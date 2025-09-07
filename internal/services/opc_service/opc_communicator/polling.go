package opc_communicator

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "opc_ua_service/internal/domain/models"
	"opc_ua_service/pkg/errors"
	"time"
)

func (o *OpcCommunicator) StartPollingForMachine(id uuid.UUID) error {
	connInfo, err := o.connector.GetConnectionInfoByUUID(id)
	if err != nil {
		return err
	}
	if connInfo == nil {
		return errors.NewNotFoundError("connection not found")
	}

	o.mu.Lock()
	if _, exists := o.pollCancelMap[id]; exists {
		o.mu.Unlock()
		return fmt.Errorf("polling already started for machine %s", id)
	}

	ctx, cancel := context.WithCancel(context.Background())
	o.pollCancelMap[id] = cancel
	o.mu.Unlock()

	interval := connInfo.Config.Config.GetTimeout()
	connInfo.IsPolled = true
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				o.logger.Info("Stopped polling for machine %s", connInfo.SessionID)
				connInfo.IsPolled = false
				return
			case <-ticker.C:
				data, err := o.ReadMachineData(id)
				if err != nil {
					o.logger.Error("Error polling machine %s: %v", connInfo.SessionID, err)
					continue
				}
				dataResponse := data.ToResponse()
				dataJSON := data.ToJSON()
				err = o.producer.Produce(context.Background(), []byte(dataResponse.MachineId), []byte(dataJSON))
				if err != nil {
					o.logger.Error("Failed to send data to Kafka", "machineId", dataResponse.MachineId, "error", err)
				}
			}
		}
	}()

	return nil
}

func (o *OpcCommunicator) StopPollingForMachine(id uuid.UUID) error {
	o.mu.Lock()
	cancel, exists := o.pollCancelMap[id]
	if !exists {
		o.mu.Unlock()
		return fmt.Errorf("polling not active for machine %s", id)
	}
	delete(o.pollCancelMap, id)
	o.mu.Unlock()

	cancel()
	o.logger.Info("Polling manually stopped for machine %s", id)
	return nil
}
