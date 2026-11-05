package simulation

import (
	"context"
	"math/rand/v2"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/chunks"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

var _ gameloops.Ticker = (*ArtifactsPlacingLoop)(nil)

const artifacts_placing_loop = "artifacts_placing_loop"
const max_new_artifacts_per_tick = 8

// ArtifactsPlacingLoop simulates artifacts's construction start.
// Afterwards, ArtifactsCompletionLoop takes care of finishing construction
// of placed artifacts.
type ArtifactsPlacingLoop struct {
	*cityhall.Deputies
	candidates []*core.Artifact
	Chunks     []core.Chunk
}

func NewArtifactPlacingLoop(city_hall *cityhall.Deputies) *ArtifactsPlacingLoop {
	candidates := makeValidCandidatesArtifacts()

	return &ArtifactsPlacingLoop{
		Deputies:   city_hall,
		candidates: candidates,
		Chunks:     chunks.FourBarsHorizontalAndFourBarsVertical[:],
	}
}

func (a *ArtifactsPlacingLoop) ID() string {
	return artifacts_placing_loop
}

func (a *ArtifactsPlacingLoop) Tick(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) (*gameloops.Update, error) {
	city, err := a.GetCityFullState(ctx, city_id)
	if err != nil {
		return nil, err
	}

	filtered_candidates := a.filterValidCandidates(city)
	if len(filtered_candidates) == 0 {
		return nil, nil
	}

	selected_candidates := make([]cityhall.BuildParams, 0, max_new_artifacts_per_tick)

	err = a.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		tiles, err := a.Cadastre.Getter.ResidentialOrIndustrialWithStats(
			task_ctx,
			city_id,
			cityhall.ResidentialOrIndustrialWithStatsParams{
				ArtifactStatus: core.Step_1_AwaitingAssignment,
				SuperKinds:     [2]core.SuperKindType{core.Industrial, core.Residential},
				Chunk:          chunk,
			},
		)
		if err != nil {
			return err
		}
		if tiles == nil {
			return nil
		}

		rand.Shuffle(len(filtered_candidates), func(i, j int) {
			filtered_candidates[i], filtered_candidates[j] = filtered_candidates[j], filtered_candidates[i]
		})

		// TODO use the same logic from artifacts_abandonment_loop, probably more performant.
		// only read from tiles_set + use a selected_candidates_set. then create slice params if only selected_candidates_set has members before making query.
		tiles_set := a.Cadastre.Mapper.Make(tiles)
		coordinates_buffer := make([]core.TilePosition, 64)
		tiles_buffer := make([]core.Tile, 64)

	tiles_loop:
		for i := range tiles {
			_, ok := tiles_set[tiles[i].Position]
			if !ok {
				continue
			}

			// FIXME this condition should not be needed
			if tiles[i].SuperKind != core.Industrial && tiles[i].SuperKind != core.Residential {
				continue
			}

		candidates_loop:
			for j := range filtered_candidates {
				candidate := filtered_candidates[j]
				candidate_relations := registry.ArtifactsRelations[candidate.Name]

				if tiles[i].SuperKind != candidate_relations.SuperKind {
					continue
				}

				coordinates, err := a.Topology.Maker.FromFixSizedArtifact(
					candidate,
					tiles[i].Position,
					coordinates_buffer,
				)
				if err != nil {
					continue
				}

				target_tiles := tiles_buffer[:len(coordinates)]
				for k := range coordinates {
					tile, ok := tiles_set[coordinates[k]]
					if !ok {
						continue candidates_loop
					}
					target_tiles[k] = tile
				}

				err = a.Permits.Validator.CheckSpatialRequirements(candidate, target_tiles)
				if err != nil {
					continue
				}

				selected_candidates = append(selected_candidates, cityhall.BuildParams{
					Anchors: slices.Repeat([]core.TilePosition{target_tiles[0].Position}, len(coordinates)),
					Artifact: core.ArtifactInstance{
						Name:   candidate.Name,
						Label:  candidate.Labels[rand.IntN(len(candidate.Labels))],
						Status: core.Step_2_AssignedAwaitingDelivery,
					},
					Coordinates:   slices.Clone(coordinates),
					SuperKindType: candidate_relations.SuperKind,
					KindType:      candidate_relations.Kind,
				})

				if len(selected_candidates) == max_new_artifacts_per_tick {
					break tiles_loop
				}

				for l := range coordinates {
					delete(tiles_set, coordinates[l])
				}

				// FIXME avoid delete, add a field in this slice and set it to true, then use a condition to check if already selected
				filtered_candidates = slices.Delete(filtered_candidates, j, j+1)
				break candidates_loop
			}
		}

		if len(selected_candidates) == 0 {
			return nil
		}

		return a.Construction.Builder.MultipleArtifacts(task_ctx, city_id, selected_candidates)
	})
	if err != nil {
		return nil, err
	}

	if len(selected_candidates) == 0 {
		return nil, nil
	}

	return &gameloops.Update{
		GameLoopName: artifacts_placing_loop,
		Body:         selected_candidates,
	}, nil
}

func (a *ArtifactsPlacingLoop) filterValidCandidates(
	city cityhall.CityFullState,
) []*core.Artifact {
	filtered_candidates := make([]*core.Artifact, 0, len(a.candidates))

	for index := range a.candidates {
		candidate := a.candidates[index]

		err := a.Permits.Validator.CheckUniqueArtifact(city, candidate)
		if err != nil {
			continue
		}

		err = a.Permits.Validator.CheckGlobalRequirements(city, candidate)
		if err != nil {
			continue
		}

		filtered_candidates = append(filtered_candidates, candidate)
	}

	return filtered_candidates
}

// makeValidCandidatesArtifacts returns a slice of `Industrial` or `Residential` artifacts.
func makeValidCandidatesArtifacts() []*core.Artifact {
	candidates_artifacts := make([]*core.Artifact, 0, len(registry.Artifacts))

	for name, artifact := range registry.Artifacts {
		artifact_super_kind := registry.ArtifactsRelations[name].SuperKind
		if (artifact_super_kind == core.Residential || artifact_super_kind == core.Industrial) &&
			!artifact.HasFreeSize() {
			candidates_artifacts = append(candidates_artifacts, artifact)
		}
	}

	return candidates_artifacts
}
