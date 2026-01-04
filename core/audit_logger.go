package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"text2sql-skill/config"
)

type AuditLogger struct {
	cfg      *config.Config
	logChan  chan *AuditEntry
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
}

type AuditEntry struct {
	Timestamp time.Time
	QueryID   string
	EventType string
	Data      map[string]interface{}
}

func NewAuditLogger(cfg *config.Config) *AuditLogger {
	logger := &AuditLogger{
		cfg:      cfg,
		logChan:  make(chan *AuditEntry, 1000),
		stopChan: make(chan struct{}),
	}

	if cfg.Audit.Enabled && cfg.Performance.AsyncProcessing {
		logger.wg.Add(1)
		go logger.processLogs()
	}

	if cfg.Audit.Storage.Type == "file" {
		os.MkdirAll(cfg.Audit.Storage.Path, 0755)
	}

	return logger
}

func (a *AuditLogger) LogEvent(queryID string, eventType string, data map[string]interface{}) {
	if !a.cfg.Audit.Enabled {
		return
	}

	entry := &AuditEntry{
		Timestamp: time.Now().UTC(),
		QueryID:   queryID,
		EventType: eventType,
		Data:      data,
	}

	if a.cfg.Performance.AsyncProcessing {
		select {
		case a.logChan <- entry:
		default:
			// Channel full, drop entry to prevent blocking
		}
	} else {
		a.processSingleEntry(entry)
	}
}

func (a *AuditLogger) processLogs() {
	defer a.wg.Done()

	for {
		select {
		case entry := <-a.logChan:
			a.processSingleEntry(entry)
		case <-a.stopChan:
			return
		}
	}
}

func (a *AuditLogger) processSingleEntry(entry *AuditEntry) {
	switch a.cfg.Audit.Storage.Type {
	case "file":
		a.writeToFile(entry)
	default:
		// For sqlite and memory storage types, use file storage as fallback
		a.writeToFile(entry)
	}
}

func (a *AuditLogger) writeToFile(entry *AuditEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()

	day := entry.Timestamp.Format("2006-01-02")
	filename := filepath.Join(a.cfg.Audit.Storage.Path, "audit_"+day+".log")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	line, _ := json.Marshal(entry)
	file.Write(append(line, '\n'))
}

func (a *AuditLogger) Close() {
	if a.cfg.Performance.AsyncProcessing {
		close(a.stopChan)
		a.wg.Wait()
	}
}
