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

var _ utilshandler.Validator = (*updatePayload)(nil)

var (
	ErrMissingSuperKind = &adapterhttp.HTTPError{
		Reason: errors.New("missing super_kind_name"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrUnknownSuperKind = &adapterhttp.HTTPError{
		Reason: errors.New("unknown super_kind_name"),
		Status: http.StatusUnprocessableEntity,
	}
)

type updatePayload struct {
	ZoneName    core.SuperKindType       `json:"super_kind_name"`
	TopLeft     utilshandler.Coordinates `json:"top_left"`
	BottomRight utilshandler.Coordinates `json:"bottom_right"`
}

// Implements Validator (payload_decoder)
func (u *updatePayload) Validate() error {
	if u.ZoneName == 0 {
		return ErrMissingSuperKind
	}

	if u.ZoneName == core.UnknownSuperKind {
		return ErrUnknownSuperKind
	}

	err := u.TopLeft.Check()
	if err != nil {
		return err
	}

	err = u.BottomRight.Check()
	if err != nil {
		return err
	}

	return nil
}

func UpdateZone(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		city_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return nil, err
		}

		payload, err := utilshandler.Decode[updatePayload](w, r)
		if err != nil {
			return nil, err
		}

		updated_tiles, err := api.UpdateZone(
			r.Context(),
			a.Deputies,
			int32(city_id),
			payload.ZoneName,
			payload.TopLeft.IntoDomainPosition(),
			payload.BottomRight.IntoDomainPosition(),
		)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Body:   updated_tiles,
			Status: http.StatusOK,
		}, nil
	}
}
