package postgres

import (
	"context"
	"log/slog"
	"maps"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.TransportsGetter = (*Store)(nil)
var _ cityhall.TransportsUpdater = (*Store)(nil)

// CityNetwork implements TransportsGetter.CityNetwork
func (s *Store) CityNetwork(
	ctx context.Context,
	city_id int32,
	traversible_artifacts []core.ArtifactType,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	artifact_ids := make([]int16, len(traversible_artifacts))

	for index, artifact_type := range traversible_artifacts {
		artifact_ids[index] = int16(artifact_type)
	}

	response, err := queries.GetTransportNetworkTiles(ctx, sqlc.GetTransportNetworkTilesParams{
		CityID:      city_id,
		ArtifactIds: artifact_ids,
	})
	if err != nil {
		return nil, getPostgresError(err)
	}

	effective_tiles := make([]core.Tile, len(response))

	for index, tile := range response {
		effective_tiles[index] = core.Tile{
			Position: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
			Artifact: core.ArtifactInstance{
				Name: tile.ArtifactID,
			},
		}
	}

	return effective_tiles, nil
}

// CityTraffic implements TransportsGetter.CityTraffic
func (s *Store) CityTraffic(
	ctx context.Context,
	city_id int32,
) ([]cityhall.TileTraffic, error) {
	queries := s.q(ctx)

	response, err := queries.GetCityTraffic(ctx, city_id)
	if err != nil {
		return nil, getPostgresError(err)
	}

	effective_traffic := make([]cityhall.TileTraffic, len(response))

	for index, tile := range response {
		effective_traffic[index] = cityhall.TileTraffic{
			FootTotal:    uint32(tile.FootValue),
			CarTotal:     uint32(tile.CarValue),
			BusTotal:     uint32(tile.BusValue),
			BicycleTotal: uint32(tile.BicycleValue),
			TrainTotal:   uint32(tile.TrainValue),
			SubwayTotal:  uint32(tile.SubwayValue),
			TilePosition: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
		}
	}

	return effective_traffic, nil
}

// ResetCityTrafficFromChunk implements TransportsUpdater.ResetCityTrafficFromChunk
func (s *Store) ResetCityTrafficFromChunk(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
	commutings map[cityhall.CommuteKey][]cityhall.Stage,
) error {
	queries := s.q(ctx)

	q_top_left := int16(chunk.TopLeft.Q)
	q_bottom_right := int16(chunk.BottomRight.Q)
	r_top_left := int16(chunk.TopLeft.R)
	r_bottom_right := int16(chunk.BottomRight.R)

	err := queries.DeleteCityTrafficFromChunk(ctx, sqlc.DeleteCityTrafficFromChunkParams{
		CityID:       city_id,
		QTopLeft:     q_top_left,
		QBottomRight: q_bottom_right,
		RTopLeft:     r_top_left,
		RBottomRight: r_bottom_right,
	})
	if err != nil {
		return getPostgresError(err)
	}

	employments, err := queries.GetWorkerJobContractsFromChunk(
		ctx,
		sqlc.GetWorkerJobContractsFromChunkParams{
			CityID:       city_id,
			QTopLeft:     q_top_left,
			QBottomRight: q_bottom_right,
			RTopLeft:     r_top_left,
			RBottomRight: r_bottom_right,
		},
	)
	if err != nil {
		return getPostgresError(err)
	}

	commute_ids_map := make(map[cityhall.CommuteKey]int32, len(employments))

	for _, item := range employments {
		key := cityhall.CommuteKey{
			WorkerAnchor: core.TilePosition{
				Q: uint8(item.WorkersAnchorQ),
				R: uint8(item.WorkersAnchorR),
			},
			EmployerAnchor: core.TilePosition{
				Q: uint8(item.EmployerAnchorQ),
				R: uint8(item.EmployerAnchorR),
			},
		}

		commute_ids_map[key] = item.ID
	}

	total_stages := 0
	for _, stages := range commutings {
		total_stages += len(stages)
	}

	if total_stages == 0 {
		return nil
	}

	batch_params := make([]sqlc.InsertTrafficBatchParams, 0, total_stages)

	for key, stages := range commutings {
		employment_id, ok := commute_ids_map[key]
		if !ok {
			slog.Error("ResetCityTrafficFromChunk employment_id not found",
				"commute key", key,
			)
			continue
		}

		path_tiles := make(map[core.TilePosition]sqlc.InsertTrafficBatchParams, len(stages))

		for _, stage := range stages {
			position := core.TilePosition{Q: stage.Q, R: stage.R}
			stage_param, already_visited := path_tiles[position]

			if !already_visited {
				stage_param = sqlc.InsertTrafficBatchParams{
					CityID:         city_id,
					EmploymentID:   employment_id,
					Q:              int16(stage.Q),
					R:              int16(stage.R),
					WorkersAnchorQ: int16(key.WorkerAnchor.Q),
					WorkersAnchorR: int16(key.WorkerAnchor.R),
				}
			}

			addTravelValueInParam(&stage_param, stage.TravelCostEffect)
			path_tiles[position] = stage_param
		}

		batch_params = append(batch_params, slices.Collect(maps.Values(path_tiles))...)
	}

	_, err = queries.InsertTrafficBatch(ctx, batch_params)
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

func addTravelValueInParam(
	traffic_param *sqlc.InsertTrafficBatchParams,
	travel_effect core.TravelCostEffect,
) {
	travel_value := int32(travel_effect.Value)

	switch travel_effect.Type {
	case core.Foot:
		traffic_param.FootValue += travel_value
	case core.Car:
		traffic_param.CarValue += travel_value
	case core.Bus:
		traffic_param.BusValue += travel_value
	case core.Bicycle:
		traffic_param.BicycleValue += travel_value
	case core.Train:
		traffic_param.TrainValue += travel_value
	case core.Subway:
		traffic_param.SubwayValue += travel_value
	}
}
