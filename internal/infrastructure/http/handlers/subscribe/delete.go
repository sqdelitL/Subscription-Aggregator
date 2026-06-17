package subscribe

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/response"
	"github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

const deleteHandlerName = "DeleteSubscribe"

func DeleteHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		paramID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			response.FailResponse(w, fmt.Errorf("sub id parse error. %v. %w", err, errs.JsonValidationError), getHandlerName)
			return
		}

		err = subscribeInteractor.Delete(ctx, uint64(id))
		if err != nil {
			response.FailResponse(w, err, getHandlerName)
			return
		}

		response.SuccessResponse(w, map[interface{}]interface{}{}, http.StatusOK)
	}
}
