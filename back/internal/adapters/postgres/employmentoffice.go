package postgres

import (
	"context"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.JobTerminator = (*Store)(nil)
var _ cityhall.JobSignator = (*Store)(nil)
var _ cityhall.JobSeeker = (*Store)(nil)

// ContractsToDelete implements JobTerminator.ContractsToDelete
func (s *Store) ContractsToDelete(
	ctx context.Context,
	city_id int32,
	params cityhall.DeleteContractsParams,
) error {
	queries := s.q(ctx)

	employer_anchor_qs, employer_anchor_rs := makeQueryCoordinates(params.EmployerAnchors)
	worker_anchor_qs, worker_anchor_rs := makeQueryCoordinates(params.WorkerAnchors)

	err := queries.DeleteJobsContracts(ctx, sqlc.DeleteJobsContractsParams{
		CityID:           city_id,
		WorkersAnchorQs:  worker_anchor_qs,
		WorkersAnchorRs:  worker_anchor_rs,
		EmployerAnchorQs: employer_anchor_qs,
		EmployerAnchorRs: employer_anchor_rs,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// ContractsToSign implements JobSignator.ContractsToSign
func (s *Store) ContractsToSign(
	ctx context.Context,
	city_id int32,
	params []core.Employment,
) error {
	queries := s.q(ctx)

	worker_artifact_ids := make([]int16, len(params))
	employer_anchor_qs := make([]int16, len(params))
	employer_anchor_rs := make([]int16, len(params))
	worker_anchor_qs := make([]int16, len(params))
	worker_anchor_rs := make([]int16, len(params))
	jobs_count := make([]int16, len(params))

	for index := range params {
		employment := &params[index]

		worker_artifact_ids[index] = int16(employment.WorkerArtifact)
		employer_anchor_qs[index] = int16(employment.EmployerAnchor.Q)
		employer_anchor_rs[index] = int16(employment.EmployerAnchor.R)
		worker_anchor_qs[index] = int16(employment.WorkerAnchor.Q)
		worker_anchor_rs[index] = int16(employment.WorkerAnchor.R)
		jobs_count[index] = int16(employment.JobsCount)
	}

	err := queries.InsertJobsContracts(ctx, sqlc.InsertJobsContractsParams{
		CityID:            city_id,
		WorkerArtifactIds: worker_artifact_ids,
		EmployerAnchorQs:  employer_anchor_qs,
		EmployerAnchorRs:  employer_anchor_rs,
		WorkersAnchorQs:   worker_anchor_qs,
		WorkersAnchorRs:   worker_anchor_rs,
		JobsCount:         jobs_count,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// SetAllCitizenAllocation implements JobSignator.SetAllCitizenAllocation
func (s *Store) SetAllCitizenAllocation(
	ctx context.Context,
	city_id int32,
	anchor_tiles []core.TilePosition,
	is_fully_allocated bool,
) error {
	queries := s.q(ctx)

	anchor_qs, anchor_rs := makeQueryCoordinates(anchor_tiles)
	allocations := slices.Repeat([]bool{is_fully_allocated}, len(anchor_tiles))

	err := queries.SetAllCitizenAllocation(ctx, sqlc.SetAllCitizenAllocationParams{
		CityID:      city_id,
		AnchorQs:    anchor_qs,
		AnchorRs:    anchor_rs,
		Allocations: allocations,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// AnchorTilesWithCitizenToAllocate implements TilesGetter.AnchorTilesWithCitizenToAllocate
// Uses two SQL requests for simplicity.
func (s *Store) AnchorTilesWithCitizenToAllocate(
	ctx context.Context,
	city_id int32,
	params cityhall.AnchorTilesWithCitizenToAllocateParams,
) (cityhall.AnchorTilesWithCitizenToAllocate, error) {
	queries := s.q(ctx)

	industrial_or_public, err := queries.GetAnchorTilesIndustrialOrPublicWithCitizenToAllocate(
		ctx,
		sqlc.GetAnchorTilesIndustrialOrPublicWithCitizenToAllocateParams{
			CityID:                 city_id,
			ArtifactStatus:         params.ArtifactStatus,
			QTopLeft:               int16(params.Chunk.TopLeft.Q),
			QBottomRight:           int16(params.Chunk.BottomRight.Q),
			RTopLeft:               int16(params.Chunk.TopLeft.R),
			RBottomRight:           int16(params.Chunk.BottomRight.R),
			SuperKindIndustrial:    params.IndustrialName,
			SuperKindPublicSystems: params.PublicSystemsName,
			AtCapacity:             params.HasAllCitizenAllocated,
		},
	)
	if err != nil {
		return cityhall.AnchorTilesWithCitizenToAllocate{}, getPostgresError(err)
	}

	residential, err := queries.GetAnchorTilesResidentialWithCitizenToAllocate(
		ctx,
		sqlc.GetAnchorTilesResidentialWithCitizenToAllocateParams{
			CityID:               city_id,
			ArtifactStatus:       params.ArtifactStatus,
			QTopLeft:             int16(params.Chunk.TopLeft.Q),
			QBottomRight:         int16(params.Chunk.BottomRight.Q),
			RTopLeft:             int16(params.Chunk.TopLeft.R),
			RBottomRight:         int16(params.Chunk.BottomRight.R),
			SuperKindResidential: params.ResidentialName,
			AtCapacity:           params.HasAllCitizenAllocated,
		},
	)
	if err != nil {
		return cityhall.AnchorTilesWithCitizenToAllocate{}, getPostgresError(err)
	}

	effective_response := cityhall.AnchorTilesWithCitizenToAllocate{
		IndustrialOrPublicSystems: make([]core.TileWithJobsCount, len(industrial_or_public)),
		Residential:               make([]core.TileWithJobsCount, len(residential)),
	}

	for index := range industrial_or_public {
		tile := &industrial_or_public[index]

		effective_response.IndustrialOrPublicSystems[index] = core.TileWithJobsCount{
			Tile: core.Tile{
				SuperKind: tile.SuperKindID,
				Artifact: core.ArtifactInstance{
					Name: tile.ArtifactID,
				},
				Position: core.TilePosition{
					Q: uint8(tile.Q),
					R: uint8(tile.R),
				},
				Anchor: core.TilePosition{
					Q: uint8(tile.AnchorQ),
					R: uint8(tile.AnchorR),
				},
			},
			JobsCount: uint8(tile.JobsCount),
		}
	}

	for index := range residential {
		tile := &residential[index]

		effective_response.Residential[index] = core.TileWithJobsCount{
			Tile: core.Tile{
				SuperKind: tile.SuperKindID,
				Artifact: core.ArtifactInstance{
					Name: tile.ArtifactID,
				},
				Position: core.TilePosition{
					Q: uint8(tile.Q),
					R: uint8(tile.R),
				},
				Anchor: core.TilePosition{
					Q: uint8(tile.AnchorQ),
					R: uint8(tile.AnchorR),
				},
			},
			JobsCount: uint8(tile.JobsCount),
		}
	}

	return effective_response, nil
}

func (s *Store) FromWorkerAnchorChunk(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) ([]core.Employment, error) {
	queries := s.q(ctx)

	response, err := queries.GetWorkerJobContractsFromChunk(
		ctx,
		sqlc.GetWorkerJobContractsFromChunkParams{
			CityID:       city_id,
			QTopLeft:     int16(chunk.TopLeft.Q),
			QBottomRight: int16(chunk.BottomRight.Q),
			RTopLeft:     int16(chunk.TopLeft.R),
			RBottomRight: int16(chunk.BottomRight.R),
		},
	)

	if err != nil {
		return nil, getPostgresError(err)
	}

	worker_employments := make([]core.Employment, len(response))

	for index, employment := range response {
		worker_employments[index] = core.Employment{
			EmployerAnchor: core.TilePosition{
				Q: uint8(employment.EmployerAnchorQ),
				R: uint8(employment.EmployerAnchorR),
			},
			WorkerAnchor: core.TilePosition{
				Q: uint8(employment.WorkersAnchorQ),
				R: uint8(employment.WorkersAnchorR),
			},
			WorkerArtifact: employment.WorkerArtifactID,
			JobsCount:      uint8(employment.JobsCount),
		}
	}

	return worker_employments, nil
}
