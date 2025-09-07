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

### Создание подключения ( POST /api/v1/connect )

```json
{
  // Тип соединения с OPC UA сервером
  // Возможные значения:
  // "certificate" - аутентификация по сертификату и ключу
  // "password"    - аутентификация по имени пользователя и паролю
  // "anonymous"   - анонимное соединение без логина и пароля
  "connectionType": "certificate",

  // Имя пользователя (используется только если connectionType = "password")
  "username": "<your_username_here>",

  // Пароль пользователя (используется только если connectionType = "password")
  "password": "<your_password_here>",

  // Файл сертификата (используется только если connectionType = "certificate")
  "certificate": "<Base64_encoded_bytes_of_certificate>",
  
  // Путь к приватному ключу (используется только если connectionType = "certificate")
  "key": "<Base64_encoded_bytes_of_key>",

  // URL OPC UA сервера
  "endpointURL": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",

  // Режим безопасности сообщений
  // Возможные значения: "None", "Sign", "SignAndEncrypt"
  "mode": "SignAndEncrypt",

  // Политика безопасности OPC UA
  // Примеры: "Basic256Sha256", "Basic256", "Basic128Rsa15", "None"
  "policy": "Basic256Sha256",

  // Таймаут опроса в секундах
  "timeout": 5,

  // Производитель станка
  "manufacturer": "Heidenhain",

  // Модель станка
  "model": "TNC640"
}
```

```json
{
  "data": {
    "UUID": "848f855f-f4c2-45f1-a216-c9fa64a6369f"
  },
  "message": "Successfully connected",
  "status": "success",
  "type": "object"
}
```

### Получить пул подключений ( GET /api/v1/connect )

```json
{
  "data": {
    "poolSize": 1,
    "connections": [
      {
        "status": "HEALTHY",
        "description": "Connection is healthy and responsive",
        "UUID": "6b0af31a-badc-418a-893e-a878e144adae",
        "sessionID": "ns=3;i=2623872058",
        "config": {
          "type": "certificate",
          "connection": {
            "EndpointURL": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
            "Certificate": "<Base64_encoded_bytes_of_certificate>",
            "Key": "<Base64_encoded_bytes_of_key>",
            "Policy": "Basic256Sha256",
            "Mode": "SignAndEncrypt",
            "Timeout": 3000000000, // в миллисекундах
            "Manufacturer": "Heidenhain",
            "Model": "TNC640"
          }
        },
        "createdAt": "2025-09-03T19:00:47.7478375+03:00",
        "lastUsed": "2025-09-03T19:00:47.7478375+03:00",
        "useCount": 1
      }
    ]
  },
  "message": "Successfully get connection pool",
  "status": "success",
  "type": "object"
}
```

### Закрыть подключение ( DELETE /api/v1/connect )

```json
{
  "UUID": "ns=3;i=3948713892"
}
```

```json
{
  "data": {
    "disconnected": true
  },
  "message": "Successfully disconnected",
  "status": "success",
  "type": "object"
}
```

### Проверить подключение ( POST /api/v1/connect/check )

```json
{
  "UUID": "ns=3;i=3948713891"
}
```

```json
{
  "data": {
    "status": "HEALTHY",
    "description": "Connection is healthy and responsive",
    "UUID": "6b0af31a-badc-418a-893e-a878e144adae",
    "sessionID": "ns=3;i=2623872058",
    "config": {
      "type": "certificate",
      "connection": {
        "EndpointURL": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
        "Certificate": "<Base64_encoded_bytes_of_certificate>",
        "Key": "<Base64_encoded_bytes_of_key>",
        "Policy": "Basic256Sha256",
        "Mode": "SignAndEncrypt",
        "Timeout": 3000000000, // в миллисекундах
        "Manufacturer": "Heidenhain",
        "Model": "TNC640"
      },
      "createdAt": "2025-09-03T19:00:47.7478375+03:00",
      "lastUsed": "2025-09-03T19:00:47.7478375+03:00",
      "useCount": 1
    }
  },
  "message": "Successfully get connection info",
  "status": "success",
  "type": "object"
}
```
### Начать сбор данных ( GET /api/v1/polling/start )
```json
{
"UUID": "848f855f-f4c2-45f1-a216-c9fa64a6369f"
}
```

```json
{
  "data": {
    "polled": true
  },
  "message": "Polling started for machine 848f855f-f4c2-45f1-a216-c9fa64a6369f",
  "status": "success",
  "type": "object"
}
```

### Остановить сбор данных ( GET /api/v1/polling/stop )
```json
{
  "UUID": "12840be9-36b2-4ecb-8243-b9d9e0952a03"
}
```

```json
{
  "data": {
    "polled": false
  },
  "message": "Polling stopped for machine 12840be9-36b2-4ecb-8243-b9d9e0952a03",
  "status": "success",
  "type": "object"
}
```

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
├── tools/build/
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
