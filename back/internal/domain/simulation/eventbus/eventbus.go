//go:generate go run -tags=generator ./scripts/generate_events.go
package eventbus

import (
	"context"
	"log/slog"
	"slices"
	"sync"
)

type EventType uint8

type Event interface {
	Type() EventType
}

// EventListener reacts to TileEvent notifications instead of running on a fixed interval.
// Unlike Ticker in gameloops, it does not operate on chunks - it's triggered by state changes.
type EventListener interface {
	ID() string
	FromEvents() []EventType
	OnChange(ctx context.Context, city_id int32, event Event) error
}

type EventBus struct {
	mutex     sync.Mutex
	listeners map[EventType][]EventListener
}

func New() *EventBus {
	return &EventBus{
		listeners: make(map[EventType][]EventListener, events_count),
	}
}

func (e *EventBus) Subscribe(listener EventListener) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, event_type := range listener.FromEvents() {
		already_subscribed := slices.ContainsFunc(
			e.listeners[event_type],
			func(subscribed_listener EventListener) bool {
				return subscribed_listener.ID() == listener.ID()
			},
		)

		if already_subscribed {
			slog.Warn(
				"skipped subscription in TilesEventBus, already subscribed",
				"listener", listener.ID(),
			)
			continue
		}

		e.listeners[event_type] = append(e.listeners[event_type], listener)
	}

}

func (e *EventBus) Publish(city_id int32, event Event) {
	e.mutex.Lock()
	listeners := slices.Clone(e.listeners[event.Type()])
	e.mutex.Unlock()

	ctx := context.Background()

	for index := range listeners {
		err := listeners[index].OnChange(ctx, city_id, event)
		if err != nil {
			slog.Error(
				"reactive_ticker: OnEvent failed",
				"city_id", city_id,
				"listener", listeners[index].ID(),
				"error", err,
			)
		}
	}
}
