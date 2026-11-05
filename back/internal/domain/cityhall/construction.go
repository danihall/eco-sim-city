package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

type BuildParams struct {
	Anchors       []core.TilePosition
	Coordinates   []core.TilePosition
	Artifact      core.ArtifactInstance
	SuperKindType core.SuperKindType
	KindType      core.KindType
}

type PathwayParams struct {
	Anchors       []core.TilePosition
	Coordinates   []core.TilePosition
	Artifact      core.ArtifactInstance
	SuperKindType core.SuperKindType
	KindType      core.KindType
}

type Builder interface {
	SingleArtifact(
		ctx context.Context,
		city_id int32,
		params BuildParams,
	) ([]core.Tile, error)
	MultipleArtifacts(
		ctx context.Context,
		city_id int32,
		params []BuildParams,
	) error
	Pathway(
		ctx context.Context,
		city_id int32,
		params PathwayParams,
	) ([]core.Tile, error)
}

type Completer interface {
	AllFromChunk(
		ctx context.Context,
		city_id int32,
		params struct {
			CurrentArtifactStatus core.ArtifactStatus
			NewArtifactStatus     core.ArtifactStatus
			Chunk                 core.Chunk
		},
	) ([]core.Tile, error)
}

type DestroyParams struct {
	ReplacementArtifact     core.ArtifactInstance
	Anchor                  core.TilePosition
	HasAllCitizensAllocated bool
}

type DestroyPublicSystemsParams struct {
	DestroyParams
	ReplacementSuperKindType core.SuperKindType
	ReplacementKindType      core.KindType
}

type Destroyer interface {
	// Destroy takes `BuildParams` as params because destroying an artifact
	// is actually like building an artifact from a data standpoint.
	// `Coordinates`, `SuperKindType` and `KindType` are ignored for this request.
	Artifact(ctx context.Context, city_id int32, destroy_params DestroyParams) ([]core.Tile, error)
	// DestroyPublicSystemsArtifact takes `BuildParams` as params because destroying a public_systems artifact
	// is actually like building an artifact from a data standpoint.
	PublicSystemsArtifact(
		ctx context.Context,
		city_id int32,
		destroy_params DestroyPublicSystemsParams,
	) ([]core.Tile, error)
}

type VacateParams struct {
	Anchors                 []core.TilePosition
	Label                   string
	Status                  core.ArtifactStatus
	HasAllCitizensAllocated bool
}

type Vacator interface {
	AbandonArtifacts(
		ctx context.Context,
		city_id int32,
		vacate_params VacateParams,
	) ([]core.Tile, error)
}

// Construction groups behaviors of **adding** or **removing** an artifact from a given tile.
// An artifact's width & height must be taken into account.
type Construction struct {
	Builder   Builder
	Completer Completer
	Destroyer Destroyer
	Vacator   Vacator
}

func NewConstruction(store Store) *Construction {
	return &Construction{
		Builder:   store,
		Destroyer: store,
		Completer: store,
		Vacator:   store,
	}
}
