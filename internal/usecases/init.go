package usecases

import (
	"opc_ua_service/internal/config"
	_ "opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/interfaces"
	_ "opc_ua_service/pkg/errors"
)

type UseCases struct {
	interfaces.ConnectionUsecase
	interfaces.PollingUsecase
}

func NewUsecases(r interfaces.Repository, s interfaces.OpcService, conf *config.Config) interfaces.Usecases {

	return &UseCases{
		NewConnectionUsecase(s, r),
		NewPollingUsecase(s, r),
	}

}
