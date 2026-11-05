package simulation

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/eventbus"
	simutils "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/utils"
)

var _ eventbus.EventListener = (*DesirabilityListener)(nil)

const listener_desirability = "listener_desirability"
const desirability_spread = 6

type DesirabilityListener struct {
	*cityhall.Deputies
	halo []cityhall.HaloTile
}

func NewDesirabilityListener(city_hall *cityhall.Deputies) *DesirabilityListener {
	return &DesirabilityListener{
		Deputies: city_hall,
		halo:     city_hall.Topology.Maker.OffsetsHalo(desirability_spread),
	}
}

func (d *DesirabilityListener) FromEvents() []eventbus.EventType {
	return []eventbus.EventType{
		eventbus.TilesChange,
	}
}

func (d *DesirabilityListener) ID() string {
	return listener_desirability
}

// OnChange implements eventbus.EventListener.OnChange
//
// see crime_loop.go, it's the same logic.
func (d *DesirabilityListener) OnChange(
	ctx context.Context,
	city_id int32,
	event eventbus.Event,
) error {
	slog.Info("msg", "onChange on city ", city_id)
	start := time.Now()
	tiles_change_event, ok := event.(eventbus.TilesChangeData)
	if !ok {
		slog.Error(
			"wrong event type given to DesirabilityListener",
			"expected", "TilesChangeData",
		)
		return nil
	}
	old_tiles := tiles_change_event.OldTiles()
	new_tiles := tiles_change_event.NewTiles()
	enacted_ordinances := tiles_change_event.EnactedOrdinances()
	if len(old_tiles) != len(new_tiles) {
		return eventbus.ErrInvalidTileEvent
	}

	base_deltas := make(map[core.TilePosition]core.Desirability, len(old_tiles))
	coalesced_deltas := make(map[core.TilePosition]core.Desirability, len(old_tiles)*6)

	for index := range old_tiles {
		old_tile := &old_tiles[index]
		new_tile := &new_tiles[index]
		old_artifact := registry.Artifacts[old_tile.Artifact.Name]
		new_artifact := registry.Artifacts[new_tile.Artifact.Name]

		old_contribution := d.getLandValueContribution(
			old_artifact,
			old_tile.Artifact.Status,
			enacted_ordinances,
		)
		new_contribution := d.getLandValueContribution(
			new_artifact,
			new_tile.Artifact.Status,
			enacted_ordinances,
		)

		delta_values := core.Desirability{
			LandValue: new_contribution.LandValue - old_contribution.LandValue,
		}

		base_deltas[old_tile.Position] = delta_values
		coalesced_deltas[old_tile.Position] = delta_values
	}

	for index := range old_tiles {
		tile := &old_tiles[index]

		base_delta_values := base_deltas[tile.Position]

		if base_delta_values.LandValue == 0 {
			continue
		}

		for _, halo_tile := range d.halo {
			neighbor := core.TilePosition{
				Q: uint8(int16(tile.Position.Q) + int16(halo_tile.Q)),
				R: uint8(int16(tile.Position.R) + int16(halo_tile.R)),
			}

			ok := d.Topology.Checker.ValidPosition(neighbor)
			if !ok {
				continue
			}

			halo_deltas := coalesced_deltas[neighbor]

			halo_deltas.LandValue += simutils.ApplyLinearAttenuation(
				base_delta_values.LandValue,
				halo_tile.Distance,
				desirability_spread,
			)

			coalesced_deltas[neighbor] = halo_deltas
		}
	}

	deltas_to_apply := make(
		[]cityhall.TileStatsDeltas[core.Desirability],
		0,
		len(coalesced_deltas),
	)

	for coordinates, deltas := range coalesced_deltas {
		if deltas.LandValue == 0 {
			continue
		}

		deltas_to_apply = append(deltas_to_apply, cityhall.TileStatsDeltas[core.Desirability]{
			TilePosition: coordinates,
			Values:       deltas,
		})
	}

	err := d.Cadastre.Updater.DesirabilityFromDeltas(ctx, city_id, deltas_to_apply)
	if err != nil {
		return err
	}

	duration := time.Since(start)
	slog.Info(
		"listener_desirability",
		"took: ", duration,
		// "deltas applied", deltas_to_apply,
	)
	return nil
}

func (d *DesirabilityListener) getLandValueContribution(
	artifact *core.Artifact,
	status core.ArtifactStatus,
	enacted_ordinances []core.OrdinanceType,
) core.Desirability {
	if status != core.Step_3_DeliveredOperational {
		return core.Desirability{}
	}

	contribution := core.Desirability{
		LandValue: artifact.SpatialEffect(core.LandValue),
	}

	for _, ordinance_effect := range artifact.OrdinanceEffects {
		if !slices.Contains(enacted_ordinances, ordinance_effect.Type) {
			continue
		}

		switch ordinance_effect.Effect {
		case core.LandValue:
			contribution.LandValue += ordinance_effect.AdditionalValue
		}
	}

	return contribution
}
