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

var _ eventbus.EventListener = (*PollutionListener)(nil)

const listener_pollution = "listener_pollution"
const pollution_spread = 5

type PollutionListener struct {
	*cityhall.Deputies
	halo []cityhall.HaloTile
}

func NewPollutionListener(city_hall *cityhall.Deputies) *PollutionListener {
	return &PollutionListener{
		Deputies: city_hall,
		halo:     city_hall.Topology.Maker.OffsetsHalo(pollution_spread),
	}
}

func (p *PollutionListener) FromEvents() []eventbus.EventType {
	return []eventbus.EventType{
		eventbus.TilesChange,
	}
}

func (p *PollutionListener) ID() string {
	return listener_pollution
}

// OnChange implements eventbus.EventListener.OnChange
//
// see crime_loop.go, it's the same logic.
func (p *PollutionListener) OnChange(
	ctx context.Context,
	city_id int32,
	event eventbus.Event,
) error {
	slog.Info("msg", "onChange on city ", city_id)
	start := time.Now()
	tiles_change_event, ok := event.(eventbus.TilesChangeData)
	if !ok {
		slog.Error(
			"wrong event type given to PollutionListener",
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

	base_deltas := make(map[core.TilePosition]core.Pollution, len(old_tiles))
	coalesced_deltas := make(map[core.TilePosition]core.Pollution, len(old_tiles)*6)

	for index := range old_tiles {
		old_tile := &old_tiles[index]
		new_tile := &new_tiles[index]
		old_artifact := registry.Artifacts[old_tile.Artifact.Name]
		new_artifact := registry.Artifacts[new_tile.Artifact.Name]

		old_contribution := p.getPollutionContribution(
			old_artifact,
			old_tile.Artifact.Status,
			enacted_ordinances,
		)
		new_contribution := p.getPollutionContribution(
			new_artifact,
			new_tile.Artifact.Status,
			enacted_ordinances,
		)

		delta_values := core.Pollution{
			Ghg1: new_contribution.Ghg1 - old_contribution.Ghg1,
		}

		base_deltas[old_tile.Position] = delta_values
		coalesced_deltas[old_tile.Position] = delta_values
	}

	for index := range old_tiles {
		tile := &old_tiles[index]

		base_delta_values := base_deltas[tile.Position]

		if base_delta_values.Ghg1 == 0 {
			continue
		}

		for _, halo_tile := range p.halo {
			neighbor := core.TilePosition{
				Q: uint8(int16(tile.Position.Q) + int16(halo_tile.Q)),
				R: uint8(int16(tile.Position.R) + int16(halo_tile.R)),
			}

			ok := p.Topology.Checker.ValidPosition(neighbor)
			if !ok {
				continue
			}

			halo_deltas := coalesced_deltas[neighbor]

			halo_deltas.Ghg1 += simutils.ApplyQuadraticAttenuation(
				base_delta_values.Ghg1,
				halo_tile.Distance,
				pollution_spread,
			)

			coalesced_deltas[neighbor] = halo_deltas
		}
	}

	deltas_to_apply := make(
		[]cityhall.TileStatsDeltas[core.Pollution],
		0,
		len(coalesced_deltas),
	)

	for coordinates, deltas := range coalesced_deltas {
		if deltas.Ghg1 == 0 {
			continue
		}

		deltas_to_apply = append(deltas_to_apply, cityhall.TileStatsDeltas[core.Pollution]{
			TilePosition: coordinates,
			Values:       deltas,
		})
	}

	err := p.Cadastre.Updater.PollutionFromDeltas(ctx, city_id, deltas_to_apply)
	if err != nil {
		return err
	}

	duration := time.Since(start)
	slog.Info(
		"listener_pollution",
		"took: ", duration,
		// "deltas applied", deltas_to_apply,
	)
	return nil
}

func (p *PollutionListener) getPollutionContribution(
	artifact *core.Artifact,
	status core.ArtifactStatus,
	enacted_ordinances []core.OrdinanceType,
) core.Pollution {
	if status != core.Step_3_DeliveredOperational {
		return core.Pollution{}
	}

	contribution := core.Pollution{
		Ghg1: artifact.SpatialEffect(core.Ghg1),
	}

	for _, ordinance_effect := range artifact.OrdinanceEffects {
		if !slices.Contains(enacted_ordinances, ordinance_effect.Type) {
			continue
		}

		switch ordinance_effect.Effect {
		case core.Ghg1:
			contribution.Ghg1 += ordinance_effect.AdditionalValue
		}
	}

	return contribution
}
