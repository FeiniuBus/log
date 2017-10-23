package log

import (
	"log/syslog"
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.Join([]string{caller.TrimmedPath(), runtime.FuncForPC(caller.PC).Name()}, ":"))
}

func New(debugLevel bool) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(strings.Join([]string{caller.TrimmedPath(), runtime.FuncForPC(caller.PC).Name()}, ":"))
	}

	if debugLevel {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return config.Build()
}

func NewLogstash(debugLevel bool, host string, port int) (*zap.Logger, error) {
	return NewLogstashWithTimeout(debugLevel, host, port, 10)
}

func NewLogstashWithTimeout(debugLevel bool, host string, port int, timeout int) (*zap.Logger, error) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeCaller = CallerEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeDuration = zapcore.SecondsDurationEncoder
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	enc := zapcore.NewJSONEncoder(cfg)
	sink, err := NewUDPSyncer(host, port, timeout)
	if err != nil {
		return nil, err
	}

	var atom zap.AtomicLevel
	if debugLevel {
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zap.NewProductionConfig()

	logger := zap.New(zapcore.NewCore(enc, zapcore.Lock(sink), atom))
	return logger, nil
}

func NewSyslog(debugLevel bool, app string) (*zap.Logger, error) {
	enc := NewSyslogEncoder(SyslogEncoderConfig{
		EncoderConfig: zapcore.EncoderConfig{
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   CallerEncoder,
		},

		Facility: syslog.LOG_LOCAL0,
		Hostname: "localhost",
		PID:      os.Getpid(),
		App:      app,
	})

	sink, err := NewSyslogSyncer()
	if err != nil {
		return nil, err
	}

	var atom zap.AtomicLevel
	if debugLevel {
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return zap.New(zapcore.NewCore(enc, zapcore.Lock(sink), atom)), nil
}
