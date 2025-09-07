package opc_service

import (
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"opc_ua_service/internal/services/opc_service/cert_manager"
	"opc_ua_service/internal/services/opc_service/opc_communicator"
	"opc_ua_service/internal/services/opc_service/opc_connector"
)

type OpcService struct {
	interfaces.CertificateManagerService
	interfaces.OpcConnectorService
	interfaces.OpcCommunicatorService
}

func NewOpcService(producer interfaces.KafkaService, logger *logging.Logger) interfaces.OpcService {
	certManager := cert_manager.NewCertificateManager(logger)
	opcConnector := opc_connector.NewOpcConnector(certManager, logger)
	opcCommunicator := opc_communicator.NewOpcCommunicator(opcConnector, producer, logger)

	return OpcService{
		certManager,
		opcConnector,
		opcCommunicator,
	}
}
