package api

import (
	"context"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// RevertZone is a user action.
func RevertZone(
	ctx context.Context,
	c *cityhall.Deputies,
	city_id int32,
	top_left core.TilePosition,
	bottom_right core.TilePosition,
) ([]core.Tile, error) {
	valid_coordinates, err := c.Topology.Maker.FromChunk(top_left, bottom_right)
	if err != nil {
		return nil, fmt.Errorf("cannot update tiles's zoning: %w", err)
	}

	revert_params := cityhall.ZoneParams{
		NewSuperKind:           core.Natural,
		NaturalSuperKind:       core.Natural,
		PublicSystemsSuperKind: core.PublicSystems,
		NewArtifactStatus:      core.Step_3_DeliveredOperational,
		OldArtifactStatus:      core.Step_1_AwaitingAssignment,
		Coordinates:            valid_coordinates,
	}

	return c.Zoning.Manager.Revert(ctx, city_id, revert_params)
}
