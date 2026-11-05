package postgres

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.Administrator = (*Store)(nil)
var _ cityhall.Reporter = (*Store)(nil)

// Create implements Administrator.Create
func (s *Store) Create(
	ctx context.Context,
	params *cityhall.CityCreationParams,
) (int32, error) {
	queries := s.q(ctx)

	city_id, err := s._createCity(ctx, queries, params)
	if err != nil {
		return 0, err
	}

	err = s._createTiles(ctx, queries, city_id, params.Tiles)
	if err != nil {
		return 0, err
	}

	err = s._createTilesSpatialEffects(ctx, queries, city_id, params.Tiles)
	if err != nil {
		return 0, err
	}

	err = s._createBudget(ctx, queries, city_id, params.Budget)
	if err != nil {
		return 0, err
	}

	err = s._createBalance(ctx, queries, city_id)
	if err != nil {
		return 0, err
	}

	err = s._addFirstStateEntry(ctx, queries, city_id, params.Tiles)
	if err != nil {
		return 0, err
	}

	return city_id, nil
}

func (s *Store) _createCity(
	ctx context.Context,
	queries_with_tx *sqlc.Queries,
	params *cityhall.CityCreationParams,
) (int32, error) {
	response, err := queries_with_tx.CreateCity(ctx, sqlc.CreateCityParams{
		Email:     params.Email,
		Password:  string(params.Password),
		Name:      params.Name,
		MayorName: params.MayorName,
	})
	if err != nil {
		return 0, getPostgresError(err)
	}

	return response, nil
}

func (s *Store) _createTiles(
	ctx context.Context,
	queries_with_tx *sqlc.Queries,
	city_id int32,
	initial_tiles *[core.TilesCount]core.Tile,
) error {
	payload := make([]sqlc.CreateCityTilesParams, len(initial_tiles))

	for index := range initial_tiles {
		initial_tile := &initial_tiles[index]

		payload[index] = sqlc.CreateCityTilesParams{
			CityID:         city_id,
			SuperKindID:    initial_tile.SuperKind,
			KindID:         initial_tile.Kind,
			ArtifactID:     initial_tile.Artifact.Name,
			AnchorQ:        int16(initial_tile.Anchor.Q),
			AnchorR:        int16(initial_tile.Anchor.R),
			Q:              int16(initial_tile.Position.Q),
			R:              int16(initial_tile.Position.R),
			ArtifactLabel:  initial_tile.Artifact.Label,
			ArtifactStatus: initial_tile.Artifact.Status,
		}
	}

	_, err := queries_with_tx.CreateCityTiles(ctx, payload)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func (s *Store) _createTilesSpatialEffects(
	ctx context.Context,
	queries_with_tx *sqlc.Queries,
	city_id int32,
	initial_tiles *[core.TilesCount]core.Tile,
) error {
	payload := make([]sqlc.CreateCityTilesSpatialEffectsParams, len(initial_tiles))

	for index := range initial_tiles {
		initial_tile := &initial_tiles[index]

		payload[index] = sqlc.CreateCityTilesSpatialEffectsParams{
			CityID:    city_id,
			Q:         int16(initial_tile.Position.Q),
			R:         int16(initial_tile.Position.R),
			Ghg1:      initial_tile.Stats.Ghg1,
			LandValue: initial_tile.Stats.LandValue,
		}
	}

	_, err := queries_with_tx.CreateCityTilesSpatialEffects(ctx, payload)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func (s *Store) _createBudget(
	ctx context.Context,
	queries_with_tx *sqlc.Queries,
	city_id int32,
	initial_budget []core.KindBudget,
) error {
	payload := make([]sqlc.CreateCityBudgetParams, len(initial_budget))

	for index := range initial_budget {
		payload[index] = sqlc.CreateCityBudgetParams{
			CityID:      city_id,
			KindID:      initial_budget[index].KindName,
			TaxRate:     int16(initial_budget[index].TaxRate),
			FundingRate: int16(initial_budget[index].FundingRate),
		}
	}

	_, err := queries_with_tx.CreateCityBudget(ctx, payload)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func (s *Store) _createBalance(
	ctx context.Context,
	queries_with_tx *sqlc.Queries,
	city_id int32,
) error {
	err := queries_with_tx.CreateCityBalance(ctx, city_id)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func (s *Store) _addFirstStateEntry(
	ctx context.Context,
	queries_with_tx *sqlc.Queries,
	city_id int32,
	initial_tiles *[core.TilesCount]core.Tile,
) error {
	var ghg_1 int32

	for index := range initial_tiles {
		initial_tile := &initial_tiles[index]

		ghg_1 += int32(initial_tile.Stats.Ghg1)
	}

	err := queries_with_tx.CreateCityHistoryFirstEntry(
		ctx,
		sqlc.CreateCityHistoryFirstEntryParams{
			CityID:    city_id,
			Ghg1Level: ghg_1,
		},
	)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// Delete implements Administrator.Delete
func (s *Store) Delete(ctx context.Context, city_id int32) error {
	queries := s.q(ctx)

	_, err := queries.DeleteCity(ctx, city_id)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// CheckIfExists implements Administrator.CheckIfExists
func (s *Store) CheckIfExists(ctx context.Context, city_id int32) (bool, error) {
	queries := s.q(ctx)

	response, err := queries.CheckIfCityExists(ctx, city_id)
	if err != nil {
		return false, err
	}

	return response, nil
}

// GetCities implements Administrator.GetCities
func (s *Store) GetCities(ctx context.Context) ([]int32, error) {
	queries := s.q(ctx)

	response, err := queries.GetCities(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CityStats implements Reporter.CityStats
func (s *Store) CityStats(ctx context.Context, city_id int32) (core.CityStats, error) {
	queries := s.q(ctx)

	response, err := queries.GetCurrentCityState(ctx, city_id)
	if err != nil {
		return core.CityStats{}, getPostgresError((err))
	}

	return core.CityStats{
		Tick:                   uint64(response.GameTick.Int64),
		ElectricityConsumption: response.ElectricityConsumption,
		WaterConsumption:       response.WaterConsumption,
		Ghg1Level:              response.Ghg1Level,
		Ghg2Level:              response.Ghg2Level,
		Ghg3Level:              response.Ghg3Level,
		ClimateChangeLevel:     response.ClimateChangeLevel.Int32,
		ChildrenCount:          response.ChildrenCount,
		AdultsCount:            response.AdultsCount,
		ElderlyCount:           response.ElderlyCount,
		CitizensCount:          response.CitizensCount.Int32,
		ChildrenEducationLevel: response.ChildrenEducationLevel,
		AdultsEducationLevel:   response.AdultsEducationLevel,
		ElderlyEducationLevel:  response.ElderlyEducationLevel,
		CitizensEducationLevel: int16(response.CitizensEducationLevel.Int32),
		ChildrenHealthLevel:    response.ChildrenHealthLevel,
		AdultsHealthLevel:      response.AdultsHealthLevel,
		ElderlyHealthLevel:     response.ElderlyHealthLevel,
		CitizensHealthLevel:    int16(response.CitizensHealthLevel.Int32),
	}, nil
}
