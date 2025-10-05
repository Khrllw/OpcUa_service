package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"opc_ua_service/internal/domain/models/connection"
)

// StartPollingByID запускает мониторинг для одного станка
// @Summary Запустить мониторинг станка
// @Description Запускает опрос OPC UA для конкретного станка по UUID
// @Tags Polling
// @Produce json
// @Param input body models.IDRequest true "ID станка"
// @Success 200 {object} swagger.PollingResponse "Мониторинг запущен"
// @Failure 400 {object} swagger.IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} swagger.NotFoundError "Данные не найдены"
// @Failure 500 {object} swagger.InternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/polling/start [get]
func (h *Handler) StartPollingByID(c *gin.Context) {
	var req connection.StartPollingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}

	eerr := h.usecase.StartPollingMachine(req)
	if eerr != nil {
		h.ErrorResponse(c, eerr, eerr.Code, eerr.Message, false)
		return
	}

	h.ResultResponse(c, fmt.Sprintf("Polling started for machine %d", req.ID), Object, connection.PollingResponse{Polled: true})
}

// StopPollingByID останавливает мониторинг для одного станка
// @Summary Остановить мониторинг станка
// @Description Останавливает опрос OPC UA для конкретного станка по UUID
// @Tags Polling
// @Produce json
// @Param input body models.IDRequest true "UUID станка"
// @Success 200 {object} swagger.PollingResponse "Мониторинг запущен"
// @Failure 400 {object} swagger.IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} swagger.NotFoundError "Данные не найдены"
// @Failure 500 {object} swagger.InternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/polling/stop [get]
func (h *Handler) StopPollingByID(c *gin.Context) {
	var req connection.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}

	eerr := h.usecase.StopPollingMachine(req.ID)
	if eerr != nil {
		h.ErrorResponse(c, eerr, eerr.Code, eerr.Message, false)
		return
	}

	h.ResultResponse(c, fmt.Sprintf("Polling stopped for machine %d", req.ID), Object, connection.PollingResponse{Polled: false})
}
