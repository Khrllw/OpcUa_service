package app

import (
	"context"
	"go.uber.org/fx"
	"net/http"
	"opc_ua_service/internal/adapters/handlers"
	"opc_ua_service/internal/adapters/repositories"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/middleware/logging"
	"opc_ua_service/internal/middleware/swagger"
	"opc_ua_service/internal/services/kafka"
	"opc_ua_service/internal/services/opc_service"

	// "opc_ua_service/internal/middleware/swagger"
	// "opc_ua_service/internal/services"
	"opc_ua_service/internal/usecases"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.LoadConfig,
		),
		LoggingModule,
		KafkaModule,
		RepositoryModule,
		ServiceModule,
		UsecaseModule,
		HttpServerModule,
	)
}

func ProvideLoggers(cfg *config.Config) *logging.Logger {
	loggerCfg := &logging.Config{
		Enabled:    cfg.Logging.Enable,
		Level:      cfg.Logging.Level,
		LogsDir:    cfg.Logging.LogsDir,
		SavingDays: IntToUint(cfg.Logging.SavingDays),
	}

	logger := logging.NewLogger(loggerCfg, "APP", cfg.App.Version)
	return logger
}

var LoggingModule = fx.Module("logging_module",
	fx.Provide(
		ProvideLoggers,
	),
	fx.Invoke(func(l *logging.Logger) {
		l.Info("Logging system initialized")
	}),
)

func InvokeHttpServer(lc fx.Lifecycle, h http.Handler) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Close()
		},
	})
}

// Swagger-конфигуратор
func NewSwaggerConfig(cfg *config.Config) *swagger.Config {
	return &swagger.Config{
		Enabled: true,
		Path:    "/swagger",
	}
}

var HttpServerModule = fx.Module("http_server_module",
	fx.Provide(
		NewSwaggerConfig,
		handlers.NewHandler,
		handlers.ProvideRouter,
	),
	fx.Invoke(InvokeHttpServer),
)

var ServiceModule = fx.Module("service_module",
	fx.Provide(opc_service.NewOpcService),
)

var KafkaModule = fx.Module("kafka_module",
	fx.Provide(kafka.NewKafka),
)

var RepositoryModule = fx.Module("postgres_module",
	fx.Provide(repositories.NewRepository),
)

var UsecaseModule = fx.Module("usecases_module",
	fx.Provide(
		usecases.NewUsecases,
	),
)

// TODO: Может быть вынести в services
func IntToUint(c int) uint {
	if c < 0 {
		panic([2]any{"a negative number", c})
	}
	return uint(c)
}
