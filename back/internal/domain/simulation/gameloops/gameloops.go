// gameloop manages the lifecycle of game-loops for each city.
// Each city gets its own set of goroutines, one per game-loop, running independently.
package gameloops

import (
	"context"
	"log/slog"
	"slices"
	"sync"
	"time"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// A year in the game takes one week in real life
// TODO does this need its own game-loop?
// Maye we should derive it from the number of game-loops
const sim_month_interval = (7 * 24 * time.Hour) / 12

type Update struct {
	GameLoopName string `json:"game_loop_name"`
	Body         any    `json:"body"`
}

type Ticker interface {
	Tick(ctx context.Context, city_id int32, chunk core.Chunk) (*Update, error)
	ID() string
}

type GameLoop struct {
	Ticker
	Chunks []core.Chunk
}

type simulatedCity struct {
	cancel     context.CancelFunc
	messages   chan any
	subscribed bool
}

type Simulation struct {
	mutex           sync.Mutex
	simulatedCities map[int32]*simulatedCity
	gameLoops       []GameLoop
}

func New() *Simulation {
	return &Simulation{
		simulatedCities: make(map[int32]*simulatedCity),
	}
}

func (s *Simulation) SetGameLoops(game_loops ...GameLoop) {
	for index := range game_loops {
		already_set := slices.ContainsFunc(
			s.gameLoops,
			func(set_game_loop GameLoop) bool {
				return set_game_loop.ID() == game_loops[index].ID()
			},
		)

		if already_set {
			slog.Warn(
				"skipped game loop append in gameLoops, already set",
				"game_loop", game_loops[index].ID(),
			)
			continue
		}

		s.gameLoops = append(s.gameLoops, game_loops[index])
	}
}

func (s *Simulation) Register(city_id int32) (<-chan any, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	city, exists := s.simulatedCities[city_id]
	if !exists {
		return nil, false
	}

	city.subscribed = true

	return city.messages, true
}

func (s *Simulation) Deregister(city_id int32) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	city, exists := s.simulatedCities[city_id]
	if !exists {
		return false
	}

	city.subscribed = false
	return true
}

// Start adds a delay between each game-loops. This is by design,
// to give the illusion of an evolving city.
func (s *Simulation) Start(city_id int32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.simulatedCities[city_id]

	if exists {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	s.simulatedCities[city_id] = &simulatedCity{
		cancel:   cancel,
		messages: make(chan any, 1),
	}

	for index, game_loop := range s.gameLoops {
		// WIP this delay is for testing purpose
		delay := time.Duration(index) * time.Second
		// delay := time.Duration(index) * time.Hour
		go s.runGameLoop(ctx, city_id, game_loop, delay)
	}
}

// runGameLoop waits for delay's end before starting its work.
func (s *Simulation) runGameLoop(
	ctx context.Context,
	city_id int32,
	game_loop GameLoop,
	delay time.Duration,
) {
	if delay > 0 {
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
			slog.Info("delay expired. game_loop can start", "city", city_id)
		}
	}

	// WIP this interval is for testing purpose
	game_loop_interval := time.Second * 8
	// game_loop_interval := (time.Hour * 24 * 7) / time.Duration(len(game_loop.Chunks))
	ticker := time.NewTicker(game_loop_interval)
	defer ticker.Stop()

	chunk_index := 0

	_run := func() {
		chunk := game_loop.Chunks[chunk_index]
		start := time.Now()
		data, err := game_loop.Tick(ctx, city_id, chunk)
		if err != nil {
			slog.Error("game_loop tick failed",
				"city_id", city_id,
				"game_loop", game_loop.ID(),
				"chunk", chunk_index,
				"error", err,
			)
		}
		duration := time.Since(start)
		slog.Info("sim tick", "game_loop", game_loop.ID(), "chunk", chunk, "duration", duration)

		if data != nil {
			s.notify(city_id, data)
		}

		chunk_index = (chunk_index + 1) % len(game_loop.Chunks)
	}

	_run()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_run()
		}
	}
}

func (s *Simulation) notify(city_id int32, data *Update) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	city, exists := s.simulatedCities[city_id]
	if !exists {
		return
	}

	if city.subscribed {
		select {
		case city.messages <- data:
		default:
		}
	}
}

func (s *Simulation) Stop(city_id int32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	city, exists := s.simulatedCities[city_id]

	if !exists {
		slog.Warn("this city does not exists", "city_id", city_id)
		return
	}

	city.cancel()
	delete(s.simulatedCities, city_id)

	slog.Info("city stopped and deleted", "city_id", city_id)
}
