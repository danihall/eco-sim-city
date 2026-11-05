package core

//go:generate stringer -type=ArtifactStatus
//go:generate go run ./scripts/generate_enums_json_methods.go -type=ArtifactStatus
//go:generate go run ./scripts/generate_mappings.go -type=ArtifactStatus -map=TILE_STATUS -prefix=TileStatus
type ArtifactStatus uint8

const (
	Step_1_AwaitingAssignment ArtifactStatus = iota + 1
	Step_2_AssignedAwaitingDelivery
	Step_3_DeliveredOperational
	// unknown
	UnknownArtifactStatus
)

type ArtifactInstance struct {
	Label  string         `json:"label"`
	Name   ArtifactType   `json:"name"`
	Status ArtifactStatus `json:"status"`
}

type TilePosition struct {
	Q uint8 `json:"q"`
	R uint8 `json:"r"`
}

type Pollution struct {
	Ghg1 int16
}

type Desirability struct {
	LandValue int16
}

type Criminality struct {
	CrimeGeneric int16
	CrimeDrug    int16
	CrimeVss     int16
}

type TileSpatialStats struct {
	Pollution
	Desirability
	Criminality
}

// A Tile contains data of the artifact placed on it,
// also the effects stemming from it artifact.
//
// The simulation should use this to evolve a city.
type Tile struct {
	Artifact               ArtifactInstance `json:"artifact"`
	Stats                  TileSpatialStats `json:"stats"`
	Position               TilePosition     `json:"position"`
	Anchor                 TilePosition     `json:"anchor"`
	SuperKind              SuperKindType    `json:"super_kind"`
	Kind                   KindType         `json:"kind"`
	HasAllCitizenAllocated bool             `json:"has_all_citizen_allocated"`
	HasStats               bool             `json:"has_stats"`
}

type PositionGetter interface {
	Position() TilePosition
}

type TileWithJobsCount struct {
	Tile
	JobsCount uint8
}

/**
 * byRowThenColumn is to be used as a callback in a slices.SortFunc.
 * It will order by R then Q (just like the ORDER BY in SQL queries).
 */
func ByRowThenColumn(a, b Tile) int {
	tile_a := a.Position
	tile_b := b.Position

	if tile_a.R != tile_b.R {
		return int(tile_a.R) - int(tile_b.R)
	}

	return int(tile_a.Q) - int(tile_b.Q)
}

func NewTile(
	super_kind SuperKindType,
	kind KindType,
	artifact ArtifactInstance,
	stats TileSpatialStats,
	position TilePosition,
	anchor TilePosition,
	has_all_citizen_allocated bool,
) Tile {
	return Tile{
		Artifact:               artifact,
		Stats:                  stats,
		Position:               position,
		Anchor:                 anchor,
		SuperKind:              super_kind,
		Kind:                   kind,
		HasAllCitizenAllocated: has_all_citizen_allocated,
		HasStats:               true,
	}
}
