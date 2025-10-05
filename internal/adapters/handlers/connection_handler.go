package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models/connection"
	connection_models "opc_ua_service/internal/domain/models/connection_types"
	"opc_ua_service/pkg/errors"
)

// AddConnection подключает клиента для мониторинга ЧПУ
// @Summary Подключение клиента
// @Description Подключение клиента для мониторинга станка: анонимно, по паролю или по сертификату. Передаётся JSON с данными для входа.
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.ConnectionRequest true "Данные для входа"
// @Success 200 {object} swagger.IDResponse "Успешное подключение"
// @Failure 400 {object} swagger.IncorrectFormatError "Неверный формат запроса или некорректные данные"
// @Failure 401 {object} swagger.IncorrectDataError "Некорректные данные"
// @Failure 500 {object} swagger.InternalServerError "Внутренняя ошибка сервера"
// @Router /connect [post]
func (h *Handler) AddConnection(c *gin.Context) {
	var req connection.ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}

	var err *errors.AppError
	resp := entities.CncMachine{}

	switch req.ConnectionType {
	case connection_models.ConnectionAnonymous:
		resp, err = h.usecase.ConnectAnonymous(req)
	case connection_models.ConnectionPassword:
		resp, err = h.usecase.ConnectWithPassword(req)
	case connection_models.ConnectionCertificate:
		resp, err = h.usecase.ConnectWithCertificate(req)
	default:
		h.BadRequest(c, fmt.Errorf("unknown connection type: %s", req.ConnectionType))
		return
	}

	if err != nil {
		h.logger.Error("Connection failed", "error", err)
		h.ErrorResponse(c, err, err.Code, err.Message, false)
		return
	}
	h.ResultResponse(c, "Successfully connected", Object, resp)
}

// CloseConnection обрабатывает запрос на отключение сессии
// @Summary Отключение сессии
// @Description Закрывает соединение OPC UA по ID
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.IDRequest true "ID для отключения"
// @Success 200 {object} swagger.DisconnectResponse "Успешное отключение"
// @Failure 400 {object} swagger.IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} swagger.NotFoundError "Данные не найдены"
// @Failure 500 {object} swagger.InternalServerError "Внутренняя ошибка сервера"
// @Router /connect [delete]
func (h *Handler) CloseConnection(c *gin.Context) {
	var req connection.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}

	state, eerr := h.usecase.DisconnectByID(req.ID)
	if eerr != nil {
		if state == nil || *state == false {
			h.ErrorResponse(c, eerr, eerr.Code, eerr.Message, false)
			return
		} else {
			h.ResultResponse(c, "Disconnected with database record delete error", Object, connection.DisconnectResponse{Disconnected: true})
		}
	}

	h.ResultResponse(c, "Successfully disconnected", Object, connection.DisconnectResponse{Disconnected: true})
}

// CheckConnection проверяет здоровье соединения по ID
// @Summary Проверка соединения
// @Description Проверяет состояние соединения OPC UA по ID
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.CheckConnectionRequest true "ID для проверки"
// @Success 200 {object} swagger.CheckConnectionResponse "Информация о подключении"
// @Failure 400 {object} swagger.IncorrectFormatError "Неверный формат запроса"
// @Failure 404 {object} swagger.NotFoundError "Данные не найдены"
// @Failure 500 {object} swagger.InternalServerError "Внутренняя ошибка сервера"
// @Router /connect/check [post]
func (h *Handler) CheckConnection(c *gin.Context) {
	var req connection.CheckConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err)
		return
	}

	connInfo, eerr := h.usecase.GetConnectionState(req.ID)
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
// @Success 200 {object} swagger.GetConnectionPoolResponse "Список активных соединений"
// @Router /connect [get]
func (h *Handler) GetConnectionPool(c *gin.Context) {
	resp := h.usecase.GetActiveConnections()
	h.ResultResponse(c, "Successfully get connection pool", Object, resp)
}
