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

var _ utilshandler.Validator = (*buildPayload)(nil)

var (
	ErrMissingArtifactName = &adapterhttp.HTTPError{
		Reason: errors.New("missing artifact_name"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrUnknownArtifactName = &adapterhttp.HTTPError{
		Reason: errors.New("unknown artifact_name"),
		Status: http.StatusUnprocessableEntity,
	}
)

type buildPayload struct {
	Anchor       utilshandler.Coordinates `json:"anchor"`
	ArtifactName core.ArtifactType        `json:"artifact_name"`
	Width        uint8                    `json:"width"`
	Height       uint8                    `json:"height"`
}

// Implements Validator (payload_decoder)
func (b *buildPayload) Validate() error {
	if b.ArtifactName == 0 {
		return ErrMissingArtifactName
	}

	if b.ArtifactName == core.UnknownArtifact {
		return ErrUnknownArtifactName
	}

	err := b.Anchor.Check()
	if err != nil {
		return err
	}

	return nil
}

func BuildArtifact(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[buildPayload](w, r)
		if err != nil {
			return nil, err
		}

		artifact_tiles, err := api.BuildArtifact(
			r.Context(),
			a.Deputies,
			a.EventBus,
			int32(city_id),
			struct {
				Anchor       core.TilePosition
				ArtifactName core.ArtifactType
				UserWidth    uint8
				UserHeight   uint8
			}{payload.Anchor.IntoDomainPosition(), payload.ArtifactName, payload.Width, payload.Height},
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
