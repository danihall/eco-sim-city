package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

type DeleteContractsParams struct {
	EmployerAnchors []core.TilePosition
	WorkerAnchors   []core.TilePosition
}

type JobTerminator interface {
	// ContractsToDelete also deletes the associated traffic tiles.
	ContractsToDelete(ctx context.Context, city_id int32, params DeleteContractsParams) error
}

type JobSignator interface {
	ContractsToSign(
		ctx context.Context,
		city_id int32,
		params []core.Employment,
	) error
	SetAllCitizenAllocation(
		ctx context.Context,
		city_id int32,
		anchor_tiles []core.TilePosition,
		is_fully_allocated bool,
	) error
}

type AnchorTilesWithCitizenToAllocateParams struct {
	ArtifactStatus         core.ArtifactStatus
	IndustrialName         core.SuperKindType
	PublicSystemsName      core.SuperKindType
	ResidentialName        core.SuperKindType
	Chunk                  core.Chunk
	HasAllCitizenAllocated bool
}

type AnchorTilesWithCitizenToAllocate struct {
	IndustrialOrPublicSystems []core.TileWithJobsCount
	Residential               []core.TileWithJobsCount
}

type JobSeeker interface {
	AnchorTilesWithCitizenToAllocate(
		ctx context.Context,
		city_id int32,
		params AnchorTilesWithCitizenToAllocateParams,
	) (AnchorTilesWithCitizenToAllocate, error)
	FromWorkerAnchorChunk(
		ctx context.Context,
		city_id int32,
		chunk core.Chunk,
	) ([]core.Employment, error)
}

type EmploymentOffice struct {
	Terminator JobTerminator
	Signator   JobSignator
	Seeker     JobSeeker
}

func NewEmploymentOffice(store Store) *EmploymentOffice {
	return &EmploymentOffice{
		Terminator: store,
		Signator:   store,
		Seeker:     store,
	}
}
