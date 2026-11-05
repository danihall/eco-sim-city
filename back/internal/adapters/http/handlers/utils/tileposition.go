package utilshandler

import (
	"errors"
	"fmt"
	"net/http"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var (
	ErrMissingTilesCoords = &adapterhttp.HTTPError{
		Reason: errors.New("missing tiles_coordinates"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrMissingPositionField = &adapterhttp.HTTPError{
		Reason: errors.New("missing field in tiles_coordinates"),
		Status: http.StatusUnprocessableEntity,
	}
)

// Coordinates mirrors core.TilePosition but with pointers for its fields,
// so we can detect if front-end has effectively sent them.
type Coordinates struct {
	Q *uint8
	R *uint8
}

type CoordinatesSlice []Coordinates

func (c *Coordinates) Check() error {
	if c.Q == nil || c.R == nil {
		return fmt.Errorf("%w, got q=%v r=%v", ErrMissingPositionField, c.Q, c.R)
	}
	return nil
}

func (c *Coordinates) IntoDomainPosition() core.TilePosition {
	return core.TilePosition{
		Q: *c.Q,
		R: *c.R,
	}
}
