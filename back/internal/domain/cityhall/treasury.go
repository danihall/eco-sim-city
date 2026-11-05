package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
)

type TreasuryWriter interface {
	Budget(ctx context.Context, city_id int32, params []core.KindBudget) error
}

type TreasuryReader interface {
	KindBudget(ctx context.Context, city_id int32, kind_name core.KindType) (core.KindBudget, error)
	AllKindsBudget(ctx context.Context, city_id int32) (map[core.KindType]core.KindBudget, error)
	CityBalance(ctx context.Context, city_id int32) (core.CityBalance, error)
}

type BudgetMaker struct{}

// Treasury manages actions over a city's budget.
type Treasury struct {
	Writer TreasuryWriter
	Reader TreasuryReader
	Maker  BudgetMaker
}

func NewTreasury(store Store) *Treasury {
	return &Treasury{
		Writer: store,
		Reader: store,
		Maker:  BudgetMaker{},
	}
}

func (b BudgetMaker) InitialBudget() []core.KindBudget {
	budget := make([]core.KindBudget, 0, len(registry.SuperKinds)*10)

	for super_kind, kinds := range registry.SuperKinds {
		for _, kind := range kinds {
			var initial_tax_rate uint8
			var initial_funding_rate uint8

			switch super_kind {
			case core.Natural:
				initial_tax_rate = 0
				initial_funding_rate = 0
			case core.PublicSystems:
				initial_tax_rate = 0
				initial_funding_rate = 80
			case core.Industrial, core.Residential:
				initial_tax_rate = 30
				initial_funding_rate = 0
			}

			budget = append(budget, core.KindBudget{
				KindName:    kind,
				TaxRate:     initial_tax_rate,
				FundingRate: initial_funding_rate,
			})
		}
	}

	return budget
}
