package core

import (
	"math"
)

// HasFreeSize returns `true` if both width and height's artifact are 0.
// NOTE any free-size artifact
func (a *Artifact) HasFreeSize() bool {
	return a.Width == 0 && a.Height == 0
}

func (a *Artifact) EffectiveSize(user_width, user_height uint8) (uint8, uint8) {
	if (a.Width == 0 && user_width == 0) || (a.Height == 0 && user_height == 0) {
		return 0, 0
	}

	var effective_width uint8
	var effective_height uint8

	if a.Width == 0 {
		effective_width = user_width
	} else {
		effective_width = a.Width
	}

	if a.Height == 0 {
		effective_height = user_height
	} else {
		effective_height = a.Height
	}

	return effective_width, effective_height
}

func (a *Artifact) HasRequirements() bool {
	if len(a.Requirements.Spatial) > 0 {
		return true
	}
	if len(a.Requirements.Global) > 0 {
		return true
	}
	return false
}

// SpatialEffect loops in given artifact's spatial effects to return desired effect value.
func (a *Artifact) SpatialEffect(effect_type EffectType) int16 {
	for _, item := range a.Effects.Spatial {
		if item.Type == effect_type {
			return item.Value
		}
	}

	return 0
}

// GlobalEffect loops in given artifact's spatial effects to return desired effect value.
func (a *Artifact) GlobalEffect(effect_type EffectType) int16 {
	for _, item := range a.Effects.Global {
		if item.Type == effect_type {
			return item.Value
		}
	}

	return 0
}

func (a *Artifact) IsTraversableBy(transport_mode EffectType) bool {
	for _, item := range a.Effects.TravelCost {
		if item.Type == transport_mode {
			return true
		}
	}

	return false
}

// TravelCostEffect loops in given artifact's spatial effects to return desired effect value.
func (a *Artifact) TravelCostEffect(effect_type EffectType) uint8 {
	for _, item := range a.Effects.TravelCost {
		if item.Type == effect_type {
			return item.Value
		}
	}

	return 0
}

// SpatialRequirement guarantees that if a specific spatial requirment was not set, a Min (0) or Max (32_767) value is used.
func (a *Artifact) SpatialRequirement(requirement_type RequirementType) int16 {
	for _, item := range a.Requirements.Spatial {
		if item.Type == requirement_type {
			return item.Value
		}
	}

	if MapRequirementsCondition[requirement_type] == Max {
		return math.MaxInt16
	}

	return 0
}

// GlobalRequirement guarantees that if a specific spatial requirment was not set, a Min (0) or Max (32_767) value is used.
func (a *Artifact) GlobalRequirement(requirement_type RequirementType) int16 {
	for _, item := range a.Requirements.Global {
		if item.Type == requirement_type {
			return item.Value
		}
	}

	if MapRequirementsCondition[requirement_type] == Max {
		return math.MaxInt16
	}

	return 0
}
