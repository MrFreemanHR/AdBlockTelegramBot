package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zaplog "adblock_bot/infrastructure/zap"
)

type VerbosityLevel uint

const (
	HighVerbosityLevel    = VerbosityLevel(3)
	NormalVerbosityLevel  = VerbosityLevel(2)
	MinimumVerbosityLevel = VerbosityLevel(1)
	ExtremeSilentLevel    = VerbosityLevel(0)
)

type EventType uint

const (
	UselessInfo = EventType(4)
	Info        = EventType(3)
	Warning     = EventType(2)
	Error       = EventType(1)
	Fatal       = EventType(0)
)

type coreLogger struct {
	ZapLogger       *zap.SugaredLogger
	currentLogLevel *VerbosityLevel
}

var currentLogger *coreLogger

func New(verbosityLevel VerbosityLevel) {
	var cfg = zap.NewProductionConfig()
	switch verbosityLevel {
	case HighVerbosityLevel, NormalVerbosityLevel:
		cfg.Level.SetLevel(zap.InfoLevel)
	case MinimumVerbosityLevel:
		cfg.Level.SetLevel(zap.WarnLevel)
	case ExtremeSilentLevel:
		cfg.Level.SetLevel(zap.PanicLevel)
	}
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	zapLogger := zaplog.New(true)
	currentLogger = &coreLogger{
		ZapLogger:       zapLogger,
		currentLogLevel: &verbosityLevel,
	}
}

func Logger() *coreLogger {
	return currentLogger
}

func (l *coreLogger) UselessInfo(msg string, opt ...any) {
	l.Event(UselessInfo, msg, opt...)
}

func (l *coreLogger) Info(msg string, opt ...any) {
	l.Event(Info, msg, opt...)
}

func (l *coreLogger) Warn(msg string, opt ...any) {
	l.Event(Warning, msg, opt...)
}

func (l *coreLogger) Error(msg string, opt ...any) {
	l.Event(Error, msg, opt...)
}

func (l *coreLogger) Fatal(msg string, opt ...any) {
	l.Event(Fatal, msg, opt...)
}

func (l *coreLogger) Event(lvl EventType, msg string, opt ...any) {
	var formattedMsg = fmt.Sprintf(msg, opt...)
	switch lvl {
	case UselessInfo:
		if *l.currentLogLevel == HighVerbosityLevel {
			l.ZapLogger.Info(formattedMsg)
		}
	case Info:
		l.ZapLogger.Info(formattedMsg)
	case Warning:
		l.ZapLogger.Warn(formattedMsg)
	case Error:
		l.ZapLogger.Error(formattedMsg)
	case Fatal:
		l.ZapLogger.Fatal(formattedMsg)
	}
	l.ZapLogger.Sync()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("01-02-06 15:04:05"))
}
