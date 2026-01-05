package log

import (
	"fmt"
	"lab1/common"
	"os"
	"time"
)

// LogModule implements the Observer interface for logging workspace events
type LogModule struct {
	logFiles map[string]*os.File
}

// NewLogModule creates a new log module
func NewLogModule() *LogModule {
	return &LogModule{
		logFiles: make(map[string]*os.File),
	}
}

// Update implements the Observer interface
func (lm *LogModule) Update(event common.WorkspaceEvent) {
	// Get or create log file for this file path
	logPath := "." + event.FilePath + ".log"
	
	file, exists := lm.logFiles[event.FilePath]
	if !exists {
		var err error
		file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
			return
		}
		lm.logFiles[event.FilePath] = file
	}
	
	// Write log entry
	timestamp := time.UnixMilli(event.Timestamp).Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, event.Type, event.Command)
	
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
	}
}

// Close closes all open log files
func (lm *LogModule) Close() {
	for _, file := range lm.logFiles {
		file.Close()
	}
}
