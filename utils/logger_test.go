package utils

import (
	"fmt"
	"os"
	"testing"
)

// TestWriteLogToFile test write log to file
func TestWriteLogToFile(t *testing.T) {
	getEnv := os.Getenv("zap_level")
	fmt.Println(getEnv)
}
