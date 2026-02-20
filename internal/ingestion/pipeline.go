package ingestion

import (
	"LogForge/internal/model"
	"LogForge/internal/storage"
	"fmt"
)

// LogChannel is a buffered channel that queues incoming logs
var LogChannel = make(chan model.LogEntry, 1000)

// StartWorker launches a background goroutine that processes logs
func StartWorker() {
	go func() {
		for logEntry := range LogChannel {

			// Save to MongoDB
			storage.SaveLog(logEntry)

			fmt.Printf(
				"[STORED] %s | %s | %s\n",
				logEntry.Service,
				logEntry.Level,
				logEntry.Message,
			)
		}

	}()
}
