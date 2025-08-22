package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"opc_ua_service/internal/interfaces"
)

type PoolingUsecase struct {
	OpcService interfaces.OpcService

	poolingCtx context.Context
	cancelFunc context.CancelFunc
	interval   time.Duration
	latestData []interfaces.MachineData
}

func NewPoolingUsecase(s interfaces.OpcService, interval time.Duration) *PoolingUsecase {
	return &PoolingUsecase{
		OpcService: s,
		interval:   interval,
	}
}

// StartPooling запускает сбор данных со всех машин
func (u *PoolingUsecase) StartPooling() error {
	conns := u.OpcService.GetAllConnectionsInfo()
	if len(conns) == 0 {
		return fmt.Errorf("no active connections available")
	}

	if u.poolingCtx != nil {
		return fmt.Errorf("pooling already started")
	}

	ctx, cancel := context.WithCancel(context.Background())
	u.poolingCtx = ctx
	u.cancelFunc = cancel

	go func() {
		ticker := time.NewTicker(u.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var allData []interfaces.MachineData
				for _, connInfo := range conns {
					// Получаем список узлов, которые нужно читать для этой машины
					nodes := connInfo.GetRelevantNodeIDs()
					if len(nodes) == 0 {
						fmt.Printf("No relevant nodes for session %s\n", connInfo.SessionID)
						continue
					}

					// Читаем данные через коммуникатор
					data, err := u.OpcService.ReadMachineData(connInfo.SessionID)
					if err != nil {
						fmt.Printf("Error reading data from session %s: %v\n", connInfo.SessionID, err)
						continue
					}
					fmt.Printf("Reading data from session %s\n\ndata: %v", connInfo.SessionID, data)
					b, err := json.MarshalIndent(data, "", "  ")
					if err != nil {
						log.Printf("JSON marshal error: %v", err)
					}
					fmt.Println(string(b))
					allData = append(allData, data)
				}

				// Обновляем последний снимок данных
				u.latestData = allData
			}
		}
	}()

	return nil
}

// StopPooling останавливает сбор данных
func (u *PoolingUsecase) StopPooling() {
	if u.cancelFunc != nil {
		u.cancelFunc()
		u.poolingCtx = nil
		u.cancelFunc = nil
	}
}
