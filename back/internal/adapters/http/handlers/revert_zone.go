package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	utilshandler "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers/utils"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

var _ utilshandler.Validator = (*revertPayload)(nil)

type revertPayload struct {
	TopLeft     utilshandler.Coordinates `json:"top_left"`
	BottomRight utilshandler.Coordinates `json:"bottom_right"`
}

// Implements Validator (payload_decoder)
func (r *revertPayload) Validate() error {
	err := r.TopLeft.Check()
	if err != nil {
		return err
	}

	err = r.BottomRight.Check()
	if err != nil {
		return err
	}

	return nil
}

func RevertZone(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[revertPayload](w, r)
		if err != nil {
			return nil, err
		}

		reverted_tiles, err := api.RevertZone(
			r.Context(),
			a.Deputies,
			int32(city_id),
			payload.TopLeft.IntoDomainPosition(),
			payload.BottomRight.IntoDomainPosition(),
		)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Body:   reverted_tiles,
			Status: http.StatusOK,
		}, nil
	}
}
