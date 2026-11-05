package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

type CityCreationParams struct {
	Email     string
	Password  []byte
	Name      string
	MayorName string
	Tiles     *[core.TilesCount]core.Tile
	Budget    []core.KindBudget
}

type Administrator interface {
	Create(ctx context.Context, params *CityCreationParams) (int32, error)
	Delete(ctx context.Context, city_id int32) error
	CheckIfExists(ctx context.Context, city_id int32) (bool, error)
	GetCities(ctx context.Context) ([]int32, error)
}

type Reporter interface {
	CityStats(ctx context.Context, city_id int32) (core.CityStats, error)
}

// Cabinet is the highest authority in the city hall. It has the power of creating or deleting a city.
// When creating a new city, initial tiles are prepared.
//
// see core/data/base_tiles.go
//
// see core/city.go
type Cabinet struct {
	Administrator Administrator
	Reporter      Reporter
}

func NewCabinet(store Store) *Cabinet {
	return &Cabinet{
		Administrator: store,
		Reporter:      store,
	}
}

// TODO add methods in Reporter interface (ElectrictyConsumption & WaterConsumption)
