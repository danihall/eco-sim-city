package simulation

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/chunks"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

var _ gameloops.Ticker = (*ArtifactsCompletionLoop)(nil)

const artifacts_completion_loop = "artifacts_completion_loop"

// ArtifactsCompletionLoop simulates artifacts previously placed
// being brought to completion. (building construction is finished).
type ArtifactsCompletionLoop struct {
	*cityhall.Deputies
	eventBus *eventbus.EventBus
	Chunks   []core.Chunk
}

func NewArtifactsCompletionLoop(
	city_hall *cityhall.Deputies,
	event_bus *eventbus.EventBus,
) *ArtifactsCompletionLoop {
	return &ArtifactsCompletionLoop{
		Deputies: city_hall,
		eventBus: event_bus,
		Chunks:   chunks.FourBarsHorizontalAndFourBarsVertical[:],
	}
}

func (a *ArtifactsCompletionLoop) ID() string {
	return artifacts_completion_loop
}

func (a *ArtifactsCompletionLoop) Tick(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) (*gameloops.Update, error) {
	// TODO write type struct for these params
	params := struct {
		CurrentArtifactStatus core.ArtifactStatus
		NewArtifactStatus     core.ArtifactStatus
		Chunk                 core.Chunk
	}{core.Step_2_AssignedAwaitingDelivery, core.Step_3_DeliveredOperational, chunk}

	completed_artifact_tiles, err := a.Construction.Completer.AllFromChunk(ctx, city_id, params)
	if err != nil {
		return nil, err
	}

	if len(completed_artifact_tiles) == 0 {
		return nil, nil
	}

	enacted_ordinances, err := a.Policies.Getter.EnactedOrdinances(ctx, city_id)
	if err != nil {
		return nil, err
	}

	old_tiles := make([]core.Tile, len(completed_artifact_tiles))
	copy(old_tiles, completed_artifact_tiles)
	for index := range old_tiles {
		old_tiles[index].Artifact.Status = core.Step_2_AssignedAwaitingDelivery
	}

	a.eventBus.Publish(
		city_id,
		eventbus.NewTilesChangeData(old_tiles, completed_artifact_tiles, enacted_ordinances),
	)

	return &gameloops.Update{
		GameLoopName: artifacts_completion_loop,
		Body:         completed_artifact_tiles,
	}, nil
}
