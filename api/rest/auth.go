package rest

import (
	"errors"
	"go-dating-app/app/dto"
	"go-dating-app/app/entity"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Registration(c echo.Context) error {
	var req dto.AuthRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: http.StatusText(http.StatusBadRequest),
		})
	}

	// Validate request.
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "please check your input",
			Errors:  h.FormatValidationErrors(err),
		})
	}

	// Create user.
	_, err := h.authService.Registration(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExists):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "email is already registered",
			})
		case errors.Is(err, entity.ErrUserInvalidEmail):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "invalid email address",
			})
		case errors.Is(err, entity.ErrUserInvalidPassword):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "password must be at least 6 characters",
			})
		default:
			h.logger.Error("auth.registration", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Message: http.StatusText(http.StatusInternalServerError),
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Registration success",
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req dto.AuthLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: http.StatusText(http.StatusBadRequest),
		})
	}

	// Validate request.
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "please check your input",
			Errors:  h.FormatValidationErrors(err),
		})
	}

	// Login.
	tokens, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			fallthrough
		case errors.Is(err, entity.ErrUserPasswordIncorrect):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "account not found or password invalid",
			})
		default:
			h.logger.Error("auth.login", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Message: http.StatusText(http.StatusInternalServerError),
			})
		}
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Login success",
		Data: dto.AuthLoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
