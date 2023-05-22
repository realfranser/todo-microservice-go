package rest

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func renderErrorResponse(c *fiber.Ctx, msg string, status int) error {
	return renderResponse(c, ErrorResponse{Error: msg}, status)
}

func renderResponse(c *fiber.Ctx, res interface{}, status int) error {

	if err := c.Status(status).JSON(&fiber.Map{
		"success": true,
		"body":    res,
	}); err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return nil
}
