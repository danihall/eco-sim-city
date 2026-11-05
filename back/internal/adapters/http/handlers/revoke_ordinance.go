package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	utilshandler "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers/utils"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

func RevokeOrdinance(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[ordinancePayload](w, r)
		if err != nil {
			return nil, err
		}

		err = api.RevokeOrdinance(r.Context(), a.Deputies, int32(city_id), payload.OrdinanceName)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Status: http.StatusNoContent,
		}, nil
	}
}
