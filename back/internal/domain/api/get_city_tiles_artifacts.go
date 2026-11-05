package api

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
)

// GetCityTilesArtifacts is a user action.
// Necessary for front-end to display latest tiles when users arrive on their city.
func GetCityTilesArtifacts(
	ctx context.Context,
	c *cityhall.Deputies,
	city_id int32,
) (core.CityTilesArtifacts, error) {
	var artifact_ids [core.TilesCount]core.ArtifactType
	var anchor_qs [core.TilesCount]uint8
	var anchor_rs [core.TilesCount]uint8
	var artifact_statuses [core.TilesCount]core.ArtifactStatus
	var super_kind_ids [core.TilesCount]core.SuperKindType

	response, err := c.Cadastre.Getter.FromCity(ctx, city_id)
	if err != nil {
		return core.CityTilesArtifacts{}, err
	}

	for index := range response {
		tile := &response[index]
		artifact := registry.Artifacts[tile.Artifact.Name]

		if tile.Position == tile.Anchor || artifact.HasFreeSize() {
			artifact_ids[index] = tile.Artifact.Name
		} else {
			artifact_ids[index] = 0
		}

		anchor_qs[index] = tile.Anchor.Q
		anchor_rs[index] = tile.Anchor.R
		artifact_statuses[index] = tile.Artifact.Status
		super_kind_ids[index] = tile.SuperKind
	}

	return core.CityTilesArtifacts{
		ArtifactIDs:      &artifact_ids,
		AnchorQs:         &anchor_qs,
		AnchorRs:         &anchor_rs,
		ArtifactStatuses: &artifact_statuses,
		SuperKindIDs:     &super_kind_ids,
	}, nil
}
