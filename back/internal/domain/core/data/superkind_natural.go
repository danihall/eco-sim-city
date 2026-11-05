//go:build generator

package data

import "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"

// SeedSuperKindNatural is the authoritative data for Natural artifacts. Update with caution.
// NOTE All Natural Artifact must have only one Label.
// NOTE Wasteland artifact has zero effect - it's the artifact's replacement for a destroyed artifact.
var SeedSuperKindNatural = core.SuperKind{
	Name:  core.Natural,
	Label: "nature",
	Kinds: []core.Kind{
		makeKindTrees(),
		makeKindUnconstructible(),
	},
}

func makeKindTrees() core.Kind {
	return core.Kind{
		Name:  core.Trees,
		Label: "arbres & forêts",
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.LandValue, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.Gingko,
				Labels:           []string{"gingko"},
				Description:      "Arbre ornemental",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -1},
						{Type: core.LandValue, Value: 2},
					},
				},
			},
			{
				Name:             core.Birch,
				Labels:           []string{"bouleau"},
				Description:      "Simple arbre",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -2},
					},
				},
			},
			{
				Name:             core.Oak,
				Labels:           []string{"chêne"},
				Description:      "Arbre massif",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -3},
					},
				},
			},
			{
				Name:             core.Paulownia,
				Labels:           []string{"paulownia"},
				Description:      "Absorbe beaucoup de CO2",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -8},
					},
					Global: []core.GlobalEffect{
						{Type: core.WaterConsumption, Value: 2},
					},
				},
			},
			{
				Name:             core.MixedForest,
				Labels:           []string{"forêt"},
				Description:      "Forêt mixte",
				CitizensCapacity: 0,
				JobLevel:         0,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -10},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Foot, Value: 6},
					},
				},
			},
		},
	}
}

// Kind Unconstructible contains artifacts that cannot be built by a player.
func makeKindUnconstructible() core.Kind {
	return core.Kind{
		Name:  core.Unconstructible,
		Label: "inconstructible",
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.LandValue, Value: 1},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.Grassland,
				Labels:           []string{"prairie"},
				Description:      "Prairie",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -2},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Foot, Value: 4},
					},
				},
			},
			{
				Name:             core.Wasteland,
				Labels:           []string{"friche"},
				Description:      "Friche",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 0},
						{Type: core.LandValue, Value: 0},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Foot, Value: 4},
					},
				},
			},
			{
				Name:             core.ShallowWater,
				Labels:           []string{"lac"},
				Description:      "Eau peu profonde",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -2},
						{Type: core.LandValue, Value: 3},
					},
				},
			},
			{
				Name:             core.DeepWater,
				Labels:           []string{"océan"},
				Description:      "Eaux profondes",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -5},
					},
				},
			},
			{
				Name:             core.Bog,
				Labels:           []string{"tourbière"},
				Description:      "Absorbe énormément de CO2",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -25},
						{Type: core.LandValue, Value: 2},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Foot, Value: 8},
					},
				},
			},
		},
	}
}
