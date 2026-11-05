package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
)

func ServerSideEvents(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		messages, err := a.Connect(r.Context(), int32(city_id))
		if err != nil {
			return nil, err
		}

		defer func() {
			err := a.Disconnect(int32(city_id))
			if err != nil {
				slog.Error(err.Error())
			}
		}()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		controller := http.NewResponseController(w)

		slog.Info("sse instanciated", "city_id", city_id)

		for {
			select {
			case <-r.Context().Done():
				slog.Info("user disconnected", "city_id", city_id)
				return nil, nil
			case data := <-messages:
				data_json, err := json.Marshal(data)
				if err != nil {
					continue
				}
				_, err = fmt.Fprintf(w, "data: %s\n\n", data_json)
				if err != nil {
					slog.Error(err.Error())
					return nil, nil
				}

				slog.Debug("msg", "data", data)
				controller.Flush()
			}
		}
	}
}
