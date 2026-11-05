package cityhall

import (
	"errors"
	"fmt"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
)

var (
	ErrUniqueArtifactAlreadyExists = &domain.Error{
		Reason: errors.New("some artifacts are unique"),
		Key:    domain.ErrNotModified,
	}
	ErrTileOccupied = &domain.Error{
		Reason: errors.New("cannot build here"),
		Key:    domain.ErrConflict,
	}
	ErrIneligible = &domain.Error{
		Reason: errors.New("artifact's requirements are not met"),
		Key:    domain.ErrConflict,
	}
)

type PermitValidator struct {
}

// Permits checks if an artifact can be built depending on artifact's requirements.
//
// see core/requirement.go
type Permits struct {
	Validator PermitValidator
}

func NewPermits() *Permits {
	return &Permits{
		Validator: PermitValidator{},
	}
}

// UserCandidate delivers permission for an artifact to be built from a user action.
// Internally, it calls multiple functions that each check for a specific set of conditions.
//   - CheckTilesOccupancy
//   - CheckUniqueArtifact
//   - CheckGlobalRequirements
//   - CheckSpatialRequirements
//
// If every tile is occupied, returns an error.
// If artifact is not free-size, returns an error if one or more tiles are occupied.
// Since user can build artifacts that has free-size, the given tiles are filtered IN-PLACE.
// Candidate returns the filtered slice of `Tile`
func (p PermitValidator) UserCandidate(
	city CityFullState,
	artifact *core.Artifact,
	tiles []core.Tile,
) ([]core.Tile, error) {
	err := p.CheckUniqueArtifact(city, artifact)
	if err != nil {
		return nil, err
	}

	unoccupied_tiles := p.MakeUnoccupiedTiles(tiles)

	if (!artifact.HasFreeSize() && len(unoccupied_tiles) < len(tiles)) ||
		len(unoccupied_tiles) == 0 {
		return nil, fmt.Errorf("%w, not enough unoccupied tiles", ErrTileOccupied)
	}

	err = p.CheckGlobalRequirements(city, artifact)
	if err != nil {
		return nil, err
	}

	err = p.CheckSpatialRequirements(artifact, unoccupied_tiles)
	if err != nil {
		return nil, err
	}

	return unoccupied_tiles, nil
}

func (p PermitValidator) CheckUniqueArtifact(
	city CityFullState,
	artifact *core.Artifact,
) error {
	if artifact.IsUnique && slices.Contains(city.Flagships, artifact.Name) {
		return fmt.Errorf(
			"%w, flagship %v already exists",
			ErrUniqueArtifactAlreadyExists,
			artifact.Name,
		)
	}

	return nil
}

// CheckGlobalRequirements returns an error if those conditions fail:
//   - If candidate artifact is unique in the whole city and is already present, permit is DENIED.
//   - If candidate artifact's global requirements (min citizen count & min education) are not met,
//     permit is DENIED.
//   - If candidate artifact's fiscal requirements (min funding rate & max tax rate) are not met,
//     permit is DENIED.
func (p PermitValidator) CheckGlobalRequirements(
	city CityFullState,
	artifact *core.Artifact,
) error {
	if len(artifact.Requirements.Global) == 0 {
		return nil
	}

	city_stats := city.Stats
	min_citizen_count := artifact.GlobalRequirement(core.MinCitizensCount)
	min_citizen_education_level := artifact.GlobalRequirement(core.MinCitizensEducationLevel)

	if city_stats.CitizensCount < int32(min_citizen_count) ||
		city_stats.CitizensEducationLevel < min_citizen_education_level {
		return fmt.Errorf(
			"%w: citizen count: %d, citizen education lv: %d but %s requires %s: %d, %s: %d",
			ErrIneligible,
			city_stats.CitizensCount,
			city_stats.CitizensEducationLevel,
			artifact.Name,
			core.MinCitizensCount,
			min_citizen_count,
			core.MinCitizensEducationLevel,
			min_citizen_education_level,
		)
	}

	kind_name := registry.ArtifactsRelations[artifact.Name].Kind
	kind_budget := city.Budget[kind_name]
	max_tax_rate := artifact.GlobalRequirement(core.MaxTaxRate)
	min_funding_rate := artifact.GlobalRequirement(core.MinFundingRate)

	if kind_budget.TaxRate > uint8(max_tax_rate) ||
		kind_budget.FundingRate < uint8(min_funding_rate) {
		return fmt.Errorf("%w: %s tax rate:%d, funding rate:%d but %s requires %s: %d, %s: %d",
			ErrIneligible,
			kind_budget.KindName,
			kind_budget.TaxRate,
			kind_budget.FundingRate,
			kind_budget.KindName,
			core.MaxTaxRate,
			max_tax_rate,
			core.MinFundingRate,
			min_funding_rate,
		)
	}

	return nil
}

// CheckSpatialRequirements returns an error if this condition fails:
//   - If more than 25% of `tile`'s spatial effects on which candidate artifact wants to be placed
//     do not meet artifact's spatial requirements, permits is DENIED.
func (p PermitValidator) CheckSpatialRequirements(
	artifact *core.Artifact,
	tiles []core.Tile,
) error {
	if len(tiles) == 0 || len(artifact.Requirements.Spatial) == 0 {
		return nil
	}

	unmet_spatial_requirements := 0

	max_ghg_1 := artifact.SpatialRequirement(core.MaxGhg1)
	max_crime_generic := artifact.SpatialRequirement(core.MaxCrimeGeneric)
	min_land_value := artifact.SpatialRequirement(core.MinLandValue)
	max_land_value := artifact.SpatialRequirement(core.MaxLandValue)
	max_crime_vss := artifact.SpatialRequirement(core.MaxCrimeVss)
	max_crime_drug := artifact.SpatialRequirement(core.MaxCrimeDrug)

	for index := range tiles {
		if !tiles[index].HasStats {
			return fmt.Errorf("missing Stats in tile")
		}

		stats := tiles[index].Stats

		switch {
		case stats.Ghg1 > max_ghg_1,
			stats.CrimeGeneric > max_crime_generic,
			stats.LandValue > max_land_value,
			stats.LandValue < min_land_value,
			stats.CrimeDrug > max_crime_drug,
			stats.CrimeVss > max_crime_vss:
			unmet_spatial_requirements++
		}
	}

	// Allow up to 25% of spatial requirments to not be met
	if unmet_spatial_requirements*4 > len(tiles) {
		return fmt.Errorf(
			"%w: too many of these requirements were not met: %s, %s, %s, %s, %s",
			ErrIneligible,
			core.MaxGhg1,
			core.MaxCrimeGeneric,
			core.MinLandValue,
			core.MaxCrimeVss,
			core.MaxCrimeDrug,
		)
	}

	return nil
}

// makeUnoccupiedTiles filters unoccupied tiles.
func (p PermitValidator) MakeUnoccupiedTiles(tiles []core.Tile) []core.Tile {
	unoccupied_tiles := make([]core.Tile, 0, len(tiles))

	for index := range tiles {
		tile := tiles[index]
		if tile.SuperKind != core.Natural &&
			tile.Artifact.Status != core.Step_1_AwaitingAssignment {
			continue
		}

		unoccupied_tiles = append(unoccupied_tiles, tile)
	}

	return unoccupied_tiles
}
