package utils

import "github.com/gofiber/fiber/v2"

func Success(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(200).JSON(fiber.Map{
		"status": 200,
		"error":  nil,
		"data":   data,
	})
}

func Error(ctx *fiber.Ctx, status int, err string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"status": status,
		"error":  err,
		"data":   nil,
	})
}
