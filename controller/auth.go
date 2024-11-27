package controller

import (
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
)

func GetPhoneNumber(c *fiber.Ctx) error {
	var author model.Author
	author.Phone = c.Params("login")
	return c.JSON(author)
}
