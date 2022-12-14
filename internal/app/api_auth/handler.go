package api_auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/internal/core/context"
	"github.com/miniyus/go-fiber/internal/core/logger"
	"github.com/miniyus/go-fiber/internal/utils"
)

type Handler interface {
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	Me(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
	RevokeToken(ctx *fiber.Ctx) error
	logger.HasLogger
}

type HandlerStruct struct {
	service Service
	logger.HasLoggerStruct
}

func NewHandler(service Service) Handler {
	return &HandlerStruct{
		service:         service,
		HasLoggerStruct: logger.HasLoggerStruct{Logger: service.GetLogger()},
	}
}

func validateSignUp(ctx *fiber.Ctx, signUp *SignUp) (bool, *api_error.ErrorResponse) {
	if err := utils.Validate(signUp); err != nil {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = err

		return false, &errRes
	}

	if signUp.Password != signUp.PasswordConfirm {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = map[string]string{
			"PasswordConfirm": "패스워드와 패스워드확인 필드가 일치하지않습니다.",
		}

		return false, &errRes
	}

	return true, nil
}

// SignUp
// @Summary Sign up
// @Description sign up
// @Tags Auth
// @Success 201 {object} SignUpResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Param request body SignUp true "sign up body"
// @Router /api/auth/register [post]
func (h *HandlerStruct) SignUp(ctx *fiber.Ctx) error {
	signUp := &SignUp{}
	err := ctx.BodyParser(signUp)
	if err != nil {
		errRes := api_error.NewValidationError(ctx)
		return errRes.Response()
	}

	if isValid, errRes := validateSignUp(ctx, signUp); !isValid {
		return errRes.Response()
	}

	result, err := h.service.SignUp(signUp)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func validateSignIn(ctx *fiber.Ctx, in *SignIn) (bool, *api_error.ErrorResponse) {
	if err := utils.Validate(in); err != nil {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = err

		return false, &errRes
	}

	return true, nil
}

// SignIn
// @Summary login
// @Description login
// @Tags Auth
// @Success 200 {object} TokenInfo
// @Failure 400 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Param request body SignIn true "login  body"
// @Router /api/auth/token [post]
func (h *HandlerStruct) SignIn(ctx *fiber.Ctx) error {
	signIn := &SignIn{}
	err := ctx.BodyParser(signIn)
	if err != nil {
		errRes := api_error.NewValidationError(ctx)

		return errRes.Response()
	}

	if isValid, errRes := validateSignIn(ctx, signIn); !isValid {
		return errRes.Response()
	}

	result, err := h.service.SignIn(signIn)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(TokenInfo{
		Token:     result.Token,
		ExpiresAt: utils.JsonTime(result.ExpiresAt),
	})
}

// Me
// @Summary get my info
// @description get login user info
// @Tags Auth
// @Success 200 {object} auth.User
// @Accept json
// @Produce json
// @Router /api/auth/me [get]
// @Security BearerAuth
func (h *HandlerStruct) Me(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(context.AuthUser).(*auth.User)
	if !ok {
		h.GetLogger().Error(user)
		return fiber.NewError(500, "Can't Load Context AuthUser")
	}

	return ctx.JSON(user)
}

// ResetPassword
// @Summary reset password
// @description reset login user's password
// @Tags Auth
// @Param request body ResetPasswordStruct true "reset password body"
// @Success 200 {object} auth.User
// @Failure 400 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/auth/password [patch]
// @Security BearerAuth
func (h *HandlerStruct) ResetPassword(ctx *fiber.Ctx) error {
	user, err := utils.GetAuthUser(ctx)
	if err != nil {
		errRes := api_error.NewFromError(ctx, err)
		return errRes.Response()
	}

	dto := &ResetPasswordStruct{}

	err = ctx.BodyParser(dto)
	if err != nil {
		errRes := api_error.NewValidationError(ctx)
		return errRes.Response()
	}

	failedFields := utils.Validate(dto)
	if failedFields != nil {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = failedFields
		return errRes.Response()
	}

	if dto.Password != dto.PasswordConfirm {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = map[string]string{
			"PasswordConfirm": "패스워드와 패스워드확인 필드가 일치하지않습니다.",
		}
		return errRes.Response()
	}

	rs, err := h.service.ResetPassword(user.Id, dto)
	if err != nil {
		errRes := api_error.NewFromError(ctx, err)
		return errRes.Response()
	}

	return ctx.JSON(rs)
}

// RevokeToken
// @Summary revoke token
// @description revoke current jwt token
// @Tags Auth
// @Success 200 {object} utils.StatusResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/auth/revoke [delete]
// @Security BearerAuth
func (h *HandlerStruct) RevokeToken(ctx *fiber.Ctx) error {
	user, err := utils.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	tokenInfo, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrForbidden
	}

	token := tokenInfo.Raw

	rs, err := h.service.RevokeToken(user.Id, token)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(utils.StatusResponse{
		Status: rs,
	})
}
