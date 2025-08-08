package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *V1) order(ctx *fiber.Ctx) error {
	orderUID := ctx.Query("order_uid")

	order, err := r.o.Order(ctx.UserContext(), orderUID)
	if err != nil {
		r.l.Error(err, "http - v1 - order")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	return ctx.Status(http.StatusOK).JSON(order)
}
