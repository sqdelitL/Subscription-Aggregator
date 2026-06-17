package subscribe

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sqdelitL/subscription-aggregator/internal/errs"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/response"
	"github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

const updateHandlerName = "UpdateSubscribe"

// UpdateHandler обновляет существующую подписку.
// @Summary      Обновление подписки
// @Description  Полностью заменяет данные подписки.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body Subscribe true "Параметры подписки (ID обязателен)"
// @Success      200  {object}  object "Пустой ответ"
// @Failure      400  {object}  response.FailView
// @Failure      404  {object}  response.FailView
// @Failure      500  {object}  response.FailView
// @Router       /subscriptions [put]
func UpdateHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var request Subscribe
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.FailResponse(w, fmt.Errorf("decode subscribe error. %v, %w", err, errs.JsonValidationError), updateHandlerName)
			return
		}

		err := request.ValidateUpdate()
		if err != nil {
			response.FailResponse(w, err, updateHandlerName)
			return
		}

		domain, err := toDomain(&request)
		if err != nil {
			response.FailResponse(w, err, updateHandlerName)
			return
		}

		err = subscribeInteractor.Update(ctx, domain)
		if err != nil {
			response.FailResponse(w, err, updateHandlerName)
			return
		}

		response.SuccessResponse(w, map[interface{}]interface{}{}, http.StatusOK)
	}
}
