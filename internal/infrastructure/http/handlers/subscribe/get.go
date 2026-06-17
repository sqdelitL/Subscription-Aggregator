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

const getHandlerName = "GetSubscribe"

// GetHandler возвращает одну подписку по ID.
// @Summary      Получение подписки
// @Description  Возвращает запись о подписке по её идентификатору.
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  Subscribe
// @Failure      400  {object}  response.FailView
// @Failure      404  {object}  response.FailView
// @Failure      500  {object}  response.FailView
// @Router       /subscriptions/{id} [get]
func GetHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		paramID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			response.FailResponse(w, fmt.Errorf("sub id parse error. %v. %w", err, errs.JsonValidationError), getHandlerName)
			return
		}

		sub, err := subscribeInteractor.Get(ctx, uint64(id))
		if err != nil {
			response.FailResponse(w, err, getHandlerName)
			return
		}

		dto, err := fromDomain(sub)
		if err != nil {
			response.FailResponse(w, err, getHandlerName)
			return
		}

		response.SuccessResponse(w, dto, http.StatusOK)
	}
}
