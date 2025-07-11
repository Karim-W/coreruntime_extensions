//
//  tracer_via_logs.go
//  coreruntime_extenstions
//
//  Created by karim-w on 11/07/2025.
//

package tracervialogs

import (
	"time"

	"github.com/karim-w/coreruntime/runtimemodules"
	"github.com/karim-w/coreruntime/traces"
	"go.uber.org/zap"
)

type tracer_impl struct{}

func (t *tracer_impl) Close() {
	zap.L().Info("Closing tracer")
	zap.L().Sync() // Ensure all logs are flushed
}

func (t *tracer_impl) Trace(tx traces.Span) error {
	zap.L().Info("Trace",
		zap.Any("operation", tx.OP),
		zap.String("trace_id", tx.TraceId),
		zap.String("span_id", tx.SpanId),
		zap.String("parent_span_id", tx.ParentId.GetOrElse("")),
		zap.String("name", tx.Name),
		zap.String("target", tx.Target),
		zap.String("invoked_at", tx.StartTime.Format(time.RFC3339)),
		zap.Int64("elapsed_ms", tx.EndTime.Sub(tx.StartTime).Milliseconds()),
		zap.Any("metadata", tx.Meta),
		zap.Bool("has_failed", tx.HasFailed),
		zap.Any("error", tx.Error),
	)

	return nil
}

var _ runtimemodules.Tracer = &tracer_impl{}
