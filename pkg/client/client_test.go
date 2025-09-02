package client

import (
	"context"
	"log"
	"testing"
	"time"

	"opc_ua_service/internal/domain/models"
)

// Для запуска теста должен быть запущен MTConnect сервис на localhost:8080
func TestFullClientWorkflow(t *testing.T) {
	// Инициализация клиента
	api := NewClient("http://localhost:8080")
	ctx := context.Background()

	// 1. Создание подключения
	log.Println("Шаг 1: Создание подключения...")
	createReq := models.ConnectionRequest{
		ConnectionType: "certificate",
		EndpointURL:    "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC",
		Key:            "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC0UCG//tSscV93zFqPlk3aag9kVpF/jqU86zRhxf6uW1fYGeFcIM+mTdCUXPRf96uOpS0oMK400vG3CGFYKoTquY5Ot8loz961EdlOOnGF02L/ZlnERO8m9FrCVmqhUkZgPawagRppgj8HeNIPGW6s8s8LSsnRGYWKRLOC49YiQKzdkcEmHhQbyrji5vMV0VhZcd99/1sfBCZM5z3Bap0J4W7d+fZBT4fjhQa6Q8Q6QLLgyKYpEMflmuTp9icT/yVMB5Pc4grLb5pcWXjxkz+WRMSQsefqIZ7BwvjNsImwCe9vaTdQ7Q4MszDmcreG0MCBqtHLFIPbsnteYi/p+m5LAgMBAAECggEAE8cCR1aCfz5ENBX/DObLN/IQKXeb/TtpK7W6CMtjHulinjpgee7DYpYAaHWK/F0PEoCoOP317adW3zY/xHWNdK17oBi6h7uHzaEyjkkMxe/6A3zMyyFV1eumv3rOYU7D7K8htG/atMiWHzyvqvulebom3tyKJwbp8HuOm68KgAvGZooOWJf5oXSL0IuRAVRn3Df6dDMC21O2tAygl3sOG6tuaqJ8G236Q7nQKnD1Jb0YOdYWk0K8qvDPAzBfDB6S/aXATKFNReBslVLuSKA/h1rYgvj/rvKzHaAehTV/iU+4RQhH1SdjIsCgJBe7cBMIjPuwCeGgmNMx9biHhAdN6QKBgQD3QvluIkS+H67h3HzZlRlRubTTQBDt36wB1xpWGQPAFFuhzPWS4+yK/z0yU4KeJXD4F3kIE42hzX+aYwRIsVeaUhGzhLE5siTzLacPD5Ool1zGoimHV+A2A2n4JaIt/PPH5aG9f8ttIgtn+yq683LB1HJUYpuZ6FyySJWTyVbxdQKBgQC6r3XCfzzLtPiMqUvC4hg472Yy4EB9VSSe7kSvYpCd05mBBJd01IikWA548wfeGxgxJ8suTilxCTQLRCT3Hn4GmBGb7gGSxOw2y8HAzovmi9CHDJ6fwEoxBV0wF/C/uGJATFm4eAr66AdNq+4CNZ3EjiI4DH7C2tQmJUYK7JoovwKBgAieAOthrluh5wpgEMnUdGlwu2iRVwWzQd9ei8BsZsEO9JKS/gv8fYXql0tltaulSmabCtDJPaph6wyKXt/Zrl/mdE95VGPaXYdMFAJmXJMHk2goxqG84kd/nvXS+e/4XNaeniBoj8Jh6VvaWQbi7SDsMn/WX+3hNznPZcccwTbxAoGAJwie378I8DLzsT2IuMPberQbs1GOSmZuFMkPFXjPciCXPRG/tU7nDy3WQNXX9EnIAicm5ZS0N41ME3r5G66FfU14iRj3vT9tgHuUFINbXyYmwMYTuKVVHfDYLkEjNoMQEA+mxtpauWGgfU4QouehCEMLxppeOtHUf/FVNt2H0jMCgYEAyK5HeBjVpLwGxo//J+GqYxLaJxnXyWQanLXdQFHda4XclgqJV1kOdFixPddBpXzEl5yqwzn/GZdr0if4QtvT7tSLLqKlK09FM3A7hYP023+1pTAFxEE5iI27bqJ/03Yt+ZVCvDk+B5LIEaDVsCMh1fbEJqYIiCMJQmAi4wqnkx4=",
		Certificate:    "MIIEpjCCA46gAwIBAgIBATANBgkqhkiG9w0BAQsFADB7MRAwDgYDVQQDDAdhcHBAYXBwMQswCQYDVQQGEwJSVTEMMAoGA1UECAwDb2JsMQ4wDAYDVQQHDAVwdW5rdDENMAsGA1UECgwEY29tcDEYMBYGCSqGSIb3DQEJARYJbWFpbEBtYWlsMRMwEQYKCZImiZPyLGQBGRYDYXBwMB4XDTI1MDgwNjEzMjMwNVoXDTI2MDgwNjEzMjMwNVowezEQMA4GA1UEAwwHYXBwQGFwcDELMAkGA1UEBhMCUlUxDDAKBgNVBAgMA29ibDEOMAwGA1UEBwwFcHVua3QxDTALBgNVBAoMBGNvbXAxGDAWBgkqhkiG9w0BCQEWCW1haWxAbWFpbDETMBEGCgmSJomT8ixkARkWA2FwcDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALRQIb/+1KxxX3fMWo+WTdpqD2RWkX+OpTzrNGHF/q5bV9gZ4Vwgz6ZN0JRc9F/3q46lLSgwrjTS8bcIYVgqhOq5jk63yWjP3rUR2U46cYXTYv9mWcRE7yb0WsJWaqFSRmA9rBqBGmmCPwd40g8ZbqzyzwtKydEZhYpEs4Lj1iJArN2RwSYeFBvKuOLm8xXRWFlx333/Wx8EJkznPcFqnQnhbt359kFPh+OFBrpDxDpAsuDIpikQx+Wa5On2JxP/JUwHk9ziCstvmlxZePGTP5ZExJCx5+ohnsHC+M2wibAJ729pN1DtDgyzMOZyt4bQwIGq0csUg9uye15iL+n6bksCAwEAAaOCATMwggEvMEYGCWCGSAGG+EIBDQQ5FjdPcGVuU1NMIEdlbmVyYXRlZCBPUEMgVUEgQ2xpZW50IEFwcGxpY2F0aW9uIENlcnRpZmljYXRlMB0GA1UdDgQWBBTEhWplbGy8bT7N+pXme9v/2zfOdTCBpQYDVR0jBIGdMIGagBTEhWplbGy8bT7N+pXme9v/2zfOdaF/pH0wezEQMA4GA1UEAwwHYXBwQGFwcDELMAkGA1UEBhMCUlUxDDAKBgNVBAgMA29ibDEOMAwGA1UEBwwFcHVua3QxDTALBgNVBAoMBGNvbXAxGDAWBgkqhkiG9w0BCQEWCW1haWxAbWFpbDETMBEGCgmSJomT8ixkARkWA2FwcIIBATAJBgNVHRMEAjAAMBMGA1UdEQQMMAqGA2FwcIIDYXBwMA0GCSqGSIb3DQEBCwUAA4IBAQACQ2jfcuFAdwVmYxSPRBUJQuTzr1dwq4krJUff1TpcpBaZR7hU4aCwfl4vjQoSqWveQC+mWoXGzMhmPfjjIGJ0UCLHeOgSuHk8/49velvYFArphRP4Y0bUowy7u5umxIEswkqsRdYSQUrUN/hX7hQc76V2t8sDsfEEojXgPmLW1HVikhRgGWnaTDbN5v3NiYTKopMKwGE3SeKD7GYEJQ+4ZN2E9rbm9a1VWurTjXbFkiCC6B+Yv8mOUrCU+sD68fnJlFWyDQuyt462dEToF5gFqRR3fxbTs+Gj4DozaVJ2egp5WefopuIWgKada2IbHbLUwLerISH4+uktWmzH8pJv",
		Manufacturer:   "Heidenhain",
		Mode:           "SignAndEncrypt",
		Model:          "TNC640",
		Policy:         "Basic256Sha256",
		Timeout:        3,
	}
	connResp, _, err := api.CreateConnection(ctx, &createReq)
	if err != nil {
		t.Fatalf("Ошибка создания подключения: %v", err)
	}
	if connResp.Status != "success" || connResp.Data.UUID == "" {
		t.Fatalf("Некорректный ответ при создании подключения: %+v", connResp)
	}
	UUID := connResp.Data.UUID
	log.Printf("Подключение создано успешно. UUID: %s\n", UUID)

	// 2. Получение списка подключений
	log.Println("Шаг 2: Получение списка всех подключений...")
	listResp, _, err := api.GetConnectionPool(ctx)
	if err != nil {
		t.Fatalf("Ошибка получения списка подключений: %v", err)
	}
	if listResp.Status != "success" || listResp.Data.PoolSize == 0 {
		t.Fatalf("Некорректный ответ при получении списка: %+v", listResp)
	}
	log.Printf("Получено %d активных подключений.\n", listResp.Data.PoolSize)

	// 3. Проверка состояния
	log.Println("Шаг 3: Проверка состояния подключения...")
	checkConnection := models.CheckConnectionRequest{UUID: UUID}
	checkResp, _, err := api.CheckConnection(ctx, &checkConnection)
	if err != nil {
		t.Fatalf("Ошибка проверки подключения: %v", err)
	}
	if checkResp.Data.Status != "HEALTHY" {
		t.Fatalf("Проверка показала нездоровое состояние: %s", checkResp.Status)
	}
	log.Printf("Состояние подключения: %s\n", checkResp.Status)

	// 4. Запуск опроса
	log.Println("Шаг 4: Запуск опроса данных...")
	pollReq := models.UUIDRequest{UUID: UUID}
	startMsg, _, err := api.StartPolling(ctx, &pollReq)
	if err != nil {
		t.Fatalf("Ошибка запуска опроса: %v", err)
	}
	log.Printf("Ответ сервера: %v\n", startMsg)

	// Даем опросу поработать
	log.Println("Ожидание 3 секунды, пока идет опрос...")
	time.Sleep(3 * time.Second)

	// 5. Остановка опроса
	log.Println("Шаг 5: Остановка опроса данных...")
	stopMsg, _, err := api.StopPolling(ctx, &pollReq)
	if err != nil {
		t.Fatalf("Ошибка остановки опроса: %v", err)
	}
	log.Printf("Ответ сервера: %v\n", stopMsg)

	// 6. Удаление подключения
	log.Println("Шаг 6: Удаление подключения...")
	deleteMsg, _, err := api.DeleteConnection(ctx, &pollReq)
	if err != nil {
		t.Fatalf("Ошибка удаления подключения: %v", err)
	}
	log.Printf("Ответ сервера: %v\n", deleteMsg)

	log.Println("Тест успешно завершен!")
}
