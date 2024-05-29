package order

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/synthao/orders/internal/domain"
	"github.com/synthao/orders/internal/middleware/sso"
	sso2 "github.com/synthao/orders/internal/module/sso"
	"net/http"
	"strconv"
	"time"
)

type CreateRequest struct {
	Sum float64 `json:"sum"`
}

type UpdateStatusRequest struct {
	Status int `json:"status"`
}

type GetListResponse struct {
	ID        int       `json:"id"`
	Status    int       `json:"status"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"created_at"`
}

type GetOneResponse struct {
	ID        int       `json:"id"`
	Status    int       `json:"status"`
	Sum       float64   `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	app     *fiber.App
	service domain.Service
}

func NewHandler(app *fiber.App, service domain.Service) *Handler {
	return &Handler{app: app, service: service}
}

func (h *Handler) InitRoutes(ssoClient *sso2.Client) {
	group := h.app.Group("/api/orders")

	group.Use(sso.New(sso.Config{Client: ssoClient}))

	h.app.Put("/api/orders/:id<int>/status", h.updateStatus)
	h.app.Delete("/api/orders/:id<int>", h.delete)
	h.app.Post("/api/orders", h.create)
	h.app.Get("/api/orders", h.list)
	h.app.Get("/api/orders/:id<int>", h.one)
}

func (h *Handler) create(ctx *fiber.Ctx) error {
	var req CreateRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": "Failed to create a record. Payload parsing error"})
	}

	id, err := h.service.Create(domain.NewOrder(req.Sum))
	if err != nil {
		return ctx.JSON(fiber.Map{"error": "Failed to create a record"})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	data, err := h.service.GetList(10, 0)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(map[string]any{
			"error": "Something went wrong while fetching data",
		})
	}

	return ctx.JSON(fromDomainToGetListResponse(data))
}

func (h *Handler) one(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	data, err := h.service.GetOne(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.SendStatus(http.StatusNotFound)
		}

		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while fetching record",
		})
	}

	return ctx.JSON(fromDomainToGetOneResponse(data))
}

func (h *Handler) updateStatus(ctx *fiber.Ctx) error {
	var req UpdateStatusRequest

	orderIDParam := ctx.Params("id")

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": "Failed to parse order id"})
	}

	err = ctx.BodyParser(&req)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": "Failed to update status. Payload parsing error"})
	}

	err = h.service.UpdateStatus(orderID, req.Status)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": "Failed to update status"})
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (h *Handler) delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return err
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.SendStatus(http.StatusNotFound)
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func fromDomainToGetOneResponse(data *domain.Order) *GetOneResponse {
	return &GetOneResponse{
		ID:        data.ID,
		Sum:       data.Sum,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
	}
}

func fromDomainToGetListResponse(data []domain.Order) []GetListResponse {
	res := make([]GetListResponse, len(data))

	for i, item := range data {
		res[i] = GetListResponse{
			ID:     item.ID,
			Sum:    item.Sum,
			Status: item.Status,
		}
	}

	return res
}
