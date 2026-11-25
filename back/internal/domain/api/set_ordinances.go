package api

import (
	"context"
	"errors"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
)

var (
	ErrInvalidOrdinanceName = &domain.Error{
		Reason: errors.New("this ordinance does not exist"),
		Key:    domain.ErrInvalid,
	}
)

// SetOrdinances is a user action.
func SetOrdinances(
	ctx context.Context,
	d *cityhall.Deputies,
	city_id int32,
	ordinances []core.OrdinanceType,
) error {
	for _, ordinance_name := range ordinances {
		_, ok := registry.Ordinances[ordinance_name]

		if !ok {
			return fmt.Errorf("%w: %v", ErrInvalidOrdinanceName, ordinance_name)
		}
	}

	return d.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		return d.Policies.Manager.Set(task_ctx, city_id, ordinances)
	})
}
