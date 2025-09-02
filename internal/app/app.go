package app

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"log"
	"net/http"
	_ "opc_ua_service/docs"
	"opc_ua_service/internal/adapters/handlers"
	"opc_ua_service/internal/adapters/producers"
	"opc_ua_service/internal/adapters/repositories"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"opc_ua_service/internal/middleware/swagger"
	"opc_ua_service/internal/services/opc_service"
	"opc_ua_service/internal/usecases"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.LoadConfig,
		),
		LoggingModule,
		RepositoryModule,
		ProducerModule,
		ServiceModule,
		UsecaseModule,
		HttpServerModule,
		fx.Invoke(InvokeRestoreConnections),
	)
}

func ProvideLoggers(cfg *config.Config) *logging.Logger {
	loggerCfg := &logging.Config{
		Enabled:    cfg.Logging.Enable,
		Level:      cfg.Logging.Level,
		LogsDir:    cfg.Logging.LogsDir,
		SavingDays: intToUint(cfg.Logging.SavingDays),
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

// InvokeHttpServer запускает HTTP-сервер
func InvokeHttpServer(lc fx.Lifecycle, h http.Handler) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("HTTP server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}

// InvokeGracefulShutdown обеспечивает корректное завершение работы сервисов
func InvokeGracefulShutdown(lc fx.Lifecycle, connector interfaces.OpcService, producer interfaces.DataProducer) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Корректное завершение работы сервисов...")
			connector.CloseAll()
			//communicator.CloseAll()
			if err := producer.Close(); err != nil {
				log.Printf("Ошибка при закрытии Kafka продюсера: %v", err)
				return err
			}
			log.Println("Все сервисы успешно остановлены.")
			return nil
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
	fx.Invoke(InvokeHttpServer, InvokeGracefulShutdown),
)

var ProducerModule = fx.Module("producer_module",
	fx.Provide(producers.NewKafkaProducer),
)

var ServiceModule = fx.Module("service_module",
	fx.Provide(opc_service.NewOpcService),
)

var RepositoryModule = fx.Module("postgres_module",
	fx.Provide(repositories.NewRepository),
)

var UsecaseModule = fx.Module("usecases_module",
	fx.Provide(
		usecases.NewUsecases,
	),
)

func intToUint(c int) uint {
	if c < 0 {
		panic([2]any{"a negative number", c})
	}
	return uint(c)
}

// InvokeRestoreConnections восстанавливает подключения и опросы при старте приложения.
func InvokeRestoreConnections(lc fx.Lifecycle, uc interfaces.Usecases, dbRepo interfaces.Repository, logger *logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Restoring connections from the database...")
			machines, err := dbRepo.GetAllCncMachines()
			if err != nil {
				logger.Error("Failed to get machine list from DB", "error", err)
				return nil
			}

			if len(machines) == 0 {
				logger.Info("No saved connections found to restore.")
				return nil
			}

			for _, machine := range machines {
				logger.Info("Attempting to restore connection", "UUID", machine.UUID, "model", machine.Model, "endpoint", machine.EndpointURL)

				connInfo, _ := uc.RestoreConnection(machine)

				if connInfo != nil {
					logger.Info("Connection restored successfully in pool", "UUID", machine.UUID)
				} else {
					logger.Warn("Connection restored in pool but is unhealthy. Will retry in background.", "UUID", machine.UUID)
				}
			}
			return nil
		},
	})
}
