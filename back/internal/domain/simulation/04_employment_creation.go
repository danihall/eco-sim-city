package simulation

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/chunks"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

var _ gameloops.Ticker = (*EmploymentLoop)(nil)

const employment_creation_loop = "employment_creation_loop"

type EmploymentLoop struct {
	*cityhall.Deputies
	Chunks []core.Chunk
}

func NewEmploymentLoop(city_hall *cityhall.Deputies) *EmploymentLoop {
	return &EmploymentLoop{
		Deputies: city_hall,
		Chunks:   chunks.FourSquares[:],
	}
}

func (e *EmploymentLoop) ID() string {
	return employment_creation_loop
}

func (e *EmploymentLoop) Tick(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) (*gameloops.Update, error) {
	err := e.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		anchor_tiles, err := e.EmploymentOffice.Seeker.AnchorTilesWithCitizenToAllocate(
			task_ctx,
			city_id,
			cityhall.AnchorTilesWithCitizenToAllocateParams{
				ArtifactStatus:         core.Step_3_DeliveredOperational,
				IndustrialName:         core.Industrial,
				PublicSystemsName:      core.PublicSystems,
				ResidentialName:        core.Residential,
				Chunk:                  chunk,
				HasAllCitizenAllocated: false,
			},
		)
		if err != nil {
			return err
		}

		if len(anchor_tiles.IndustrialOrPublicSystems) == 0 || len(anchor_tiles.Residential) == 0 {
			return nil
		}

		cap := len(anchor_tiles.IndustrialOrPublicSystems) + len(anchor_tiles.Residential)

		new_employments := make([]core.Employment, 0, cap)
		tiles_with_all_citizen_allocated := make([]core.TilePosition, 0, cap)

		for worker_index := range anchor_tiles.Residential {
			worker_tile := anchor_tiles.Residential[worker_index]
			worker_artifact := registry.Artifacts[worker_tile.Artifact.Name]
			allocable_workers_count := worker_artifact.CitizensCapacity - worker_tile.JobsCount
			worker_min_job_level := registry.SocialClasses[registry.ArtifactsRelations[worker_artifact.Name].Kind].MinJobLevel

			for employer_index := range anchor_tiles.IndustrialOrPublicSystems {
				employer_tile := &anchor_tiles.IndustrialOrPublicSystems[employer_index]
				employer_artifact := registry.Artifacts[employer_tile.Artifact.Name]
				job_offers_count := employer_artifact.CitizensCapacity - employer_tile.JobsCount

				if allocable_workers_count == 0 {
					break
				}

				if job_offers_count == 0 {
					continue
				}

				if employer_artifact.JobLevel >= worker_min_job_level {
					contracts_count := min(allocable_workers_count, job_offers_count)

					new_employments = append(new_employments, core.Employment{
						EmployerAnchor: employer_tile.Anchor,
						WorkerAnchor:   worker_tile.Anchor,
						WorkerArtifact: worker_artifact.Name,
						JobsCount:      contracts_count,
					})

					employer_tile.JobsCount += contracts_count
					allocable_workers_count -= contracts_count
					job_offers_count -= contracts_count

					if job_offers_count == 0 {
						tiles_with_all_citizen_allocated = append(
							tiles_with_all_citizen_allocated,
							employer_tile.Anchor,
						)
					}

					if allocable_workers_count == 0 {
						tiles_with_all_citizen_allocated = append(
							tiles_with_all_citizen_allocated,
							worker_tile.Anchor,
						)
					}
				}
			}
		}

		if len(new_employments) == 0 {
			return nil
		}

		err = e.EmploymentOffice.Signator.ContractsToSign(task_ctx, city_id, new_employments)
		if err != nil {
			return err
		}

		err = e.EmploymentOffice.Signator.SetAllCitizenAllocation(
			task_ctx,
			city_id,
			tiles_with_all_citizen_allocated,
			true,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
