package core

//go:generate stringer -type=SuperKindType
//go:generate go run ./scripts/generate_enums_json_methods.go -type=SuperKindType
//go:generate go run ./scripts/generate_mappings.go -type=SuperKindType -map=SUPER_KINDS -prefix=SuperKind
type SuperKindType uint8

const (
	Natural SuperKindType = iota + 1
	PublicSystems
	Industrial
	Residential
	// unkown
	UnknownSuperKind
)

// A SuperKind represents one of the 4 "zone type" used by the simulation.
//
// Industrial and Residential SuperKind's artifacts should not be built directly by a user,
// they should be placed by the simulation.
type SuperKind struct {
	Name  SuperKindType
	Label string
	Kinds []Kind
}
