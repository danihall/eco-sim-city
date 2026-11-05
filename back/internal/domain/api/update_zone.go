package api

import (
	"context"
	"errors"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var (
	ErrSuperKindZoningInvalid = &domain.Error{
		Reason: errors.New("zone's super_kind invalid, must be 'industrial' or 'residential'"),
		Key:    domain.ErrInvalid,
	}
)

// UpdateZone is a user action.
// UpdateZone should only update a tile's super_kind given the following conditions:
//
// - tile is of super_kind core.Natural.
// - tile is not of super_kind core.PublicSystems.
// - tile has no Artifact (Step_1_AwaitingAssignment)
func UpdateZone(
	ctx context.Context,
	c *cityhall.Deputies,
	city_id int32,
	super_kind_type core.SuperKindType,
	top_left core.TilePosition,
	bottom_right core.TilePosition,
) ([]core.Tile, error) {
	if super_kind_type != core.Industrial && super_kind_type != core.Residential {
		return nil, fmt.Errorf(
			"%w. received '%v'",
			ErrSuperKindZoningInvalid,
			super_kind_type.String(),
		)
	}

	valid_coordinates, err := c.Topology.Maker.FromChunk(top_left, bottom_right)
	if err != nil {
		return nil, fmt.Errorf("cannot update tiles's zoning: %w", err)
	}

	update_params := cityhall.ZoneParams{
		NewSuperKind:           super_kind_type,
		NaturalSuperKind:       core.Natural,
		PublicSystemsSuperKind: core.PublicSystems,
		NewArtifactStatus:      core.Step_1_AwaitingAssignment,
		OldArtifactStatus:      core.Step_1_AwaitingAssignment,
		Coordinates:            valid_coordinates,
	}

	return c.Zoning.Manager.Update(ctx, city_id, update_params)
}
