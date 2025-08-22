package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// StartMonitoring запускает воркер проверки здоровья соединений
// @Summary Запустить мониторинг пула соединений
// @Description Запускает периодическую проверку состояния всех соединений в пуле OPC UA
// @Tags Connection
// @Produce json
// @Success 200 {object} gin.H "Мониторинг запущен"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /connections/monitor/start [get]
func (h *Handler) StartPooling(c *gin.Context) {
	if err := h.usecase.StartPooling(); err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "monitoring started"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "monitoring is already running"})
	}
}

// StopMonitoring останавливает воркер проверки здоровья соединений
// @Summary Остановить мониторинг пула соединений
// @Description Останавливает периодическую проверку состояния всех соединений в пуле OPC UA
// @Tags Connection
// @Produce json
// @Success 200 {object} gin.H "Мониторинг остановлен"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /connections/monitor/stop [get]
func (h *Handler) StopPooling(c *gin.Context) {
	h.usecase.StopPooling()
	c.JSON(http.StatusOK, gin.H{"status": "monitoring stopped"})

}
