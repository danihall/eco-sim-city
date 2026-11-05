//go:build generator

package data

import (
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// SeedSuperKindIndustrial is the authoritative data for Industrial artifacts. Update with caution.
var SeedSuperKindIndustrial = core.SuperKind{
	Name:  core.Industrial,
	Label: "industries",
	Kinds: []core.Kind{
		makeKindPetrochemical(),
		makeKindMetallurgy(),
		makeKindBuilding(),
		makeKindTextiles(),
		makeKindAgrofood(),
		makeKindDigital(),
		makeKindCommercialB2B(),
		makeKindCommercialB2CLowEnd(),
		makeKindCommercialB2CMediumEnd(),
		makeKindCommercialB2CHighEnd(),
	},
}

func makeKindPetrochemical() core.Kind {
	return core.Kind{
		Name:  core.Petrochemical,
		Label: "pétrochimie",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 50},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 2},
			},
			Global: []core.GlobalEffect{
				{Type: core.AdultsHealth, Value: -1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.ChemicalPlant,
				Labels:           []string{"Toca-Cola"},
				Description:      "production d'agilité et de composés chimiques de toute sortes!",
				CitizensCapacity: 9,
				JobLevel:         0,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 100},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 200},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.BudgetCost, Value: -10},
					},
				},
			},
			{
				Name:             core.PlasticsFactory,
				Labels:           []string{"Dark Matter"},
				Description:      "du plastique en pagaille",
				CitizensCapacity: 9,
				JobLevel:         1,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 120},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 250},
						{Type: core.ElectricityConsumption, Value: 12},
						{Type: core.WaterConsumption, Value: 30},
						{Type: core.BudgetCost, Value: -12},
					},
				},
			},
		},
	}
}

func makeKindMetallurgy() core.Kind {
	return core.Kind{
		Name:  core.Metallurgy,
		Label: "métallurgie",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 50},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 2},
			},
			Global: []core.GlobalEffect{
				{Type: core.AdultsHealth, Value: -1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.SteelMill,
				Labels:           []string{"ArcelorMitemps"},
				Description:      "l'homme a créé le feu. ArcelorMitemps a créé l'acier",
				CitizensCapacity: 12,
				JobLevel:         0,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 200},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 200},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.BudgetCost, Value: -20},
					},
				},
			},
			{
				Name:             core.AluminiumFactory,
				Labels:           []string{"Alumin & co"},
				Description:      "de l'aluminium pour tous",
				CitizensCapacity: 6,
				JobLevel:         1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 120},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 250},
						{Type: core.ElectricityConsumption, Value: 20},
						{Type: core.WaterConsumption, Value: 30},
						{Type: core.BudgetCost, Value: -12},
					},
				},
			},
		},
	}
}

func makeKindBuilding() core.Kind {
	return core.Kind{
		Name:  core.Construction,
		Label: "construction",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 15},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 50},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 3},
			},
			Global: []core.GlobalEffect{
				{Type: core.AdultsHealth, Value: -1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.CementFactory,
				Labels:           []string{"Vinsi"},
				Description:      "du ciment pas piqué des hannetons",
				CitizensCapacity: 12,
				JobLevel:         0,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 150},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 200},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 30},
						{Type: core.BudgetCost, Value: -30},
					},
				},
			},
			{
				Name:             core.LumberYard,
				Labels:           []string{"ArcelorBois"},
				Description:      "nous faisons du bois solide comme de l'acier",
				CitizensCapacity: 4,
				JobLevel:         0,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 40},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 30},
						{Type: core.Ghg3, Value: 100},
						{Type: core.ElectricityConsumption, Value: 10},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.BudgetCost, Value: -12},
					},
				},
			},
		},
	}
}

func makeKindTextiles() core.Kind {
	return core.Kind{
		Name:  core.Textiles,
		Label: "textiles",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 50},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 2},
			},
			Global: []core.GlobalEffect{
				{Type: core.AdultsHealth, Value: -1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.TextilesFactory,
				Labels:           []string{"De Furrsac"},
				Description:      "la griffe et le style",
				CitizensCapacity: 15,
				JobLevel:         0,
				Width:            4,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 150},
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 150},
						{Type: core.Ghg3, Value: 150},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 34},
						{Type: core.BudgetCost, Value: -32},
					},
				},
			},
		},
	}
}

func makeKindAgrofood() core.Kind {
	return core.Kind{
		Name:  core.Agrofood,
		Label: "agroalimentaire",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 50},
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
		Artifacts: []core.Artifact{
			{
				Name:             core.FoodProcessingPlant,
				Labels:           []string{"Pizza Tut"},
				Description:      "Production d'aliments transformés",
				CitizensCapacity: 4,
				JobLevel:         0,
				Width:            2,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 100},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 80},
						{Type: core.Ghg3, Value: 100},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.BudgetCost, Value: -8},
					},
				},
			},
			{
				Name:             core.BiologicalFarm,
				Labels:           []string{"Le potager d'Alex", "Naturelia"},
				Description:      "Agriculture bio",
				CitizensCapacity: 4,
				JobLevel:         1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 4},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: 10},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.ChildrenHealth, Value: 2},
						{Type: core.AdultsHealth, Value: 2},
						{Type: core.ElderlyHealth, Value: 2},
						{Type: core.BudgetCost, Value: 2},
					},
				},
			},
		},
	}
}

func makeKindDigital() core.Kind {
	return core.Kind{
		Name:  core.Digital,
		Label: "numérique",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxGhg1, Value: 8},
				{Type: core.MaxCrimeGeneric, Value: 2},
				{Type: core.MinLandValue, Value: 8},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MinCitizensEducationLevel, Value: 4},
				{Type: core.MaxTaxRate, Value: 20},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeVss, Value: 1},
			},
			Global: []core.GlobalEffect{
				{Type: core.CrimeTaxFraud, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.Esn,
				Labels:           []string{"ClapGemini", "Sorla Steria"},
				Description:      "Entreprise de services numériques",
				CitizensCapacity: 4,
				JobLevel:         3,
				Width:            3,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 2},
						{Type: core.LandValue, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 200},
						{Type: core.Ghg3, Value: 100},
						{Type: core.ElectricityConsumption, Value: 20},
						{Type: core.WaterConsumption, Value: 3},
						{Type: core.BudgetCost, Value: -16},
					},
				},
			},
			{
				Name:             core.VideoGameStartup,
				Labels:           []string{"Nintendro"},
				Description:      "nous faisons des jeux vidéo garantis sans additifs",
				CitizensCapacity: 2,
				JobLevel:         4,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.LandValue, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 200},
						{Type: core.Ghg3, Value: 150},
						{Type: core.ElectricityConsumption, Value: 10},
						{Type: core.WaterConsumption, Value: 3},
						{Type: core.BudgetCost, Value: -12},
					},
				},
			},
			{
				Name:             core.DataCenter,
				Labels:           []string{"Big Data Corp"},
				Description:      "nous fournissons de quoi stocker vos données en toute légalité",
				CitizensCapacity: 3,
				JobLevel:         1,
				Width:            7,
				Height:           5,
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensEducationLevel, Value: 2},
					},
				},
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 4},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 300},
						{Type: core.Ghg3, Value: 200},
						{Type: core.ElectricityConsumption, Value: 80},
						{Type: core.WaterConsumption, Value: 90},
						{Type: core.BudgetCost, Value: -45},
					},
				},
			},
			{
				Name:             core.AiStartup,
				Labels:           []string{"Pimped AI"},
				Description:      "nous créons de la drogue virtuelle",
				CitizensCapacity: 2,
				JobLevel:         6,
				Width:            2,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeDrug, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 100},
						{Type: core.ElectricityConsumption, Value: 24},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.ChildrenEducation, Value: -16},
						{Type: core.AdultsEducation, Value: -8},
						{Type: core.ElderlyEducation, Value: -2},
						{Type: core.BudgetCost, Value: -30},
						{Type: core.CrimeTaxFraud, Value: 6},
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

func makeKindCommercialB2B() core.Kind {
	return core.Kind{
		Name:  core.CommercialB2B,
		Label: "commerce b2b",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 40},
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
		Artifacts: []core.Artifact{
			{
				Name:             core.SmallPlant,
				Labels:           []string{"Michel & frères"},
				Description:      "fournitures en tout genres",
				CitizensCapacity: 2,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 10},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 30},
						{Type: core.Ghg3, Value: 50},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 5},
						{Type: core.BudgetCost, Value: -4},
					},
				},
			},
			{
				Name:             core.ApplianceManufacturer,
				Labels:           []string{"Moulinez", "Ubik"},
				Description:      "Fabricant d'électroménager",
				CitizensCapacity: 3,
				JobLevel:         0,
				Width:            1,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 10},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 30},
						{Type: core.Ghg3, Value: 80},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -6},
					},
				},
			},
			{
				Name:             core.CarFactory,
				Labels:           []string{"Vroum Vroum Corp"},
				Description:      "vroum vroum font les grosses voitures",
				CitizensCapacity: 8,
				JobLevel:         0,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 100},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 150},
						{Type: core.Ghg3, Value: 350},
						{Type: core.ElectricityConsumption, Value: 20},
						{Type: core.WaterConsumption, Value: 15},
						{Type: core.BudgetCost, Value: -40},
						{Type: core.CrimeTaxFraud, Value: 1},
					},
				},
			},
			{
				Name:             core.LogisticsDepot,
				Labels:           []string{"Amazone Inc"},
				Description:      "livré avant même que vous cliquiez",
				CitizensCapacity: 3,
				JobLevel:         1,
				Width:            2,
				Height:           4,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 20},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 80},
						{Type: core.Ghg3, Value: 200},
						{Type: core.ElectricityConsumption, Value: 10},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -12},
						{Type: core.CrimeTaxFraud, Value: 2},
					},
				},
			},
			{
				Name:             core.RoboticFactory,
				Labels:           []string{"IRobots"},
				Description:      "Production de robots pour l'industrie",
				CitizensCapacity: 3,
				JobLevel:         5,
				Width:            3,
				Height:           3,
				Requirements: core.Requirements{
					Spatial: []core.SpatialRequirement{
						{Type: core.MaxLandValue, Value: 30},
					},
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensEducationLevel, Value: 2},
						{Type: core.MinCitizensCount, Value: 200},
						{Type: core.MaxTaxRate, Value: 10},
					},
				},
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 20},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 100},
						{Type: core.Ghg3, Value: 200},
						{Type: core.ElectricityConsumption, Value: 30},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.BudgetCost, Value: -15},
					},
				},
			},
			{
				Name:             core.DroneFactory,
				Labels:           []string{"Narvalo Inc", "Prostate & Gamble"},
				Description:      "Fabrication de drones",
				CitizensCapacity: 4,
				JobLevel:         1,
				Width:            1,
				Height:           4,
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensCount, Value: 200},
						{Type: core.MaxTaxRate, Value: 10},
					},
				},
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 12},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 60},
						{Type: core.Ghg3, Value: 80},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.BudgetCost, Value: -6},
					},
				},
			},
		},
	}
}

func makeKindCommercialB2CLowEnd() core.Kind {
	return core.Kind{
		Name:  core.CommercialB2CLowEnd,
		Label: "commerce b2c bas de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxLandValue, Value: 10},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MaxTaxRate, Value: 50},
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
		Artifacts: []core.Artifact{
			{
				Name:             core.Pub,
				Labels:           []string{"La grosse caisse", "Le Balto", "Chez le Naf"},
				Description:      "bar",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 12},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: -1},
						{Type: core.AdultsHealth, Value: -1},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.Lavomatic,
				Labels:           []string{"Lavomatic"},
				Description:      "lavomatique",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: 2},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.FastFood,
				Labels:           []string{"Los Pollos"},
				Description:      "fast-food",
				CitizensCapacity: 2,
				JobLevel:         0,
				Width:            1,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 20},
						{Type: core.Ghg3, Value: 10},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 10},
						{Type: core.BudgetCost, Value: -2},
					},
				},
			},
			{
				Name:             core.Pharmacy,
				Labels:           []string{"pharmacie du coin"},
				Description:      "Pharmacie",
				CitizensCapacity: 1,
				JobLevel:         2,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 40},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.ElderlyHealth, Value: 1},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.CarShop,
				Labels:           []string{"Auto Boutique"},
				Description:      "Concessionnaire auto",
				CitizensCapacity: 1,
				JobLevel:         1,
				Width:            2,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 3},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 40},
						{Type: core.Ghg3, Value: 60},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.AdultsHealth, Value: -1},
						{Type: core.BudgetCost, Value: -3},
					},
				},
			},
			{
				Name:             core.BikeShop,
				Labels:           []string{"Zivélo"},
				Description:      "Magasin/réparateur de vélo",
				CitizensCapacity: 1,
				JobLevel:         1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: -2},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 1},
						{Type: core.Ghg1, Value: -1},
					},
				},
			},
			{
				Name:             core.SmartphoneFactory,
				Labels:           []string{"Zapple", "Amsung", "MinitelCo"},
				Description:      "Production de smartphones",
				CitizensCapacity: 8,
				JobLevel:         1,
				Width:            4,
				Height:           4,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 16},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 30},
						{Type: core.Ghg3, Value: 40},
						{Type: core.ElectricityConsumption, Value: 16},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.AdultsHealth, Value: -1},
						{Type: core.BudgetCost, Value: -16},
					},
				},
			},
			{
				Name:             core.DeliveryCenter,
				Labels:           []string{"Velideroo"},
				Description:      "Centre d'acheminement et de livraison",
				CitizensCapacity: 8,
				JobLevel:         0,
				Width:            4,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 6},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 24},
						{Type: core.Ghg3, Value: 190},
						{Type: core.ElectricityConsumption, Value: 12},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.AdultsHealth, Value: -2},
						{Type: core.BudgetCost, Value: -6},
					},
				},
			},
		},
	}
}

func makeKindCommercialB2CMediumEnd() core.Kind {
	return core.Kind{
		Name:  core.CommercialB2CMediumEnd,
		Label: "commerce b2c milieu de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxGhg1, Value: 12},
				{Type: core.MinLandValue, Value: 1},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MinCitizensEducationLevel, Value: 1},
				{Type: core.MaxTaxRate, Value: 35},
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
		Artifacts: []core.Artifact{
			{
				Name:             core.Hotel,
				Labels:           []string{"Ibisc hotel"},
				Description:      "hôtel moyen de gamme",
				CitizensCapacity: 4,
				JobLevel:         1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 3},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 26},
						{Type: core.Ghg3, Value: 6},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -4},
					},
				},
			},
			{
				Name:             core.Supermarket,
				Labels:           []string{"Daily Nonop"},
				Description:      "supermarché du coin",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 12},
						{Type: core.Ghg3, Value: 4},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 5},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.Restaurant,
				Labels:           []string{"Chez Bart"},
				Description:      "restaurant local",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.Cinema,
				Labels:           []string{"Mk3"},
				Description:      "salle de cinéma",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            1,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.BudgetCost, Value: -2},
					},
				},
			},
			{
				Name:             core.Theater,
				Labels:           []string{"Bouffes du nord"},
				Description:      "salle de théâtre",
				CitizensCapacity: 2,
				JobLevel:         3,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.ElectricityConsumption, Value: 3},
						{Type: core.WaterConsumption, Value: 3},
					},
				},
			},
			{
				Name:             core.ConcertHall,
				Labels:           []string{"Taratata"},
				Description:      "salle de concert",
				CitizensCapacity: 2,
				JobLevel:         2,
				Width:            1,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 15},
						{Type: core.Ghg3, Value: 10},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 3},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
		},
	}
}

func makeKindCommercialB2CHighEnd() core.Kind {
	return core.Kind{
		Name:  core.CommercialB2CHighEnd,
		Label: "commerce b2c haut de gamme",
		Requirements: core.Requirements{
			Spatial: []core.SpatialRequirement{
				{Type: core.MaxGhg1, Value: 4},
				{Type: core.MaxCrimeGeneric, Value: 3},
				{Type: core.MinLandValue, Value: 4},
			},
			Global: []core.GlobalRequirement{
				{Type: core.MinCitizensEducationLevel, Value: 4},
				{Type: core.MaxTaxRate, Value: 15},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.LandValue, Value: 3},
				{Type: core.CrimeGeneric, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.CosmeticSurgeryClinic,
				Labels:           []string{"Clinique Pimp my booty"},
				Description:      "clinique de chirurgie esthétique",
				CitizensCapacity: 1,
				JobLevel:         5,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 5},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -10},
					},
				},
			},
			{
				Name:             core.GourmetRestaurant,
				Labels:           []string{"La tour de cristal"},
				Description:      "restaurant recommandé par le guide Michelin",
				CitizensCapacity: 3,
				JobLevel:         3,
				Width:            2,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 3},
						{Type: core.Ghg3, Value: 2},
						{Type: core.ElectricityConsumption, Value: 3},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -5},
					},
				},
			},
			{
				Name:             core.ClothesStore,
				Labels:           []string{"Chein"},
				Description:      "vêtements haut de gamme",
				CitizensCapacity: 3,
				JobLevel:         1,
				Width:            2,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.Ghg3, Value: 50},
						{Type: core.ElectricityConsumption, Value: 3},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.BudgetCost, Value: -8},
					},
				},
			},
			{
				Name:             core.BigBank,
				Labels:           []string{"Too big to fail"},
				Description:      "banque d'investissement",
				CitizensCapacity: 3,
				JobLevel:         5,
				Width:            4,
				Height:           4,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.Ghg3, Value: 6},
						{Type: core.ElectricityConsumption, Value: 150},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -40},
						{Type: core.CrimeTaxFraud, Value: 40},
					},
				},
			},
			{
				Name:             core.InsuranceCompany,
				Labels:           []string{"AXXXXXA", "Assure & co"},
				Description:      "Compagnie d'assurances",
				CitizensCapacity: 2,
				JobLevel:         4,
				Width:            4,
				Height:           4,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.LandValue, Value: 6},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 6},
						{Type: core.Ghg3, Value: 6},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.BudgetCost, Value: -32},
						{Type: core.CrimeTaxFraud, Value: 16},
					},
				},
			},
		},
	}
}
