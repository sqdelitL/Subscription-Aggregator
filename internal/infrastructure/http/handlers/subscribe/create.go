package subscribe

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sqdelitL/subscription-aggregator/internal/errs"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/response"
	"github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

const createHandlerName = "CreateSubscribe"

func CreateHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var request Subscribe
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.FailResponse(w, fmt.Errorf("decode subscribe error. %v, %w", err, errs.JsonValidationError), createHandlerName)
			return
		}

		err := request.ValidateCreate()
		if err != nil {
			response.FailResponse(w, err, createHandlerName)
			return
		}

		domain, err := toDomain(&request)
		if err != nil {
			response.FailResponse(w, err, createHandlerName)
			return
		}

		id, err := subscribeInteractor.Create(ctx, domain)
		if err != nil {
			response.FailResponse(w, err, createHandlerName)
			return
		}

		response.SuccessResponse(w, map[interface{}]interface{}{
			"id": id,
		}, http.StatusCreated)
	}
}
