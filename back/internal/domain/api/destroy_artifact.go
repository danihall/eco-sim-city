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
	ErrCommonAnchorsButDifferentSuperKind = &domain.Error{
		Reason: errors.New("tiles with common anchors should have common super_kind"),
		Key:    domain.ErrConflict,
	}
)

// DestroyArtifact is a user action.
// A destroyed artifact will always be replaced by Wasteland artifact.
// We only need to check one valid coordinate because adapter will take care of destroying
// all related tiles via anchor_q & anchor_r
func DestroyArtifact(
	ctx context.Context,
	d *cityhall.Deputies,
	event_bus *eventbus.EventBus,
	city_id int32,
	anchor core.TilePosition,
) ([]core.Tile, error) {
	ok := d.Topology.Checker.ValidPosition(anchor)
	if !ok {
		return nil, fmt.Errorf(
			"tile is not destroyable: %w: %v",
			cityhall.WarnNoActionableCoordinates, anchor,
		)
	}

	var tiles []core.Tile
	var destroyed_tiles []core.Tile

	err := d.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		var err error
		// NOTE this query could be avoided if front-end sends all anchored tiles with their ZoneName
		// Or if we don't use checkAllTilesSameSuperKind
		tiles, err = d.Cadastre.Getter.FromAnchor(task_ctx, city_id, anchor)
		if err != nil {
			return err
		}

		err = checkAllTilesSameSuperKind(tiles)
		if err != nil {
			return err
		}

		super_kind_name := tiles[0].SuperKind
		wasteland := registry.Artifacts[core.Wasteland]
		artifact_to_destroy := registry.Artifacts[tiles[0].Artifact.Name]

		destroy_params := cityhall.DestroyParams{
			ReplacementArtifact: core.ArtifactInstance{
				Name:   wasteland.Name,
				Label:  wasteland.Labels[0],
				Status: core.Step_1_AwaitingAssignment,
			},
			Anchor:                  anchor,
			HasAllCitizensAllocated: false,
		}

		if artifact_to_destroy.IsUnique {
			err = d.Flagship.Manager.Remove(task_ctx, city_id, artifact_to_destroy.Name)
			if err != nil {
				return err
			}
		}

		if super_kind_name == core.PublicSystems {
			destroy_public_sys_params := cityhall.DestroyPublicSystemsParams{
				DestroyParams:            destroy_params,
				ReplacementSuperKindType: core.Natural,
				ReplacementKindType:      core.Unconstructible,
			}

			destroyed_tiles, err = d.Construction.Destroyer.PublicSystemsArtifact(
				task_ctx,
				city_id,
				destroy_public_sys_params,
			)
		} else {
			destroyed_tiles, err = d.Construction.Destroyer.Artifact(
				task_ctx,
				city_id,
				destroy_params,
			)
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	enacted_ordinances, err := d.Policies.Getter.EnactedOrdinances(ctx, city_id)
	if err != nil {
		return nil, err
	}

	event_bus.Publish(
		city_id,
		eventbus.NewTilesChangeData(tiles, destroyed_tiles, enacted_ordinances),
	)

	return destroyed_tiles, nil
}

// checkAllTilesSameSuperKind is implemented to ensure no incoherent data is in DB, but in practice, this should never happen.
func checkAllTilesSameSuperKind(tiles []core.Tile) error {
	for index := range tiles {
		if index == 0 {
			continue
		}

		previous_tiles := &tiles[index-1]
		current_tile := &tiles[index]

		if current_tile.SuperKind != previous_tiles.SuperKind {
			return fmt.Errorf(
				"%w, got %v (q=%d, r=%d) and %v (q=%d, r=%d). Common anchor is anchor_q=%d, anchor_q=%d",
				ErrCommonAnchorsButDifferentSuperKind,
				previous_tiles.SuperKind,
				previous_tiles.Position.Q,
				previous_tiles.Position.R,
				current_tile.SuperKind,
				current_tile.Position.Q,
				current_tile.Position.R,
				current_tile.Anchor.Q,
				current_tile.Anchor.R,
			)
		}
	}

	return nil
}
