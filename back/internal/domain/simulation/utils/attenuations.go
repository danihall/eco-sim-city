package simutils

// ApplyLinearAttenuation scales a delta value based on its distance from the source,
// within a halo of given spread.
//
// At distance 0 the full delta applies, at maximum distance (spread) the delta resolves to zero.
func ApplyLinearAttenuation(delta int16, distance uint8, spread int8) int16 {
	attenuation_factor := int16(spread) - int16(distance)

	return (delta * attenuation_factor) / int16(spread)
}

// ApplyQuadraticAttenuation scales agressively a delta value based on its distance from the source,
// within a halo of the given spread.
//
// At distance 0 the full delta applies, then it quickly degrades as it moves away from the source.
func ApplyQuadraticAttenuation(delta int16, distance uint8, spread int8) int16 {
	attenuation_factor := int32(spread) - int32(distance)
	_spread := int32(spread)

	return int16(
		(int32(delta) * attenuation_factor * attenuation_factor) / (_spread * _spread),
	)
}
