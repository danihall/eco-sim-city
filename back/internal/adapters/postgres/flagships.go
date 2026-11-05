package postgres

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.FlagshipManager = (*Store)(nil)
var _ cityhall.FlagshipGetter = (*Store)(nil)

// Add implements FlagshipManager.Add
func (s *Store) Add(
	ctx context.Context,
	city_id int32,
	artifact_name core.ArtifactType,
) error {
	queries := s.q(ctx)

	err := queries.AddUniqueArtifact(ctx, sqlc.AddUniqueArtifactParams{
		CityID:     city_id,
		ArtifactID: artifact_name,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// Remove implements FlagshipManager.Remove
func (s *Store) Remove(
	ctx context.Context,
	city_id int32,
	artifact_name core.ArtifactType,
) error {
	queries := s.q(ctx)

	err := queries.RemoveUniqueArtifact(ctx, sqlc.RemoveUniqueArtifactParams{
		CityID:     city_id,
		ArtifactID: artifact_name,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// All implements FlagshipGetter.All
func (s *Store) All(
	ctx context.Context,
	city_id int32,
) ([]core.ArtifactType, error) {
	queries := s.q(ctx)

	response, err := queries.GetUniqueArtifacts(ctx, city_id)
	if err != nil {
		return nil, getPostgresError(err)
	}

	return response, nil
}
