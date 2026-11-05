package api

import (
	"context"
	"errors"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/data"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEMailConflict = &domain.Error{
		Reason: errors.New("email already exists"),
		Key:    domain.ErrConflict,
	}
)

// CreateFullCity is a user action.
// TODO maybe use an event to start simulating city instead of doing it in this function
func CreateFullCity(
	ctx context.Context,
	c *cityhall.Deputies,
	sim *gameloops.Simulation,
	email, password, name, mayor_name string,
) (int32, error) {
	var new_city int32

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	params := &cityhall.CityCreationParams{
		Email:     email,
		Password:  hash,
		Name:      name,
		MayorName: mayor_name,
		Tiles:     &data.InitialTiles,
		Budget:    c.Treasury.Maker.InitialBudget(),
	}

	err = c.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		new_city, err = c.Cabinet.Administrator.Create(task_ctx, params)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	sim.Start(new_city)

	return new_city, nil
}
