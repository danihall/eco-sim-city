package cityhall

import (
	"context"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
)

type CommuteKey struct {
	WorkerAnchor   core.TilePosition
	EmployerAnchor core.TilePosition
}

type Stage struct {
	core.TilePosition
	core.TravelCostEffect
}

type TileTraffic struct {
	FootTotal    uint32
	CarTotal     uint32
	BusTotal     uint32
	BicycleTotal uint32
	TrainTotal   uint32
	SubwayTotal  uint32
	core.TilePosition
}

type TransportsGetter interface {
	// CityTransportNetwork returns tiles that are part of the transport network's city.
	// Each tile in the slice has the TilePosition and the ArtifactType fields.
	// Transport network tiles should have ArtifactType included in registry.TransportNetworkArtifacts:
	//
	// 	- core.Road
	// 	- core.Rail
	// 	- core.SubwayStation
	// 	- core.BusStop
	//  - core.TrainStation
	// 	- core.PedestrianBicycleWay
	CityNetwork(
		ctx context.Context,
		city_id int32,
		traversible_artifacts []core.ArtifactType,
	) ([]core.Tile, error)
	// CityTraffic returns all traffic for a given transport mode.
	CityTraffic(
		ctx context.Context,
		city_id int32,
	) ([]TileTraffic, error)
}

type TransportsMapper struct{}

// MakeNetwork returns a map of transport network tiles keyed by their position,
// and a slice of subway sation positions.
func (t TransportsMapper) MakeNetwork(
	tiles []core.Tile,
) (map[core.TilePosition]core.Tile, []core.TilePosition) {
	tiles_set := make(map[core.TilePosition]core.Tile, len(tiles))
	subway_tiles := make([]core.TilePosition, 0, len(tiles))

	for index := range tiles {
		tile := tiles[index]

		if !slices.Contains(registry.TransportNetworkArtifacts, tile.Artifact.Name) {
			continue
		}

		tiles_set[tile.Position] = tile

		if tile.Artifact.Name == core.SubwayStation {
			subway_tiles = append(subway_tiles, tile.Position)
		}
	}

	return tiles_set, subway_tiles
}

type TransportsUpdater interface {
	ResetCityTrafficFromChunk(
		ctx context.Context,
		city_id int32,
		chunk core.Chunk,
		commutings map[CommuteKey][]Stage,
	) error
}

type Transportation struct {
	Getter  TransportsGetter
	Mapper  TransportsMapper
	Updater TransportsUpdater
}

func NewTransportation(store Store) *Transportation {
	return &Transportation{
		Getter:  store,
		Mapper:  TransportsMapper{},
		Updater: store,
	}
}
