package core

//go:generate go run -tags=generator ./scripts/generate_artifacts_size.go
//go:generate stringer -type=ArtifactType
//go:generate go run ./scripts/generate_enums_json_methods.go -type=ArtifactType
//go:generate go run ./scripts/generate_mappings.go -type=ArtifactType -map=ARTIFACTS -prefix=Artifact
type ArtifactType uint8

const (
	Gingko ArtifactType = iota + 1
	Birch
	Oak
	Paulownia
	MixedForest
	ShallowWater
	Grassland
	Wasteland
	DeepWater
	Bog
	WaterPumps
	WaterTower
	CoalPlant
	OilPlant
	GazPlant
	SolarPlant
	WindTurbine
	NuclearPlant
	FusionPlant
	School
	HighSchool
	College
	UniversityCampus
	Library
	Museum
	RadioCenter
	Road
	Rail
	SubwayStation
	BusStop
	TrainStation
	PedestrianBicycleWay
	Nursery
	MedicalCenter
	Hospital
	RetirementHome
	PoliceKiosk
	PoliceStation
	FireStation
	CityJail
	Courthouse
	MeteoCenter
	SmallPark
	MediumPark
	LargePark
	CommunityGarden
	Arboretum
	Fountain
	Playground
	Stadium
	SmallPlaza
	MediumPlaza
	LargePlaza
	Landfill
	WasteIncinerator
	RecyclingCenter
	ChemicalPlant
	PlasticsFactory
	SteelMill
	AluminiumFactory
	CementFactory
	LumberYard
	TextilesFactory
	FoodProcessingPlant
	BiologicalFarm
	Esn
	VideoGameStartup
	DataCenter
	AiStartup
	SmallPlant
	ApplianceManufacturer
	CarFactory
	LogisticsDepot
	RoboticFactory
	DroneFactory
	Pub
	Lavomatic
	FastFood
	Pharmacy
	CarShop
	BikeShop
	SmartphoneFactory
	DeliveryCenter
	Hotel
	Supermarket
	Restaurant
	Cinema
	Theater
	ConcertHall
	CosmeticSurgeryClinic
	GourmetRestaurant
	ClothesStore
	BigBank
	InsuranceCompany
	HouseLow
	ResidentialBuildingLow
	LargeResidentialTowerLow
	HouseMedium
	ResidentialBuildingMedium
	LargeResidentialTowerMedium
	HouseHigh
	ResidentialBuildingHigh
	LargeResidentialTowerHigh
	HouseUltraHigh
	ResidentialBuildingUltraHigh
	LargeResidentialTowerUltraHigh
	// unknown
	UnknownArtifact
)

// An Artifact is a building or a natural element placed on the city's grid.
// Artifacts can spread on multiple tiles (Width & Height).
// Artifacts always belong to a parent Kind.
type Artifact struct {
	Name             ArtifactType
	Labels           []string
	Description      string
	CitizensCapacity uint8
	JobLevel         uint8
	Width            uint8 // if 0, means it has a free width
	Height           uint8 // if 0, means it has a free height
	IsUnique         bool
	Requirements     Requirements
	Effects          Effects
	OrdinanceEffects []OrdinanceEffect
}

type ArtifactRelation struct {
	SuperKind SuperKindType
	Kind      KindType
}
