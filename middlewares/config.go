package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salahfarzin/roja-shop/configs"
)

func Config(configs configs.Configs) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("configs", configs)

		return ctx.Next()
	}
}
