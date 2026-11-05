package api

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

// DeleteCity is a user action.
// TODO maybe use an event to stop simulating city instead of doing it in this function
func DeleteCity(
	ctx context.Context,
	c *cityhall.Deputies,
	s *gameloops.Simulation,
	city_id int32,
) error {
	s.Stop(city_id)
	return c.Cabinet.Administrator.Delete(ctx, city_id)
}
