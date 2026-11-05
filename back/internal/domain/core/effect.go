package core

//go:generate stringer -type=EffectType
type EffectType uint8

const (
	// Spatial
	Ghg1 EffectType = iota + 1
	LandValue
	CrimeGeneric
	CrimeDrug
	CrimeVss
	// Global
	Ghg2
	Ghg3
	ElectricityConsumption
	WaterConsumption
	ChildrenEducation
	AdultsEducation
	ElderlyEducation
	ChildrenHealth
	AdultsHealth
	ElderlyHealth
	BudgetCost
	CrimeTaxFraud
	// Travel Cost
	Car
	Bicycle
	Bus
	Foot
	Subway
	Train
	// unknown
	UnknownEffect
)

type SpatialEffect struct {
	Type  EffectType
	Value int16
}

type GlobalEffect struct {
	Type  EffectType
	Value int16
}

type TravelCostEffect struct {
	Type  EffectType
	Value uint8
}

// Effects represent the impact of an Artifact.
// For any Artifact, each effect can be inherited from its parent Kind's effect.
//
// NOTE TravelCost.Foot, if not set on a Kind or Artifact means that the tile is not traversible by foot.
//
// Same goes for:
//   - TravelCost.Car
//   - TravelCost.Bicycle
//   - TravelCost.Bus
//   - TravelCost.Subway
//   - TravelCost.Train
//
// for their implied mode of transport.
type Effects struct {
	Spatial    []SpatialEffect
	Global     []GlobalEffect
	TravelCost []TravelCostEffect
}
