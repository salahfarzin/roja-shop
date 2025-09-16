package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/salahfarzin/roja-shop/configs"
	"github.com/salahfarzin/roja-shop/services"
	"github.com/salahfarzin/roja-shop/types"
)

var (
	soldCount int
)

type Product interface {
	Index(ctx *fiber.Ctx) error
	Store(ctx *fiber.Ctx) error
	Sell(ctx *fiber.Ctx) error
}

type product struct {
	service services.Product
}

func NewProduct(service services.Product) Product {
	return &product{service: service}
}

func (p *product) Index(ctx *fiber.Ctx) error {
	products, err := p.service.GetAll()
	if err != nil {
		return errorResponse(ctx, fiber.StatusInternalServerError, "could not fetch products from database")
	}

	serverURL := ctx.BaseURL()
	for i := range products {
		if path := products[i].File.Path; path != "" {
			if len(path) > 0 && path[0] == '/' {
				fmt.Println("Image path:", serverURL+products[i].File.Path)
				products[i].File.Path = serverURL + products[i].File.Path
			} else {
				products[i].File.Path = serverURL + "/" + path
			}
		}
	}

	return ctx.JSON(products)
}

func (p *product) Store(ctx *fiber.Ctx) error {
	// Parse form fields
	title := ctx.FormValue("title")
	discountStr := ctx.FormValue("discount")
	priceStr := ctx.FormValue("price")
	oldPriceStr := ctx.FormValue("oldPrice")
	brand := ctx.FormValue("brand")
	inventoryStr := ctx.FormValue("inventory")
	description := ctx.FormValue("description")
	// Parse details and styleNotes as JSON if provided
	details := make(map[string]string)
	styleNotes := make(map[string]string)
	_ = ctx.BodyParser(&details)    // Optionally parse from form or JSON
	_ = ctx.BodyParser(&styleNotes) // Optionally parse from form or JSON

	if title == "" || priceStr == "" {
		return errorResponse(ctx, fiber.StatusBadRequest, "title and price are required")
	}

	price, _ := strconv.ParseFloat(priceStr, 64)
	oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)
	inventory, _ := strconv.Atoi(inventoryStr)
	discount, _ := strconv.ParseFloat(discountStr, 64)

	// Handle image upload
	file, _, err := services.UploadService.SaveFile(ctx, "image", configs.New())
	if err != nil && err != fiber.ErrUnprocessableEntity {
		if errors.Is(err, services.ErrNotImage) {
			return errorResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
		return errorResponse(ctx, fiber.StatusInternalServerError, "failed to save image")
	}

	productID := uuid.New().String()
	product := types.Product{
		ID:          productID,
		Brand:       brand,
		Title:       title,
		Inventory:   inventory,
		Price:       price,
		Discount:    discount,
		OldPrice:    oldPrice,
		Description: description,
		Details:     details,
		StyleNotes:  styleNotes,
		File:        file,
	}

	_, err = p.service.Store(product, file)
	if err != nil {
		return errorResponse(ctx, fiber.StatusInternalServerError, "failed to store product")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product created",
		"product": product,
	})
}

func (p *product) Sell(ctx *fiber.Ctx) error {
	var body map[string]any
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	uuid, ok := body["uuid"].(string)
	if !ok || uuid == "" {
		return errorResponse(ctx, fiber.StatusBadRequest, "uuid is required")
	}

	soldCount++
	return ctx.JSON(fiber.Map{"message": "received", "uuid": uuid, "count": soldCount})
}

func errorResponse(ctx *fiber.Ctx, status int, msg string) error {
	return ctx.Status(status).JSON(struct {
		Error string `json:"message"`
	}{Error: msg})
}
