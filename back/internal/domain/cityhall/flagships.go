package cityhall

import (
	"context"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

type FlagshipManager interface {
	Add(
		ctx context.Context,
		city_id int32,
		artifact_name core.ArtifactType,
	) error
	Remove(
		ctx context.Context,
		city_id int32,
		artifact_name core.ArtifactType,
	) error
}

type FlagshipGetter interface {
	All(
		ctx context.Context,
		city_id int32,
	) ([]core.ArtifactType, error)
}

// Flagships manages actions over unique artifacts for a city.
type Flagship struct {
	Manager FlagshipManager
	Getter  FlagshipGetter
}

func NewFlagship(store Store) *Flagship {
	return &Flagship{
		Manager: store,
		Getter:  store,
	}
}
