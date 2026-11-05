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

var _ eventbus.EventListener = (*CriminalityListener)(nil)

const listener_criminality = "listener_criminality"
const criminality_spread = 4

type CriminalityListener struct {
	*cityhall.Deputies
	halo []cityhall.HaloTile
}

func NewCriminalityListener(city_hall *cityhall.Deputies) *CriminalityListener {
	return &CriminalityListener{
		Deputies: city_hall,
		halo:     city_hall.Topology.Maker.OffsetsHalo(criminality_spread),
	}
}

func (c *CriminalityListener) FromEvents() []eventbus.EventType {
	return []eventbus.EventType{
		eventbus.TilesChange,
	}
}

func (c *CriminalityListener) ID() string {
	return listener_criminality
}

// OnChange implements eventbus.EventListener.OnChange
//
// When an artifact becomes operational (Step_3_DeliveredOperational), its crime effects
// (CrimeGeneric, CrimeDrug, CrimeVss) are diffused to surrounding tiles.
// When an artifact stops being operational (Step_1_AwaitingAssignment | Step_2_AssignedAwaitingDelivery),
// its crime effects are removed.
//
// The delta is computed by comparing each tile's OLD contribution (based on its old
// artifact & old status) against its NEW contribution (based on its new artifact & new status).
// This makes the computation correct regardless of whether the artifact identity changed
// (build/destroy) or stayed the same (completion/abandonment/rebuilding the same artifact).
//
// Crime effect application is two-fold:
//
//   - iterate on given artifacts's tiles to apply crime effects.
//   - iterate again to apply crime effects on tiles halo given `criminality_spread`.
func (c *CriminalityListener) OnChange(
	ctx context.Context,
	city_id int32,
	event eventbus.Event,
) error {
	slog.Info("msg", "onChange on city ", city_id)
	start := time.Now()
	tiles_change_event, ok := event.(eventbus.TilesChangeData)
	if !ok {
		slog.Error(
			"wrong event type given to CriminalityListener",
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

	base_deltas := make(map[core.TilePosition]core.Criminality, len(old_tiles))
	coalesced_deltas := make(map[core.TilePosition]core.Criminality, len(old_tiles)*6)

	for index := range old_tiles {
		old_tile := &old_tiles[index]
		new_tile := &new_tiles[index]
		old_artifact := registry.Artifacts[old_tile.Artifact.Name]
		new_artifact := registry.Artifacts[new_tile.Artifact.Name]

		old_contribution := c.getCrimeContribution(
			old_artifact,
			old_tile.Artifact.Status,
			enacted_ordinances,
		)
		new_contribution := c.getCrimeContribution(
			new_artifact,
			new_tile.Artifact.Status,
			enacted_ordinances,
		)

		delta_values := core.Criminality{
			CrimeGeneric: new_contribution.CrimeGeneric - old_contribution.CrimeGeneric,
			CrimeDrug:    new_contribution.CrimeDrug - old_contribution.CrimeDrug,
			CrimeVss:     new_contribution.CrimeVss - old_contribution.CrimeVss,
		}

		base_deltas[old_tile.Position] = delta_values
		coalesced_deltas[old_tile.Position] = delta_values
	}

	for index := range old_tiles {
		tile := &old_tiles[index]

		base_delta_values := base_deltas[tile.Position]

		if base_delta_values.CrimeGeneric == 0 &&
			base_delta_values.CrimeDrug == 0 &&
			base_delta_values.CrimeVss == 0 {
			continue
		}

		for _, halo_tile := range c.halo {
			neighbor := core.TilePosition{
				Q: uint8(int16(tile.Position.Q) + int16(halo_tile.Q)),
				R: uint8(int16(tile.Position.R) + int16(halo_tile.R)),
			}

			ok := c.Topology.Checker.ValidPosition(neighbor)
			if !ok {
				continue
			}

			halo_deltas := coalesced_deltas[neighbor]

			halo_deltas.CrimeGeneric += simutils.ApplyLinearAttenuation(
				base_delta_values.CrimeGeneric,
				halo_tile.Distance,
				criminality_spread,
			)
			halo_deltas.CrimeDrug += simutils.ApplyLinearAttenuation(
				base_delta_values.CrimeDrug,
				halo_tile.Distance,
				criminality_spread,
			)
			halo_deltas.CrimeVss += simutils.ApplyLinearAttenuation(
				base_delta_values.CrimeVss,
				halo_tile.Distance,
				criminality_spread,
			)

			coalesced_deltas[neighbor] = halo_deltas
		}
	}

	deltas_to_apply := make(
		[]cityhall.TileStatsDeltas[core.Criminality],
		0,
		len(coalesced_deltas),
	)

	for coordinates, deltas := range coalesced_deltas {
		if deltas.CrimeGeneric == 0 && deltas.CrimeDrug == 0 && deltas.CrimeVss == 0 {
			continue
		}

		deltas_to_apply = append(deltas_to_apply, cityhall.TileStatsDeltas[core.Criminality]{
			TilePosition: coordinates,
			Values:       deltas,
		})
	}

	err := c.Cadastre.Updater.CriminalityFromDeltas(ctx, city_id, deltas_to_apply)
	if err != nil {
		return err
	}

	duration := time.Since(start)
	slog.Info(
		"crime_listener",
		"took: ", duration,
		// "deltas applied", deltas_to_apply,
	)
	return nil
}

func (c *CriminalityListener) getCrimeContribution(
	artifact *core.Artifact,
	status core.ArtifactStatus,
	enacted_ordinances []core.OrdinanceType,
) core.Criminality {
	if status != core.Step_3_DeliveredOperational {
		return core.Criminality{}
	}

	contribution := core.Criminality{
		CrimeGeneric: artifact.SpatialEffect(core.CrimeGeneric),
		CrimeDrug:    artifact.SpatialEffect(core.CrimeDrug),
		CrimeVss:     artifact.SpatialEffect(core.CrimeVss),
	}

	for _, ordinance_effect := range artifact.OrdinanceEffects {
		if !slices.Contains(enacted_ordinances, ordinance_effect.Type) {
			continue
		}

		switch ordinance_effect.Effect {
		case core.CrimeGeneric:
			contribution.CrimeGeneric += ordinance_effect.AdditionalValue
		case core.CrimeDrug:
			contribution.CrimeDrug += ordinance_effect.AdditionalValue
		case core.CrimeVss:
			contribution.CrimeVss += ordinance_effect.AdditionalValue
		}
	}

	return contribution
}
