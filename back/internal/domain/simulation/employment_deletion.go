package simulation

import (
	"context"
	"log/slog"
	"time"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
)

// FIXME this should be a game_loop running after 04_employment_creations
// because deleting an employment must be repercuted on traffic and then on pollution.
// So having a listener for this would lead to a chain of events.
// Thats too complicated so let's embrace the fuzzy nature of the city sim
// and have the employment destruction logic in a game_loop before traffic_loop.
var _ eventbus.EventListener = (*EmploymentDeletionListener)(nil)

const listener_employment_deletion = "listener_employment_deletion"

type EmploymentDeletionListener struct {
	*cityhall.Deputies
}

func NewEmploymentDeletionListener(city_hall *cityhall.Deputies) *EmploymentDeletionListener {
	return &EmploymentDeletionListener{
		Deputies: city_hall,
	}
}

func (e *EmploymentDeletionListener) FromEvents() []eventbus.EventType {
	return []eventbus.EventType{
		eventbus.TilesChange,
	}
}

func (e *EmploymentDeletionListener) ID() string {
	return listener_employment_deletion
}

// OnChange implements eventbus.EventListener.OnChange
//
// When artifacts are abandonned/destroyed, existing job contracts between employers and resident
// related to these artifacts must be deleted.
func (e *EmploymentDeletionListener) OnChange(
	ctx context.Context,
	city_id int32,
	event eventbus.Event,
) error {
	slog.Info("msg", "onChange on city ", city_id)
	start := time.Now()
	tiles_change_event, ok := event.(eventbus.TilesChangeData)
	if !ok {
		slog.Error(
			"wrong event type given to EmploymentDeletionListener",
			"expected", "TilesChangeData",
		)
		return nil
	}
	new_tiles := tiles_change_event.NewTiles()

	params := cityhall.DeleteContractsParams{
		EmployerAnchors: make([]core.TilePosition, 0, len(new_tiles)),
		WorkerAnchors:   make([]core.TilePosition, 0, len(new_tiles)),
	}

	for index := range new_tiles {
		tile := &new_tiles[index]

		if tile.Artifact.Status == core.Step_1_AwaitingAssignment && tile.Anchor == tile.Position {
			switch tile.SuperKind {
			case core.Residential:
				params.WorkerAnchors = append(params.WorkerAnchors, tile.Anchor)
			case core.Industrial, core.PublicSystems:
				params.EmployerAnchors = append(params.EmployerAnchors, tile.Anchor)
			}
		}
	}

	if len(params.EmployerAnchors) == 0 && len(params.WorkerAnchors) == 0 {
		return nil
	}

	err := e.EmploymentOffice.Terminator.ContractsToDelete(ctx, city_id, params)
	if err != nil {
		return err
	}

	duration := time.Since(start)
	slog.Info("listener_employment_deletion", "took: ", duration, "contracts deleted", params)
	return nil
}
