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

// DeleteHandler удаляет подписку по ID.
// @Summary      Удаление подписки
// @Description  Удаляет запись о подписке.
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  object "Пустой ответ"
// @Failure      400  {object}  response.FailView
// @Failure      404  {object}  response.FailView
// @Failure      500  {object}  response.FailView
// @Router       /subscriptions/{id} [delete]
func DeleteHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		paramID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			response.FailResponse(w, fmt.Errorf("sub id parse error. %v. %w", err, errs.JsonValidationError), deleteHandlerName)
			return
		}

		err = subscribeInteractor.Delete(ctx, uint64(id))
		if err != nil {
			response.FailResponse(w, err, deleteHandlerName)
			return
		}

		response.SuccessResponse(w, map[interface{}]interface{}{}, http.StatusOK)
	}
}
