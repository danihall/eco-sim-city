package chunks

import "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"

// SixteenSquares splits the grid into squares of 32*32 tiles,
// from top and left to bottom and right.
var SixteenSquares = [16]core.Chunk{
	{
		TopLeft:     core.TilePosition{Q: 0, R: 0},
		BottomRight: core.TilePosition{Q: 31, R: 31},
	},
	{
		TopLeft:     core.TilePosition{Q: 32, R: 0},
		BottomRight: core.TilePosition{Q: 63, R: 31},
	},
	{
		TopLeft:     core.TilePosition{Q: 64, R: 0},
		BottomRight: core.TilePosition{Q: 95, R: 31},
	},
	{
		TopLeft:     core.TilePosition{Q: 96, R: 0},
		BottomRight: core.TilePosition{Q: 127, R: 31},
	},

	{
		TopLeft:     core.TilePosition{Q: 0, R: 32},
		BottomRight: core.TilePosition{Q: 31, R: 63},
	},
	{
		TopLeft:     core.TilePosition{Q: 32, R: 32},
		BottomRight: core.TilePosition{Q: 63, R: 63},
	},
	{
		TopLeft:     core.TilePosition{Q: 64, R: 32},
		BottomRight: core.TilePosition{Q: 95, R: 63},
	},
	{
		TopLeft:     core.TilePosition{Q: 96, R: 32},
		BottomRight: core.TilePosition{Q: 127, R: 63},
	},

	{
		TopLeft:     core.TilePosition{Q: 0, R: 64},
		BottomRight: core.TilePosition{Q: 31, R: 95},
	},
	{
		TopLeft:     core.TilePosition{Q: 32, R: 64},
		BottomRight: core.TilePosition{Q: 63, R: 95},
	},
	{
		TopLeft:     core.TilePosition{Q: 64, R: 64},
		BottomRight: core.TilePosition{Q: 95, R: 95},
	},
	{
		TopLeft:     core.TilePosition{Q: 96, R: 64},
		BottomRight: core.TilePosition{Q: 127, R: 95},
	},

	{
		TopLeft:     core.TilePosition{Q: 0, R: 96},
		BottomRight: core.TilePosition{Q: 31, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 32, R: 96},
		BottomRight: core.TilePosition{Q: 63, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 64, R: 96},
		BottomRight: core.TilePosition{Q: 95, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 96, R: 96},
		BottomRight: core.TilePosition{Q: 127, R: 127},
	},
}
