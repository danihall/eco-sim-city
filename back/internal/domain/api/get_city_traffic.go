package api

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

type CityTrafficArrays struct {
	Foot    *[core.TilesCount]uint32
	Car     *[core.TilesCount]uint32
	Bus     *[core.TilesCount]uint32
	Bicycle *[core.TilesCount]uint32
	Train   *[core.TilesCount]uint32
	Subway  *[core.TilesCount]uint32
}

func GetCityTraffic(
	ctx context.Context,
	d *cityhall.Deputies,
	city_id int32,
) (CityTrafficArrays, error) {
	var foot [core.TilesCount]uint32
	var car [core.TilesCount]uint32
	var bus [core.TilesCount]uint32
	var bicycle [core.TilesCount]uint32
	var train [core.TilesCount]uint32
	var subway [core.TilesCount]uint32

	traffic_tiles, err := d.Transportation.Getter.CityTraffic(ctx, city_id)
	if err != nil {
		return CityTrafficArrays{}, err
	}

	for _, tile := range traffic_tiles {
		tile_index := int(tile.Q) + (int(tile.R) * 128)

		foot[tile_index] = tile.FootTotal
		car[tile_index] = tile.CarTotal
		bus[tile_index] = tile.BusTotal
		bicycle[tile_index] = tile.BicycleTotal
		train[tile_index] = tile.TrainTotal
		subway[tile_index] = tile.SubwayTotal
	}

	return CityTrafficArrays{
		Foot:    &foot,
		Car:     &car,
		Bus:     &bus,
		Bicycle: &bicycle,
		Train:   &train,
		Subway:  &subway,
	}, nil
}
