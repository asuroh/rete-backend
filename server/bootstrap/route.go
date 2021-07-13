package bootstrap

import (
	"retel-backend/pkg/logruslogger"
	api "retel-backend/server/handler"
	"retel-backend/server/middleware"

	chimiddleware "github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RegisterRoutes ...
func (boot *Bootup) RegisterRoutes() {
	handlerType := api.Handler{
		DB:         boot.DB,
		EnvConfig:  boot.EnvConfig,
		Validate:   boot.Validator,
		Translator: boot.Translator,
		ContractUC: &boot.ContractUC,
		Jwe:        boot.Jwe,
		Jwt:        boot.Jwt,
	}

	mJwt := middleware.VerifyMiddlewareInit{
		ContractUC: &boot.ContractUC,
	}

	boot.R.Route("/v1", func(r chi.Router) {
		// Define a limit rate to 1000 requests per IP per request.
		rate, _ := limiter.NewRateFromFormatted("1000-S")
		store, _ := sredis.NewStoreWithOptions(boot.ContractUC.Redis, limiter.StoreOptions{
			Prefix:   "limiter_global",
			MaxRetry: 3,
		})
		rateMiddleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
		r.Use(rateMiddleware.Handler)

		// Logging setup
		r.Use(chimiddleware.RequestID)
		r.Use(logruslogger.NewStructuredLogger(boot.EnvConfig["LOG_FILE_PATH"], boot.EnvConfig["LOG_DEFAULT"], boot.ContractUC.ReqID))
		r.Use(chimiddleware.Recoverer)

		// API
		r.Route("/api", func(r chi.Router) {
			userHandler := api.UserHandler{Handler: handlerType}
			r.Route("/user", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Post("/login", userHandler.LoginHandler)
					r.Post("/register", userHandler.RegisterHandler)
				})
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Get("/token", userHandler.GetByTokenHandler)
					r.Put("/token", userHandler.UpdateByTokenHandler)
					r.Post("/profile", userHandler.UploadImageHandler)
				})
			})

			productHandler := api.ProductHandler{Handler: handlerType}
			r.Route("/product", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Get("/", productHandler.GetAllHandler)
					r.Get("/id/{id}", productHandler.GetByIDHandler)
				})
			})

			categoryHandler := api.CategoryHandler{Handler: handlerType}
			r.Route("/category", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Get("/", categoryHandler.GetAllHandler)
				})
			})

			userCartHandler := api.UserCartHandler{Handler: handlerType}
			r.Route("/cart", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Get("/", userCartHandler.GetAllHandler)
					r.Post("/", userCartHandler.CheckoutHandler)
					r.Put("/id/{id}", userCartHandler.UpdateHandler)
					r.Delete("/id/{id}", userCartHandler.DeleteHandler)
				})
			})
		})
	})
}
