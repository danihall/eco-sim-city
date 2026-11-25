package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

func GetAllOrdinances(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		all_ordinances, err := api.GetAllOrdinances(r.Context(), a.Deputies, int32(city_id))
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Body:   all_ordinances,
			Status: http.StatusOK,
		}, nil
	}
}
