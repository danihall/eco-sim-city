package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	utilshandler "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers/utils"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ utilshandler.Validator = (*buildPathwayPayload)(nil)

type buildPathwayPayload struct {
	ArtifactName core.ArtifactType          `json:"artifact_name"`
	Coordinates  []utilshandler.Coordinates `json:"coordinates"`
}

// Implements Validator (payload_decoder)
func (b *buildPathwayPayload) Validate() error {
	if b.ArtifactName == 0 {
		return ErrMissingArtifactName
	}

	if b.ArtifactName == core.UnknownArtifact {
		return ErrUnknownArtifactName
	}

	for index := range b.Coordinates {
		err := b.Coordinates[index].Check()
		if err != nil {
			return err
		}
	}

	return nil
}

func BuildPathway(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[buildPathwayPayload](w, r)
		if err != nil {
			return nil, err
		}

		coordinates := make([]core.TilePosition, len(payload.Coordinates))

		for index, position := range payload.Coordinates {
			coordinates[index] = position.IntoDomainPosition()
		}

		artifact_tiles, err := api.BuildPathway(
			r.Context(),
			a.Deputies,
			a.EventBus,
			int32(city_id),
			payload.ArtifactName,
			coordinates,
		)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Body:   artifact_tiles,
			Status: http.StatusOK,
		}, nil
	}
}
