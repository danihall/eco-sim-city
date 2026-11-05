package handlers

import (
	"errors"
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	utilshandler "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers/utils"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
)

var _ utilshandler.Validator = (*ordinancePayload)(nil)

var (
	ErrMissingOrdinanceName = &adapterhttp.HTTPError{
		Reason: errors.New("missing ordinance_name"),
		Status: http.StatusUnprocessableEntity,
	}
)

type ordinancePayload struct {
	OrdinanceName core.OrdinanceType `json:"ordinance_name"`
}

// Implements Validator (payload_decoder)
func (o *ordinancePayload) Validate() error {
	if o.OrdinanceName == 0 {
		return ErrMissingOrdinanceName
	}

	return nil
}

func EnactOrdinance(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[ordinancePayload](w, r)
		if err != nil {
			return nil, err
		}

		err = api.EnactOrdinance(r.Context(), a.Deputies, int32(city_id), payload.OrdinanceName)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Status: http.StatusNoContent,
		}, nil
	}
}
