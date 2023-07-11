package handlers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// validateStruct validate struct by `go-playground/validator` tag
func (h *handler) validateStruct(s interface{}) error {
	var fields []string
	if err := h.validator.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.StructNamespace(), err.Tag(), err.Param())
			fields = append(fields, err.StructNamespace())
		}
		return &errorResponse{
			Code:    fiber.StatusBadRequest,
			Message: fmt.Sprintf("validate error on fields `%s`", strings.Join(fields, ",")),
			Err:     err,
		}
	}
	return nil
}
