package routes

import (
	"github.com/dw-parameter-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initVoucherRoutes(router fiber.Router) {
	r := router.Group("/voucher")
	h := handlers.NewVoucherHandler()

	r.Post("/all", func(c *fiber.Ctx) error {
		return h.GetAllVoucher(c)
	})

	r.Get("/:id", func(c *fiber.Ctx) error {
		return h.GetVoucherByID(c)
	})

	r.Post("/code", func(c *fiber.Ctx) error {
		return h.GetVoucherByCode(c)
	})

	r.Post("/", func(c *fiber.Ctx) error {
		return h.Create(c)
	})

	r.Delete("/", func(c *fiber.Ctx) error {
		return h.Deactivate(c)
	})

	r.Put("/", func(c *fiber.Ctx) error {
		return h.Update(c)
	})

}
