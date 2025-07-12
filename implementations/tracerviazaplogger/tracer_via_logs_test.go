//
//  tracer_via_logs_test.go
//  coreruntime_extenstions
//
//  Created by karim-w on 11/07/2025.
//

package tracerviazaplogger_test

import (
	"testing"
	"time"

	"github.com/karim-w/coreruntime/traces"
	"github.com/karim-w/coreruntime_extensions/implementations/tracerviazaplogger"
	"go.uber.org/zap"
)

func TestTracerImpl(t *testing.T) {
	zap.ReplaceGlobals(
		zap.NewExample(),
	)
	// Initialize the tracer
	tracer := tracerviazaplogger.Initialize()

	// Create a sample trace
	tx := traces.Span{
		OP:        0, // Example operation code
		TraceId:   "12345",
		SpanId:    "67890",
		Name:      "Test Span",
		Target:    "Test Target",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(100 * time.Millisecond),
		Meta:      map[string]interface{}{"key": "value"},
		HasFailed: false,
		Error:     nil,
	}

	// Trace the sample span
	err := tracer.Trace(tx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tracer.Close() // Ensure to close the tracer
}
