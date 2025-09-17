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

	saleRepo := repositories.NewSaleRepo(dbConn)
	productRepo := repositories.NewProduct(dbConn)
	saleService := services.NewSale(saleRepo, productRepo)
	productHandler := handlers.NewProduct(services.NewProduct(productRepo), saleService)

	// Product routes
	api.Get("/products", productHandler.Index)
	api.Post("/products", productHandler.Create)
	api.Put("/products/:id", productHandler.Update)
	api.Patch("/products/:id", productHandler.Update)
	api.Post("/products/sell/:id", productHandler.Sell)

	// Catch-all route for undefined paths (404)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not Found",
		})
	})
}
