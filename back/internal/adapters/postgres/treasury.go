package postgres

import (
	"context"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ cityhall.TreasuryReader = (*Store)(nil)
var _ cityhall.TreasuryWriter = (*Store)(nil)

// KindBudget implements TreasuryReader.KindBudget
func (s *Store) KindBudget(
	ctx context.Context,
	city_id int32,
	kind_name core.KindType,
) (core.KindBudget, error) {
	queries := s.q(ctx)

	response, err := queries.GetCityKindBudget(ctx, sqlc.GetCityKindBudgetParams{
		CityID: city_id,
		KindID: kind_name,
	})
	if err != nil {
		return core.KindBudget{}, getPostgresError(err)
	}

	return core.KindBudget{
		TaxRate:     uint8(response.TaxRate),
		FundingRate: uint8(response.FundingRate),
	}, nil
}

// AllKindsBudget implements TreasuryReader.AllKindsBudget
func (s *Store) AllKindsBudget(
	ctx context.Context,
	city_id int32,
) (map[core.KindType]core.KindBudget, error) {
	queries := s.q(ctx)

	response, err := queries.GetCityAllKindsBudget(ctx, city_id)
	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("no budget found for city_id=%v", city_id)
	}

	city_kinds_budget := make(map[core.KindType]core.KindBudget, len(response))
	for index := range response {
		kind_budget := &response[index]

		city_kinds_budget[kind_budget.KindID] = core.KindBudget{
			TaxRate:     uint8(kind_budget.TaxRate),
			FundingRate: uint8(kind_budget.FundingRate),
		}
	}

	return city_kinds_budget, nil
}

// CityBalance implements TreasuryReader.CityBalance
func (s *Store) CityBalance(ctx context.Context, city_id int32) (core.CityBalance, error) {
	queries := s.q(ctx)

	response, err := queries.GetCityBalance(ctx, city_id)
	if err != nil {
		return core.CityBalance{}, err
	}

	city_balance := core.CityBalance{
		Revenues: response.Revenues,
		Expenses: response.Expenses,
	}

	return city_balance, nil
}

// TODO impl dat shit
func (s *Store) Budget(ctx context.Context, city_id int32, params []core.KindBudget) error {
	return nil
}
