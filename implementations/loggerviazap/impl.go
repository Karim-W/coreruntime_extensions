//
//  impl.go
//  coreruntime_extenstions
//
//  Created by karim-w on 10/07/2025.
//

package loggerviazap

import (
	"os"

	"github.com/karim-w/coreruntime/runtimemodules"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOptions struct {
	IncludeCaller     bool
	IncludeStacktrace bool
	ShowPid           bool
	OutErrtoStdout    bool
}

type logger_impl struct {
	logger *zap.Logger
}

func (l *logger_impl) Error(msg string, args ...runtimemodules.LoggerField) {
	l.logger.Error(msg, _Assert_to_Zap_Fields(args...)...)
}

func (l *logger_impl) Info(msg string, args ...runtimemodules.LoggerField) {
	l.logger.Info(msg, _Assert_to_Zap_Fields(args...)...)
}

func (l *logger_impl) Warn(msg string, args ...runtimemodules.LoggerField) {
	l.logger.Warn(msg, _Assert_to_Zap_Fields(args...)...)
}

func (l *logger_impl) With(args ...runtimemodules.LoggerField) runtimemodules.Logger {
	return &logger_impl{
		logger: l.logger.With(_Assert_to_Zap_Fields(args...)...),
	}
}

func Initialize(opt LoggerOptions) runtimemodules.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()

	encoderCfg.TimeKey = "ts"

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     !opt.IncludeCaller,
		DisableStacktrace: !opt.IncludeStacktrace,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
	}

	if opt.ShowPid {
		config.InitialFields = map[string]any{
			"pid": os.Getpid(),
		}
	}

	if opt.OutErrtoStdout {
		config.ErrorOutputPaths = []string{
			"stdout",
		}
	} else {
		config.ErrorOutputPaths = []string{
			"stderr",
		}
	}

	l := zap.Must(config.Build())

	l.Info("Logger initialized")

	zap.ReplaceGlobals(l)

	return &logger_impl{
		logger: l,
	}
}

func _Assert_to_Zap_Fields(
	args ...runtimemodules.LoggerField,
) (res []zap.Field) {
	res = make([]zap.Field, 0, len(args))

	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}

	return
}
