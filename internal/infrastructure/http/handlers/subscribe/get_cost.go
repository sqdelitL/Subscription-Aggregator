package subscribe

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sqdelitL/subscription-aggregator/internal/domain/port/repository"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/response"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/util"
	"github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

const getCostHandlerName = "GetCostHandler"

// GetCostHandler подсчитывает суммарную стоимость подписок за период.
// @Summary      Суммарная стоимость подписок
// @Description  Вычисляет общую стоимость всех подписок за указанный период с возможной фильтрацией по пользователю и/или сервису.
// @Tags         subscriptions
// @Produce      json
// @Param        start_date   query     string  true   "Начало периода (MM-YYYY)"
// @Param        end_date     query     string  false  "Конец периода (MM-YYYY), если не указан – период открытый"
// @Param        user_id      query     string  false  "ID пользователя (UUID)"
// @Param        service_name query     string  false  "Название сервиса"
// @Success      200  {object}  object "Итоговая сумма"
// @Failure      400  {object}  response.FailView
// @Failure      500  {object}  response.FailView
// @Router       /subscriptions/cost [get]
func GetCostHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		startDate := r.URL.Query().Get("start_date")
		if startDate == "" {
			response.FailResponse(w, fmt.Errorf("start_date required query param. %w", errs.JsonValidationError), getCostHandlerName)
			return
		}

		start, err := util.ParseMMYYYY(startDate)
		if err != nil {
			response.FailResponse(w, fmt.Errorf("start date format error. %w", errs.JsonValidationError), getCostHandlerName)
			return
		}

		endDate := r.URL.Query().Get("end_date")
		var end *time.Time
		if endDate != "" {
			parsed, err := util.ParseMMYYYY(endDate)
			if err != nil {
				response.FailResponse(w, fmt.Errorf("end date format error. %w", errs.JsonValidationError), getCostHandlerName)
				return
			}
			end = &parsed
		}

		userIDParam := r.URL.Query().Get("user_id")
		var userID *uuid.UUID
		if userIDParam != "" {
			id, err := uuid.Parse(userIDParam)
			if err != nil {
				response.FailResponse(w, fmt.Errorf("user id format error. %w", errs.JsonValidationError), getCostHandlerName)
				return
			}
			userID = &id
		}

		serviceNameParam := r.URL.Query().Get("service_name")
		var serviceName *string
		if serviceNameParam != "" {
			serviceName = &serviceNameParam
		}

		subs, err := subscribeInteractor.GetTotalCost(ctx, repository.CostFilter{
			UserID:      userID,
			Start:       start,
			End:         end,
			ServiceName: serviceName,
		})
		if err != nil {
			response.FailResponse(w, err, getCostHandlerName)
			return
		}

		response.SuccessResponse(w, subs, http.StatusOK)
	}
}
