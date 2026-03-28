package main

import (
	"net/http"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/auth"
	"example.com/ecommerce/internal/orders"
	"example.com/ecommerce/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProduct)
	r.Get("/products/{id}", productHandler.FindProductById)

	authService := auth.NewService(repo.New(app.db), tokenAuth)
	authHandler := auth.NewHandler(authService)
	r.Post("/users/register", authHandler.RegisterUser)
	r.Post("/users/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(auth.JWTAuthenticator)

		r.Post("/products", productHandler.CreateProduct)
		r.Put("/products/{id}", productHandler.UpdateProduct)
		r.Delete("/products/{id}", productHandler.DeleteProduct)
	})

	orderService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderService)
	r.Post("/orders", ordersHandler.PlaceOrder)

	return r
}
