package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/salahfarzin/roja-shop/configs"
	"github.com/salahfarzin/roja-shop/services"
	"github.com/salahfarzin/roja-shop/types"
	// ...existing code...
)

type Product interface {
	Index(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Sell(ctx *fiber.Ctx) error
}

type product struct {
	service     services.Product
	saleService services.Sale
}

func NewProduct(service services.Product, saleService services.Sale) Product {
	return &product{service: service, saleService: saleService}
}

func (p *product) Index(ctx *fiber.Ctx) error {
	// Pagination: ?page=1
	page := 1
	limit := 6
	if pStr := ctx.Query("page"); pStr != "" {
		if pInt, err := strconv.Atoi(pStr); err == nil && pInt > 0 {
			page = pInt
		}
	}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	offset := (page - 1) * limit
	products, err := p.service.GetAll(limit, offset)
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

func (p *product) Create(ctx *fiber.Ctx) error {
	return p.handleUpsert(ctx, "", true)
}

func (p *product) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, fiber.StatusBadRequest, "missing product ID")
	}

	return p.handleUpsert(ctx, id, false)
}

func (p *product) handleUpsert(ctx *fiber.Ctx, id string, isCreate bool) error {
	// Parse all fields from form values (multipart/form-data)
	title := ctx.FormValue("title")
	brand := ctx.FormValue("brand")
	description := ctx.FormValue("description")
	priceStr := ctx.FormValue("price")
	oldPriceStr := ctx.FormValue("old_price")
	discountStr := ctx.FormValue("discount")
	inventoryStr := ctx.FormValue("inventory")
	var detailsPtr *map[string]string
	var styleNotesPtr *map[string]string
	details := make(map[string]string)
	styleNotes := make(map[string]string)

	// Optionally parse details and styleNotes as JSON strings if provided
	detailsStr := ctx.FormValue("details")
	if detailsStr != "" {
		_ = json.Unmarshal([]byte(detailsStr), &details)
		detailsPtr = &details
	}
	styleNotesStr := ctx.FormValue("style_notes")
	if styleNotesStr != "" {
		_ = json.Unmarshal([]byte(styleNotesStr), &styleNotes)
		styleNotesPtr = &styleNotes
	}

	price, _ := strconv.ParseFloat(priceStr, 64)
	oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)
	discount, _ := strconv.ParseFloat(discountStr, 64)
	inventory, _ := strconv.Atoi(inventoryStr)

	var oldPricePtr *float64
	if oldPriceStr != "" {
		oldPricePtr = &oldPrice
	}
	var discountPtr *float64
	if discountStr != "" {
		discountPtr = &discount
	}

	if title == "" || price == 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "title and price are required")
	}

	fileHeader, err := ctx.FormFile("image")
	var file *types.File
	if isCreate {
		if err != nil || fileHeader == nil {
			return errorResponse(ctx, fiber.StatusBadRequest, "image is required for product creation")
		}
		file, _, err = services.UploadService.SaveFile(ctx, fileHeader, configs.New())
		if err != nil && err != fiber.ErrUnprocessableEntity {
			if errors.Is(err, services.ErrNotImage) {
				return errorResponse(ctx, fiber.StatusBadRequest, err.Error())
			}
			return errorResponse(ctx, fiber.StatusInternalServerError, "failed to save image")
		}
	} else if fileHeader != nil {
		file, _, err = services.UploadService.SaveFile(ctx, fileHeader, configs.New())
		if err != nil && err != fiber.ErrUnprocessableEntity {
			if errors.Is(err, services.ErrNotImage) {
				return errorResponse(ctx, fiber.StatusBadRequest, err.Error())
			}
			return errorResponse(ctx, fiber.StatusInternalServerError, "failed to save image")
		}
	}

	// If updating and no new file is uploaded, keep the old file
	if !isCreate && file == nil {
		prod, err := p.service.GetOne(id)
		if err == nil && prod != nil {
			file = prod.File
		}
	}

	productID := id
	if isCreate {
		productID = uuid.New().String()
	}

	product := types.Product{
		ID:          productID,
		Brand:       brand,
		Title:       title,
		Inventory:   inventory,
		Price:       price,
		OldPrice:    oldPricePtr,
		Discount:    discountPtr,
		Description: description,
		Details:     detailsPtr,
		StyleNotes:  styleNotesPtr,
		File:        file,
	}

	if isCreate {
		_, err = p.service.Create(product, file)
	} else {
		err = p.service.Update(id, product)
	}

	action := "update"
	if isCreate {
		action = "create"
	}
	if err != nil {
		return errorResponse(ctx, fiber.StatusInternalServerError, fmt.Sprintf("failed to %s product", action))
	}

	product.File.Path = ctx.BaseURL() + product.File.Path

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product " + action + "d",
		"product": product,
	})
}

func (p *product) Sell(ctx *fiber.Ctx) error {
	var body struct {
		ProductID string `json:"uuid"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}
	if body.ProductID == "" {
		return errorResponse(ctx, fiber.StatusBadRequest, "product_id is required")
	}

	quantity := 1
	product, err := p.saleService.Sell(body.ProductID, quantity)
	if err != nil {
		if err == services.ErrInsufficientInventory {
			return errorResponse(ctx, fiber.StatusBadRequest, "insufficient inventory")
		}
		return errorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "sale recorded", "product": product})
}

func errorResponse(ctx *fiber.Ctx, status int, msg string) error {
	return ctx.Status(status).JSON(struct {
		Error string `json:"message"`
	}{Error: msg})
}
