package rest

import (
	"errors"
	"go-dating-app/app/dto"
	"go-dating-app/app/service"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	logger      *slog.Logger
	validator   *validator.Validate
	authService service.AuthService
}

func FormatValidationErrors(errs error) []dto.ValidationErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(errs, &ve) {
		responses := make([]dto.ValidationErrorResponse, len(ve))
		for i, e := range ve {
			responses[i] = dto.ValidationErrorResponse{
				Field:   e.Field(),
				Message: e.Tag(),
			}
		}
		return responses
	}

	return nil
}

func NewHandler(logger *slog.Logger, authService service.AuthService) Handler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return Handler{
		logger:      logger,
		validator:   validate,
		authService: authService,
	}
}
