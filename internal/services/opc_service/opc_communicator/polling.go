package opc_communicator

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
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

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopped polling for machine %s", connInfo.SessionID)
				return
			case <-ticker.C:
				nodes := connInfo.GetRelevantNodeIDs()
				if len(nodes) == 0 {
					log.Printf("No relevant nodes for machine %s", connInfo.SessionID)
					continue
				}

				data, err := o.ReadMachineData(id)
				if err != nil {
					log.Printf("Error polling machine %s: %v", connInfo.SessionID, err)
					continue
				}

				log.Printf("Polled data from machine %s:", connInfo.SessionID)
				log.Printf("%s", data.ToJSON())
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
	log.Printf("Polling manually stopped for machine %s", id)
	return nil
}
