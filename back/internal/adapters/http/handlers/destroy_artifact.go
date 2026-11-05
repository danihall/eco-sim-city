package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	utilshandler "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers/utils"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

var _ utilshandler.Validator = (*destroyPayload)(nil)

type destroyPayload struct {
	Position utilshandler.Coordinates `json:"position"`
}

// Implements Validator (payload_decoder)
func (d *destroyPayload) Validate() error {
	err := d.Position.Check()
	if err != nil {
		return err
	}

	return nil
}

func DestroyArtifact(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[destroyPayload](w, r)
		if err != nil {
			return nil, err
		}

		deleted_tiles, err := api.DestroyArtifact(
			r.Context(),
			a.Deputies,
			a.EventBus,
			int32(city_id),
			payload.Position.IntoDomainPosition(),
		)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Body:   deleted_tiles,
			Status: http.StatusOK,
		}, nil
	}
}
