package v1

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (r *V1) orderJSON(ctx *fiber.Ctx) error {
	orderUID := ctx.Query("order_uid")

	if orderUID == "" {
		return errorResponse(ctx, http.StatusBadRequest, "order_uid required")
	}

	order, err := r.o.Order(ctx.UserContext(), orderUID)
	if err != nil {
		r.l.Error(err, "http - v1 - orderJSON")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problenms") // TODO: норм ошибку, проверка на sql.ErrNoRows
	}

	return ctx.Status(http.StatusOK).JSON(order)
}

func (r *V1) orderHTML(ctx *fiber.Ctx) error {
	orderUID := ctx.Query("order_uid")

	// path to html files
	tmplPath := filepath.Join("docs", "html")

	// show `order search`
	if orderUID == "" {
		t, err := template.ParseFiles(filepath.Join(tmplPath, "order_form.html"))
		if err != nil {
			return errorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
		ctx.Type("html")
		return t.Execute(ctx.Response().BodyWriter(), nil)
	}

	// get order
	order, err := r.o.Order(ctx.UserContext(), orderUID)
	if err != nil {
		r.l.Error(err, "http - v1 - order")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems") // TODO: норм ошибку, проверка на sql.ErrNoRows
	}

	pretty, _ := json.MarshalIndent(order, "", "  ")

	// show `order info`
	t, err := template.ParseFiles(filepath.Join(tmplPath, "order_info.html"))
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	data := map[string]interface{}{
		"Order":      order,
		"PrettyJSON": string(pretty),
	}
	ctx.Type("html")
	return t.Execute(ctx.Response().BodyWriter(), data)
}
