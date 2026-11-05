package app

import (
	"context"
	"errors"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

var (
	ErrCityNotFound = &domain.Error{
		Reason: errors.New("city does not exist"),
		Key:    domain.ErrNotFound,
	}
	ErrCityNotSimulated = &domain.Error{
		Reason: errors.New("city not registered in simulated cities"),
		Key:    domain.ErrNotFound,
	}
)

type App struct {
	*cityhall.Deputies
	Sim      *gameloops.Simulation
	EventBus *eventbus.EventBus
}

func Boot(store cityhall.Store) *App {
	sim := gameloops.New()
	event_bus := eventbus.New()

	// city hall actors
	city_hall := cityhall.New(store)

	// game loops
	// artifacts_placing_loop := simulation.NewArtifactPlacingLoop(city_hall)
	// artifacts_completion_loop := simulation.NewArtifactsCompletionLoop(city_hall, event_bus)
	// artifacts_abandonment_loop := simulation.NewArtifactAbandonmentLoop(city_hall, event_bus)
	// employment_loop := simulation.NewEmploymentLoop(city_hall)
	// traffic_loop := simulation.NewTrafficLoop(city_hall)

	sim.SetGameLoops(
	// gameloops.GameLoop{
	// 	Ticker: artifacts_placing_loop,
	// 	Chunks: artifacts_placing_loop.Chunks,
	// },
	// gameloops.GameLoop{
	// 	Ticker: artifacts_completion_loop,
	// 	Chunks: artifacts_completion_loop.Chunks,
	// },
	// gameloops.GameLoop{
	// 	Ticker: artifacts_abandonment_loop,
	// 	Chunks: artifacts_abandonment_loop.Chunks,
	// },
	// gameloops.GameLoop{
	// 	Ticker: employment_loop,
	// 	Chunks: employment_loop.Chunks,
	// },
	// gameloops.GameLoop{
	// 	Ticker: traffic_loop,
	// 	Chunks: traffic_loop.Chunks,
	// },
	)

	// listeners tile-change event
	criminality_listener := simulation.NewCriminalityListener(city_hall)
	desirability_listener := simulation.NewDesirabilityListener(city_hall)
	pollution_listener := simulation.NewPollutionListener(city_hall)
	employment_deletion_listener := simulation.NewEmploymentDeletionListener(city_hall)

	event_bus.Subscribe(criminality_listener)
	event_bus.Subscribe(desirability_listener)
	event_bus.Subscribe(pollution_listener)
	event_bus.Subscribe(employment_deletion_listener)

	return &App{
		Deputies: city_hall,
		Sim:      sim,
		EventBus: event_bus,
	}
}

func (a *App) StartAllCitiesSimulation(ctx context.Context) error {
	cities_ids, err := a.Cabinet.Administrator.GetCities(ctx)
	if err != nil {
		return err
	}

	for _, city_id := range cities_ids {
		a.Sim.Start(city_id)
	}

	return nil
}

func (a *App) Connect(ctx context.Context, city_id int32) (<-chan any, error) {
	exists, err := a.Cabinet.Administrator.CheckIfExists(ctx, city_id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("%w %d", ErrCityNotFound, city_id)
	}

	messages, ok := a.Sim.Register(city_id)
	if !ok {
		return nil, fmt.Errorf("%w %d", ErrCityNotSimulated, city_id)
	}

	return messages, nil
}

func (a *App) Disconnect(city_id int32) error {
	ok := a.Sim.Deregister(city_id)
	if !ok {
		return fmt.Errorf("%w %d", ErrCityNotSimulated, city_id)
	}

	return nil
}
