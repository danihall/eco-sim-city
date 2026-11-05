package core

func (e EffectType) IsPublicTransitMode() bool {
	return e == Bus || e == Train || e == Subway
}
