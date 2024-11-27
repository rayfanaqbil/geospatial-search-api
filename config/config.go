package config

import (
	"gocroot/helper"

	"github.com/gofiber/fiber/v2"
)

var IPPort, Net = helper.GetAddress()

var Iteung = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "GoCroot",
	AppName:       "Golang Change Root",
	Network:       Net,
}
