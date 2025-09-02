package opc_connector

/*
// RestoreConnection восстанавливает подключение из БД в пул памяти.
// Эта функция теперь всегда успешна, даже если эндпоинт недоступен,
// помечая такое соединение как IsHealthy: false.
func (oc *OpcConnector) RestoreConnection(machine entities.CncMachine) (*models.ConnectionInfo, error) {
	req := models.ConnectionRequest{
		EndpointURL:  machine.EndpointURL,
		Model:        machine.Model,
		Manufacturer: machine.Manufacturer,
	}

	if machine.ConnectionType == "certificate" {
		connection, err := oc.CreateCertificateConnection()
		if err != nil {
			return nil, err
		}
	}

	// Создаем базовый объект подключения, по умолчанию нездоровый
	connInfo := createConnectionInfo(machine.UUID, "unknown", req, machine.Manufacturer)
	connInfo.IsHealthy = false

	probeURL := strings.TrimSuffix(machine.EndpointURL, "/") + "/probe"
	xmlData, err := client.FetchXML(probeURL)
	if err != nil {
		oc.logger.Warn("Failed to get /probe for session. Connection will be restored as unhealthy.", "sessionID", machine.UUID, "error", err)
	} else {
		var devices models.MTConnectDevices
		if err := xml.Unmarshal(xmlData, &devices); err != nil {
			s.logger.Warn("Failed to parse /probe for session.", "sessionID", machine.UUID, "error", err)
		} else if len(devices.Devices) == 0 {
			s.logger.Warn("No devices found in /probe for session.", "sessionID", machine.UUID)
		} else {
			targetDevice, err := findTargetDevice(&devices, machine.Model, machine.Manufacturer, machine.EndpointURL)
			if err != nil {
				s.logger.Warn("Target device not found for session.", "sessionID", machine.UUID, "error", err)
			} else {
				// Все успешно, обновляем информацию и статус
				connInfo.MachineID = targetDevice.Name
				connInfo.Config.Manufacturer = targetDevice.Description.Manufacturer
				connInfo.IsHealthy = true
				if err := s.pollingMgr.LoadMetadataForEndpoint(machine.EndpointURL); err != nil {
					s.logger.Warn("Failed to load metadata for endpoint.", "endpoint", machine.EndpointURL, "error", err)
				}
			}
		}
	}

	oc.mu.Lock()
	defer oc.mu.Unlock()
	oc.connections[(machine.UUID).toString()] = connInfo

	return connInfo, nil
}

*/
