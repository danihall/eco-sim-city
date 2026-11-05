package simulation

import (
	"context"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/chunks"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

var _ gameloops.Ticker = (*ArtifactsAbandonmentLoop)(nil)

const artifacts_abandonment_loop = "artifacts_abandonment_loop"
const max_abandoned_artifacts_per_tick = 8
const abandoned_label = "abandonné"

// ArtifactsAbandonmentLoop simulates artifacts that get abandonned
// due to their requirements not being met.
type ArtifactsAbandonmentLoop struct {
	*cityhall.Deputies
	eventBus *eventbus.EventBus
	Chunks   []core.Chunk
}

func NewArtifactAbandonmentLoop(
	city_hall *cityhall.Deputies,
	event_bus *eventbus.EventBus,
) *ArtifactsAbandonmentLoop {
	return &ArtifactsAbandonmentLoop{
		Deputies: city_hall,
		eventBus: event_bus,
		Chunks:   chunks.FourBarsHorizontalAndFourBarsVertical[:],
	}
}

func (a *ArtifactsAbandonmentLoop) ID() string {
	return artifacts_abandonment_loop
}

func (a *ArtifactsAbandonmentLoop) Tick(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) (*gameloops.Update, error) {
	city, err := a.GetCityFullState(ctx, city_id)
	if err != nil {
		return nil, err
	}

	var abandonned_tiles []core.Tile

	err = a.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		tiles, err := a.Cadastre.Getter.ResidentialOrIndustrialWithStats(
			task_ctx,
			city_id,
			cityhall.ResidentialOrIndustrialWithStatsParams{
				ArtifactStatus: core.Step_3_DeliveredOperational,
				SuperKinds:     [2]core.SuperKindType{core.Industrial, core.Residential},
				Chunk:          chunk,
			},
		)
		if err != nil {
			return err
		}
		if len(tiles) == 0 {
			return nil
		}

		tiles_set := a.Cadastre.Mapper.Make(tiles)
		coordinates_buffer := make([]core.TilePosition, 64)
		tiles_buffer := make([]core.Tile, 64)
		anchors_to_abandon_set := make(
			map[core.TilePosition]core.ArtifactType,
			max_abandoned_artifacts_per_tick,
		)

	tiles_loop:
		for i := range tiles {
			tile := &tiles[i]
			_, exists := anchors_to_abandon_set[tile.Anchor]
			if exists {
				continue
			}

			artifact := registry.Artifacts[tile.Artifact.Name]

			err = a.Permits.Validator.CheckGlobalRequirements(city, artifact)
			if err != nil {
				anchors_to_abandon_set[tile.Anchor] = artifact.Name
				if len(anchors_to_abandon_set) == max_abandoned_artifacts_per_tick {
					break
				}
				continue
			}

			coordinates, err := a.Topology.Maker.FromFixSizedArtifact(
				artifact,
				tiles[i].Anchor,
				coordinates_buffer,
			)
			if err != nil {
				return fmt.Errorf("in %s: FromFixSizedArtifact failed", artifacts_abandonment_loop)
			}

			target_tiles := tiles_buffer[:len(coordinates)]
			for k := range coordinates {
				target_tile, ok := tiles_set[coordinates[k]]
				if !ok {
					continue tiles_loop
				}
				target_tiles[k] = target_tile
			}

			err = a.Permits.Validator.CheckSpatialRequirements(artifact, target_tiles)
			if err != nil {
				anchors_to_abandon_set[tile.Anchor] = artifact.Name
				if len(anchors_to_abandon_set) == max_abandoned_artifacts_per_tick {
					break
				}
			}
		}

		if len(anchors_to_abandon_set) == 0 {
			return nil
		}

		anchors_to_abandon_slice := make([]core.TilePosition, 0, len(anchors_to_abandon_set))

		for anchor := range anchors_to_abandon_set {
			anchors_to_abandon_slice = append(anchors_to_abandon_slice, anchor)
		}

		abandonned_tiles, err = a.Construction.Vacator.AbandonArtifacts(
			task_ctx,
			city_id,
			cityhall.VacateParams{
				Anchors:                 anchors_to_abandon_slice,
				Label:                   abandoned_label,
				Status:                  core.Step_1_AwaitingAssignment,
				HasAllCitizensAllocated: false,
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(abandonned_tiles) == 0 {
		return nil, nil
	}

	old_tiles := make([]core.Tile, len(abandonned_tiles))
	copy(old_tiles, abandonned_tiles)
	for index := range old_tiles {
		old_tiles[index].Artifact.Status = core.Step_3_DeliveredOperational
	}

	a.eventBus.Publish(
		city_id,
		eventbus.NewTilesChangeData(old_tiles, abandonned_tiles, city.EnactedOrdinances),
	)

	type debug struct {
		Artifact core.ArtifactType
		Position core.TilePosition
	}

	debug_abandonned_tiles := make([]debug, 0, len(abandonned_tiles))
	for _, item := range abandonned_tiles {

		debug_abandonned_tiles = append(debug_abandonned_tiles, debug{
			Artifact: item.Artifact.Name,
			Position: item.Position,
		})
	}

	return &gameloops.Update{
		GameLoopName: artifacts_abandonment_loop,
		Body:         abandonned_tiles,
	}, nil
}
