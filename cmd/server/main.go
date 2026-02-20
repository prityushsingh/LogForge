package main

import (
	"net/http"
	"strconv"
	"time"

	"LogForge/internal/ingestion"
	"LogForge/internal/model"
	"LogForge/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize DB and worker FIRST
	storage.InitMongo()
	ingestion.StartWorker()

	// -------------------------
	// Health endpoint
	// -------------------------
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "logforge",
		})
	})

	// -------------------------
	// POST /logs — ingest logs
	// -------------------------
	r.POST("/logs", func(c *gin.Context) {
		var logEntry model.LogEntry

		if err := c.ShouldBindJSON(&logEntry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid log format",
			})
			return
		}

		if logEntry.Level == "" || logEntry.Service == "" || logEntry.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing required fields",
			})
			return
		}

		// Send to ingestion pipeline
		ingestion.LogChannel <- logEntry

		c.JSON(http.StatusAccepted, gin.H{
			"status": "log accepted",
		})
	})

	// -------------------------
	// GET /logs — query logs
	// -------------------------
	r.GET("/logs", func(c *gin.Context) {

		service := c.Query("service")
		level := c.Query("level")

		limitStr := c.DefaultQuery("limit", "50")

		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 {
			limit = 50
		}

		// Parse time range
		fromStr := c.Query("from")
		toStr := c.Query("to")

		var fromTime, toTime time.Time

		if fromStr != "" {
			fromTime, _ = time.Parse(time.RFC3339, fromStr)
		}

		if toStr != "" {
			toTime, _ = time.Parse(time.RFC3339, toStr)
		}

		logs, err := storage.QueryLogs(service, level, limit, fromTime, toTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to query logs",
			})
			return
		}

		c.JSON(http.StatusOK, logs)
	})

	// -------------------------
	// Start server
	// -------------------------
	r.Run(":8080")
}
