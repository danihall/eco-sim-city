package handlers

import (
	"errors"
	"net/http"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	utilshandler "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers/utils"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
)

var _ utilshandler.Validator = (*cityPayload)(nil)

var (
	ErrMissingEmail = &adapterhttp.HTTPError{
		Reason: errors.New("missing email"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrMissingPassword = &adapterhttp.HTTPError{
		Reason: errors.New("missing password"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrPasswordTooShort = &adapterhttp.HTTPError{
		Reason: errors.New("password too short. minimum's length is 4"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrPasswordTooLong = &adapterhttp.HTTPError{
		Reason: errors.New("password too long lol"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrMissingCityName = &adapterhttp.HTTPError{
		Reason: errors.New("missing city's name"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrMissingMayorName = &adapterhttp.HTTPError{
		Reason: errors.New("missing city's mayor's name"),
		Status: http.StatusUnprocessableEntity,
	}
)

type cityPayload struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	MayorName string `json:"mayor_name"`
}

// Implements Validator (payload_decoder)
func (c *cityPayload) Validate() error {
	switch {
	case c.Email == "":
		return ErrMissingEmail
	case c.Password == "":
		return ErrMissingPassword
	case len(c.Password) < 5:
		return ErrPasswordTooShort
	case len(c.Password) > 72:
		return ErrPasswordTooLong
	case c.Name == "":
		return ErrMissingCityName
	case c.MayorName == "":
		return ErrMissingMayorName
	}

	return nil
}

func CreateFullCity(a *app.App) adapterhttp.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*adapterhttp.Response, error) {
		payload, err := utilshandler.Decode[cityPayload](w, r)
		if err != nil {
			return nil, err
		}

		new_city, err := api.CreateFullCity(
			r.Context(),
			a.Deputies,
			a.Sim,
			payload.Email,
			payload.Password,
			payload.Name,
			payload.MayorName,
		)
		if err != nil {
			return nil, err
		}

		return &adapterhttp.Response{
			Body:   new_city,
			Status: http.StatusCreated,
		}, nil
	}
}
