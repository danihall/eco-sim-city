package core

//go:generate stringer -type=RequirementType
type RequirementType uint8

const (
	// Spatial
	MaxGhg1 RequirementType = iota
	MinLandValue
	MaxLandValue
	MaxCrimeGeneric
	MaxCrimeDrug
	MaxCrimeVss
	// Global
	MinCitizensCount
	MinCitizensEducationLevel
	MaxTaxRate
	MinFundingRate
)

type ConditionType uint8

const (
	Min ConditionType = iota
	Max
)

var MapRequirementsCondition = map[RequirementType]ConditionType{
	MaxGhg1:                   Max,
	MinLandValue:              Min,
	MaxLandValue:              Max,
	MaxCrimeGeneric:           Max,
	MaxCrimeDrug:              Max,
	MaxCrimeVss:               Max,
	MinCitizensCount:          Min,
	MinCitizensEducationLevel: Min,
	MaxTaxRate:                Max,
	MinFundingRate:            Min,
}

type SpatialRequirement struct {
	Type  RequirementType
	Value int16
}

type GlobalRequirement struct {
	Type  RequirementType
	Value int16
}

// Requirements are the counter-part of Effects.
// For an Artifact to be placed on a tile, Requirements of the Artifact should be met.
type Requirements struct {
	Spatial []SpatialRequirement
	Global  []GlobalRequirement
}
