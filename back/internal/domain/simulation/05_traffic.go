package simulation

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/chunks"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/simulation/gameloops"
)

var _ gameloops.Ticker = (*TrafficLoop)(nil)

var (
	ErrTransportMode = &domain.Error{
		Reason: errors.New("No transport mode found. This should not happen"),
		Key:    domain.ErrNotFound,
	}
)

const traffic_loop = "traffic_loop"
const max_travel_steps = 448
const max_walking_steps = 4
const move_back_penalty = 50
const closest_bias float32 = 0.005
const intersection_bias float32 = 0.5

type commute struct {
	stagesBuffer         []cityhall.Stage
	currentMode          core.EffectType
	workersCount         uint8
	hasTookPublicTransit bool
}

// newCommute sets currentMode to core.Foot for start.
func newCommute(
	stages_buffer []cityhall.Stage,
	workers_count uint8,
) *commute {
	return &commute{
		stagesBuffer:         stages_buffer[:0],
		currentMode:          core.Foot,
		workersCount:         workers_count,
		hasTookPublicTransit: false,
	}
}

type commuteContext struct {
	transportNetwork map[core.TilePosition]core.Tile
	subwayNetwork    []core.TilePosition
	ordinancesEffect map[cityhall.PolicyEffectKey]int16
	employerAnchor   core.TilePosition
	workerAnchor     core.TilePosition
	socialClass      core.SocialClass
}

type next struct {
	distance   uint
	position   core.TilePosition
	travelCost uint16
	found      bool
}

// TrafficLoop updates traffic values along the transport network
// from a worker's tile to its employer's tile.
//
// Each worker trip is calculated following these steps:
//
//   - walking distance is 4 tiles.
//   - find if employer is accessible by walking, (excluding worker and employer tiles).
//   - else, start loop and attempt to reach employer stage by stage.
type TrafficLoop struct {
	*cityhall.Deputies
	Chunks   []core.Chunk
	haloWalk []cityhall.HaloTile
}

func NewTrafficLoop(city_hall *cityhall.Deputies) *TrafficLoop {
	return &TrafficLoop{
		Deputies: city_hall,
		Chunks:   chunks.SixteenSquares[:],
		haloWalk: city_hall.Topology.Maker.OffsetsHalo(max_walking_steps),
	}
}

func (t *TrafficLoop) ID() string {
	return traffic_loop
}

func (t *TrafficLoop) Tick(
	ctx context.Context,
	city_id int32,
	chunk core.Chunk,
) (*gameloops.Update, error) {
	data, err := t.GetCityDataForTrafficCompute(ctx, city_id, chunk)
	if err != nil {
		return nil, err
	}

	if len(data.WorkerEmployments) == 0 {
		return nil, nil
	}

	err = t.TaskForce.Executor.Task(ctx, func(task_ctx context.Context) error {
		transport_network, subway_network := t.Transportation.Mapper.MakeNetwork(data.TransportNetwork)
		ordinances_effects := t.Policies.Mapper.MakeEffectiveTravelEffects(data.EnactedOrdinances)

		walk_buffer := make([]core.TilePosition, 0, max_walking_steps+1)
		stages_buffer := make([]cityhall.Stage, 0, max_travel_steps)
		all_commutings := make(map[cityhall.CommuteKey][]cityhall.Stage, len(data.WorkerEmployments))

		for _, item := range data.WorkerEmployments {
			context := commuteContext{
				transportNetwork: transport_network,
				subwayNetwork:    subway_network,
				ordinancesEffect: ordinances_effects,
				employerAnchor:   item.EmployerAnchor,
				workerAnchor:     item.WorkerAnchor,
				socialClass:      registry.SocialClasses[registry.ArtifactsRelations[item.WorkerArtifact].Kind],
			}
			key := cityhall.CommuteKey{
				WorkerAnchor:   context.workerAnchor,
				EmployerAnchor: context.employerAnchor,
			}
			commute := newCommute(stages_buffer, item.JobsCount)

			distance_to_work := t.Topology.Computor.Distance(context.workerAnchor, context.employerAnchor)

			if distance_to_work <= 1 {
				continue
			}

			// employer is accessible by walking, try a direct line to employer_anchor.
			if distance_to_work <= max_walking_steps {
				direct_path := t.Topology.Maker.Line(context.workerAnchor, context.employerAnchor, walk_buffer)
				if direct_path != nil {
					// exclude worker_anchor and employer_anchor
					commute.push(direct_path[1 : len(direct_path)-1]...)
					all_commutings[key] = slices.Clone(commute.stagesBuffer)
					continue
				}
			}

			// Could not walk directly to employer tile, let's try to use the transport network.
			walk_to_transport_path := t.getPathFromWorkerToNearestTransport(context, walk_buffer)
			if walk_to_transport_path == nil {
				continue
			}

			is_reachable := t.checkIfEmployerReachableByTransport(context)
			if !is_reachable {
				continue
			}

			// exclude worker_anchor
			commute.push(walk_to_transport_path[1:]...)

			start_artifact := registry.Artifacts[context.transportNetwork[commute.lastPosition()].Artifact.Name]
			_, transport_mode, found := t.getNewTransport(context, commute, start_artifact)
			if !found {
				slog.Error("error in", traffic_loop, ErrTransportMode)
				continue
			}

			commute.setMode(transport_mode)
			step := 0

			// find tile by tile the stages to reach any valid target.
			for step <= max_travel_steps {
				next_stage := t.getNextStage(context, commute)
				if !next_stage.found {
					break
				}

				if next_stage.position == context.employerAnchor {
					all_commutings[key] = slices.Clone(commute.stagesBuffer)
					break
				}

				commute.push(next_stage.position)
				step += int(next_stage.travelCost)

				if next_stage.distance <= max_walking_steps {
					final_walk := t.tryCommuteEnding(context, commute, next_stage.distance, walk_buffer)
					if final_walk == nil {
						continue
					}

					commute.setMode(core.Foot)
					commute.push(final_walk[1 : len(final_walk)-1]...)
					all_commutings[key] = slices.Clone(commute.stagesBuffer)
					break
				}
			}
		}

		slog.Debug("traffic", "all_commutings", all_commutings)

		// FIXME put a condition before calling ResetCityTrafficFromChunk
		return t.Transportation.Updater.ResetCityTrafficFromChunk(
			task_ctx,
			city_id,
			chunk,
			all_commutings,
		)
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// getPathFromWorkerToNearestTransport finds a direct line
// between a worker's anchor and the nearest transport tile.
//
// (see TransportNetworkArtifacts in registry).
//
// NOTE the returned path must contain at least 2 elements,
// the source and a transport tile position.
func (t *TrafficLoop) getPathFromWorkerToNearestTransport(
	context commuteContext,
	path_buffer []core.TilePosition,
) []core.TilePosition {
	for _, tile := range t.haloWalk {
		position := core.TilePosition{
			Q: uint8(int16(context.workerAnchor.Q) + int16(tile.Q)),
			R: uint8(int16(context.workerAnchor.R) + int16(tile.R)),
		}

		_, ok := context.transportNetwork[position]
		if !ok {
			continue
		}

		path := t.Topology.Maker.Line(context.workerAnchor, position, path_buffer)
		if path == nil {
			continue
		}

		return path
	}

	return nil
}

// checkIfEmployerReachableByTransport does not return the
// potential paths in order to optimize and give more freedom
// for a travel to find a suitable ending of commute.
func (t *TrafficLoop) checkIfEmployerReachableByTransport(context commuteContext) bool {
	for _, tile := range t.haloWalk {
		position := core.TilePosition{
			Q: uint8(int16(context.employerAnchor.Q) + int16(tile.Q)),
			R: uint8(int16(context.employerAnchor.R) + int16(tile.R)),
		}

		_, ok := context.transportNetwork[position]
		if !ok {
			continue
		}

		ok = t.Topology.Checker.HasLine(context.employerAnchor, position)
		if !ok {
			continue
		}

		return true
	}

	return false
}

// getNextStage is the main way of reaching the employer position.
// It first checks for a special case: current mode of transport is Subway
// and a station getting close to the employer can be found.
// Else, it looks for the 6 neighboring tiles and goes to the most suitable
// depending on a variety of condition (reduces distance to employer, social class, ordinances).
// This "algorithm" is heavily inspired by how SimCity-2000 probably did implement traffic.
// So some randomness is incorporatd when chossing the most suitable next stage of commute.
func (t *TrafficLoop) getNextStage(
	context commuteContext,
	commute *commute,
) next {
	last_position := commute.lastPosition()
	current_distance := t.Topology.Computor.Distance(last_position, context.employerAnchor)

	if commute.currentMode == core.Subway {
		next_station := t.teleportToNextSubwayStation(last_position, current_distance, context)

		if next_station.found {
			return next_station
		}
	}

	transport_change := next{}
	closest := next{}
	equal := next{}
	farthest := next{}
	var new_transport_mode core.EffectType

	current_q := int8(last_position.Q)
	current_r := int8(last_position.R)

	for _, offset := range cityhall.NeighborOffsets {
		neighbor := core.TilePosition{
			Q: uint8(current_q + offset.Q),
			R: uint8(current_r + offset.R),
		}

		if neighbor == context.employerAnchor {
			return next{0, neighbor, 0, true}
		}

		tile, ok := context.transportNetwork[neighbor]
		if !ok {
			continue
		}

		artifact := registry.Artifacts[tile.Artifact.Name]
		new_distance := t.Topology.Computor.Distance(neighbor, context.employerAnchor)
		travel_cost := getTravelCost(context, artifact, commute.currentMode)

		if travel_cost == 0 {
			change_cost, mode, ok := t.getNewTransport(context, commute, artifact)
			if !ok {
				continue
			}

			if transport_change.found && transport_change.travelCost < change_cost {
				continue
			}

			new_transport_mode = mode
			transport_change = next{new_distance, neighbor, change_cost, true}
			continue
		}

		if !closest.found && new_distance < current_distance {
			closest = next{new_distance, neighbor, travel_cost, true}
			continue
		}

		if !equal.found && new_distance == current_distance {
			equal = next{new_distance, neighbor, travel_cost, true}
			continue
		}

		if !farthest.found && new_distance > current_distance {
			if commute.isPreviousStagePosition(neighbor) {
				travel_cost += move_back_penalty
			}
			farthest = next{new_distance, neighbor, travel_cost, true}
		}
	}

	switch {
	case closest.found && rand.Float32() > closest_bias:
		return closest
	case transport_change.found && rand.Float32() > float32(transport_change.travelCost)/10:
		commute.setMode(new_transport_mode)
		return transport_change
	case equal.found && rand.Float32() > intersection_bias:
		return equal
	case farthest.found:
		return farthest
	default:
		return next{}
	}
}

// getNewTransport applies SocialClass travel effect values
// to find out a suitable mode of transport with which a travel can occur.
//
// Most of the time, it will be core.Car, it's up to the player
// to build infrastructure that offer a variety of transports modes such as
// core.BusStop, core.SubwayStation, core.PedestrianBicycleWay or core.TrainStation.
//
// This is by design, since players are likely to build lots of roads
// instead of creating a very dense public transport system (budget constraints).
func (t *TrafficLoop) getNewTransport(
	context commuteContext,
	commute *commute,
	artifact *core.Artifact,
) (uint16, core.EffectType, bool) {
	var best_travel_cost uint16
	var new_transport core.EffectType
	var found bool

	for _, travel_effect := range artifact.Effects.TravelCost {
		if travel_effect.Type == core.Car && commute.hasTookPublicTransit {
			continue
		}

		ok := canSwitchTransport(travel_effect.Type, artifact.Name)
		if !ok {
			continue
		}

		cost := getTravelCost(context, artifact, travel_effect.Type)
		if cost == 0 {
			continue
		}

		if !found || cost < best_travel_cost {
			best_travel_cost = cost
			new_transport = travel_effect.Type
			found = true
		}
	}

	return best_travel_cost, new_transport, found
}

// teleportToNextSubwayStation is called for dealing with a special case,
// when commuting using the subway network. There aren't any subway lines per se,
// we just attempt to find the station closest to the employer tile.
func (t *TrafficLoop) teleportToNextSubwayStation(
	current_station core.TilePosition,
	current_distance uint,
	context commuteContext,
) next {
	best_distance := current_distance
	var best_position core.TilePosition
	var found bool

	for _, position := range context.subwayNetwork {
		if position == current_station {
			continue
		}

		if position == context.employerAnchor {
			return next{0, position, 0, true}
		}

		new_distance := t.Topology.Computor.Distance(position, context.employerAnchor)

		if new_distance < best_distance {
			best_distance = new_distance
			best_position = position
			found = true
		}
	}

	if !found {
		return next{}
	}

	travel_cost := getTravelCost(context, registry.Artifacts[core.SubwayStation], core.Subway)
	return next{best_distance, best_position, travel_cost, true}
}

// tryCommuteEnding is called in the parent loop if we're at walking distance from the
// employer tile. Some randomness is used to decide if one should offboard transport network
// and finish the commute by foot.
func (t *TrafficLoop) tryCommuteEnding(
	context commuteContext,
	commute *commute,
	distance uint,
	walk_buffer []core.TilePosition,
) []core.TilePosition {
	last_position := commute.lastPosition()
	artifact_name := context.transportNetwork[last_position].Artifact.Name

	if !canSwitchTransport(commute.currentMode, artifact_name) {
		return nil
	}

	final_path := t.Topology.Maker.Line(last_position, context.employerAnchor, walk_buffer)
	if final_path == nil {
		return nil
	}

	if rand.Float32()*2 > float32(uint(context.socialClass.AdditionalTravelCostFoot)+distance)/10 {
		return final_path
	}

	return nil
}

// canSwitchTransport tells if one can board/offboard a public transport.
//
//   - core.Bus requires a core.busStop
//   - core.Train requires a core.TrainStation
//   - core.Subway requires a core.SubwayStation
//
// example:
// artifact core.Road has core.Bus as travel effect but to use the bus,
// one must first take the bus at a bus stop.
// Same goes for going out of a bus.
func canSwitchTransport(
	travel_effect_type core.EffectType,
	artifact_name core.ArtifactType,
) bool {
	switch travel_effect_type {
	case core.Bus:
		return artifact_name == core.BusStop
	case core.Train:
		return artifact_name == core.TrainStation
	case core.Subway:
		return artifact_name == core.SubwayStation
	default:
		return true
	}
}

// getTravelCost adds 3 values together:
//
//   - base travel effect cost of the given artifact
//   - social class additional cost for the given travel effect
//   - ordinance additional value for the given travel effect
//
// base travel effect cost + social class additional cost is always > 0,
// but ordinance additional value can lower the number.
//
// Note that core.Car is favored in all SocialClasses.
// To give a chance to more vertuous modes of transport, one must
// build the corresponding infrastructure:
// core.BusStop, core.SubwayStation, core.PedestrianBicycleWay or core.TrainStation
//
// At minimum, 1 is returned if artifact is traversable by the given transport mode.
// 0 is returned if artifact is NOT traversable by the given transport mode.
func getTravelCost(
	context commuteContext,
	artifact *core.Artifact,
	transport_mode core.EffectType,
) uint16 {
	base_cost := artifact.TravelCostEffect(transport_mode)
	if base_cost == 0 {
		return 0
	}

	ordinance_impact := context.ordinancesEffect[cityhall.PolicyEffectKey{
		ArtifactType: artifact.Name,
		EffectType:   transport_mode,
	}]

	var social_class_cost uint8

	switch transport_mode {
	case core.Car:
		social_class_cost = context.socialClass.AdditionalTravelCostCar
	case core.Bicycle:
		social_class_cost = context.socialClass.AdditionalTravelCostBicycle
	case core.Bus:
		social_class_cost = context.socialClass.AdditionalTravelCostBus
	case core.Foot:
		social_class_cost = context.socialClass.AdditionalTravelCostFoot
	case core.Subway:
		social_class_cost = context.socialClass.AdditionalTravelCostSubway
	case core.Train:
		social_class_cost = context.socialClass.AdditionalTravelCostTrain
	}

	return uint16(max(1, int16(base_cost)+int16(social_class_cost)+ordinance_impact))
}

// push spreads and insert entries in commute stages buffer.
// currentMode must be updated before if it has changed.
func (c *commute) push(positions ...core.TilePosition) {
	for _, position := range positions {
		c.stagesBuffer = append(c.stagesBuffer, cityhall.Stage{
			TilePosition: position,
			TravelCostEffect: core.TravelCostEffect{
				Type:  c.currentMode,
				Value: c.workersCount,
			},
		})
	}
}

func (c *commute) setMode(mode core.EffectType) {
	c.currentMode = mode

	if mode.IsPublicTransitMode() {
		c.hasTookPublicTransit = true
	}
}

func (c *commute) lastPosition() core.TilePosition {
	return c.stagesBuffer[len(c.stagesBuffer)-1].TilePosition
}

func (c *commute) isPreviousStagePosition(position core.TilePosition) bool {
	return len(c.stagesBuffer) > 1 && c.stagesBuffer[len(c.stagesBuffer)-2].TilePosition == position
}
