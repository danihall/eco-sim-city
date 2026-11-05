package chunks

import (
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

// FourSquares contains 4 chunks.
// Each chunk occupies 2/3 of the whole grid.
var FourSquares = [4]core.Chunk{
	{
		TopLeft:     core.TilePosition{Q: 0, R: 0},
		BottomRight: core.TilePosition{Q: 85, R: 85},
	},
	{
		TopLeft:     core.TilePosition{Q: 42, R: 0},
		BottomRight: core.TilePosition{Q: 127, R: 85},
	},
	{
		TopLeft:     core.TilePosition{Q: 0, R: 42},
		BottomRight: core.TilePosition{Q: 85, R: 127},
	},
	{
		TopLeft:     core.TilePosition{Q: 42, R: 42},
		BottomRight: core.TilePosition{Q: 127, R: 127},
	},
}
