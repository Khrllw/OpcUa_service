<div align="center">

# OpcUa Service

![Go](https://img.shields.io/badge/Go-1.19%2B-00ADD8?logo=go&logoColor=white)
![OPC UA](https://img.shields.io/badge/OPC%20UA-Supported-00ADD8?logo=opc-foundation&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-Integrated-00ADD8?logo=apachekafka&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-00ADD8?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-00ADD8?logo=docker&logoColor=white)

*Сервис для сбора данных по протоколу OPC UA, их отправки в Apache Kafka и управления через REST API*

</div> 

## 🚀 Возможности

- 🚀 **Потоковая передача в Kafka**: Все данные со станков в реальном времени отправляются в топик Apache Kafka для дальнейшей обработки и аналитики
- 🕹️ **Управляемый опрос**: Запускайте и останавливайте мониторинг для каждого станка индивидуально через REST API с настраиваемым интервалом
- 💾 **Персистентность**: Состояния подключений и опроса сохраняются в базе данных PostgreSQL, что позволяет автоматически восстанавливать их после перезапуска сервиса.
- 🌐 **REST API**: Удобный HTTP API для получения актуальных данных, проверки доступности станков и управления процессами опроса
- 🐳 **Простота развертывания**: Готовая конфигурация docker-compose.yml для быстрого запуска Apache Kafka и сопутствующих сервисов
- 🎛️ **Веб-интерфейс для Kafka**: Встроенный Kafka UI для удобного просмотра топиков и сообщений
- 🔧 **Универсальность**: Автоматическое извлечение и кэширование метаинформации из /probe для корректной интерпретации данных с различных станков

## 🗂 Архитектура

```
+-------------+    HTTP/REST    +----------------+    OPC UA Binary     +------------------+
│             │ <-------------> │                 │ <-----------------> │                  │
│   Клиент    │      JSON       │       API       │       asyncua       │  OPC UA Сервер   │
│             │                 │                 │                     │                  │
+-------------+                 +----------------+                      +------------------+
```

## 📦 Установка и Запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/Khrllw/OpcUa_service.git
cd OpcUa_service
```

## 🔌 API

## Создание подключения

```http
POST /api/v1/connect
```

```bash
curl -X POST "/api/v1/connect"
```

```json
{
  "certificate": "certs/new_app/opcua-client-application-cert-20250806-132304.der",
  "key": "certs/new_app/opcua-client-application-cert-20250806-132304.key",
  "endpointURL": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
  "mode": "SignAndEncrypt",
  "policy": "Basic256Sha256",
  "connectionType": "certificate"
}
```

```json
{
  "status": "OK",
  "token": "",
  "connectionInfo": {
    "sessionId": "ns=3;i=2069698690",
    "config": {
      "EndpointURL": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
      "Certificate": "certs/new_app/opcua-client-application-cert-20250806-132304.der",
      "Key": "certs/new_app/opcua-client-application-cert-20250806-132304.key",
      "Policy": "Basic256Sha256",
      "Mode": "SignAndEncrypt",
      "Timeout": 30000000000,
      "Manufacturer": "Heidenhain",
      "Model": "TNC640"
    },
    "createdAt": "2025-08-25T12:50:15.867731+03:00",
    "lastUsed": "2025-08-25T12:50:15.867731+03:00",
    "useCount": 1,
    "isHealthy": true
  }
}
```

## Получить пул подключений

```http
GET /api/v1/connect
```

```bash
curl -X GET "http://localhost:8080/api/v1/connect"
```

```json
{
  "connections": [
    {
      "sessionId": "ns=3;i=2069698690",
      "config": {
        "EndpointURL": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
        "Certificate": "certs/new_app/opcua-client-application-cert-20250806-132304.der",
        "Key": "certs/new_app/opcua-client-application-cert-20250806-132304.key",
        "Policy": "Basic256Sha256",
        "Mode": "SignAndEncrypt",
        "Timeout": 30000000000,
        "Manufacturer": "Heidenhain",
        "Model": "TNC640"
      },
      "createdAt": "2025-08-25T12:50:15.867731+03:00",
      "lastUsed": "2025-08-25T12:50:15.867731+03:00",
      "useCount": 1,
      "isHealthy": true
    }
  ],
  "poolSize": 1,
  "status": "ok"
}
```

## Закрыть подключение

```http
DELETE /api/v1/connect
```

```bash
curl -X DELETE "http://localhost:8080/api/v1/connect"
```

```json
{
  "sessionID": "ns=3;i=3948713892"
}
```

```json
{
  "message": "Session ns=3;i=2069698690 disconnected successfully"
}
```

## Проверить подключение

```http
POST /api/v1/connect/check
```

```bash
curl -X POST "http://localhost:8080/api/v1/connect/check"
```

```json
{
  "sessionID": "ns=3;i=2069698691"
}
```

```json
{
  "details": {
    "created_at": "2025-08-25T12:52:26.5091654+03:00",
    "endpoint": "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
    "is_healthy": true,
    "last_used": "2025-08-25T12:52:26.5091654+03:00",
    "session_id": "ns=3;i=2069698691",
    "use_count": 1
  },
  "sessionID": "ns=3;i=2069698691",
  "status": "healthy"
}
```

## Начать сбор данных

```http
GET /api/v1/polling/start
```

```bash
curl -X GET "http://localhost:8080/api/v1/polling/start"
```

```json
{
  "status": "monitoring started"
}
```

## Остановить сбор данных

```http
GET /api/v1/polling/stop
```

```bash
curl -X GET "http://localhost:8080/api/v1/polling/stop"
```

```json
{
  "status": "monitoring stopped"
}
```



## 🔧 Структура проекта

```

OpcUa_service/
│
├── 📁 certs/ # Директория для сертификатов безопасности
│ ├── 📁 new_app/ # Сертификаты клиентского приложения
│ │ ├── 📄 certificate.pem # Публичный сертификат приложения
│ │ ├── 📄 private_key.pem # Приватный ключ приложения
│ │ └── 📄 trust_list.pem # Доверенные сертификаты
│ └── 📁 new_server/ # Сертификаты серверной части
│ ├── 📄 server_cert.pem # Серверный сертификат
│ └── 📄 server_key.pem # Приватный ключ сервера
│
├── 📁 cmd/ # Точки входа приложения
│ └── 📁 app/
│ └── 📄 main.go # Основной entry point приложения
│
├── 📁 docs/ # Документация проекта
│
├── 📁 internal/ # Внутренняя логика приложения
│ │
│ ├── 📁 adapters/ # Адаптеры для внешних систем
│ │ ├── 📁 handlers/ # HTTP обработчики REST API endpoints
│ │ │ ├── 📄 connection_handler.go # Обработчик подключений к OPC UA серверам
│ │ │ ├── 📄 init.go # Инициализация HTTP роутеров и обработчиков
│ │ │ ├── 📄 polling_handler.go # Обработчик опроса данных с OPC UA серверов
│ │ │ ├── 📄 response.go # Утилиты для формирования HTTP ответов
│ │ │ └── 📄 middleware.go # Промежуточное ПО для HTTP обработки
│ │ │
│ │ └── 📁 repositories/ # Репозитории для работы с persistence layer
│ │ ├── 📁 cnc_machine/ # Репозитории для данных станков с ЧПУ
│ │ │ ├── 📄 func.go # Интерфейс репозитория станков
│ │ │ └── 📄 init.go # In-memory реализация репозитория (для тестов)
│ │ └── 📄 init.go # Инициализация репозиториев
│ │       
│ ├── 📁 app/ # Конфигурация и инициализация приложения
│ │ └── 📄 app.go # Сборка и конфигурация DI контейнера приложения
│ │
│ ├── 📁 config/ # Конфигурационные параметры приложения
│ │ └── 📄 config.go # Загрузка и валидация конфигурации из env переменных
│ │
│ ├── 📁 domain/ # Доменные модели и бизнес-сущности
│ │ ├── 📁 entities/ # Core бизнес-сущности
│ │ │ └── 📄 cnc_machine.go # Сущность станка
│ │ │
│ │ └── 📁 models/ # DTO модели для API и persistence
│ │ ├── 📁 machine_models/ # Модели данных для станков
│ │ │ └── 📄 machine_TNC640.go # Специфичная модель для станков TNC 640
│ │ │
│ │ ├── 📁 opc_custom/ # Кастомные OPC UA data types
│ │ │ ├── 📄 cutter.go # Модель данных фрезы/инструмента
│ │ │ ├── 📄 execution_stack.go # Модель стека выполнения ЧПУ программ
│ │ │ └── 📄 tool_data.go # Модель данных инструмента
│ │ │
│ │ ├── 📄 cnc_machine.go # Модель станка для API
│ │ ├── 📄 connection.go # Модель OPC UA подключения
│ │ └── 📄 opc.go # Базовые OPC UA модели
│ │
│ ├── 📁 interfaces/ # Интерфейсы (контракты) для инверсии зависимостей
│ │ ├── 📄 kafka.go # Интерфейсы для Kafka producers/consumers
│ │ ├── 📄 machine_data.go # Интерфейсы для работы с данными станков
│ │ ├── 📄 repository.go # Базовые интерфейсы репозиториев
│ │ ├── 📄 uasecase.go # Интерфейсы use cases
│ │ └── 📄 services.go # Интерфейсы сервисного слоя
│ │
│ ├── 📁 middleware/ # HTTP middleware компоненты
│ │ ├── 📁 logging/ # Middleware для логирования
│ │ │ └── 📄 logger.go # Structured logging middleware
│ │ │
│ │ └── 📁 swagger/ # Swagger/OpenAPI документация
│ │ └── 📄 swagger.go # Генерация и обслуживание Swagger UI
│ │
│ ├── 📁 services/ # Бизнес-сервисы и application services
│ │ ├── 📁 kafka/ # Сервисы для работы с Apache Kafka
│ │ │ ├── 📄 init.go # Инициализация Kafka producers/consumers
│ │ │ └── 📄 kafka.go # Реализация Kafka клиента
│ │ │
│ │ └── 📁 opc_service/ # Сервисы для работы с OPC UA
│ │ ├── 📁 cert_manager/ # Управление OPC UA сертификатами
│ │ │ ├── 📄 cert.go # Базовые операции с сертификатами
│ │ │ ├── 📄 cert_analyzer.go # Анализ и валидация сертификатов
│ │ │ └── 📄 key_cert.go # Управление ключами и сертификатами
│ │ │
│ │ ├── 📁 opc_communicator/ # Коммуникация с OPC UA серверами
│ │ │ └── 📄 communicator.go # Основной клиент для OPC UA communication
│ │ │
│ │ └── 📁 opc_connector/ # Подключение и управление OPC UA сессиями
│ │ └── 📄 connector.go # Менеджер подключений и пул сессий
│ │
│ └── 📁 usecases/ # Use Cases (бизнес-сценарии приложения)
│ ├── 📄 connection_usecase.go # Сценарии управления OPC UA подключениями
│ ├── 📄 init.go # Инициализация use cases
│ └── 📄 polling_usecase.go # Сценарии опроса данных с OPC UA серверов
│
├── 📁 logs/ # Логи приложения по дням
│ └── 📄 { _date_ }.log
│
├── 📁 pkg/ # Вспомогательные пакеты для повторного использования
│ └── 📁 errors/ # Кастомные ошибки и error handling
│ └── 📄 errors.go # Определение кастомных ошибок и утилиты обработки
│
├── 📄 .env.example # Пример переменных окружения
├── 📄 .gitignore # Git ignore правила
├── 📄 requirements.txt # Зависимости Python
└── 📄 README.md # Документация проекта

```
