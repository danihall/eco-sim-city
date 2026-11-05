// city_hall assemble all actors of city's domain as deputies.
//
// Each deputy fulfills a role analogous to a city's councelor.
// From a domain-driven-design perspective, mimicking how a city's council operates
// feels cognitively natural and provides a meaningful way to express
// the many sub-domains of a city's simulation.
package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"golang.org/x/sync/errgroup"
)

type Store interface {
	Administrator
	Reporter
	Builder
	Destroyer
	Completer
	Vacator
	JobTerminator
	JobSignator
	JobSeeker
	FlagshipManager
	FlagshipGetter
	PolicyManager
	PolicyGetter
	TaskForceExecutor
	TreasuryWriter
	TreasuryReader
	TilesGetter
	TilesUpdater
	TransportsGetter
	TransportsUpdater
	ZoneManager
}

type Deputies struct {
	Cabinet          *Cabinet
	Cadastre         *Cadastre
	Construction     *Construction
	EmploymentOffice *EmploymentOffice
	Flagship         *Flagship
	Permits          *Permits
	Policies         *Policies
	TaskForce        *TaskForce
	Topology         *Topology
	Transportation   *Transportation
	Treasury         *Treasury
	Zoning           *Zoning
}

func New(store Store) *Deputies {
	return &Deputies{
		Cabinet:          NewCabinet(store),
		Cadastre:         NewCadastre(store),
		Construction:     NewConstruction(store),
		EmploymentOffice: NewEmploymentOffice(store),
		Flagship:         NewFlagship(store),
		Permits:          NewPermits(),
		Policies:         NewPolicies(store),
		TaskForce:        NewTaskForce(store),
		Topology:         NewTopology(),
		Transportation:   NewTransportation(store),
		Treasury:         NewTreasury(store),
		Zoning:           NewZoning(store),
	}
}

type CityFullState struct {
	Flagships         []core.ArtifactType
	EnactedOrdinances []core.OrdinanceType
	Stats             core.CityStats
	Budget            map[core.KindType]core.KindBudget
}

func (d *Deputies) GetCityFullState(
	ctx context.Context,
	city_id int32,
) (CityFullState, error) {
	group, group_ctx := errgroup.WithContext(ctx)

	var existing_flagships []core.ArtifactType
	var enacted_ordinances []core.OrdinanceType
	var city_stats core.CityStats
	var city_budget map[core.KindType]core.KindBudget

	{
		group.Go(func() error {
			response, err := d.Flagship.Getter.All(group_ctx, city_id)
			if err != nil {
				return err
			}

			existing_flagships = response
			return nil
		})

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

		group.Go(func() error {
			response, err := d.Treasury.Reader.AllKindsBudget(group_ctx, city_id)
			if err != nil {
				return err
			}
			city_budget = response
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return CityFullState{}, err
	}

	return CityFullState{
		Flagships:         existing_flagships,
		EnactedOrdinances: enacted_ordinances,
		Stats:             city_stats,
		Budget:            city_budget,
	}, nil
}

type CityDataForTrafficCompute struct {
	TransportNetwork  []core.Tile
	WorkerEmployments []core.Employment
	EnactedOrdinances []core.OrdinanceType
}

// GetCityDataForTrafficCompute TransportNetwork field contains these artifacts:
//
//   - core.Road
//   - core.Rail
//   - core.SubwayStation
//   - core.BusStop
//   - core.TrainStation
//   - core.PedestrianBicycleWay
func (d *Deputies) GetCityDataForTrafficCompute(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) (CityDataForTrafficCompute, error) {
	group, group_ctx := errgroup.WithContext(ctx)

	var transport_network_tiles []core.Tile
	var worker_employments []core.Employment
	var enacted_ordinances []core.OrdinanceType

	{
		group.Go(func() error {
			response, err := d.Transportation.Getter.CityNetwork(
				group_ctx,
				city_id,
				registry.TransportNetworkArtifacts,
			)
			if err != nil {
				return err
			}

			transport_network_tiles = response
			return nil
		})

		group.Go(func() error {
			response, err := d.EmploymentOffice.Seeker.FromWorkerAnchorChunk(
				group_ctx,
				city_id,
				chunk,
			)
			if err != nil {
				return err
			}

			worker_employments = response
			return nil
		})

		group.Go(func() error {
			response, err := d.Policies.Getter.EnactedOrdinances(group_ctx, city_id)
			if err != nil {
				return err
			}

			enacted_ordinances = response
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return CityDataForTrafficCompute{}, err
	}

	return CityDataForTrafficCompute{
		TransportNetwork:  transport_network_tiles,
		WorkerEmployments: worker_employments,
		EnactedOrdinances: enacted_ordinances,
	}, nil
}
