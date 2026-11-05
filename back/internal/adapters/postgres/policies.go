package postgres

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.PolicyManager = (*Store)(nil)
var _ cityhall.PolicyGetter = (*Store)(nil)

// Enact implements PolicyManager.Enact
func (s *Store) Enact(
	ctx context.Context,
	city_id int32,
	ordinance_name core.OrdinanceType,
) error {
	queries := s.q(ctx)

	err := queries.EnactOrdinance(ctx, sqlc.EnactOrdinanceParams{
		CityID:      city_id,
		OrdinanceID: ordinance_name,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// Revoke implements PolicyManager.Revoke
func (s *Store) Revoke(
	ctx context.Context,
	city_id int32,
	ordinance_name core.OrdinanceType,
) error {
	queries := s.q(ctx)

	err := queries.RevokeOrdinance(ctx, sqlc.RevokeOrdinanceParams{
		CityID:      city_id,
		OrdinanceID: ordinance_name,
	})
	if err != nil {
		return getPostgresError(err)
	}

	return nil
}

// EnactedOrdinances implements PolicyGetter.EnactedOrdinances
func (s *Store) EnactedOrdinances(
	ctx context.Context,
	city_id int32,
) ([]core.OrdinanceType, error) {
	queries := s.q(ctx)

	response, err := queries.GetEnactedOrdinances(ctx, city_id)
	if err != nil {
		return nil, getPostgresError(err)
	}

	return response, nil
}
