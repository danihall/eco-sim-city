package postgres

import (
	"context"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"github.com/jackc/pgx/v5"
)

var _ cityhall.ZoneManager = (*Store)(nil)

// Update implements ZoneManager.Update
func (s *Store) Update(
	ctx context.Context,
	city_id int32,
	p cityhall.ZoneParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	qs := make([]int16, len(p.Coordinates))
	rs := make([]int16, len(p.Coordinates))

	for index := range p.Coordinates {
		qs[index] = int16(p.Coordinates[index].Q)
		rs[index] = int16(p.Coordinates[index].R)
	}

	response, err := queries.UpdateTilesZones(ctx, sqlc.UpdateTilesZonesParams{
		NewSuperKindID:           p.NewSuperKind,
		NewArtifactStatus:        p.NewArtifactStatus,
		CityID:                   city_id,
		NaturalSuperKindID:       p.NaturalSuperKind,
		PublicSystemsSuperKindID: p.PublicSystemsSuperKind,
		OldArtifactStatus:        p.OldArtifactStatus,
		Qs:                       qs,
		Rs:                       rs,
	})
	if err != nil {
		return nil, getPostgresError(err)
	}
	if len(response) == 0 {
		return nil, fmt.Errorf("postgres: UpdateTilesZones: %w", pgx.ErrNoRows)
	}

	updated_tiles := make([]core.Tile, len(response))

	for index := range response {
		tile := response[index]

		updated_tiles[index] = core.Tile{
			Artifact: core.ArtifactInstance{
				Status: tile.ArtifactStatus,
			},
			Position: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
		}
	}

	return updated_tiles, nil
}

// Revert implements ZoneManager.Revert
func (s *Store) Revert(
	ctx context.Context,
	city_id int32,
	p cityhall.ZoneParams,
) ([]core.Tile, error) {
	queries := s.q(ctx)

	qs := make([]int16, len(p.Coordinates))
	rs := make([]int16, len(p.Coordinates))

	for index := range p.Coordinates {
		qs[index] = int16(p.Coordinates[index].Q)
		rs[index] = int16(p.Coordinates[index].R)
	}

	response, err := queries.RevertTilesZones(ctx, sqlc.RevertTilesZonesParams{
		NewSuperKindID:           p.NewSuperKind,
		NewArtifactStatus:        p.NewArtifactStatus,
		CityID:                   city_id,
		NaturalSuperKindID:       p.NaturalSuperKind,
		PublicSystemsSuperKindID: p.PublicSystemsSuperKind,
		OldArtifactStatus:        p.OldArtifactStatus,
		Qs:                       qs,
		Rs:                       rs,
	})
	if err != nil {
		return nil, getPostgresError(err)
	}
	if len(response) == 0 {
		return nil, getPostgresError(err)
	}

	reverted_tiles := make([]core.Tile, len(response))

	for index := range response {
		tile := response[index]
		reverted_tiles[index] = core.Tile{
			Artifact: core.ArtifactInstance{
				Status: tile.ArtifactStatus,
			},
			Position: core.TilePosition{
				Q: uint8(tile.Q),
				R: uint8(tile.R),
			},
		}
	}

	return reverted_tiles, nil
}
