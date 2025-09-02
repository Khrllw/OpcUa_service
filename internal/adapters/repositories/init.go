package repositories

import (
	"fmt"
	"log"
	"opc_ua_service/internal/adapters/repositories/anonymous_connection"
	"opc_ua_service/internal/adapters/repositories/certificate_connection"
	"opc_ua_service/internal/adapters/repositories/cnc_machine"
	"opc_ua_service/internal/adapters/repositories/password_connection"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/interfaces"
)

type Repository struct {
	interfaces.CncMachineRepository
	interfaces.CertificateConnectionRepository
	interfaces.PasswordConnectionRepository
	interfaces.AnonymousConnectionRepository
}

func NewRepository(cfg *config.Config) (interfaces.Repository, error) {
	//logger := logging.NewModuleLogger("ADAPTER", "POSTGRES", parentLogger)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Вывод в stdout
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Порог для медленных запросов
			LogLevel:                  logger.Info,            // Уровень логирования (Info - все запросы)
			IgnoreRecordNotFoundError: true,                   // Игнорировать ошибки "запись не найдена"
			Colorful:                  true,                   // Цветной вывод
		},
	)

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Выполнение автомиграций
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("ошибка выполнения автомиграций: %w", err)

	}

	return &Repository{
		cnc_machine.NewCncMachineRepository(db),
		certificate_connection.NewCertificateConnectionRepository(db),
		password_connection.NewPasswordConnectionRepository(db),
		anonymous_connection.NewAnonymousConnectionRepository(db),
	}, nil

}

// autoMigrate - выполнение автомиграций для моделей
func autoMigrate(db *gorm.DB) error {

	// Удаляем таблицы в правильном порядке зависимостей
	tables := []string{
		"certificate_connection",
		"anonymous_connection",
		"password_connection",
		"cnc_machine",
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	// Создаем таблицы
	models := []interface{}{
		&entities.CncMachine{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return nil
}
