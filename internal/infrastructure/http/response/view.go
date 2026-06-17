package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sqdelitL/subscription-aggregator/internal/errs"
)

type FailView struct {
	Error string `json:"error"`
}

func SuccessResponse(w http.ResponseWriter, data interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode success response",
			"error", err,
		)
	}
}

func FailResponse(w http.ResponseWriter, err error, handlerName string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err == nil {
		slog.Error("FailResponse called with nil error", "handler", handlerName)
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(FailView{Error: "internal server error"})
		if err != nil {
			slog.Error("failed to encode fail response: %v", "err", err)
			return
		}
		return
	}

	statusCode := mapTransportCode(err)

	slog.Warn("request handling failed", "handler", handlerName, "error", err, "status", statusCode)

	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(FailView{Error: err.Error()}); encodeErr != nil {
		slog.Error("failed to encode error response", "error", encodeErr)
	}
}

func mapTransportCode(err error) int {
	switch {
	case errors.Is(err, errs.InternalError):
		return http.StatusInternalServerError
	case errors.Is(err, errs.NegativeSubscribePriceError):
		return http.StatusBadRequest
	case errors.Is(err, errs.WrongSubscribeDateIntervalError):
		return http.StatusBadRequest
	case errors.Is(err, errs.NotFoundSubscribeError):
		return http.StatusNotFound
	case errors.Is(err, errs.JsonValidationError):
		return http.StatusBadRequest
	default:
		slog.Error("unknown custom error", "error", err)
		return http.StatusInternalServerError
	}
}
