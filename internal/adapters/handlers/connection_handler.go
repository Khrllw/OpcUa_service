package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opc_ua_service/internal/domain/models"
)

// AddConnection аутентифицирует клиента для мониторинга ЧПУ
// @Summary Аутентификация клиента
// @Description Аутентификация клиента для мониторинга станка: анонимно, по паролю или по сертификату. Передаётся JSON с данными для входа.
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body models.ConnectionRequest true "Данные для входа"
// @Success 200 {object} models.ConnectionAuthResponse "Успешная аутентификация"
// @Failure 400 {object} gin.H "Неверный формат запроса или неизвестный тип аутентификации"
// @Failure 401 {object} gin.H "Неверные учетные данные"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /connect [post]
func (h *Handler) AddConnection(c *gin.Context) {
	var req models.ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Error decoding request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	h.logger.Info("Auth attempt", "authType", req.ConnectionType, "username", req.Username)

	var token string
	var err error

	connInfo := &models.ConnectionInfo{}
	switch req.ConnectionType {
	case "anonymous":
		token, err = h.usecase.LoginClientAnonymous()
	case "password":
		token, err = h.usecase.LoginClientPassword(req)
	case "certificate":
		connInfo, err = h.usecase.ConnectByCert(req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown auth type"})
		return
	}

	if err != nil {
		h.logger.Error("Auth failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	inf := models.ToResponse(connInfo)

	c.JSON(http.StatusOK, models.ConnectionAuthResponse{
		Status:         models.StatusOK,
		Token:          token,
		ConnectionInfo: &inf,
	})
}

// GetConnectionPool возвращает текущий пул открытых соединений
// @Summary Получить пул соединений
// @Description Возвращает список активных соединений в пуле OPC UA
// @Tags Connection
// @Produce json
// @Success 200 {object} gin.H "Список активных соединений"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /connect [get]
func (h *Handler) GetConnectionPool(c *gin.Context) {
	// Получаем все активные соединения из usecase/pool
	pool := h.usecase.GetActiveConnections() // возвращает []*models.ConnectionInfo

	// Преобразуем все в респонс
	var resp []models.ConnectionInfoResponse
	for _, connInfo := range pool {
		if connInfo != nil {
			resp = append(resp, models.ToResponse(connInfo))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"poolSize":    len(resp),
		"connections": resp, // массив объектов ConnectionInfoResponse
	})
}

// CloseConnection обрабатывает DELETE запрос на отключение сессии по JSON
// @Summary Отключение сессии
// @Description Закрывает соединение OPC UA по sessionID. sessionID передаётся в теле запроса в формате JSON.
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body DisconnectRequest true "Session ID для отключения"
// @Success 200 {object} gin.H "Сообщение об успешном отключении"
// @Failure 400 {object} gin.H "Отсутствует или некорректный sessionID"
// @Failure 500 {object} gin.H "Не удалось отключить сессию"
// @Router /connections [delete]
func (h *Handler) CloseConnection(c *gin.Context) {
	var req models.DisconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid sessionID"})
		return
	}

	// Вызываем usecase для закрытия соединения
	err := h.usecase.DisconnectBySessionID(req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect session: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session " + req.SessionID + " disconnected successfully"})
}

// CheckConnectionHandler проверяет здоровье соединения по sessionID
// @Summary Проверка соединения
// @Description Проверяет состояние соединения OPC UA по sessionID
// @Tags Connection
// @Accept json
// @Produce json
// @Param input body CheckConnectionRequest true "Session ID для проверки"
// @Success 200 {object} gin.H "Статус соединения"
// @Failure 400 {object} gin.H "Отсутствует или некорректный sessionID"
// @Failure 404 {object} gin.H "Соединение не найдено"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /connect/check [post]
func (h *Handler) CheckConnection(c *gin.Context) {
	var req models.CheckConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid sessionID"})
		return
	}

	// Получаем информацию о соединении
	connInfo, err := h.usecase.GetConnectionStats(req.SessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found: " + err.Error()})
		return
	}

	// Проверяем состояние
	status := "healthy"
	if isHealthy, ok := connInfo["is_healthy"].(bool); ok && !isHealthy {
		status = "unhealthy"
	}

	c.JSON(http.StatusOK, gin.H{
		"sessionID": req.SessionID,
		"status":    status,
		"details":   connInfo,
	})
}
