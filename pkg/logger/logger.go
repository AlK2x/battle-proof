package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*AppLogger)(nil)

type Encoder string

const (
	Console Encoder = "console"
	JSON    Encoder = "json"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger interface {
	Log(msg string)
	Error(err error, args ...any)
	Warning(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Config struct {
	level   Level
	encoder Encoder
	writer  io.Writer
}

var mapErrorLevel = map[Level]zapcore.Level{
	DebugLevel: zap.DebugLevel,
	InfoLevel:  zap.InfoLevel,
	WarnLevel:  zap.WarnLevel,
	ErrorLevel: zap.ErrorLevel,
}

func NewLogger(config Config) Logger {
	level := mapErrorLevel[config.level]

	var encoder zapcore.Encoder
	if config.encoder == Console {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(config.writer), zap.NewAtomicLevelAt(level))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &AppLogger{
		logger:        logger,
		sugaredLogger: logger.Sugar(),
	}
}

type AppLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

func (l *AppLogger) WithName(name string) Logger {
	l.logger.Named(name)
	l.sugaredLogger.Named(name)
	return l
}

func (l *AppLogger) Log(msg string) {
	l.logger.Log(zap.InfoLevel, msg)
}

func (l *AppLogger) Error(err error, args ...any) {
	l.sugaredLogger.Errorw(err.Error(), args)
}

func (l *AppLogger) Warning(msg string, args ...any) {
	l.sugaredLogger.Warnw(msg, args)
}

func (l *AppLogger) Info(msg string, args ...any) {
	l.sugaredLogger.Infow(msg, args)
}

func (l *AppLogger) Debug(msg string, args ...any) {
	l.sugaredLogger.Debugw(msg, args)
}
