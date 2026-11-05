package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
)

type PolicyManager interface {
	Enact(
		ctx context.Context,
		city_id int32,
		ordinance_name core.OrdinanceType,
	) error
	Revoke(
		ctx context.Context,
		city_id int32,
		ordinance_name core.OrdinanceType,
	) error
}

type PolicyGetter interface {
	EnactedOrdinances(ctx context.Context, city_id int32) ([]core.OrdinanceType, error)
}

type PolicyMapper struct{}

type PolicyEffectKey struct {
	core.ArtifactType
	core.EffectType
}

func (p PolicyMapper) MakeEffectiveTravelEffects(
	enacted_ordinances []core.OrdinanceType,
) map[PolicyEffectKey]int16 {
	enacted_ordinances_map := make(map[core.OrdinanceType]struct{}, len(enacted_ordinances))
	effective_ordinances_travel_effects := make(map[PolicyEffectKey]int16)

	for _, enacted_ordinance := range enacted_ordinances {
		switch enacted_ordinance {
		case core.FreeBicycles:
			enacted_ordinances_map[enacted_ordinance] = struct{}{}
		}
	}

	for _, transport_network_artifact := range registry.TransportNetworkArtifacts {
		artifact := registry.Artifacts[transport_network_artifact]

		for _, ordinance_effect := range artifact.OrdinanceEffects {
			_, ok := enacted_ordinances_map[ordinance_effect.Type]
			if !ok {
				continue
			}

			switch ordinance_effect.Type {
			case core.FreeBicycles:
				key := PolicyEffectKey{artifact.Name, ordinance_effect.Effect}
				effective_ordinances_travel_effects[key] = ordinance_effect.AdditionalValue
			}
		}
	}

	return effective_ordinances_travel_effects
}

// Policies keeps record of enacted ordinances.
// Also tell which ordinances have been enacted by a user.
//
// see core/ordinance.go
type Policies struct {
	Manager PolicyManager
	Getter  PolicyGetter
	Mapper  PolicyMapper
}

func NewPolicies(store Store) *Policies {
	return &Policies{
		Manager: store,
		Getter:  store,
		Mapper:  PolicyMapper{},
	}
}
