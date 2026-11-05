package postgres

import (
	"context"
	"fmt"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.Builder = (*Store)(nil)
var _ cityhall.Completer = (*Store)(nil)
var _ cityhall.Destroyer = (*Store)(nil)
var _ cityhall.Vacator = (*Store)(nil)

// SingleArtifact implements cityhall.Builder.SingleArtifact
func (s *Store) SingleArtifact(
	ctx context.Context,
	city_id int32,
	params cityhall.BuildParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	qs, rs := makeQueryCoordinates(params.Coordinates)
	anchor_qs, anchor_rs := makeQueryCoordinates(params.Anchors)

	response, err := queries.BuildArtifact(ctx, sqlc.BuildArtifactParams{
		CityID:         city_id,
		SuperKindID:    params.SuperKindType,
		KindID:         params.KindType,
		ArtifactID:     params.Artifact.Name,
		ArtifactLabel:  params.Artifact.Label,
		ArtifactStatus: params.Artifact.Status,
		Qs:             qs,
		Rs:             rs,
		AnchorQs:       anchor_qs,
		AnchorRs:       anchor_rs,
	})
	if err != nil {
		return nil, getPostgresError(err)
	}

	updated_tiles := make([]core.Tile, len(response))
	for index := range response {
		tile := &response[index]
		updated_tiles[index] = core.Tile{
			SuperKind: tile.SuperKindID,
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
		}
	}

	return updated_tiles, nil
}

// MultipleArtifacts implements cityhall.Builder.MultipleArtifacts
func (s *Store) MultipleArtifacts(
	ctx context.Context,
	city_id int32,
	params []cityhall.BuildParams,
) error {
	queries := s.q(ctx)

	batch_params := sqlc.BuildMultipleArtifactsParams{
		CityID:           city_id,
		Qs:               []int16{},
		Rs:               []int16{},
		SuperKindIds:     []int16{},
		KindIds:          []int16{},
		ArtifactIds:      []int16{},
		ArtifactLabels:   []string{},
		ArtifactStatuses: []int16{},
		AnchorQs:         []int16{},
		AnchorRs:         []int16{},
	}

	expected_updated_rows_count := 0
	for index := range params {
		item := params[index]

		qs, rs := makeQueryCoordinates(item.Coordinates)
		anchor_qs, anchor_rs := makeQueryCoordinates(item.Anchors)
		super_kind_id := item.SuperKindType
		kind_id := item.KindType
		artifact_id := item.Artifact.Name

		batch_params.Qs = slices.Concat(batch_params.Qs, qs)
		batch_params.Rs = slices.Concat(batch_params.Rs, rs)
		batch_params.SuperKindIds = slices.Concat(
			batch_params.SuperKindIds,
			slices.Repeat([]int16{int16(super_kind_id)}, len(qs)),
		)
		batch_params.KindIds = slices.Concat(
			batch_params.KindIds,
			slices.Repeat([]int16{int16(kind_id)}, len(qs)),
		)
		batch_params.ArtifactIds = slices.Concat(
			batch_params.ArtifactIds,
			slices.Repeat([]int16{int16(artifact_id)}, len(qs)),
		)
		batch_params.ArtifactLabels = slices.Concat(
			batch_params.ArtifactLabels,
			slices.Repeat([]string{item.Artifact.Label}, len(qs)),
		)
		batch_params.ArtifactStatuses = slices.Concat(
			batch_params.ArtifactStatuses,
			slices.Repeat([]int16{int16(item.Artifact.Status)}, len(qs)),
		)
		batch_params.AnchorQs = slices.Concat(batch_params.AnchorQs, anchor_qs)
		batch_params.AnchorRs = slices.Concat(batch_params.AnchorRs, anchor_rs)
		expected_updated_rows_count += len(qs)
	}

	response, err := queries.BuildMultipleArtifacts(ctx, batch_params)
	if err != nil {
		return err
	}

	if len(response) != expected_updated_rows_count {
		return fmt.Errorf("mismatch between number of artifacts given and tiles updated")
	}

	return nil
}

// Pathway implements cityhall.Builder.Pathway
func (s *Store) Pathway(
	ctx context.Context,
	city_id int32,
	params cityhall.PathwayParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	qs, rs := makeQueryCoordinates(params.Coordinates)
	anchor_qs, anchors_rs := makeQueryCoordinates(params.Anchors)

	response, err := queries.BuildPathway(ctx, sqlc.BuildPathwayParams{
		CityID:                  city_id,
		SuperKindID:             params.SuperKindType,
		KindID:                  params.KindType,
		ArtifactID:              params.Artifact.Name,
		ArtifactLabel:           params.Artifact.Label,
		CompletedArtifactStatus: params.Artifact.Status,
		AnchorQs:                anchor_qs,
		AnchorRs:                anchors_rs,
		Qs:                      qs,
		Rs:                      rs,
	})
	if err != nil {
		return nil, getPostgresError(err)
	}

	updated_tiles := make([]core.Tile, len(response))
	for index := range response {
		tile := &response[index]

		updated_tiles[index] = core.Tile{
			Artifact: core.ArtifactInstance{
				Label:  tile.ArtifactLabel,
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
		}
	}

	return updated_tiles, nil
}

// AllFromChunk implements cityhall.Completer.AllFromChunk
func (s *Store) AllFromChunk(
	ctx context.Context,
	city_id int32,
	params struct {
		CurrentArtifactStatus core.ArtifactStatus
		NewArtifactStatus     core.ArtifactStatus
		Chunk                 core.Chunk
	},
) ([]core.Tile, error) {
	queries := s.q(ctx)

	response, err := queries.CompleteArtifactsFromChunk(
		ctx,
		sqlc.CompleteArtifactsFromChunkParams{
			CityID:                city_id,
			CurrentArtifactStatus: params.CurrentArtifactStatus,
			NewArtifactStatus:     params.NewArtifactStatus,
			QTopLeft:              int16(params.Chunk.TopLeft.Q),
			QBottomRight:          int16(params.Chunk.BottomRight.Q),
			RTopLeft:              int16(params.Chunk.TopLeft.R),
			RBottomRight:          int16(params.Chunk.BottomRight.R),
		},
	)
	if err != nil {
		return nil, getPostgresError(err)
	}

	updated_tiles := make([]core.Tile, len(response))
	for index := range response {
		tile := &response[index]

		updated_tiles[index] = core.Tile{
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
		}
	}

	return updated_tiles, nil
}

// Artifact implements cityhall.Destroyer.Artifact
func (s *Store) Artifact(
	ctx context.Context,
	city_id int32,
	destroy_params cityhall.DestroyParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	params := sqlc.DestroyArtifactParams{
		ReplacementArtifactID:     destroy_params.ReplacementArtifact.Name,
		ReplacementArtifactLabel:  destroy_params.ReplacementArtifact.Label,
		ReplacementArtifactStatus: destroy_params.ReplacementArtifact.Status,
		CityID:                    city_id,
		AnchorQ:                   int16(destroy_params.Anchor.Q),
		AnchorR:                   int16(destroy_params.Anchor.R),
	}

	response, err := queries.DestroyArtifact(ctx, params)
	if err != nil {
		return nil, getPostgresError(err)
	}

	destroyed_tiles := make([]core.Tile, len((response)))

	for index := range response {
		tile := &response[index]
		destroyed_tiles[index] = core.Tile{
			SuperKind: tile.SuperKindID,
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
		}
	}

	return destroyed_tiles, nil
}

// PublicSystemsArtifact implements cityhall.Destroyer.PublicSystemsArtifact
func (s *Store) PublicSystemsArtifact(
	ctx context.Context,
	city_id int32,
	destroy_params cityhall.DestroyPublicSystemsParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	params := sqlc.DestroyPublicSystemsArtifactParams{
		ReplacementSuperKindID:    destroy_params.ReplacementSuperKindType,
		ReplacementKindID:         destroy_params.ReplacementKindType,
		ReplacementArtifactID:     destroy_params.ReplacementArtifact.Name,
		ReplacementArtifactLabel:  destroy_params.ReplacementArtifact.Label,
		ReplacementArtifactStatus: destroy_params.ReplacementArtifact.Status,
		CityID:                    city_id,
		AnchorQ:                   int16(destroy_params.Anchor.Q),
		AnchorR:                   int16(destroy_params.Anchor.R),
	}

	response, err := queries.DestroyPublicSystemsArtifact(ctx, params)
	if err != nil {
		return nil, getPostgresError(err)
	}

	destroyed_tiles := make([]core.Tile, len((response)))

	for index := range response {
		tile := &response[index]
		destroyed_tiles[index] = core.Tile{
			SuperKind: tile.SuperKindID,
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
		}
	}

	return destroyed_tiles, nil
}

// AbandonArtifacts implements cityhall.Vacator.AbandonArtifacts
func (s *Store) AbandonArtifacts(
	ctx context.Context,
	city_id int32,
	params cityhall.VacateParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	anchor_qs, anchor_rs := makeQueryCoordinates(params.Anchors)

	response, err := queries.VacateMultipleArtifacts(
		ctx,
		sqlc.VacateMultipleArtifactsParams{
			ArtifactLabel:          params.Label,
			ArtifactStatus:         params.Status,
			HasAllCitizenAllocated: params.HasAllCitizensAllocated,
			CityID:                 city_id,
			AnchorQs:               anchor_qs,
			AnchorRs:               anchor_rs,
		},
	)
	if err != nil {
		return nil, err
	}

	abandoned_tiles := make([]core.Tile, len(response))

	for index := range response {
		tile := &response[index]
		abandoned_tiles[index] = core.Tile{
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
		}
	}

	return abandoned_tiles, nil
}
