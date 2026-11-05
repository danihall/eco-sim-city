package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

func GetCityTraffic(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		snapshot, err := api.GetCityTraffic(r.Context(), a.Deputies, int32(city_id))
		if err != nil {
			return nil, err
		}

		binary_payload := []any{
			*snapshot.Foot,
			*snapshot.Car,
			*snapshot.Bus,
			*snapshot.Bicycle,
			*snapshot.Train,
			*snapshot.Subway,
		}

		return &adapterhttp.Response{
			Body: adapterhttp.BinaryBody{
				Data: binary_payload,
			},
			Status: http.StatusOK,
		}, nil
	}
}
