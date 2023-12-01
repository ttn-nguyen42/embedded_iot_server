package privateapi

import (
	"github.com/gofiber/fiber/v2"
)


func GlobalErrorHandler() func(c *fiber.Ctx, err error) error {
	return nil
}
