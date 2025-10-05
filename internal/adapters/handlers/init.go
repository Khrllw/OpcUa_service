package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	"net/http"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"opc_ua_service/internal/middleware/swagger"
)

type Handler struct {
	logger  *logging.Logger
	usecase interfaces.Usecases
}

// NewHandler создает новый экземпляр Handler со всеми зависимостями
func NewHandler(usecase interfaces.Usecases, parentLogger *logging.Logger, service interfaces.OpcService) *Handler {
	handlerLogger := parentLogger.WithPrefix("HANDLER")
	handlerLogger.Info("Handler initialized",
		"component", "GENERAL",
	)
	return &Handler{
		logger:  handlerLogger,
		usecase: usecase,
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler, cfg *config.Config, swagCfg *swagger.Config) http.Handler {
	gin.SetMode(cfg.App.GinMode)
	r := gin.Default()

	// Swagger-роутер
	swagger.Setup(r, swagCfg)

	// LoggerMiddleware
	r.Use(LoggingMiddleware(h.logger))

	// Общая группа для API
	baseRouter := r.Group("/api/v1")

	// Подключение
	connectGroup := baseRouter.Group("/connect")
	connectGroup.POST("/", h.AddConnection)        // Добавить соединение
	connectGroup.GET("/", h.GetConnectionPool)     // Получить пул открытых соединений
	connectGroup.DELETE("/", h.CloseConnection)    // Закрыть соединение
	connectGroup.POST("/check", h.CheckConnection) // Проверить соединение

	// Мониторинг
	pollingGroup := baseRouter.Group("/polling")
	pollingGroup.GET("/start", h.StartPollingByID) // Начать опрос
	pollingGroup.GET("/stop", h.StopPollingByID)   // Остановить опрос

	return r
}
