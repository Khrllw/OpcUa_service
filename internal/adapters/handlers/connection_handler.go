package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/pkg/errors"
)

// AddConnection подключает клиента для мониторинга ЧПУ
// @Summary Подключение клиента
// @Description Подключение клиента для мониторинга станка: анонимно, по паролю или по сертификату. Передаётся JSON с данными для входа.
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.ConnectionRequest true "Данные для входа"
// @Success 200 {object} models.UUIDResponseSwagger "Успешное подключение"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса или некорректные данные"
// @Failure 401 {object} IncorrectDataError "Некорректные данные"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /connect [post]
func (h *Handler) AddConnection(c *gin.Context) {
	var req models.ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}

	var err *errors.AppError
	var resp models.UUIDResponse

	switch req.ConnectionType {
	case "anonymous":
		resp, err = h.usecase.ConnectAnonymous(req)
	case "password":
		resp, err = h.usecase.ConnectWithPassword(req)
	case "certificate":
		resp, err = h.usecase.ConnectWithCertificate(req)
	default:
		h.BadRequest(c, fmt.Errorf("unknown connection type: %s", req.ConnectionType))
		return
	}

	if err != nil {
		h.logger.Error("Connection failed", "error", err)
		h.ErrorResponse(c, err, err.Code, err.Message, true)
		return
	}
	h.ResultResponse(c, "Successfully connected", Object, resp)
}

// CloseConnection обрабатывает запрос на отключение сессии
// @Summary Отключение сессии
// @Description Закрывает соединение OPC UA по UUID
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.DisconnectRequest true "UUID для отключения"
// @Success 200 {object} models.DisconnectResponseSwagger "Успешное отключение"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} NotFoundError "Данные не найдены"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /connect [delete]
func (h *Handler) CloseConnection(c *gin.Context) {
	var req models.DisconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}
	id, err := uuid.Parse(req.UUID)
	if err != nil {
		h.BadRequest(c, fmt.Errorf("incorrect UUID: %s", req.UUID))
	}

	// Вызываем usecase для закрытия соединения
	state, eerr := h.usecase.DisconnectByUUID(id)
	if eerr != nil {
		if state == nil || *state == false {
			h.ErrorResponse(c, err, eerr.Code, eerr.Message, true)
			return
		} else {
			h.ResultResponse(c, "Disconnected with database record delete error", Object, models.DisconnectResponse{Disconnected: true})
		}
	}

	h.ResultResponse(c, "Successfully disconnected", Object, models.DisconnectResponse{Disconnected: true})
}

// CheckConnection проверяет здоровье соединения по UUID
// @Summary Проверка соединения
// @Description Проверяет состояние соединения OPC UA по UUID
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.CheckConnectionRequest true "UUID для проверки"
// @Success 200 {object} models.CheckConnectionResponseSwagger "Информация о подключении"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} NotFoundError "Данные не найдены"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /connect/check [post]
func (h *Handler) CheckConnection(c *gin.Context) {
	var req models.CheckConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}
	id, err := uuid.Parse(req.UUID)
	if err != nil {
		h.BadRequest(c, fmt.Errorf("incorrect UUID: %s", req.UUID))
	}

	// Получаем информацию о соединении
	connInfo, eerr := h.usecase.GetConnectionState(id)
	if eerr != nil {
		h.ErrorResponse(c, eerr, eerr.Code, eerr.Message, false)
		return
	}

	h.ResultResponse(c, "Successfully get connection info", Object, connInfo)
}

// GetConnectionPool возвращает текущий пул открытых соединений
// @Summary Получить пул соединений
// @Description Возвращает список активных соединений в пуле OPC UA
// @Tags Connection
// @Produce json
// @Success 200 {object} models.GetConnectionPoolResponseSwagger "Список активных соединений"
// @Router /connect [get]
func (h *Handler) GetConnectionPool(c *gin.Context) {
	resp := h.usecase.GetActiveConnections()
	h.ResultResponse(c, "Successfully get connection pool", Object, resp)
}
