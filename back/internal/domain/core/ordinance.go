package core

//go:generate stringer -type=OrdinanceType
//go:generate go run ./scripts/generate_enums_json_methods.go -type=OrdinanceType
type OrdinanceType uint8

const (
	DrugPrevention OrdinanceType = iota + 1
	PassCulture
	VSSCampaign
	FreeBicycles
	// unknown
	UnknownOrdinance
)

// An Ordinance, if enacted, must be taken into account when computing the Effects of an Artifact.
// It acts as an additional value for an Effect.
type Ordinance struct {
	Name        OrdinanceType
	Label       string
	Description string
}

type OrdinanceEffect struct {
	Type            OrdinanceType
	Effect          EffectType
	AdditionalValue int16
}
