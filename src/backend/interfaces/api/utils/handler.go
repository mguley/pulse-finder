package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Handler struct{}

// NewHandler initializes and returns a new Handler instance.
func NewHandler() *Handler { return &Handler{} }

// requestBodyLimit specifies the maximum size of the request body (1MB).
const requestBodyLimit = int64(1_048_576)

// ErrInvalidIdParameter is an error for invalid ID parameters in requests.
var ErrInvalidIdParameter = errors.New("invalid id parameter")

// ExtractId extracts the "id" parameter from the HTTP request context and parses it into an int64.
// It returns an error if the parameter is missing or invalid.
func (h *Handler) ExtractId(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	idParam := params.ByName("id")

	if idParam == "" {
		return 0, fmt.Errorf("%w: missing 'id' parameter", ErrInvalidIdParameter)
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("%w: '%s' is not a valid integer ID", ErrInvalidIdParameter, idParam)
	}

	return id, nil
}

// WriteJson writes a JSON response to the client with the specified status code and headers.
// If data is nil, it sends an empty response.
func (h *Handler) WriteJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	h.setDefaultHeaders(w, headers)
	w.WriteHeader(status)

	if data == nil {
		return nil
	}

	// Marshal the data into JSON format and write it to the response.
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return fmt.Errorf("unable to marshal JSON response: %w", err)
	}

	if _, err = w.Write(append(payload, '\n')); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return fmt.Errorf("unable to write JSON response: %w", err)
	}

	return nil
}

// ReadJson reads and parses JSON data from the request body into the provided data structure.
// It enforces a body size limit and disallows unknown fields.
func (h *Handler) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, requestBodyLimit)

	// Create a JSON decoder that disallows unknown fields to enforce strict parsing.
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the request body into the provided data structure.
	if err := decoder.Decode(data); err != nil {
		return h.handleJsonDecodeError(err)
	}

	// Check for any additional, unexpected data in the request body.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON object")
	}

	return nil
}

// GetQueryString retrieves a string value from the query string or returns the provided default value.
func (h *Handler) GetQueryString(qs url.Values, key, value string) string {
	if !qs.Has(key) {
		return value
	}
	return qs.Get(key)
}

// GetQueryInt retrieves an integer value from the query string or returns the provided default value.
func (h *Handler) GetQueryInt(qs url.Values, key string, value int) int {
	if !qs.Has(key) {
		return value
	}
	if i, err := strconv.Atoi(qs.Get(key)); err == nil {
		return i
	}
	return value
}

// handleJSONDecodeError handles various types of errors returned during JSON decoding.
func (h *Handler) handleJsonDecodeError(err error) error {
	var (
		syntaxError           *json.SyntaxError
		unmarshalTypeError    *json.UnmarshalTypeError
		invalidUnmarshalError *json.InvalidUnmarshalError
		maxBytesError         *http.MaxBytesError
	)

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("body contains unknown field %s", fieldName)
	case errors.As(err, &maxBytesError):
		return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return err
	}
}

// setDefaultHeaders applies default headers to the response, allowing any provided headers to override them.
func (h *Handler) setDefaultHeaders(w http.ResponseWriter, headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
