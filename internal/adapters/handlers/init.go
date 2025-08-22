package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/swaggo/files"
	"net/http"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"opc_ua_service/internal/middleware/swagger"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	logger  *logging.Logger
	usecase interfaces.Usecases
	service interfaces.OpcService
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
		service: service,
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler, cfg *config.Config, swagCfg *swagger.Config) http.Handler {
	r := gin.Default()

	// Swagger-роутер
	swagger.Setup(r, swagCfg)

	// Logger
	r.Use(LoggingMiddleware(h.logger))

	// Общая группа для API
	baseRouter := r.Group("/api/v1")

	// Подключение
	connectGroup := baseRouter.Group("/connect")
	connectGroup.POST("/", h.AddConnection)
	connectGroup.GET("/", h.GetConnectionPool)
	connectGroup.DELETE("/", h.CloseConnection)
	connectGroup.POST("/check", h.CheckConnection)

	// Мониторинг
	poolingGroup := baseRouter.Group("/pooling")
	poolingGroup.GET("/start", h.StartPooling)
	poolingGroup.GET("/stop", h.StopPooling)

	return r
}
