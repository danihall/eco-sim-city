package api

import (
	"context"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"golang.org/x/sync/errgroup"
)

// EnactOrdinance is a user action.
func GetAllOrdinances(
	ctx context.Context,
	d *cityhall.Deputies,
	city_id int32,
) ([]cityhall.OrdinanceWithState, error) {
	group, group_ctx := errgroup.WithContext(ctx)

	var enacted_ordinances []core.OrdinanceType
	var city_stats core.CityStats

	{
		group.Go(func() error {
			response, err := d.Policies.Getter.EnactedOrdinances(group_ctx, city_id)
			if err != nil {
				return err
			}

			enacted_ordinances = response
			return nil
		})

		group.Go(func() error {
			response, err := d.Cabinet.Reporter.CityStats(group_ctx, city_id)
			if err != nil {
				return err
			}
			city_stats = response
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return nil, err
	}

	all_ordinances := make([]cityhall.OrdinanceWithState, 0, len(registry.Ordinances))

	for name, ordinance := range registry.Ordinances {
		all_ordinances = append(all_ordinances, cityhall.OrdinanceWithState{
			Name:        name,
			Label:       ordinance.Label,
			Description: ordinance.Description,
			TotalCost:   int64(ordinance.CostPerCitizen) * int64(city_stats.CitizensCount),
			IsEnacted:   slices.Contains(enacted_ordinances, name),
		})
	}

	return all_ordinances, nil
}
