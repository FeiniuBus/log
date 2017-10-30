package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int8

const (
	Debug LogLevel = iota - 1
	Info
	Warn
	Error
	DPanic
	Panic
	Fatal
)

func (level LogLevel) cast() zapcore.Level {
	switch level {
	case Debug:
		return zap.DebugLevel
	case Info:
		return zap.InfoLevel
	case Warn:
		return zap.WarnLevel
	case Error:
		return zap.ErrorLevel
	case DPanic:
		return zap.DPanicLevel
	case Panic:
		return zap.PanicLevel
	case Fatal:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func NewAt(level LogLevel) (*Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = CallerEncoder
	config.Level = zap.NewAtomicLevelAt(level.cast())

	_log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{log: _log}, nil
}

func NewLogstashAt(level LogLevel, host string, port int) (*Logger, error) {
	return NewLogstashWithTimeoutAt(level, host, port, 10)
}

func NewLogstashWithTimeoutAt(level LogLevel, host string, port, timeout int) (*Logger, error) {
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

	atom := zap.NewAtomicLevelAt(level.cast())
	_log := zap.New(zapcore.NewCore(enc, sink, atom), zap.AddCaller(), zap.AddStacktrace(atom))
	return &Logger{log: _log}, nil
}
