package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/api_error"
	"net/http"
)

type Handler interface {
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	//Revoke(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) *HandlerStruct {
	return &HandlerStruct{service: service}
}

func validateSignUp(signUp *SignUp) (bool, *api_error.ErrorResponse) {
	if err := api_error.Validate(signUp); err != nil {
		errRes := &api_error.ErrorResponse{
			Status:       "error",
			Code:         fiber.StatusBadRequest,
			Message:      http.StatusText(fiber.StatusBadRequest),
			FailedFields: err,
		}

		return false, errRes
	}

	if signUp.Password != signUp.PasswordConfirm {
		errRes := &api_error.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: http.StatusText(fiber.StatusBadRequest),
			FailedFields: map[string]string{
				"PasswordConfirm": "패스워드와 패스워드 확인 필드가 같지 않습니다.",
			},
		}

		return false, errRes
	}

	return true, nil
}

func (h *HandlerStruct) SignUp(ctx *fiber.Ctx) error {
	signUp := &SignUp{}
	err := ctx.BodyParser(signUp)
	if err != nil {
		errRes := api_error.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: http.StatusText(fiber.StatusBadRequest),
		}
		return errRes.Response(ctx)
	}

	if isValid, errRes := validateSignUp(signUp); !isValid {
		return errRes.Response(ctx)
	}

	result, err := h.service.SignUp(signUp)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"token":      result.Token,
		"expires_in": result.ExpiresAt,
	})
}

func (h *HandlerStruct) SignIn(ctx *fiber.Ctx) error {
	signIn := &SignIn{}
	err := ctx.BodyParser(signIn)
	if err != nil {
		errRes := api_error.ErrorResponse{
			Code:         fiber.StatusBadRequest,
			Message:      http.StatusText(fiber.StatusBadRequest),
			FailedFields: api_error.Validate(err),
		}

		return errRes.Response(ctx)
	}

	result, err := h.service.SignIn(signIn)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"token":      result.Token,
		"expires_in": result.ExpiresAt,
	})
}

//func (h *HandlerStruct) Revoke(ctx *fiber.Ctx) error {
//
//}
