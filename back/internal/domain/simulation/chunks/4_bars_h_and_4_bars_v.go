package chunks

import (
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// FourBarsHorizontalAndFourBarsVertical contains 8 chunks:
//
//   - 4 chunks spanning the whole grid's width (Horizontal bands) and some padding
//   - 4 chunks spanning the whole grid's height (Vertical bands) + some padding
var FourBarsHorizontalAndFourBarsVertical = [8]core.Chunk{
	{
		TopLeft:     core.TilePosition{Q: 0, R: 0},
		BottomRight: core.TilePosition{Q: 127, R: 35},
	},
	{
		TopLeft:     core.TilePosition{Q: 0, R: 28},
		BottomRight: core.TilePosition{Q: 127, R: 67},
	},
	{
		TopLeft:     core.TilePosition{Q: 0, R: 60},
		BottomRight: core.TilePosition{Q: 127, R: 99},
	},
	{
		TopLeft:     core.TilePosition{Q: 0, R: 92},
		BottomRight: core.TilePosition{Q: 127, R: 127},
	},

	{
		TopLeft:     core.TilePosition{Q: 0, R: 0},
		BottomRight: core.TilePosition{Q: 35, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 28, R: 0},
		BottomRight: core.TilePosition{Q: 67, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 60, R: 0},
		BottomRight: core.TilePosition{Q: 99, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 92, R: 0},
		BottomRight: core.TilePosition{Q: 127, R: 127},
	},
}
