// @title OpcUaService API
// @version 1.0.0
// @description API для работы с протоколом OPC UA
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1
package main

import "opc_ua_service/internal/app"

func main() {
	app.New().Run()
}
