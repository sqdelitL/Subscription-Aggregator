package http

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/handlers/subscribe"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/middleware"
	usecase "github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

type Router struct {
	interactor *usecase.Interactor
}

func NewRouter(subscribeInteractor *usecase.Interactor) *Router {
	return &Router{
		interactor: subscribeInteractor,
	}
}

func (rt *Router) SetupChi() *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.LoggerMiddleware)
		r.Use(chimiddleware.Timeout(3 * time.Second))
		
		r.Post("/subscriptions", subscribe.CreateHandler(rt.interactor))
		r.Get("/subscriptions", subscribe.GetAllHandler(rt.interactor))
		r.Get("/subscriptions/cost", subscribe.GetCostHandler(rt.interactor))
		r.Get("/subscriptions/{id}", subscribe.GetHandler(rt.interactor))
		r.Put("/subscriptions", subscribe.UpdateHandler(rt.interactor))
		r.Delete("/subscriptions/{id}", subscribe.DeleteHandler(rt.interactor))
	})

	return r
}
