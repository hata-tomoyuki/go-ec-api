package main

import (
	"net/http"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/address"
	"example.com/ecommerce/internal/auth"
	"example.com/ecommerce/internal/carts"
	"example.com/ecommerce/internal/categories"
	"example.com/ecommerce/internal/orders"
	"example.com/ecommerce/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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

	productService := products.NewService(queries, app.db)
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProduct)
	r.Get("/products/{id}", productHandler.FindProductById)

	categoryService := categories.NewService(queries, app.db)
	categoryHandler := categories.NewHandler(categoryService)
	r.Get("/categories", categoryHandler.ListCategories)
	r.Get("/categories/{id}", categoryHandler.FindCategoryById)
	r.Get("/categories/{id}/products", categoryHandler.ListProductsByCategory)

	authService := auth.NewService(app.db, queries, tokenAuth)
	authHandler := auth.NewHandler(authService)

	// 認証エンドポイント: IP ごとに 5 リクエスト/分
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, 1*time.Minute))
		r.Post("/auth/register", authHandler.RegisterUser)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
	})

	orderService := orders.NewService(queries, app.db)
	ordersHandler := orders.NewHandler(orderService)

	cartsService := carts.NewService(queries)
	cartsHandler := carts.NewHandler(cartsService)

	addressService := address.NewService(queries)
	addressHandler := address.NewHandler(addressService)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(auth.JWTAuthenticator(queries))

		// 認証済みユーザー全員
		r.Post("/auth/logout", authHandler.Logout)
		r.Get("/users/me", authHandler.GetMe)
		r.Put("/users/me", authHandler.UpdateMe)
		r.Put("/users/me/password", authHandler.UpdatePassword)

		r.Post("/orders", ordersHandler.PlaceOrder)
		r.Get("/orders", ordersHandler.ListOrdersByCustomerID)
		r.Get("/orders/{id}", ordersHandler.FindOrderById)
		r.Put("/orders/{id}/cancel", ordersHandler.CancelOrder)

		r.Post("/cart", cartsHandler.CreateCart)
		r.Post("/cart/items", cartsHandler.AddItemToCart)
		r.Get("/cart", cartsHandler.ShowCartItems)
		r.Put("/cart/items/{id}", cartsHandler.UpdateCartItemQuantity)
		r.Delete("/cart/items/{id}", cartsHandler.RemoveItemFromCart)
		r.Delete("/cart", cartsHandler.ClearCart)

		r.Get("/addresses", addressHandler.ListAddresses)
		r.Post("/addresses", addressHandler.CreateAddress)
		r.Get("/addresses/{id}", addressHandler.FindAddressById)
		r.Put("/addresses/{id}", addressHandler.UpdateAddress)
		r.Delete("/addresses/{id}", addressHandler.DeleteAddress)

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

			r.Get("/admin/orders", ordersHandler.ListAllOrders)
			r.Get("/admin/orders/{id}", ordersHandler.FindOrderByIdAdmin)
			r.Put("/admin/orders/{id}/status", ordersHandler.UpdateOrderStatus)
		})
	})

	return r
}
