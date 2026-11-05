package core

//go:generate stringer -type=KindType
//go:generate go run ./scripts/generate_enums_json_methods.go -type=KindType
type KindType uint8

const (
	Trees KindType = iota + 1
	Unconstructible
	WaterSystems
	PowerSystems
	EducationSystems
	TransportSystems
	HealthcareSystems
	SafetySystems
	ParksAndRecSystems
	GarbageSystems
	Petrochemical
	Metallurgy
	Construction
	Textiles
	Agrofood
	Digital
	CommercialB2B
	CommercialB2CLowEnd
	CommercialB2CMediumEnd
	CommercialB2CHighEnd
	ResidentialLowEnd
	ResidentialMediumEnd
	ResidentialHighEnd
	ResidentialUltraHighEnd
	// unknown
	UnknownKind
)

// A Kind is a child of a SuperKind.
type Kind struct {
	Name         KindType
	Label        string
	Effects      Effects
	Requirements Requirements
	Artifacts    []Artifact
	*SocialClass
}

// A SocialClass applies only to each `Residential` SuperKind's Kind.
type SocialClass struct {
	MinJobLevel                 uint8
	AdditionalTravelCostFoot    uint8
	AdditionalTravelCostCar     uint8
	AdditionalTravelCostBus     uint8
	AdditionalTravelCostBicycle uint8
	AdditionalTravelCostTrain   uint8
	AdditionalTravelCostSubway  uint8
}

type KindBudget struct {
	KindName    KindType
	TaxRate     uint8
	FundingRate uint8
}
