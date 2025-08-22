package opc_service

import (
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/services/opc_service/cert_manager"
	"opc_ua_service/internal/services/opc_service/opc_communicator"
	"opc_ua_service/internal/services/opc_service/opc_connector"
)

type OpcService struct {
	interfaces.CertificateManagerService
	interfaces.OpcConnectorService
	interfaces.OpcCommunicatorService
}

func NewOpcService() interfaces.OpcService {
	certManager := cert_manager.NewCertificateManager()
	opcConnector := opc_connector.NewOpcConnector(certManager)
	opcCommunicator := opc_communicator.NewOpcCommunicator()

	return OpcService{
		certManager,
		opcConnector,
		opcCommunicator,
	}
}
