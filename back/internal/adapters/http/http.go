package http

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
)

type Response struct {
	Body   any
	Status int
}

type BinaryBody struct {
	Data []any
}

type HTTPError struct {
	Reason error
	Status int
}

// Error implements error interface
func (h *HTTPError) Error() string {
	return h.Reason.Error()
}

var (
	ErrPathInvalid = &HTTPError{
		Reason: errors.New("invalid path"),
		Status: http.StatusNotFound,
	}
	ErrPathInvalidNumber = &HTTPError{
		Reason: errors.New("invalid number in path"),
		Status: http.StatusBadRequest,
	}
)

func _getError(err error) *HTTPError {
	var num_error *strconv.NumError
	if errors.As(err, &num_error) {
		return ErrPathInvalidNumber
	}

	var payload_error *HTTPError
	if errors.As(err, &payload_error) {
		return &HTTPError{
			Reason: err,
			Status: payload_error.Status,
		}
	}

	var domain_error *domain.Error
	if errors.As(err, &domain_error) {
		return &HTTPError{
			Reason: err,
			Status: _mapErrKeyToStatusCode(domain_error.Key),
		}
	}

	return &HTTPError{Reason: err, Status: http.StatusInternalServerError}
}

func _mapErrKeyToStatusCode(key domain.ErrKey) int {
	switch key {
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrInvalid:
		return http.StatusBadRequest
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrNotModified:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

type CustomHandler func(w http.ResponseWriter, r *http.Request) (*Response, error)

// ServeHTTP satisfies http.Handler interface
func (c CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := c(w, r)

	if err != nil {
		api_error := _getError(err)

		slog.Error("HTTP request failed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", api_error.Status,
			"error", err,
		)

		http.Error(w, api_error.Reason.Error(), api_error.Status)
		return
	}

	if response == nil {
		return
	}

	if response.Body == nil {
		w.WriteHeader(response.Status)
		return
	}

	binary_body, ok := response.Body.(BinaryBody)
	if ok {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(response.Status)

		for _, item := range binary_body.Data {
			err := binary.Write(w, binary.LittleEndian, item)
			if err != nil {
				slog.Error("Failed to write binary item to network", "error", err)
				return
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response.Body)
}

func StrictURLMiddleware(h http.Handler) CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) (*Response, error) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")

		start := time.Now()
		if strings.Contains(r.RequestURI, "//") {
			return nil, ErrPathInvalid
		}

		h.ServeHTTP(w, r)
		duration := time.Since(start)
		slog.Info("api", "request URI", r.RequestURI, "response time", duration)
		return nil, nil
	}
}
