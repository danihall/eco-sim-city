package api

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// EnactOrdinance is a user action.
func EnactOrdinance(
	ctx context.Context,
	d *cityhall.Deputies,
	city_id int32,
	ordinance_name core.OrdinanceType,
) error {
	return d.Policies.Manager.Enact(ctx, city_id, ordinance_name)
}
