package main

import (
	"net/http"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/auth"
	"example.com/ecommerce/internal/categories"
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

	queries := repo.New(app.db)

	productService := products.NewService(queries)
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProduct)
	r.Get("/products/{id}", productHandler.FindProductById)

	categoryService := categories.NewService(queries)
	categoryHandler := categories.NewHandler(categoryService)
	r.Get("/categories", categoryHandler.ListCategories)
	r.Get("/categories/{id}", categoryHandler.FindCategoryById)

	authService := auth.NewService(queries, tokenAuth)
	authHandler := auth.NewHandler(authService)
	r.Post("/auth/register", authHandler.RegisterUser)
	r.Post("/auth/login", authHandler.Login)

	orderService := orders.NewService(queries, app.db)
	ordersHandler := orders.NewHandler(orderService)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(auth.JWTAuthenticator(queries))

		// 認証済みユーザー全員
		r.Post("/auth/logout", authHandler.Logout)
		r.Post("/orders", ordersHandler.PlaceOrder)
		r.Get("/users/me", authHandler.GetMe)

		// 管理者のみ
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAdmin)

			r.Post("/products", productHandler.CreateProduct)
			r.Put("/products/{id}", productHandler.UpdateProduct)
			r.Delete("/products/{id}", productHandler.DeleteProduct)

			r.Post("/categories", categoryHandler.CreateCategories)
			r.Put("/categories/{id}", categoryHandler.UpdateCategory)
			r.Delete("/categories/{id}", categoryHandler.DeleteCategory)
			r.Post("/categories/{id}/products", categoryHandler.AddProductToCategory)
			r.Delete("/categories/{id}/products/{productId}", categoryHandler.RemoveProductFromCategory)
		})
	})

	return r
}
