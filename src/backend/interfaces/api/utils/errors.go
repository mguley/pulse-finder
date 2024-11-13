package utils

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Errors is a structured error handler with logging capabilities.
type Errors struct {
	Logger  *slog.Logger
	Handler *Handler
}

// NewErrors initializes a new Errors instance with a given logger.
func NewErrors(l *slog.Logger, h *Handler) *Errors {
	return &Errors{Logger: l, Handler: h}
}

// LogError logs the request method, URL, and error message for better context.
func (e *Errors) LogError(r *http.Request, err error) {
	e.Logger.Error("Error occurred", "error", err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

// ErrorResponse sends a JSON error message with the specified status and logs it if necessary.
func (e *Errors) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	e.LogError(r, fmt.Errorf("%v", message))
	e.logAndSend(w, r, status, map[string]any{"error": message})
}

// ServerErrorResponse logs a server error and sends a 500 Internal Server Error response.
func (e *Errors) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.LogError(r, err)
	e.ErrorResponse(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// NotFoundResponse sends a 404 Not Found response.
func (e *Errors) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	e.ErrorResponse(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// MethodNotAllowedResponse sends a 405 Method Not Allowed response.
func (e *Errors) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	e.ErrorResponse(w, r, http.StatusMethodNotAllowed,
		fmt.Sprintf("the %s method is not supported for this resource", r.Method))
}

// FailedValidationResponse sends a 422 Unprocessable Entity response with validation errors.
func (e *Errors) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	e.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// InvalidCredentialsResponse sends a 401 Unauthorized response for failed authentication attempts.
func (e *Errors) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	e.ErrorResponse(w, r, http.StatusUnauthorized, "invalid authentication credentials")
}

// InvalidAuthenticationTokenResponse sends a 401 Unauthorized response for invalid or missing tokens.
func (e *Errors) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	e.ErrorResponse(w, r, http.StatusUnauthorized, "invalid or missing authentication token")
}

// AuthenticationRequiredResponse sends a 401 Unauthorized response indicating that authentication is required.
func (e *Errors) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	e.ErrorResponse(w, r, http.StatusUnauthorized, "you must be authenticated to access this resource")
}

// InactiveAccountResponse sends a 403 Forbidden response indicating that the account needs activation.
func (e *Errors) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	e.ErrorResponse(w, r, http.StatusForbidden, "your user account must be activated to access this resource")
}

// NotPermittedResponse sends a 403 Forbidden response indicating insufficient permissions.
func (e *Errors) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	e.ErrorResponse(w, r, http.StatusForbidden, "your user account doesn't have permission to access this resource")
}

// logAndSend sends the JSON response and handles any errors that occur during writing.
func (e *Errors) logAndSend(w http.ResponseWriter, r *http.Request, status int, payload map[string]any) {
	if err := e.Handler.WriteJson(w, status, payload, nil); err != nil {
		e.LogError(r, err)
	}
}
