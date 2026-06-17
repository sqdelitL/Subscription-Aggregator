package subscribe

import (
	"net/http"

	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/response"
	"github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

const getAllHandlerName = "GetAllHandler"

// GetAllHandler возвращает список всех подписок.
// @Summary      Список подписок
// @Description  Возвращает все существующие записи о подписках.
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   Subscribe
// @Failure      500  {object}  response.FailView
// @Router       /subscriptions [get]
func GetAllHandler(subscribeInteractor *subscribe.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		subs, err := subscribeInteractor.List(ctx)
		if err != nil {
			response.FailResponse(w, err, getAllHandlerName)
			return
		}

		dto, err := fromDomains(subs)
		if err != nil {
			response.FailResponse(w, err, getAllHandlerName)
			return
		}

		response.SuccessResponse(w, dto, http.StatusOK)
	}
}
