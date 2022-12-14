package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/logger"
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
	"strconv"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	logger.HasLogger
}

type HandlerStruct struct {
	service Service
	logger.HasLoggerStruct
}

func NewHandler(s Service) Handler {
	return &HandlerStruct{
		service: s,
		HasLoggerStruct: logger.HasLoggerStruct{
			Logger: s.GetLogger(),
		},
	}
}

// Create
// @Summary create search
// @description create search
// @Tags Search
// @param request body CreateSearch true "create search"
// @Success 201 {object} Response
// @Failure 400 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search [post]
// @Security BearerAuth
func (h *HandlerStruct) Create(c *fiber.Ctx) error {
	dto := &CreateSearch{}

	err := c.BodyParser(dto)
	if err != nil {
		res := api_error.NewValidationError(c)
		return res.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	created, err := h.service.Create(dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// All
// @Summary create search
// @description create search
// @Tags Search
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {object} ResponseAll
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/all [get]
// @Security BearerAuth
func (h *HandlerStruct) All(c *fiber.Ctx) error {
	page, err := utils.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	all, err := h.service.All(page)

	if err != nil {
		return err
	}

	return c.JSON(ResponseAll{
		Paginator: all,
		Data:      all.Data.([]entity.Search),
	})
}

// Get
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @Success 200 {object} Response
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [get]
// @Security BearerAuth
func (h *HandlerStruct) Get(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}
	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	search, err := h.service.Find(uint(pk), user.Id)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(search)
}

// Update
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @param request body UpdateSearch true "update search"
// @Success 200 {object} Response
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [put]
// @Security BearerAuth
func (h *HandlerStruct) Update(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}
	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	dto := &UpdateSearch{}
	err = c.BodyParser(dto)
	if err != nil {
		errRes := api_error.NewValidationError(c)
		return errRes.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	exists, err := h.service.Find(uint(pk), user.Id)
	if exists == nil || exists.Id == 0 {
		return fiber.ErrNotFound
	}

	updated, err := h.service.Update(uint(pk), user.Id, dto)
	if err != nil {
		return err
	}

	return c.JSON(updated)
}

// Patch
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @param request body PatchSearch true "update search"
// @Success 200 {object} Response
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [patch]
// @Security BearerAuth
func (h *HandlerStruct) Patch(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}
	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	dto := &PatchSearch{}
	err = c.BodyParser(dto)
	if err != nil {
		errRes := api_error.NewValidationError(c)
		return errRes.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	exists, err := h.service.Find(uint(pk), user.Id)
	if exists == nil || exists.Id == 0 {
		return fiber.ErrNotFound
	}

	patch, err := h.service.Patch(uint(pk), user.Id, dto)
	if err != nil {
		return err
	}

	return c.JSON(patch)
}

// Delete
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @Success 204 {object} utils.StatusResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [delete]
// @Security BearerAuth
func (h *HandlerStruct) Delete(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}
	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	rs, err := h.service.Delete(uint(pk), user.Id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.StatusResponse{
		Status: rs,
	})
}
