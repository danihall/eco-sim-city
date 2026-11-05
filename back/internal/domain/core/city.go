package core

type CityStats struct {
	Tick uint64

	ElectricityConsumption int32
	WaterConsumption       int32

	Ghg1Level          int32 // globally
	Ghg2Level          int32 // globally
	Ghg3Level          int32 // globally
	ClimateChangeLevel int32 // Ghg1Level + Ghg2Level + Ghg3Level

	ChildrenCount int32 // globally
	AdultsCount   int32 // globally
	ElderlyCount  int32 // globally
	CitizensCount int32 // ChildrenCount + AdultsCount + ElderlyCount

	ChildrenEducationLevel int16 // globally
	AdultsEducationLevel   int16 // globally
	ElderlyEducationLevel  int16 // globally
	CitizensEducationLevel int16 // ChildrenEducationLevel + AdultsEducationLevel + ElderlyEducationLevel

	ChildrenHealthLevel int16 // globally
	AdultsHealthLevel   int16 // globally
	ElderlyHealthLevel  int16 // globally
	CitizensHealthLevel int16 // ChildrenHealthLevel + AdultsHealthLevel + ElderlyHealthLevel
}

type CityBalance struct {
	Revenues int32
	Expenses int32
}

type CityTilesArtifacts struct {
	ArtifactIDs      *[TilesCount]ArtifactType
	AnchorQs         *[TilesCount]uint8
	AnchorRs         *[TilesCount]uint8
	ArtifactStatuses *[TilesCount]ArtifactStatus
	SuperKindIDs     *[TilesCount]SuperKindType
}

const TilesCount uint16 = 128 * 128
