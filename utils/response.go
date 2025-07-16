package utils

import "github.com/gofiber/fiber/v2"

// JSONResponse standardizes successful API JSON responses.
// It takes a Fiber context, HTTP status code, a message string, and data (can be nil).
func JSONResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

// ErrorResponse standardizes API error responses.
// It takes a Fiber context, HTTP status code, and an error message string.
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}