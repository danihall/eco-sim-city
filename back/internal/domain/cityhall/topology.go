package cityhall

import (
	"errors"
	"fmt"
	"math"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/data"
)

var (
	ErrArtifactSizeInvalid = &domain.Error{
		Reason: errors.New("no valid size for given artifact"),
		Key:    domain.ErrInvalid,
	}
	ErrInvalidChunkCoordinates = &domain.Error{
		Reason: errors.New("top_left corner must be above and left from bottom_right"),
		Key:    domain.ErrInvalid,
	}
	WarnNoActionableCoordinates = &domain.Error{
		Reason: errors.New("no coordinate is on ground or in-bounds"),
		Key:    domain.ErrNotModified,
	}
)

type CoordinatesComputor struct{}

type CoordinateChecker struct {
	groundTiles map[core.TilePosition]struct{}
	computor    CoordinatesComputor
}

type CoordinatesMaker struct {
	computor CoordinatesComputor
	checker  CoordinateChecker
}

// Topology deals with coordinates, geometrically speaking.
// It holds a private record of tiles that are on ground, and by definition in-bound.
type Topology struct {
	Computor CoordinatesComputor
	Checker  CoordinateChecker
	Maker    CoordinatesMaker
}

func NewTopology() *Topology {
	ground_tiles := make(map[core.TilePosition]struct{}, core.TilesCount)

	for index := range data.InitialTiles {
		base_tile := &data.InitialTiles[index]
		if base_tile.Artifact.Name == core.DeepWater ||
			base_tile.Artifact.Name == core.ShallowWater {
			continue
		}

		ground_tiles[base_tile.Position] = struct{}{}
	}

	coordinates_computor := CoordinatesComputor{}
	coordinates_checker := CoordinateChecker{ground_tiles, coordinates_computor}
	coordinate_maker := CoordinatesMaker{coordinates_computor, coordinates_checker}

	return &Topology{
		Computor: coordinates_computor,
		Checker:  coordinates_checker,
		Maker:    coordinate_maker,
	}
}

type HaloTile struct {
	Q, R     int8
	Distance uint8
}

type Neighbor struct {
	Q int8
	R int8
}

// see https://www.redblobgames.com/grids/hexagons/#neighbors
var NeighborOffsets = [6]Neighbor{
	{1, 0},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{0, -1},
	{1, -1},
}

// Position checks if given position if on-ground and in-bounds.
func (c CoordinateChecker) ValidPosition(position core.TilePosition) bool {
	if position.Q > 127 || position.R > 127 {
		return false
	}

	_, on_ground := c.groundTiles[position]

	return on_ground
}

func (c CoordinateChecker) HasLine(origin, destination core.TilePosition) bool {
	distance := c.computor.Distance(origin, destination)

	if distance == 0 {
		return false
	}

	if distance == 1 {
		return c.ValidPosition(destination)
	}

	origin_q := float64(origin.Q)
	origin_r := float64(origin.R)
	destination_q := float64(destination.Q)
	destination_r := float64(destination.R)

	for step := 0; step < int(distance)+1; step++ {
		current_point := float64(step) / float64(distance)

		q := origin_q + (destination_q-origin_q)*current_point
		r := origin_r + (destination_r-origin_r)*current_point
		position := roundToNearestTile(q, r)

		if !c.ValidPosition(position) {
			return false
		}
	}

	return true
}

// FromArtifact returns `coordinates` that are *on ground*,
// given a specific artifact and its anchor.
// (user_width, user_height params can be 0).
func (c CoordinatesMaker) FromArtifact(
	artifact *core.Artifact,
	anchor core.TilePosition,
	user_width, user_height uint8,
	buffer []core.TilePosition,
) ([]core.TilePosition, error) {
	effective_width, effective_height := artifact.EffectiveSize(user_width, user_height)
	if effective_width == 0 || effective_height == 0 {
		return nil, fmt.Errorf(
			"%w, got width=%d height=%d",
			ErrArtifactSizeInvalid,
			effective_width,
			effective_height,
		)
	}

	coordinates := getBuffer(buffer, int(effective_width*effective_height))

	has_free_size := artifact.HasFreeSize()
	for h := range effective_height {
		for w := range effective_width {
			position := core.TilePosition{
				Q: uint8(int(anchor.Q) + int(w)),
				R: uint8(int(anchor.R) + int(h)),
			}

			if !c.checker.ValidPosition(position) {
				if !has_free_size {
					return nil, WarnNoActionableCoordinates
				}
				continue
			}

			coordinates = append(coordinates, position)
		}
	}

	if len(coordinates) == 0 {
		return nil, WarnNoActionableCoordinates
	}
	return coordinates, nil
}

// FromFixSizedArtifact returns `coordinates` that are *on ground*,
// given an artifact that has no free size and its anchor.
func (c CoordinatesMaker) FromFixSizedArtifact(
	artifact *core.Artifact,
	anchor core.TilePosition,
	buffer []core.TilePosition,
) ([]core.TilePosition, error) {
	effective_width := artifact.Width
	effective_height := artifact.Height

	coordinates := getBuffer(buffer, int(effective_width*effective_height))

	for h := range effective_height {
		for w := range effective_width {
			position := core.TilePosition{
				Q: uint8(int(anchor.Q) + int(w)),
				R: uint8(int(anchor.R) + int(h)),
			}

			if !c.checker.ValidPosition(position) {
				return nil, WarnNoActionableCoordinates
			}

			coordinates = append(coordinates, position)
		}
	}

	return coordinates, nil
}

// FilteredTiles returns coordinates that are *on ground* from given tiles.
// It does not return an error if tiles length is different from coordinates length.
func (c CoordinatesMaker) FilteredTiles(tiles []core.Tile) []core.TilePosition {
	coordinates := make([]core.TilePosition, 0, len(tiles))

	for index := range tiles {
		position := tiles[index].Position

		if !c.checker.ValidPosition(position) {
			continue
		}

		coordinates = append(coordinates, position)
	}

	return coordinates
}

// FromChunk returns `coordinates` that are *on ground*,
// given a chunk (defined by it top-left and bottom-right corners).
// NOTE top-left corner must be sound, meaning it must really be
// the top-left corner of the chunk.
func (c CoordinatesMaker) FromChunk(
	top_left core.TilePosition,
	bottom_right core.TilePosition,
) ([]core.TilePosition, error) {
	if top_left.Q > bottom_right.Q || top_left.R > bottom_right.R {
		return nil, ErrInvalidChunkCoordinates
	}

	width := int(bottom_right.Q-top_left.Q) + 1
	height := int(bottom_right.R-top_left.R) + 1
	coordinates := make([]core.TilePosition, 0, width*height)

	for r := range height {
		for q := range width {
			position := core.TilePosition{
				Q: uint8(int(top_left.Q) + q),
				R: uint8(int(top_left.R) + r),
			}

			if !c.checker.ValidPosition(position) {
				continue
			}

			coordinates = append(coordinates, position)
		}
	}

	if len(coordinates) == 0 {
		return nil, WarnNoActionableCoordinates
	}

	return coordinates, nil
}

// OffsetsHalo generates all tile's RELATIVE OFFSETS positions
// within `spread` distance from a center.
//
// For example:
//   - Spread = 1, 1 ring around the center.
//   - Spread = 2, 2 rings with the second circling around the first.
//   - Spread = 3, 3 rings with the third one circling around the second.
//   - and so on.
//
// OffsetsHalo(1) -> 6 tiles
// OffsetsHalo(2) -> 18 tiles
// OffsetsHalo(3) -> 36 tiles
// OffsetsHalo(4) -> 60 tiles
//
// NOTE: The positions can be off-ground (water tiles)
// and out-of-bounds since they're relative offsets.
// But they're guaranteed to be sorted from the inner ring toward the outer ring.
//
// formula taken from https://www.redblobgames.com/grids/hexagons/#range
// FIXED make the halo sorted to have rings circling their lower rings.
// FIXED include the outermost ring, it is more intuitive
func (c CoordinatesMaker) OffsetsHalo(spread int) []HaloTile {
	if spread <= 0 || spread > math.MaxInt8 {
		return nil
	}

	halo := make([]HaloTile, 0, 3*spread*(spread+1))

	for distance := 1; distance <= spread; distance++ {
		var q int8
		var r = int8(-distance)

		for _, hex := range NeighborOffsets {
			for step := 0; step < distance; step++ {
				halo = append(halo, HaloTile{
					Q:        q,
					R:        r,
					Distance: uint8(distance),
				})

				q += hex.Q
				r += hex.R
			}
		}
	}

	return halo
}

// Line returns the coordinates of the direct line between an origin position
// and a destination position.
//
// It includes the origin and destination.
//
// If the line encounters a out-of-bounds of off-ground tile,
// the line is considered invalid.
//
// see https://www.redblobgames.com/grids/hexagons/#line-drawing
func (c CoordinatesMaker) Line(
	origin, destination core.TilePosition,
	buffer []core.TilePosition,
) []core.TilePosition {
	distance := c.computor.Distance(origin, destination)

	if distance == 0 {
		return nil
	}

	line := getBuffer(buffer, int(distance)+1)

	if distance == 1 {
		if !c.checker.ValidPosition(destination) {
			return nil
		}
		line = append(line, origin, destination)
		return line
	}

	origin_q := float64(origin.Q)
	origin_r := float64(origin.R)
	destination_q := float64(destination.Q)
	destination_r := float64(destination.R)

	for step := 0; step < int(distance)+1; step++ {
		current_point := float64(step) / float64(distance)

		q := origin_q + (destination_q-origin_q)*current_point
		r := origin_r + (destination_r-origin_r)*current_point
		position := roundToNearestTile(q, r)

		if !c.checker.ValidPosition(position) {
			return nil
		}

		line = append(line, position)
	}

	return line
}

// Distance formula taken from https://www.redblobgames.com/grids/hexagons/#distances
func (c CoordinatesComputor) Distance(origin, destination core.TilePosition) uint {
	q_diff := int(origin.Q) - int(destination.Q)
	r_diff := int(origin.R) - int(destination.R)

	return uint((abs(q_diff) + abs(r_diff) + abs(q_diff+r_diff)) / 2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// roundToNearestTile takes floating coordinates and gives the closest tile.
// see https://www.redblobgames.com/grids/hexagons/#rounding
func roundToNearestTile(q, r float64) core.TilePosition {
	s := -q - r

	rounded_q := math.Round(q)
	rounded_r := math.Round(r)
	rounded_s := math.Round(s)

	q_diff := math.Abs(rounded_q - q)
	r_diff := math.Abs(rounded_r - r)
	s_diff := math.Abs(rounded_s - s)

	if q_diff > r_diff && q_diff > s_diff {
		rounded_q = -rounded_r - rounded_s
	} else if r_diff > s_diff {
		rounded_r = -rounded_q - rounded_s
	}

	return core.TilePosition{
		Q: uint8(rounded_q),
		R: uint8(rounded_r),
	}
}

func getBuffer(buffer []core.TilePosition, required_capacity int) []core.TilePosition {
	if cap(buffer) >= required_capacity {
		return buffer[:0]
	}

	return make([]core.TilePosition, 0, required_capacity)
}
