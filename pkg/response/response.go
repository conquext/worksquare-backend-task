package response

import (
	"housing-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *fiber.Ctx, message string, err error) error {
	errorInfo := &models.ErrorInfo{
		Code:    fiber.StatusBadRequest,
		Message: message,
	}
	if err != nil {
		errorInfo.Details = err.Error()
	}

	return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
		Success: false,
		Error:   errorInfo,
	})
}

func Unauthorized(c *fiber.Ctx, message string, err error) error {
	errorInfo := &models.ErrorInfo{
		Code:    fiber.StatusUnauthorized,
		Message: message,
	}
	if err != nil {
		errorInfo.Details = err.Error()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
		Success: false,
		Error:   errorInfo,
	})
}

func NotFound(c *fiber.Ctx, message string, err error) error {
	errorInfo := &models.ErrorInfo{
		Code:    fiber.StatusNotFound,
		Message: message,
	}
	if err != nil {
		errorInfo.Details = err.Error()
	}

	return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
		Success: false,
		Error:   errorInfo,
	})
}

func Conflict(c *fiber.Ctx, message string, err error) error {
	errorInfo := &models.ErrorInfo{
		Code:    fiber.StatusConflict,
		Message: message,
	}
	if err != nil {
		errorInfo.Details = err.Error()
	}

	return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
		Success: false,
		Error:   errorInfo,
	})
}

func InternalServerError(c *fiber.Ctx, message string, err error) error {
	errorInfo := &models.ErrorInfo{
		Code:    fiber.StatusInternalServerError,
		Message: message,
	}
	if err != nil {
		errorInfo.Details = err.Error()
	}

	return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
		Success: false,
		Error:   errorInfo,
	})
}

func ValidationError(c *fiber.Ctx, message string, errors []models.ValidationError) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
		Success: false,
		Error: &models.ErrorInfo{
			Code:    fiber.StatusUnprocessableEntity,
			Message: message,
		},
		Data: models.ValidationErrorResponse{
			Message: message,
			Errors:  errors,
		},
	})
}
