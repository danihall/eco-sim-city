package handlers

import (
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

func GetCityTilesArtifacts(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		snapshot, err := api.GetCityTilesArtifacts(r.Context(), a.Deputies, int32(city_id))
		if err != nil {
			return nil, err
		}

		binary_payload := []any{
			*snapshot.ArtifactIDs,
			*snapshot.AnchorQs,
			*snapshot.AnchorRs,
			*snapshot.ArtifactStatuses,
			*snapshot.SuperKindIDs,
		}

		return &adapterhttp.Response{
			Body: adapterhttp.BinaryBody{
				Data: binary_payload,
			},
			Status: http.StatusOK,
		}, nil
	}
}
