package postgres

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.PolicyManager = (*Store)(nil)
var _ cityhall.PolicyGetter = (*Store)(nil)

// Set implements PolicyManager.Set
func (s *Store) Set(
	ctx context.Context,
	city_id int32,
	ordinances []core.OrdinanceType,
) error {
	queries := s.q(ctx)

	err := queries.RevokeOrdinances(ctx, city_id)
	if err != nil {
		return getPostgresError(err)
	}

	if len(ordinances) == 0 {
		return nil
	}

	ordinance_ids := make([]int16, len(ordinances))

	for index, id := range ordinances {
		ordinance_ids[index] = int16(id)
	}

	err = queries.EnactOrdinances(ctx, sqlc.EnactOrdinancesParams{
		CityID:       city_id,
		OrdinanceIds: ordinance_ids,
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
