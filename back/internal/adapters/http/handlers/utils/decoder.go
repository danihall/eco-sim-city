package utilshandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	adapterhttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
)

var (
	ErrNoPayload = &adapterhttp.HTTPError{
		Reason: errors.New("missing payload"),
		Status: http.StatusBadRequest,
	}
	ErrPayloadTooLarge = &adapterhttp.HTTPError{
		Reason: errors.New("request payload too large"),
		Status: http.StatusRequestEntityTooLarge,
	}
	ErrPayloadWrongType = &adapterhttp.HTTPError{
		Reason: errors.New("wrong type in payload"),
		Status: http.StatusUnprocessableEntity,
	}
	ErrPayloadUnkownField = &adapterhttp.HTTPError{
		Reason: errors.New("json: unknown field"),
		Status: http.StatusUnprocessableEntity,
	}
)

type Validator interface {
	Validate() error
}

func Decode[P any](w http.ResponseWriter, r *http.Request) (*P, error) {
	payload := new(P)

	r.Body = http.MaxBytesReader(w, r.Body, 10_000)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(payload)
	if err != nil {
		return nil, getPayloadError(err)
	}

	checkable_payload, ok := any(payload).(Validator)
	if ok {
		err = checkable_payload.Validate()
		if err != nil {
			return nil, err
		}
	}

	return payload, nil
}

func getPayloadError(err error) error {
	var unmarshall_type_error *json.UnmarshalTypeError
	var max_bytes_error *http.MaxBytesError

	switch {
	case errors.Is(err, io.EOF):
		return ErrNoPayload
	case errors.As(err, &unmarshall_type_error):
		message := fmt.Errorf(
			"%w: field '%v', expected %v, received %v",
			ErrPayloadWrongType,
			unmarshall_type_error.Field,
			unmarshall_type_error.Type,
			unmarshall_type_error.Value)

		return &adapterhttp.HTTPError{
			Reason: message,
			Status: ErrPayloadWrongType.Status,
		}
	case errors.As(err, &max_bytes_error):
		return ErrPayloadTooLarge
	case strings.HasPrefix(err.Error(), "json: unknown field"):
		return ErrPayloadUnkownField
	}

	return err
}
