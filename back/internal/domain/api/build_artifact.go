package api

import (
	"context"
	"errors"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
)

var (
	ErrArtifactInvalid = &domain.Error{
		Reason: errors.New("this artifact does not exist"),
		Key:    domain.ErrInvalid,
	}
	ErrSuperKindArtifactInvalid = &domain.Error{
		Reason: errors.New(
			"artifact to build has invalid super_kind. Must be 'natural' or 'public_systems'",
		),
		Key: domain.ErrInvalid,
	}
	ErrInvalidAnchor = &domain.Error{
		Reason: errors.New(
			"artifact has no valid anchor",
		),
		Key: domain.ErrInvalid,
	}
)

// BuildArtifact is a user action.
// It involves 3 successive steps to ensure artifact can be built:
//
//   - Ask *topology* to generate in-bounds and on-ground `coordinates`.
//   - Get city full stats from city_hall
//   - Ask *cadastre* to fetch tiles at given `coordinates`.
//   - Ask *permits* to validate construction.
func BuildArtifact(
	ctx context.Context,
	d *cityhall.Deputies,
	event_bus *eventbus.EventBus,
	city_id int32,
	params struct {
		Anchor       core.TilePosition
		ArtifactName core.ArtifactType
		UserWidth    uint8
		UserHeight   uint8
	},
) ([]core.Tile, error) {
	artifact, ok := registry.Artifacts[params.ArtifactName]
	if !ok {
		return nil, fmt.Errorf("%v: %w", params.ArtifactName, ErrArtifactInvalid)
	}

	artifact_super_kind := registry.ArtifactsRelations[params.ArtifactName].SuperKind
	if artifact_super_kind != core.Natural && artifact_super_kind != core.PublicSystems {
		return nil, fmt.Errorf(
			"%w. Received '%v'",
			ErrSuperKindArtifactInvalid,
			artifact_super_kind,
		)
	}

	coordinates, err := d.Topology.Maker.FromArtifact(
		artifact,
		params.Anchor,
		params.UserWidth, params.UserHeight,
		nil,
	)
	if err != nil {
		return nil, err
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

		tiles_to_update, err = d.Permits.Validator.UserCandidate(city, artifact, tiles)
		if err != nil {
			return err
		}

		// This is just to get []core.TilePosition, oh well.
		coordinates_to_update := d.Topology.Maker.FilteredTiles(tiles_to_update)

		if artifact.IsUnique {
			err = d.Flagship.Manager.Add(task_ctx, city_id, params.ArtifactName)
			if err != nil {
				return err
			}
		}

		anchors := make([]core.TilePosition, len(coordinates_to_update))
		if artifact.HasFreeSize() {
			copy(anchors, coordinates_to_update)
		} else {
			for index := range coordinates_to_update {
				anchors[index] = params.Anchor
			}
		}

		updated_tiles, err = d.Construction.Builder.SingleArtifact(
			task_ctx,
			city_id,
			cityhall.BuildParams{
				Anchors: anchors,
				Artifact: core.ArtifactInstance{
					Name:   artifact.Name,
					Label:  artifact.Labels[0],
					Status: core.Step_3_DeliveredOperational,
				},
				Coordinates:   coordinates_to_update,
				SuperKindType: artifact_super_kind,
				KindType:      registry.ArtifactsRelations[params.ArtifactName].Kind,
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
