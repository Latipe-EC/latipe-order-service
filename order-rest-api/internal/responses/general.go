package responses

import (
	"github.com/gofiber/fiber/v2"
)

type General struct {
	Status    int         `json:"-"`
	Code      int         `json:"code"`
	ErrorCode string      `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

func (g *General) JSON(c *fiber.Ctx) error {
	return c.Status(g.Status).JSON(g)
}
