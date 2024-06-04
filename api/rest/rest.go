package rest

import (
	"errors"
	"go-dating-app/app/dto"
	"go-dating-app/app/service"
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTrans "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	logger      *slog.Logger
	validator   *validator.Validate
	authService service.AuthService
}

func (h *Handler) FormatValidationErrors(errs error) []dto.ValidationErrorResponse {
	// Validation translation.
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTrans.RegisterDefaultTranslations(h.validator, trans)

	// Transform validation errors.
	var ve validator.ValidationErrors
	if errors.As(errs, &ve) {
		responses := make([]dto.ValidationErrorResponse, len(ve))
		for i, e := range ve {
			responses[i] = dto.ValidationErrorResponse{
				Field:   e.Field(),
				Message: e.Translate(trans),
			}
		}
		return responses
	}

	return nil
}

func NewHandler(logger *slog.Logger, authService service.AuthService) Handler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// RegisterTagNameFunc is used to customize the name of the field with tag taken from json instead struct property name.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return Handler{
		logger:      logger,
		validator:   validate,
		authService: authService,
	}
}
