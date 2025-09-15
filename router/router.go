package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/salahfarzin/roja-shop/handlers"
	"github.com/salahfarzin/roja-shop/repositories"
	"github.com/salahfarzin/roja-shop/services"
)

func New(app *fiber.App, dbConn *sql.DB) {
	// Serve static files from public/images at /images
	app.Static("/images", "public/images")
	app.Static("/uploads", "storage/uploads")

	api := app.Group("/api/v1")

	productService := services.NewProduct(repositories.NewProduct(dbConn))
	productHandler := handlers.NewProduct(productService)

	// Product routes
	api.Get("/products", productHandler.Index)
	api.Post("/products/sell", productHandler.Sell)
	api.Post("/products", productHandler.Store)
}
