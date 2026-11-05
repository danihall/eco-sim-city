//go:build generator

package data

import (
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// SeedSuperKindResidential is the authoritative data for Residential artifacts. Update with caution.
// NOTE Each Residential Kind has an associated Social Class.
var SeedSuperKindResidential = core.SuperKind{
	Name:  core.Residential,
	Label: "résidences",
	Kinds: []core.Kind{
		makeKindResidentialLowEnd(),
		makeKindResidentialMediumEnd(),
		makeKindResidentialHighEnd(),
		makeKindResidentialUltraHighEnd(),
	},
}

func makeKindResidentialLowEnd() core.Kind {
	return core.Kind{
		Name:  core.ResidentialLowEnd,
		Label: "résidentiel bas de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		SocialClass: &core.SocialClass{
			MinJobLevel:                 0,
			AdditionalTravelCostBus:     1,
			AdditionalTravelCostSubway:  2,
			AdditionalTravelCostTrain:   3,
			AdditionalTravelCostCar:     4,
			AdditionalTravelCostBicycle: 5,
			AdditionalTravelCostFoot:    6,
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.HouseLow,
				Labels:           []string{"petite maison"},
				Description:      "petite maison modeste mais fonctionnelle.",
				CitizensCapacity: 1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 1},
						{Type: core.Ghg3, Value: 10},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 1},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.ResidentialBuildingLow,
				Labels:           []string{"simples appartements"},
				Description:      "immeuble d'appartements basiques, dense et peu coûteux.",
				CitizensCapacity: 6,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeDrug, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.Ghg3, Value: 20},
						{Type: core.ElectricityConsumption, Value: 3},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.BudgetCost, Value: -3},
					},
				},
			},
			{
				Name:             core.LargeResidentialTowerLow,
				Labels:           []string{"grande tour d'habitation"},
				Description:      "grande tour d'habitation. Forte densité, conditions de vie précaires.",
				CitizensCapacity: 10,
				Width:            3,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeDrug, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 8},
						{Type: core.Ghg3, Value: 30},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.BudgetCost, Value: -4},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.DrugPrevention,
						Effect:          core.CrimeDrug,
						AdditionalValue: -1,
					},
				},
			},
		},
	}
}

func makeKindResidentialMediumEnd() core.Kind {
	return core.Kind{
		Name:  core.ResidentialMediumEnd,
		Label: "résidentiel milieu de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxGhg1, Value: 60},
				{Type: core.MaxLandValue, Value: 20},
				{Type: core.MinLandValue, Value: 1},
				{Type: core.MaxCrimeGeneric, Value: 15},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 2},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		SocialClass: &core.SocialClass{
			MinJobLevel:                 1,
			AdditionalTravelCostSubway:  1,
			AdditionalTravelCostCar:     2,
			AdditionalTravelCostTrain:   3,
			AdditionalTravelCostBus:     4,
			AdditionalTravelCostBicycle: 5,
			AdditionalTravelCostFoot:    6,
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.HouseMedium,
				Labels:           []string{"jolie maison"},
				Description:      "maison confortable pour famille de classe moyenne.",
				CitizensCapacity: 1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 1},
						{Type: core.Ghg3, Value: 10},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 1},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.ResidentialBuildingMedium,
				Labels:           []string{"simples appartements"},
				Description:      "immeuble d'appartements spacieux",
				CitizensCapacity: 6,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.Ghg3, Value: 20},
						{Type: core.ElectricityConsumption, Value: 3},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.BudgetCost, Value: -4},
					},
				},
			},
			{
				Name:             core.LargeResidentialTowerMedium,
				Labels:           []string{"grande tour d'habitation"},
				Description:      "tour résidentielle accessible à la classe moyenne.",
				CitizensCapacity: 10,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeDrug, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: 30},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 6},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.DrugPrevention,
						Effect:          core.CrimeDrug,
						AdditionalValue: -1,
					},
				},
			},
		},
	}
}

func makeKindResidentialHighEnd() core.Kind {
	return core.Kind{
		Name:  core.ResidentialHighEnd,
		Label: "résidentiel haut de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxGhg1, Value: 30},
				{Type: core.MaxLandValue, Value: 40},
				{Type: core.MinLandValue, Value: 8},
				{Type: core.MaxCrimeGeneric, Value: 6},
				{Type: core.MaxCrimeDrug, Value: 5},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
			Global: []core.GlobalEffect{
				{Type: core.CrimeTaxFraud, Value: 2},
			},
		},
		SocialClass: &core.SocialClass{
			MinJobLevel:                 3,
			AdditionalTravelCostCar:     1,
			AdditionalTravelCostBicycle: 3,
			AdditionalTravelCostFoot:    4,
			AdditionalTravelCostTrain:   5,
			AdditionalTravelCostSubway:  7,
			AdditionalTravelCostBus:     9,
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.HouseHigh,
				Labels:           []string{"villa"},
				Description:      "villa spacieuse haut de gamme",
				CitizensCapacity: 1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: 20},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.BudgetCost, Value: -3},
					},
				},
			},
			{
				Name:             core.ResidentialBuildingHigh,
				Labels:           []string{"grands appartements"},
				Description:      "immeuble de standing",
				CitizensCapacity: 4,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: 40},
						{Type: core.ElectricityConsumption, Value: 5},
						{Type: core.WaterConsumption, Value: 5},
						{Type: core.BudgetCost, Value: -6},
					},
				},
			},
			{
				Name:             core.LargeResidentialTowerHigh,
				Labels:           []string{"grande tour d'habitation haut de gamme"},
				Description:      "tour résidentielle de prestige avec services inclus",
				CitizensCapacity: 6,
				Width:            3,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeDrug, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 20},
						{Type: core.Ghg3, Value: 50},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.BudgetCost, Value: -10},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.DrugPrevention,
						Effect:          core.CrimeDrug,
						AdditionalValue: -1,
					},
				},
			},
		},
	}
}

func makeKindResidentialUltraHighEnd() core.Kind {
	return core.Kind{
		Name:  core.ResidentialUltraHighEnd,
		Label: "résidentiel ultra-haut de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxGhg1, Value: 10},
				{Type: core.MinLandValue, Value: 20},
				{Type: core.MaxCrimeDrug, Value: 1},
				{Type: core.MaxCrimeGeneric, Value: 1},
			},
		},
		Effects: core.Effects{
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
			Global: []core.GlobalEffect{
				{Type: core.CrimeTaxFraud, Value: 10},
			},
		},
		SocialClass: &core.SocialClass{
			MinJobLevel:                 5,
			AdditionalTravelCostCar:     0,
			AdditionalTravelCostTrain:   8,
			AdditionalTravelCostBicycle: 9,
			AdditionalTravelCostFoot:    12,
			AdditionalTravelCostSubway:  15,
			AdditionalTravelCostBus:     17,
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.HouseUltraHigh,
				Labels:           []string{"manoir"},
				Description:      "manoir luxueux avec parc privatif, réservé à une élite fortunée.",
				CitizensCapacity: 1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeGeneric, Value: -2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 20},
						{Type: core.Ghg3, Value: 60},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.BudgetCost, Value: -6},
					},
				},
			},
			{
				Name:             core.ResidentialBuildingUltraHigh,
				Labels:           []string{"appartements de luxe"},
				Description:      "appartements de grand luxe",
				CitizensCapacity: 3,
				Width:            3,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeGeneric, Value: -1},
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 20},
						{Type: core.Ghg3, Value: 90},
						{Type: core.ElectricityConsumption, Value: 5},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.BudgetCost, Value: -14},
					},
				},
			},
			{
				Name:             core.LargeResidentialTowerUltraHigh,
				Labels:           []string{"penthouse"},
				Description:      "penthouses sompteux pour ultra-riches",
				CitizensCapacity: 3,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeGeneric, Value: -1},
						{Type: core.CrimeVss, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 40},
						{Type: core.Ghg3, Value: 150},
						{Type: core.ElectricityConsumption, Value: 12},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.BudgetCost, Value: -30},
						{Type: core.CrimeTaxFraud, Value: 24},
					},
				},
			},
		},
	}
}
