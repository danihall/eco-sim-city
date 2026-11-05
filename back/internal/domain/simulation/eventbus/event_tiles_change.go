package eventbus

import (
	"errors"
	"log/slog"
	"slices"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var (
	ErrInvalidTileEvent = &domain.Error{
		Reason: errors.New("old tiles and new tiles have different length"),
		Key:    domain.ErrInvalid,
	}
)

type TilesChangeData struct {
	oldTiles          []core.Tile
	newTiles          []core.Tile
	enactedOrdinances []core.OrdinanceType
}

func NewTilesChangeData(
	old_tiles []core.Tile,
	new_tiles []core.Tile,
	enacted_ordinances []core.OrdinanceType,
) TilesChangeData {
	slices.SortFunc(old_tiles, core.ByRowThenColumn)
	slices.SortFunc(new_tiles, core.ByRowThenColumn)

	slog.Debug("new tiles change data", "old_tiles", old_tiles, "new_tiles", new_tiles)

	return TilesChangeData{
		oldTiles:          old_tiles,
		newTiles:          new_tiles,
		enactedOrdinances: enacted_ordinances,
	}
}

func (t TilesChangeData) OldTiles() []core.Tile {
	return t.oldTiles
}

func (t TilesChangeData) NewTiles() []core.Tile {
	return t.newTiles
}

func (t TilesChangeData) EnactedOrdinances() []core.OrdinanceType {
	return t.enactedOrdinances
}
