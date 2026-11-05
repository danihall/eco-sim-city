package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// ResidentialOrIndustrialWithStatsParams is only used internally and is not exposed to the user.
type ResidentialOrIndustrialWithStatsParams struct {
	ArtifactStatus core.ArtifactStatus
	SuperKinds     [2]core.SuperKindType
	Chunk          core.Chunk
}

type TilesGetter interface {
	FromCoordinatesFullPropsWithStats(
		ctx context.Context,
		city_id int32,
		coordinates []core.TilePosition,
	) ([]core.Tile, error)
	FromAnchor(
		ctx context.Context,
		city_id int32,
		anchor core.TilePosition,
	) ([]core.Tile, error)
	SpatialStats(
		ctx context.Context,
		city_id int32,
		coordinates []core.TilePosition,
	) ([]core.TileSpatialStats, error)
	ResidentialOrIndustrialWithStats(
		ctx context.Context,
		city_id int32,
		params ResidentialOrIndustrialWithStatsParams,
	) ([]core.Tile, error)
	FromCity(
		ctx context.Context,
		city_id int32,
	) (*[core.TilesCount]core.Tile, error)
}

type TilesMapper struct{}

// Make returns a map of all tiles's given data keyed by their position. Useful for fast look-ups.
func (t TilesMapper) Make(tiles []core.Tile) map[core.TilePosition]core.Tile {
	tiles_set := make(map[core.TilePosition]core.Tile, len(tiles))

	for index := range tiles {
		tile := tiles[index]
		tiles_set[tile.Position] = tile
	}

	return tiles_set
}

type deltas interface {
	core.Pollution | core.Desirability | core.Criminality
}

type TileStatsDeltas[D deltas] struct {
	Values D
	core.TilePosition
}

func (t TileStatsDeltas[D]) Position() core.TilePosition {
	return t.TilePosition
}

type TilesUpdater interface {
	CriminalityFromDeltas(
		ctx context.Context,
		city_id int32,
		deltas_to_apply []TileStatsDeltas[core.Criminality],
	) error
	DesirabilityFromDeltas(
		ctx context.Context,
		city_id int32,
		deltas_to_apply []TileStatsDeltas[core.Desirability],
	) error
	PollutionFromDeltas(
		ctx context.Context,
		city_id int32,
		deltas_to_apply []TileStatsDeltas[core.Pollution],
	) error
}

// Cadastre deals with whatever is on any given city's tiles.
//
// I didn't know but "cadastre" is an english word.
//
// see core/tile.go
type Cadastre struct {
	Getter  TilesGetter
	Mapper  TilesMapper
	Updater TilesUpdater
}

func NewCadastre(store Store) *Cadastre {
	return &Cadastre{
		Getter:  store,
		Mapper:  TilesMapper{},
		Updater: store,
	}
}
