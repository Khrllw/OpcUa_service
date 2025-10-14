package producers

import "time"

// cleanerWorker удаляет из БД обработанные записи пакетами
func (p *KafkaProducer) cleanerWorker() {
	p.logger.Info("Cleaner worker started")
	ticker := time.NewTicker(p.publishTimeout) // можно реже, чем публикация
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for {
				ids, err := p.usecase.GetProcessedPollDataIDs(p.batchSize)
				if err != nil {
					p.logger.Warn("Failed to fetch processed IDs", "error", err)
					break
				}
				if len(ids) == 0 {
					break
				}

				if err := p.usecase.DeletePollDataBatch(ids); err != nil {
					p.logger.Warn("Failed to delete processed poll data", "error", err)
					break
				}
			}
		case <-p.stopChan:
			p.logger.Info("Cleaner worker stopped")
			return
		}
	}
}
