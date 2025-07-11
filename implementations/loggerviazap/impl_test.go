//
//  impl_test.go
//  coreruntime_extenstions
//
//  Created by karim-w on 11/07/2025.
//

package loggerviazap_test

import (
	"testing"

	"github.com/karim-w/coreruntime/runtimemodules"
	"github.com/karim-w/coreruntime_extensions/implementations/loggerviazap"
	"go.uber.org/zap"
)

func TestLoggerImpl(t *testing.T) {
	// Initialize the logger with default options
	logger := loggerviazap.Initialize(loggerviazap.LoggerOptions{
		IncludeCaller:     true,
		IncludeStacktrace: true,
		ShowPid:           true,
		OutErrtoStdout:    false,
	})

	// Test logging methods
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Test with additional fields
	logger.With(runtimemodules.Field(
		"key1", "value1",
	)).Info("This is an info message with fields", runtimemodules.Field("key2", "value2"))

	zap.L().Sync() // Ensure all logs are flushed before the test ends
}
