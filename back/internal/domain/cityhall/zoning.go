// zoning groups the behaviors of **updating** or **reverting** tiles's zones.
package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

type ZoneParams struct {
	NewSuperKind           core.SuperKindType
	NaturalSuperKind       core.SuperKindType
	PublicSystemsSuperKind core.SuperKindType
	NewArtifactStatus      core.ArtifactStatus
	OldArtifactStatus      core.ArtifactStatus
	Coordinates            []core.TilePosition
}

type ZoneManager interface {
	Update(ctx context.Context, city_id int32, update_params ZoneParams) ([]core.Tile, error)
	Revert(ctx context.Context, city_id int32, revert_params ZoneParams) ([]core.Tile, error)
}

// Zoning manages actions over a city's tiles zoning (SuperKind).
//
// see core/superkind.go
type Zoning struct {
	Manager ZoneManager
}

func NewZoning(store Store) *Zoning {
	return &Zoning{
		Manager: store,
	}
}
