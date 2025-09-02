package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
)

// StartPollingByUUID запускает мониторинг для одного станка
// @Summary Запустить мониторинг станка
// @Description Запускает опрос OPC UA для конкретного станка по UUID
// @Tags Polling
// @Produce json
// @Param input body models.UUIDRequest true "UUID станка"
// @Success 200 {object} models.PollingResponseSwagger "Мониторинг запущен"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} NotFoundError "Данные не найдены"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/polling/start [get]
func (h *Handler) StartPollingByUUID(c *gin.Context) {
	var req models.UUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}
	id, err := uuid.Parse(req.UUID)
	if err != nil {
		h.BadRequest(c, fmt.Errorf("incorrect UUID: %s", req.UUID))
		return
	}

	eerr := h.usecase.StartPollingMachine(id)
	if eerr != nil {
		h.ErrorResponse(c, eerr, eerr.Code, eerr.Message, false)
		return
	}

	h.ResultResponse(c, fmt.Sprintf("Polling started for machine %s", id), Object, models.PollingResponse{Polled: true})
}

// StopPollingByUUID останавливает мониторинг для одного станка
// @Summary Остановить мониторинг станка
// @Description Останавливает опрос OPC UA для конкретного станка по UUID
// @Tags Polling
// @Produce json
// @Param input body models.UUIDRequest true "UUID станка"
// @Success 200 {object} models.PollingResponseSwagger "Мониторинг запущен"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} NotFoundError "Данные не найдены"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /api/v1/polling/stop [get]
func (h *Handler) StopPollingByUUID(c *gin.Context) {
	var req models.UUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}
	id, err := uuid.Parse(req.UUID)
	if err != nil {
		h.BadRequest(c, fmt.Errorf("incorrect UUID: %s", req.UUID))
		return
	}

	eerr := h.usecase.StopPollingMachine(id)
	if eerr != nil {
		h.ErrorResponse(c, eerr, eerr.Code, eerr.Message, false)
		return
	}

	h.ResultResponse(c, fmt.Sprintf("Polling stopped for machine %s", id), Object, models.PollingResponse{Polled: false})
}
