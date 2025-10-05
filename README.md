<div align="center">

# OpcUa Service

![Go](https://img.shields.io/badge/Go-1.19%2B-00ADD8?logo=go&logoColor=white)
![OPC UA](https://img.shields.io/badge/OPC%20UA-Supported-00ADD8?logo=opc-foundation&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-Integrated-00ADD8?logo=apachekafka&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-00ADD8?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-00ADD8?logo=docker&logoColor=white)

*Сервис для сбора данных по протоколу OPC UA, их отправки в Apache Kafka и управления через REST API*

</div> 

### 🚀 Возможности

- 🚀 **Потоковая передача в Kafka**: Все данные со станков в реальном времени отправляются в топик Apache Kafka для
  дальнейшей обработки и аналитики
- 🕹️ **Управляемый опрос**: Запускайте и останавливайте мониторинг для каждого станка индивидуально через REST API с
  настраиваемым интервалом
- 💾 **Персистентность**: Состояния подключений и опроса сохраняются в базе данных PostgreSQL, что позволяет
  автоматически восстанавливать их после перезапуска сервиса.
- 🌐 **REST API**: Удобный HTTP API для получения актуальных данных, проверки доступности станков и управления процессами
  опроса
- 🐳 **Простота развертывания**: Готовая конфигурация docker-compose.yml для быстрого запуска Apache Kafka и
  сопутствующих сервисов
- 🎛️ **Веб-интерфейс для Kafka**: Встроенный Kafka UI для удобного просмотра топиков и сообщений
- 🔧 **Универсальность**: Автоматическое извлечение и кэширование метаинформации из /probe для корректной интерпретации
  данных с различных станков


### ⚙️ Архитектура


```text
┌──────────────────────────┐            ┌──────────────────────────┐            ┌──────────────────────────┐ 
│         REST API         ├──────── >  │       OPC UA Сервис      ├───────── > │       OPC UA Сервер      │
│        (Gin-Gonic)       │            │         (Go App)         │ < ─────────┤     (Binary Protocol)    │
└──────────────────────────┘            └────────────┬───────┬─────┘            └──────────────────────────┘
             ^                                       │       └────────────────────────────────┐
             │                                       v                ( Опрос )               v
┌────────────┴─────────────┐            ┌──────────────────────────┐            ┌──────────────────────────┐ 
│       Пользователь /     │            │        PostgreSQL        │            │       Apache Kafka       │
│          Система         │            │  (Состояния подключений) │            │   (Потоковая обработка)  │
└──────────────────────────┘            └──────────────────────────┘            └──────────────────────────┘ 
```

<div align="center">

## 📦 Установка
</div>

### 1. Клонирование репозитория

```bash
git clone https://github.com/Khrllw/OpcUa_service.git
cd OpcUa_service
```

### 2. Конфигурация приложения

Откройте файл .env и при необходимости измените его

```dotenv
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=medapp

# App
APP_PORT=8080
GIN_MODE=debug

# Kafka
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=opc-data

# Logger
LOGGER_ENABLE=true
LOGGER_LOGS_DIR=./logs
LOGGER_LOG_LEVEL=DEBUG
LOGGER_SAVING_DAYS=7
```

### 3. Запуск Apache Kafka

```bash
docker-compose up
```

После запуска [Веб-интерфейс Kafka](http://localhost:8081)

Либо просмотреть сообщения сервера можно в реальном времени командой:<br>
`docker-compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic opc_data`

### 4. Запуск приложения

```
# Windows
./build/windows_opc_ua.exe

# Linux
./build/linux_opc_ua

# MacOS
./build/macos_opc_ua

# Golang
go run cmd/app/main.go
```


<div align="center">

## **🌐 API ENDPOINTS**
</div>

### 1. Создание подключения
  POST /api/v1/connect 

### 2. Получить пул подключений
  GET /api/v1/connect

### 3. Закрыть подключение
  DELETE /api/v1/connect

### 4. Проверить подключение
  POST /api/v1/connect/check

### 5. Начать сбор данных
  GET /api/v1/polling/start

### 6. Остановить сбор данных
  GET /api/v1/polling/stop


<div align="center">

## 🗂️ Структура проекта
</div>

```
OpcUa_service/
│ 
├── cmd/app/                     
│       └── 📄 main.go                 # Главная точка входа приложения
├── 📁 docs/                           # Документация проекта
├── internal/
│   ├── 📁 app/                        # Сборка и запуск приложения с помощью Fx для DI
│   ├── 📁 config/                     # Логика загрузки конфигурации из .env
│   ├── adapters/ 
│   │   ├── 📁 handlers/               # Обработчики HTTP-запросов (слой API на Gin)
│   │   └── 📁 repositories/           # Реализации репозиториев (PostgreSQL)
│   ├── 📁 domain/                     # Основные бизнес-сущности (entities) и модели (models)  
    ├── 📁 interfaces/                 # Go-интерфейсы для всех слоев (контракты)       
│   ├── middleware/
│   │   ├── 📁 logging/                # Логирование
│   │   └── 📁 swagger/                # Swagger/OpenAPI документация
│   ├── services/ 
│   │   ├── 📁 kafka/                  # Продюсер для Apache Kafka
│   │   └── opc_service/ 
│   │       ├── 📁 cert_manager/       # Управление OPC UA сертификатами
│   │       ├── 📁 opc_connector/      # Подключение и управление OPC UA сессиями
│   │       └── 📁 opc_communicator/   # Коммуникация с OPC UA серверами
│   └── 📁 usecases/                   # Сценарии использования, связывающие API и сервисный слой
├── logs/ 
│   └── 📄 { _date_ }.log              # Логи приложения по дням
├── pkg/
│   ├── 📁 client/                     # Клиентская библиотека для API
│   ├── 📁 opc_custom/                 # Зарегистрированные OPC UA структуры
│   └── 📁 machine_models/             # Поддерживаемые модели ЧПУ 
├── tools
│   ├── codegen/                        
│   │   ├── 📄 opc_types.yaml          # Конфиг кастомных структур OPC UA
│   │   └── 📄 codegen.go              # Скрипт для генерации кастомных структур OPC UA
│   └── build/
│       └── 📄 build.go                # Скрипт для сборки исполняемых файлов
├── 📁 build/                          # Папка с готовыми исполняемыми файлами
├── 📄 .env                            # Файл конфигурации
├── 📄 docker-compose.yml              # Файл для запуска Kafka и Kafka-UI
├── 📄 LICENSE
└── 📄 README.md
```

## 🆘 Поддержка

- 🐛 [Создайте issue](https://github.com/Khrllw/OpcUa_service/issues)
- 📧 Напишите на email: khrllw@gmail.com

## 📝 Лицензия

Проект распространяется под [лицензией MIT](LICENSE)

Copyright (c) 2025 khrllw
