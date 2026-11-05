package api

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
)

var (
	ErrArtifactPathwayInvalid = &domain.Error{
		Reason: errors.New("this artifact is not of a valid transport network's type"),
		Key:    domain.ErrInvalid,
	}
	ErrPathwayInvalidCoordinates = &domain.Error{
		Reason: errors.New("all pathway's coordinates are invalid"),
		Key:    domain.ErrInvalid,
	}
)

// BuildPathway is a user action.
func BuildPathway(
	ctx context.Context,
	d *cityhall.Deputies,
	event_bus *eventbus.EventBus,
	city_id int32,
	artifact_name core.ArtifactType,
	coordinates []core.TilePosition,
) ([]core.Tile, error) {
	artifact, ok := registry.Artifacts[artifact_name]
	if !ok {
		return nil, fmt.Errorf("%v: %w", artifact_name, ErrArtifactInvalid)
	}

	if !slices.Contains(registry.TransportNetworkArtifacts, artifact_name) {
		return nil, fmt.Errorf("%v: %w", artifact_name, ErrArtifactPathwayInvalid)
	}

	city, err := d.GetCityFullState(ctx, city_id)
	if err != nil {
		return nil, err
	}

	var tiles_to_update []core.Tile
	var updated_tiles []core.Tile

	err = d.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		tiles, err := d.Cadastre.Getter.FromCoordinatesFullPropsWithStats(
			task_ctx,
			city_id,
			coordinates,
		)
		if err != nil {
			return err
		}

		tiles_to_update = d.Permits.Validator.MakeUnoccupiedTiles(tiles)
		if len(tiles_to_update) == 0 {
			return ErrPathwayInvalidCoordinates
		}

		// This is just to get []core.TilePosition, oh well.
		coordinates_to_update := d.Topology.Maker.FilteredTiles(tiles_to_update)

		if len(coordinates_to_update) == 0 {
			return ErrPathwayInvalidCoordinates
		}

		updated_tiles, err = d.Construction.Builder.Pathway(
			task_ctx,
			city_id,
			cityhall.PathwayParams{
				Anchors:     coordinates_to_update,
				Coordinates: coordinates_to_update,
				Artifact: core.ArtifactInstance{
					Name:   artifact.Name,
					Label:  artifact.Labels[0],
					Status: core.Step_3_DeliveredOperational,
				},
				SuperKindType: registry.ArtifactsRelations[artifact_name].SuperKind,
				KindType:      registry.ArtifactsRelations[artifact_name].Kind,
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	event_bus.Publish(
		city_id,
		eventbus.NewTilesChangeData(tiles_to_update, updated_tiles, city.EnactedOrdinances),
	)

	return updated_tiles, nil
}
