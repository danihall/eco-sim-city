//go:build generator

package data

import (
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// SeedSuperKindPublicSystems is the authoritative data for Public Systems artifacts. Update with caution.
// NOTE All Public Systems Artifact must have only one Label.
var SeedSuperKindPublicSystems = core.SuperKind{
	Name:  core.PublicSystems,
	Label: "services publics",
	Kinds: []core.Kind{
		makeKindWaterSystems(),
		makeKindPowerSystems(),
		makeKindEducationSystems(),
		makeKindTransportSystems(),
		makeKindHealthcareSystems(),
		makeSafetySystems(),
		makeParksAndRecSystems(),
		makeGarbageSystems(),
	},
}

func makeKindWaterSystems() core.Kind {
	return core.Kind{
		Name:  core.WaterSystems,
		Label: "production d'eau",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.Ghg1, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.WaterPumps,
				Labels:           []string{"pompe à eau"},
				Description:      "Fournit de l'eau",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: -100},
						{Type: core.BudgetCost, Value: 1},
					},
				},
			},
			{
				Name:             core.WaterTower,
				Labels:           []string{"château d'eau"},
				Description:      "Stocke de l'eau",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.WaterConsumption, Value: -20},
						{Type: core.BudgetCost, Value: 1},
					},
				},
			},
		},
	}
}

func makeKindPowerSystems() core.Kind {
	return core.Kind{
		Name:  core.PowerSystems,
		Label: "production d'électricité",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.CrimeGeneric, Value: 1},
			},
			Global: []core.GlobalEffect{
				{Type: core.Ghg2, Value: 4},
				{Type: core.Ghg3, Value: 2},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.CoalPlant,
				Labels:           []string{"centrale à charbon"},
				Description:      "Fournit de l'électricité. Pollue beaucoup",
				CitizensCapacity: 8,
				JobLevel:         0,
				Width:            4,
				Height:           4,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 300},
						{Type: core.LandValue, Value: -4},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 30},
						{Type: core.ElectricityConsumption, Value: -400},
						{Type: core.AdultsHealth, Value: -3},
						{Type: core.BudgetCost, Value: 6},
					},
				},
			},
			{
				Name:             core.OilPlant,
				Labels:           []string{"centrale à pétrole"},
				Description:      "Fournit de l'électricité. Pollue un peu moins qu'une centrale à charbon",
				CitizensCapacity: 8,
				JobLevel:         0,
				Width:            4,
				Height:           4,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 250},
						{Type: core.LandValue, Value: -4},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 40},
						{Type: core.ElectricityConsumption, Value: -500},
						{Type: core.AdultsHealth, Value: -2},
						{Type: core.BudgetCost, Value: 12},
					},
				},
			},
			{
				Name:             core.GazPlant,
				Labels:           []string{"centrale à gaz"},
				Description:      "Fournit de l'électricité. Pollue moins qu'une centrale à pétrole",
				CitizensCapacity: 5,
				JobLevel:         0,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 200},
						{Type: core.LandValue, Value: -2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 30},
						{Type: core.ElectricityConsumption, Value: -450},
						{Type: core.AdultsHealth, Value: -1},
						{Type: core.BudgetCost, Value: 12},
					},
				},
			},
			{
				Name:             core.SolarPlant,
				Labels:           []string{"centrale solaire"},
				Description:      "Ne génère pas de pollution directe mais fournit moins d'électricité",
				CitizensCapacity: 2,
				JobLevel:         0,
				Width:            2,
				Height:           3,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 20},
						{Type: core.ElectricityConsumption, Value: -100},
						{Type: core.BudgetCost, Value: 8},
					},
				},
			},
			{
				Name:             core.WindTurbine,
				Labels:           []string{"éolienne"},
				Description:      "Ne génère pas de pollution et prend peu de place mais fournit peu d'électricité",
				CitizensCapacity: 0,
				JobLevel:         1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 5},
						{Type: core.ElectricityConsumption, Value: -20},
						{Type: core.BudgetCost, Value: 4},
					},
				},
			},
			{
				Name:             core.NuclearPlant,
				Labels:           []string{"centrale à fission nucléaire"},
				Description:      "Génère très peu de pollution et fournit beaucoup d'électricité",
				CitizensCapacity: 8,
				JobLevel:         2,
				Width:            4,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.LandValue, Value: -4},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 45},
						{Type: core.ElectricityConsumption, Value: -1000},
						{Type: core.WaterConsumption, Value: 200},
						{Type: core.BudgetCost, Value: 40},
					},
				},
			},
			{
				Name:             core.FusionPlant,
				Labels:           []string{"centrale à fusion nucléaire"},
				Description:      "Fournit beaucoup d'électricité sans générer de pollution",
				CitizensCapacity: 4,
				JobLevel:         4,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 40},
						{Type: core.ElectricityConsumption, Value: -2000},
						{Type: core.WaterConsumption, Value: 300},
						{Type: core.BudgetCost, Value: 80},
					},
				},
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensCount, Value: 2_000},
					},
				},
			},
		},
	}
}

func makeKindEducationSystems() core.Kind {
	return core.Kind{
		Name:  core.EducationSystems,
		Label: "éducation publique",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
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
				Name:             core.School,
				Labels:           []string{"école primaire"},
				Description:      "Éducation de base pour les enfants",
				CitizensCapacity: 3,
				JobLevel:         1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.Ghg3, Value: 2},
						{Type: core.ElectricityConsumption, Value: 10},
						{Type: core.WaterConsumption, Value: 40},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.BudgetCost, Value: 4},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.PassCulture,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 2,
					},
					{
						Type:            core.VSSCampaign,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 1,
					},
				},
			},
			{
				Name:             core.HighSchool,
				Labels:           []string{"lycée"},
				Description:      "Éducation plus avancée pour les jeunes",
				CitizensCapacity: 4,
				JobLevel:         2,
				Width:            3,
				Height:           2,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 3},
						{Type: core.Ghg3, Value: 3},
						{Type: core.ElectricityConsumption, Value: 10},
						{Type: core.WaterConsumption, Value: 50},
						{Type: core.ChildrenEducation, Value: 2},
						{Type: core.AdultsEducation, Value: 1},
						{Type: core.BudgetCost, Value: 6},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.PassCulture,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 1,
					},
					{
						Type:            core.VSSCampaign,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 2,
					},
					{
						Type:            core.VSSCampaign,
						Effect:          core.AdultsEducation,
						AdditionalValue: 2,
					},
				},
			},
			{
				Name:             core.College,
				Labels:           []string{"université"},
				Description:      "Enseignement supérieur pour les jeunes adultes",
				CitizensCapacity: 6,
				JobLevel:         3,
				Width:            3,
				Height:           4,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 6},
						{Type: core.Ghg3, Value: 6},
						{Type: core.ElectricityConsumption, Value: 20},
						{Type: core.WaterConsumption, Value: 50},
						{Type: core.AdultsEducation, Value: 2},
						{Type: core.ElderlyEducation, Value: 1},
						{Type: core.BudgetCost, Value: 14},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.VSSCampaign,
						Effect:          core.AdultsEducation,
						AdditionalValue: 1,
					},
				},
			},
			{
				Name:             core.UniversityCampus,
				Labels:           []string{"école universitaire de recherche"},
				Description:      "Enseignement supérieur avancé pour les adultes",
				CitizensCapacity: 3,
				JobLevel:         4,
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensCount, Value: 1_000},
						{Type: core.MinCitizensEducationLevel, Value: 4},
						{Type: core.MinFundingRate, Value: 30},
					},
				},
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 3},
						{Type: core.Ghg3, Value: 4},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.AdultsEducation, Value: 4},
						{Type: core.ElderlyEducation, Value: 3},
						{Type: core.BudgetCost, Value: 4},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.VSSCampaign,
						Effect:          core.AdultsEducation,
						AdditionalValue: 1,
					},
				},
			},
			{
				Name:             core.Library,
				Labels:           []string{"bibliothèque"},
				Description:      "Bibliothèque municipale pour tous publics",
				CitizensCapacity: 1,
				JobLevel:         1,
				Width:            2,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 1},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 30},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.AdultsEducation, Value: 1},
						{Type: core.ElderlyEducation, Value: 1},
						{Type: core.BudgetCost, Value: 2},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.PassCulture,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 2,
					},
				},
			},
			{
				Name:             core.Museum,
				Labels:           []string{"musée"},
				Description:      "Expositions et conférences pour tous publics",
				CitizensCapacity: 2,
				JobLevel:         2,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.Ghg3, Value: 2},
						{Type: core.ElectricityConsumption, Value: 30},
						{Type: core.WaterConsumption, Value: 40},
						{Type: core.ChildrenEducation, Value: 2},
						{Type: core.AdultsEducation, Value: 2},
						{Type: core.ElderlyEducation, Value: 3},
						{Type: core.BudgetCost, Value: 16},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.PassCulture,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 2,
					},
				},
			},
			{
				Name:             core.RadioCenter,
				Labels:           []string{"Centre radio"},
				Description:      "Émissions radio pour tous publics",
				CitizensCapacity: 6,
				JobLevel:         2,
				Width:            3,
				Height:           2,
				IsUnique:         true,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 30},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.AdultsEducation, Value: 1},
						{Type: core.ElderlyEducation, Value: 1},
						{Type: core.BudgetCost, Value: 12},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.VSSCampaign,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 1,
					},
					{
						Type:            core.VSSCampaign,
						Effect:          core.AdultsEducation,
						AdditionalValue: 1,
					},
					{
						Type:            core.VSSCampaign,
						Effect:          core.ElderlyEducation,
						AdditionalValue: 1,
					},
				},
			},
		},
	}
}

func makeKindTransportSystems() core.Kind {
	return core.Kind{
		Name:  core.TransportSystems,
		Label: "transports",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
		Effects: core.Effects{
			Global: []core.GlobalEffect{
				{Type: core.Ghg2, Value: 1},
				{Type: core.Ghg3, Value: 4},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 3},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.Road,
				Labels:           []string{"route"},
				Description:      "Permet aux voitures, bus, vélos et piétons de circuler",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.BudgetCost, Value: 1},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Car, Value: 1},
						{Type: core.Bus, Value: 1},
						{Type: core.Bicycle, Value: 3},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.FreeBicycles,
						Effect:          core.Bicycle,
						AdditionalValue: -1,
					},
				},
			},
			{
				Name:             core.Rail,
				Labels:           []string{"voie ferrée"},
				Description:      "Permet aux trains de circuler sur le réseau ferré",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.BudgetCost, Value: 1},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Train, Value: 1},
					},
				},
			},
			{
				Name:             core.SubwayStation,
				Labels:           []string{"station de métro"},
				Description:      "Chaque station fait partie du réseau de métro souterrain",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeVss, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.ElectricityConsumption, Value: 3},
						{Type: core.BudgetCost, Value: 20},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Subway, Value: 1},
					},
				},
			},
			{
				Name:             core.BusStop,
				Labels:           []string{"arrêt de bus"},
				Description:      "Chaque station fait partie du réseau de bus. Les bus circulent sur les routes",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.BudgetCost, Value: 1},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Bus, Value: 1},
					},
				},
			},
			{
				Name:             core.TrainStation,
				Labels:           []string{"gare"},
				Description:      "Chaque gare fait partie du réseau ferré. Les trains circulent sur les voies ferrées",
				CitizensCapacity: 2,
				JobLevel:         0,
				Width:            2,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeDrug, Value: 2},
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.ElectricityConsumption, Value: 5},
						{Type: core.BudgetCost, Value: 8},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Train, Value: 1},
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
			{
				Name:             core.PedestrianBicycleWay,
				Labels:           []string{"voie piétonne"},
				Description:      "Voie réservée aux piétons, cyclistes et bus",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: 1},
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: 1},
						{Type: core.BudgetCost, Value: 1},
					},
					TravelCost: []core.TravelCostEffect{
						{Type: core.Foot, Value: 2},
						{Type: core.Bicycle, Value: 1},
						{Type: core.Bus, Value: 3},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.FreeBicycles,
						Effect:          core.Bicycle,
						AdditionalValue: -1,
					},
				},
			},
		},
	}
}

func makeKindHealthcareSystems() core.Kind {
	return core.Kind{
		Name:  core.HealthcareSystems,
		Label: "santé publique",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 10},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.LandValue, Value: 2},
			},
			Global: []core.GlobalEffect{
				{Type: core.Ghg3, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.Nursery,
				Labels:           []string{"crèche"},
				Description:      "La crèche municipale permet de garder les bambins pendant la journée",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 1},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 1},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.BudgetCost, Value: 1},
					},
				},
			},
			{
				Name:             core.MedicalCenter,
				Labels:           []string{"centre médical"},
				Description:      "Le centre médical permet aux habitants de se soigner",
				CitizensCapacity: 1,
				JobLevel:         2,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.ChildrenHealth, Value: 2},
						{Type: core.AdultsHealth, Value: 2},
						{Type: core.ElderlyHealth, Value: 2},
						{Type: core.BudgetCost, Value: 5},
					},
				},
			},
			{
				Name:             core.Hospital,
				Labels:           []string{"hôpital"},
				Description:      "Un hôpital augmente significativement le santé des habitants",
				CitizensCapacity: 5,
				JobLevel:         3,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: 3},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 20},
						{Type: core.Ghg3, Value: 4},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.ChildrenHealth, Value: 4},
						{Type: core.AdultsHealth, Value: 4},
						{Type: core.ElderlyHealth, Value: 5},
						{Type: core.BudgetCost, Value: 30},
					},
				},
			},
			{
				Name:             core.RetirementHome,
				Labels:           []string{"maison de retraite"},
				Description:      "Garde les personnes âgées en bonne santé",
				CitizensCapacity: 4,
				JobLevel:         2,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 20},
						{Type: core.Ghg3, Value: 1},
						{Type: core.ElectricityConsumption, Value: 2},
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.ElderlyHealth, Value: 1},
						{Type: core.BudgetCost, Value: 6},
					},
				},
			},
		},
	}
}

func makeSafetySystems() core.Kind {
	return core.Kind{
		Name:  core.SafetySystems,
		Label: "sécurité publique",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.LandValue, Value: 1},
				{Type: core.CrimeGeneric, Value: -3},
			},
			Global: []core.GlobalEffect{
				{Type: core.Ghg3, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.PoliceKiosk,
				Labels:           []string{"poste de police"},
				Description:      "Un petit poste de police local diminue la criminalité aux alentours",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeVss, Value: -1},
						{Type: core.CrimeDrug, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.ElectricityConsumption, Value: 1},
						{Type: core.WaterConsumption, Value: 1},
						{Type: core.BudgetCost, Value: 1},
					},
				},
			},
			{
				Name:             core.PoliceStation,
				Labels:           []string{"commissariat"},
				Description:      "Un commissariat diminue la criminalité sur une large zone",
				CitizensCapacity: 3,
				JobLevel:         1,
				Width:            2,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeGeneric, Value: -5},
						{Type: core.CrimeVss, Value: -1},
						{Type: core.CrimeDrug, Value: -3},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.ElectricityConsumption, Value: 5},
						{Type: core.WaterConsumption, Value: 5},
						{Type: core.BudgetCost, Value: 8},
					},
				},
			},
			{
				Name:             core.FireStation,
				Labels:           []string{"caserne de pompiers"},
				Description:      "Réduit les dégâts liés aux incendies et intervient en premiers secours",
				CitizensCapacity: 1,
				JobLevel:         1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeGeneric, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.Ghg3, Value: 2},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.AdultsHealth, Value: 1},
						{Type: core.ElderlyHealth, Value: 1},
						{Type: core.BudgetCost, Value: 10},
					},
				},
			},
			{
				Name:             core.CityJail,
				Labels:           []string{"prison"},
				Description:      "Réduit légèrement la criminalité dans toute la ville",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            3,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: -10},
						{Type: core.CrimeGeneric, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 3},
						{Type: core.AdultsHealth, Value: -2},
						{Type: core.BudgetCost, Value: 6},
					},
				},
			},
			{
				Name:             core.Courthouse,
				Labels:           []string{"tribunal"},
				Description:      "Un tribunal réduit surtout la fraude fiscale mais pas toute la criminalité immédiate",
				CitizensCapacity: 5,
				JobLevel:         3,
				Width:            3,
				Height:           3,
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensCount, Value: 500},
					},
				},
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: 4},
						{Type: core.CrimeGeneric, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 4},
						{Type: core.ElectricityConsumption, Value: 8},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.BudgetCost, Value: 30},
						{Type: core.CrimeTaxFraud, Value: -128},
					},
				},
			},
			{
				Name:             core.MeteoCenter,
				Labels:           []string{"CMN"},
				Description:      "Centre Météorologique National",
				CitizensCapacity: 6,
				JobLevel:         4,
				Width:            3,
				Height:           3,
				IsUnique:         true,
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensCount, Value: 1_000},
					},
				},
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.ElectricityConsumption, Value: 12},
						{Type: core.WaterConsumption, Value: 8},
						{Type: core.ChildrenHealth, Value: 4},
						{Type: core.AdultsHealth, Value: 4},
						{Type: core.ElderlyHealth, Value: 4},
						{Type: core.BudgetCost, Value: 70},
					},
				},
			},
		},
	}
}

func makeParksAndRecSystems() core.Kind {
	return core.Kind{
		Name:  core.ParksAndRecSystems,
		Label: "parcs et divertissements",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
		Effects: core.Effects{
			Spatial: []core.SpatialEffect{
				{Type: core.LandValue, Value: 3},
			},
			Global: []core.GlobalEffect{
				{Type: core.ElectricityConsumption, Value: 1},
				{Type: core.WaterConsumption, Value: 1},
			},
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.SmallPark,
				Labels:           []string{"square"},
				Description:      "Augmente légèrement la valeur du terrain aux alentours. Réduit aussi légèrement la pollution",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -2},
					},
					Global: []core.GlobalEffect{
						{Type: core.BudgetCost, Value: 2},
					},
				},
			},
			{
				Name:             core.MediumPark,
				Labels:           []string{"jardin public"},
				Description:      "Augmente la valeur du terrain aux alentours. Réduit aussi légèrement la pollution",
				CitizensCapacity: 1,
				JobLevel:         0,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -3},
					},
					Global: []core.GlobalEffect{
						{Type: core.WaterConsumption, Value: 2},
						{Type: core.BudgetCost, Value: 4},
					},
				},
			},
			{
				Name:             core.LargePark,
				Labels:           []string{"grand parc"},
				Description:      "Augmente significativement la valeur du terrain aux alentours. Réduit aussi la pollution",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -6},
						{Type: core.LandValue, Value: 8},
					},
					Global: []core.GlobalEffect{
						{Type: core.ElectricityConsumption, Value: 4},
						{Type: core.WaterConsumption, Value: 3},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.AdultsHealth, Value: 1},
						{Type: core.ElderlyHealth, Value: 1},
						{Type: core.BudgetCost, Value: 9},
					},
				},
			},
			{
				Name:             core.CommunityGarden,
				Labels:           []string{"jardin collectif"},
				Description:      "Réduit la pollution directe et indirecte",
				CitizensCapacity: 1,
				JobLevel:         1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -1},
						{Type: core.LandValue, Value: 1},
						{Type: core.CrimeGeneric, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: -3},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.AdultsHealth, Value: 1},
						{Type: core.ElderlyHealth, Value: 1},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.AdultsEducation, Value: 1},
						{Type: core.ElderlyEducation, Value: 1},
						{Type: core.BudgetCost, Value: 1},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.VSSCampaign,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 1,
					},
				},
			},
			{
				Name:             core.Arboretum,
				Labels:           []string{"arboretum"},
				Description:      "Réduit la pollution directe et indirecte",
				CitizensCapacity: 0,
				JobLevel:         0,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -7},
						{Type: core.LandValue, Value: 4},
					},
					Global: []core.GlobalEffect{
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.AdultsEducation, Value: 1},
						{Type: core.ElderlyEducation, Value: 1},
						{Type: core.BudgetCost, Value: 5},
					},
				},
				OrdinanceEffects: []core.OrdinanceEffect{
					{
						Type:            core.PassCulture,
						Effect:          core.ChildrenEducation,
						AdditionalValue: 2,
					},
				},
			},
			{
				Name:             core.Fountain,
				Labels:           []string{"fontaine publique"},
				Description:      "Très utile lors d'épisodes caniculaires",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -1},
						{Type: core.LandValue, Value: 4},
					},
					Global: []core.GlobalEffect{
						{Type: core.WaterConsumption, Value: 6},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.ElderlyHealth, Value: 1},
						{Type: core.BudgetCost, Value: 2},
					},
				},
			},
			{
				Name:             core.Playground,
				Labels:           []string{"aire de jeux"},
				Description:      "À destination d'un jeune public",
				CitizensCapacity: 0,
				JobLevel:         0,
				Width:            1,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.LandValue, Value: 2},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.ChildrenEducation, Value: 1},
						{Type: core.ChildrenHealth, Value: 1},
						{Type: core.BudgetCost, Value: 4},
					},
				},
			},
			{
				Name:             core.Stadium,
				Labels:           []string{"stade"},
				Description:      "Accueille l'équipe locale et génère un peu de revenus",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 1},
						{Type: core.CrimeGeneric, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 8},
						{Type: core.Ghg3, Value: 12},
						{Type: core.ElectricityConsumption, Value: 6},
						{Type: core.WaterConsumption, Value: 4},
						{Type: core.BudgetCost, Value: -6},
					},
				},
				Requirements: core.Requirements{
					Global: []core.GlobalRequirement{
						{Type: core.MinCitizensCount, Value: 700},
					},
				},
			},
			{
				Name:             core.SmallPlaza,
				Labels:           []string{"petite plaza"},
				Description:      "petite place publique urbaine",
				CitizensCapacity: 1,
				JobLevel:         1,
				Width:            1,
				Height:           1,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 1},
						{Type: core.Ghg3, Value: 1},
						{Type: core.BudgetCost, Value: 1},
					},
				},
			},
			{
				Name:             core.MediumPlaza,
				Labels:           []string{"plaza moyenne"},
				Description:      "place publique urbaine moyenne",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            2,
				Height:           1,
				Effects: core.Effects{
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 1},
						{Type: core.Ghg3, Value: 1},
						{Type: core.BudgetCost, Value: -1},
					},
				},
			},
			{
				Name:             core.LargePlaza,
				Labels:           []string{"grande plaza"},
				Description:      "grande place publique urbaine",
				CitizensCapacity: 3,
				JobLevel:         1,
				Width:            2,
				Height:           2,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.CrimeGeneric, Value: 1},
						{Type: core.CrimeVss, Value: 1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 6},
						{Type: core.Ghg3, Value: 6},
						{Type: core.BudgetCost, Value: -5},
					},
				},
			},
		},
	}
}

func makeGarbageSystems() core.Kind {
	return core.Kind{
		Name:  core.GarbageSystems,
		Label: "gestion des déchets",
		Requirements: core.Requirements{
			Global: []core.GlobalRequirement{
				{Type: core.MinFundingRate, Value: 20},
			},
		},
		Effects: core.Effects{
			TravelCost: []core.TravelCostEffect{
				{Type: core.Foot, Value: 4},
			},
		},
		Artifacts: []core.Artifact{
			{
				Name:             core.Landfill,
				Labels:           []string{"déchetterie"},
				Description:      "Stocke les déchets générés par la ville. Les déchets entassés disparaissent au fil du temps",
				CitizensCapacity: 1,
				JobLevel:         0,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 30},
						{Type: core.LandValue, Value: -15},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg3, Value: -1},
						{Type: core.ChildrenHealth, Value: -1},
						{Type: core.AdultsHealth, Value: -1},
						{Type: core.ElderlyHealth, Value: -1},
					},
				},
			},
			{
				Name:             core.WasteIncinerator,
				Labels:           []string{"incinérateur de déchets"},
				Description:      "Génère de l'électricité à partir de déchets mais pollue beaucoup",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 250},
						{Type: core.LandValue, Value: -8},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 2},
						{Type: core.Ghg3, Value: -6},
						{Type: core.ChildrenHealth, Value: -2},
						{Type: core.AdultsHealth, Value: -2},
						{Type: core.ElderlyHealth, Value: -2},
						{Type: core.ElectricityConsumption, Value: -80},
						{Type: core.BudgetCost, Value: 6},
					},
				},
			},
			{
				Name:             core.RecyclingCenter,
				Labels:           []string{"centre de recyclage"},
				Description:      "Réduit légèrement la pollution indirecte",
				CitizensCapacity: 2,
				JobLevel:         1,
				Width:            3,
				Height:           3,
				Effects: core.Effects{
					Spatial: []core.SpatialEffect{
						{Type: core.Ghg1, Value: 2},
						{Type: core.LandValue, Value: -1},
					},
					Global: []core.GlobalEffect{
						{Type: core.Ghg2, Value: 10},
						{Type: core.Ghg3, Value: -48},
						{Type: core.ElectricityConsumption, Value: 20},
						{Type: core.WaterConsumption, Value: 20},
						{Type: core.BudgetCost, Value: 30},
					},
				},
			},
		},
	}
}
