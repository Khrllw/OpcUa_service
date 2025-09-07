package repositories

import (
	"fmt"
	"log"
	"os"
	"time"

	"opc_ua_service/internal/adapters/repositories/anonymous_connection"
	"opc_ua_service/internal/adapters/repositories/certificate_connection"
	"opc_ua_service/internal/adapters/repositories/cnc_machine"
	"opc_ua_service/internal/adapters/repositories/password_connection"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	interfaces.CncMachineRepository
	interfaces.CertificateConnectionRepository
	interfaces.PasswordConnectionRepository
	interfaces.AnonymousConnectionRepository
}

func NewRepository(cfg *config.Config, appLogger *logging.Logger) (interfaces.Repository, error) {
	// Шаг 1: Проверка и создание БД при необходимости
	if err := ensureDatabaseExists(cfg, appLogger); err != nil {
		return nil, err
	}

	// Шаг 2: Подключение к основной БД
	appDb, err := connectToAppDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to main database: %w", err)
	}

	// Шаг 3: Автоматическая миграция
	if err := autoMigrate(appDb); err != nil {
		return nil, fmt.Errorf("error during migration: %w", err)
	}

	// Шаг 4: Инициализация репозиториев
	return &Repository{
		CncMachineRepository:            cnc_machine.NewCncMachineRepository(appDb),
		CertificateConnectionRepository: certificate_connection.NewCertificateConnectionRepository(appDb),
		PasswordConnectionRepository:    password_connection.NewPasswordConnectionRepository(appDb),
		AnonymousConnectionRepository:   anonymous_connection.NewAnonymousConnectionRepository(appDb),
	}, nil
}

// ensureDatabaseExists проверяет существование БД и создает её при необходимости
func ensureDatabaseExists(cfg *config.Config, log *logging.Logger) error {
	serviceDB, err := gorm.Open(postgres.Open(buildDSN(cfg, "postgres")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to service database: %w", err)
	}
	defer closeDB(serviceDB)

	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)"
	if err := serviceDB.Raw(query, cfg.Database.DBName).Scan(&exists).Error; err != nil {
		return fmt.Errorf("database check error: %w", err)
	}

	if !exists {
		log.Info("Database not found. Creating...", "db_name", cfg.Database.DBName)
		createQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName)
		if err := serviceDB.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Info("Database created successfully.", "db_name", cfg.Database.DBName)
	} else {
		log.Info("Database already exists.", "db_name", cfg.Database.DBName)
	}

	return nil
}

// connectToAppDatabase подключается к основной БД
func connectToAppDatabase(cfg *config.Config) (*gorm.DB, error) {
	logWriter := log.New(os.Stdout, "\r\n", log.LstdFlags)
	dbLogger := logger.New(logWriter, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	return gorm.Open(postgres.Open(buildDSN(cfg, cfg.Database.DBName)), &gorm.Config{
		Logger: dbLogger,
	})
}

// autoMigrate выполняет миграцию и очистку зависимых таблиц
func autoMigrate(db *gorm.DB) error {
	// Удаление таблиц в порядке зависимостей
	tablesToDrop := []string{
		"certificate_connection",
		"anonymous_connection",
		"password_connection",
		"cnc_machine",
	}

	for _, table := range tablesToDrop {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	models := []interface{}{
		&entities.CncMachine{},
		&entities.CertificateConnection{},
		&entities.PasswordConnection{},
		&entities.AnonymousConnection{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	return nil
}

// buildDSN формирует строку подключения
func buildDSN(cfg *config.Config, dbName string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		dbName,
		cfg.Database.Port,
	)
}

// closeDB закрывает соединение с БД
func closeDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
}
