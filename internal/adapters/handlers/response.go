package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"opc_ua_service/pkg/errors"
)

// Константы типов ответа
const (
	Object = "object" // Используется когда ответ содержит один объект
	Array  = "array"  // Используется когда ответ содержит массив объектов
	Empty  = "empty"  // Используется когда ответ не содержит данных
)

// Response - интерфейс для стандартных ответов API
type Response interface {
	ErrorResponse(c *gin.Context, err error, statusCode int, message string, showError bool)
	ResultResponse(c *gin.Context, message string, dataType string, data interface{})
	BadRequest(c *gin.Context, err error)
}

// ErrorResponse - возвращает стандартизированный ответ с ошибкой
func (h *Handler) ErrorResponse(c *gin.Context, err error, statusCode int, message string, showError bool) {
	errorMessage := message
	if showError && err != nil {
		errorMessage = message + ": " + err.Error()
	}

	c.JSON(statusCode, gin.H{
		"status": "error",
		"error": gin.H{
			"code":    statusCode,
			"message": errorMessage,
		},
	})
}

// ResultResponse - возвращает JSON с данными
func (h *Handler) ResultResponse(c *gin.Context, message string, dataType string, data interface{}) {
	response := gin.H{
		"status":  "success",
		"message": message,
		"type":    dataType,
	}

	if data != nil {
		response["data"] = data
	}

	c.JSON(http.StatusOK, response)
}

// BadRequest - возвращает ошибку 400
func (h *Handler) BadRequest(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
}
