package postgres

import (
	"context"
	"fmt"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.TilesGetter = (*Store)(nil)
var _ cityhall.TilesUpdater = (*Store)(nil)

// FromCoordinatesFullPropsWithStats implements TilesGetter.FromCoordinatesFullPropsWithStats
func (s *Store) FromCoordinatesFullPropsWithStats(
	ctx context.Context,
	city_id int32,
	coordinates []core.TilePosition,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	qs, rs := makeQueryCoordinates(coordinates)

	response, err := queries.GetTilesAllPropsWithStats(
		ctx,
		sqlc.GetTilesAllPropsWithStatsParams{
			CityID: city_id,
			Qs:     qs,
			Rs:     rs,
		},
	)
	if err != nil {
		return nil, getPostgresError(err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf(
			"no tiles found. city_id=%v at coordinates q=%v, r=%v",
			city_id,
			qs,
			rs,
		)
	}

	effective_tiles := make([]core.Tile, len(response))

	for index := range response {
		tile := &response[index]

		effective_tiles[index] = core.NewTile(
			tile.SuperKindID,
			tile.KindID,
			core.ArtifactInstance{
				Name:   tile.ArtifactID,
				Label:  tile.ArtifactLabel,
				Status: tile.ArtifactStatus,
			},
			core.TileSpatialStats{
				Pollution: core.Pollution{
					Ghg1: tile.Ghg1,
				},
				Desirability: core.Desirability{
					LandValue: tile.LandValue,
				},
				Criminality: core.Criminality{
					CrimeGeneric: tile.CrimeGeneric,
					CrimeDrug:    tile.CrimeDrug,
					CrimeVss:     tile.CrimeVss,
				},
			},
			core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
			core.TilePosition{
				Q: uint8(tile.AnchorQ),
				R: uint8(tile.AnchorR),
			},
			tile.HasAllCitizenAllocated,
		)
	}

	return effective_tiles, nil
}

// FromAnchor implements TilesGetter.FromAnchor
func (s *Store) FromAnchor(
	ctx context.Context,
	city_id int32,
	anchor core.TilePosition,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	response, err := queries.GetTilesByAnchor(ctx, sqlc.GetTilesByAnchorParams{
		CityID: city_id,
		Q:      int16(anchor.Q),
		R:      int16(anchor.R),
	})
	if err != nil {
		return nil, getPostgresError(err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf(
			"no tiles found. city_id=%v at coordinates q=%v, r=%v",
			city_id,
			anchor.Q,
			anchor.R,
		)
	}

	effective_tiles := make([]core.Tile, len(response))

	for index := range response {
		tile := &response[index]
		effective_tiles[index] = core.Tile{
			SuperKind: tile.SuperKindID,
			Kind:      tile.KindID,
			Artifact: core.ArtifactInstance{
				Name:   tile.ArtifactID,
				Label:  tile.ArtifactLabel,
				Status: tile.ArtifactStatus,
			},
			Position: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
			Anchor: core.TilePosition{
				Q: uint8(tile.AnchorQ),
				R: uint8(tile.AnchorR),
			},
		}
	}

	return effective_tiles, nil
}

// SpatialStats implements TilesGetter.SpatialStats
func (s *Store) SpatialStats(
	ctx context.Context,
	city_id int32,
	coordinates []core.TilePosition,
) ([]core.TileSpatialStats, error) {
	queries := s.q(ctx)

	qs, rs := makeQueryCoordinates(coordinates)

	response, err := queries.GetTilesSpatialEffects(
		ctx,
		sqlc.GetTilesSpatialEffectsParams{
			CityID: city_id,
			Qs:     qs,
			Rs:     rs,
		},
	)
	if err != nil {
		return nil, getPostgresError(err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf(
			"no tiles found. city_id=%v at coordinates q=%v, r=%v",
			city_id,
			qs,
			rs,
		)
	}

	effective_spatial_values := make([]core.TileSpatialStats, len(response))
	for index := range response {
		tile := &response[index]
		effective_spatial_values[index] = core.TileSpatialStats{
			Pollution: core.Pollution{
				Ghg1: tile.Ghg1,
			},
			Desirability: core.Desirability{
				LandValue: tile.LandValue,
			},
			Criminality: core.Criminality{
				CrimeGeneric: tile.CrimeGeneric,
				CrimeDrug:    tile.CrimeDrug,
				CrimeVss:     tile.CrimeVss,
			},
		}
	}

	return effective_spatial_values, nil
}

// ResidentialOrIndustrialWithStats implements TilesGetter.ResidentialOrIndustrialWithStats
func (s *Store) ResidentialOrIndustrialWithStats(
	ctx context.Context,
	city_id int32,
	params cityhall.ResidentialOrIndustrialWithStatsParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	response, err := queries.GetTilesResidentialOrIndustrialWithEffects(
		ctx,
		sqlc.GetTilesResidentialOrIndustrialWithEffectsParams{
			CityID:         city_id,
			ArtifactStatus: params.ArtifactStatus,
			QTopLeft:       int16(params.Chunk.TopLeft.Q),
			QBottomRight:   int16(params.Chunk.BottomRight.Q),
			RTopLeft:       int16(params.Chunk.TopLeft.R),
			RBottomRight:   int16(params.Chunk.BottomRight.R),
			SuperKindIds: []int16{
				int16(params.SuperKinds[0]),
				int16(params.SuperKinds[1]),
			},
		},
	)
	if err != nil {
		return nil, getPostgresError(err)
	}

	if len(response) == 0 {
		return nil, nil
	}

	effective_response := make([]core.Tile, len(response))
	for index := range response {
		tile := &response[index]

		effective_response[index] = core.Tile{
			SuperKind: tile.SuperKindID,
			Kind:      tile.KindID,
			Artifact: core.ArtifactInstance{
				Name: tile.ArtifactID,
			},
			Stats: core.TileSpatialStats{
				Pollution: core.Pollution{
					Ghg1: tile.Ghg1,
				},
				Desirability: core.Desirability{
					LandValue: tile.LandValue,
				},
				Criminality: core.Criminality{
					CrimeGeneric: tile.CrimeGeneric,
					CrimeDrug:    tile.CrimeDrug,
					CrimeVss:     tile.CrimeVss,
				},
			},
			Position: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
			Anchor: core.TilePosition{
				Q: uint8(tile.AnchorQ),
				R: uint8(tile.AnchorR),
			},
			HasStats: true,
		}
	}

	return effective_response, nil
}

// FromCity implements TilesGetter.FromCity
func (s *Store) FromCity(
	ctx context.Context,
	city_id int32,
) (*[core.TilesCount]core.Tile, error) {
	queries := s.q(ctx)

	response, err := queries.GetTilesFromCity(ctx, city_id)
	if err != nil {
		return nil, getPostgresError(err)
	}

	response_length := len(response)

	if response_length != int(core.TilesCount) {
		return nil, fmt.Errorf("abnormal tiles count for city, got: %d", response_length)
	}

	var city_tiles [core.TilesCount]core.Tile
	for index := range response {
		tile := &response[index]

		city_tiles[index] = core.Tile{
			Artifact: core.ArtifactInstance{
				Name:   tile.ArtifactID,
				Status: tile.ArtifactStatus,
			},
			Position: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
			Anchor: core.TilePosition{
				Q: uint8(tile.AnchorQ),
				R: uint8(tile.AnchorR),
			},
			SuperKind: tile.SuperKindID,
		}
	}

	return &city_tiles, nil
}

// TilesCriminality implements TilesDeltasUpdater.TilesCriminality
func (s *Store) CriminalityFromDeltas(
	ctx context.Context,
	city_id int32,
	deltas_to_apply []cityhall.TileStatsDeltas[core.Criminality],
) error {
	queries := s.q(ctx)

	slices.SortFunc(deltas_to_apply, byTilePosition)

	params := sqlc.UpdateTilesCriminalityParams{
		CityID:             city_id,
		Qs:                 make([]int16, len(deltas_to_apply)),
		Rs:                 make([]int16, len(deltas_to_apply)),
		CrimeGenericDeltas: make([]int16, len(deltas_to_apply)),
		CrimeDrugDeltas:    make([]int16, len(deltas_to_apply)),
		CrimeVssDeltas:     make([]int16, len(deltas_to_apply)),
	}

	for index := range deltas_to_apply {
		delta := &deltas_to_apply[index]

		params.Qs[index] = int16(delta.Q)
		params.Rs[index] = int16(delta.R)
		params.CrimeGenericDeltas[index] = delta.Values.CrimeGeneric
		params.CrimeDrugDeltas[index] = delta.Values.CrimeDrug
		params.CrimeVssDeltas[index] = delta.Values.CrimeVss
	}

	err := queries.UpdateTilesCriminality(ctx, params)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func (s *Store) DesirabilityFromDeltas(
	ctx context.Context,
	city_id int32,
	deltas_to_apply []cityhall.TileStatsDeltas[core.Desirability],
) error {
	queries := s.q(ctx)

	slices.SortFunc(deltas_to_apply, byTilePosition)

	params := sqlc.UpdateTilesDesirabilityParams{
		CityID:          city_id,
		Qs:              make([]int16, len(deltas_to_apply)),
		Rs:              make([]int16, len(deltas_to_apply)),
		LandValueDeltas: make([]int16, len(deltas_to_apply)),
	}

	for index := range deltas_to_apply {
		delta := &deltas_to_apply[index]

		params.Qs[index] = int16(delta.Q)
		params.Rs[index] = int16(delta.R)
		params.LandValueDeltas[index] = delta.Values.LandValue
	}

	err := queries.UpdateTilesDesirability(ctx, params)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func (s *Store) PollutionFromDeltas(
	ctx context.Context,
	city_id int32,
	deltas_to_apply []cityhall.TileStatsDeltas[core.Pollution],
) error {
	queries := s.q(ctx)

	slices.SortFunc(deltas_to_apply, byTilePosition)

	params := sqlc.UpdateTilesPollutionParams{
		CityID:     city_id,
		Qs:         make([]int16, len(deltas_to_apply)),
		Rs:         make([]int16, len(deltas_to_apply)),
		Ghg1Deltas: make([]int16, len(deltas_to_apply)),
	}

	for index := range deltas_to_apply {
		delta := &deltas_to_apply[index]

		params.Qs[index] = int16(delta.Q)
		params.Rs[index] = int16(delta.R)
		params.Ghg1Deltas[index] = delta.Values.Ghg1
	}

	err := queries.UpdateTilesPollution(ctx, params)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}
