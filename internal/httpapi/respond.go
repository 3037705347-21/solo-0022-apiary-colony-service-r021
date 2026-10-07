package httpapi

import (
	"encoding/json"
	"net/http"

	"example.com/apiary-colony-service/internal/domain"
)

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type errorBody struct {
	Error errorDetail `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	code := domain.CodeOf(err)
	detail := errorDetail{Code: string(code), Message: err.Error()}
	if domainErr, ok := domain.AsError(err); ok {
		detail.Message = domainErr.Message
		detail.Field = domainErr.Field
	}
	writeJSON(w, statusFor(code), errorBody{Error: detail})
}

func statusFor(code domain.ErrorCode) int {
	switch code {
	case domain.CodeInvalidRequest:
		return http.StatusBadRequest
	case domain.CodeNotFound:
		return http.StatusNotFound
	case domain.CodeConflict, domain.CodeInvalidState:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
