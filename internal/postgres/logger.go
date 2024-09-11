package postgres

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
)

type Logger struct {
	Level zerolog.Level
	Log   zerolog.Logger
}

func NewLogLevel(level string) Logger {
	l := Logger{}

	switch level {
	case "info":
		l.Level = zerolog.InfoLevel
	case "warn":
		l.Level = zerolog.WarnLevel
	case "debug":
		l.Level = zerolog.DebugLevel
	case "error":
		l.Level = zerolog.ErrorLevel
	default:
		l.Level = zerolog.Disabled
	}

	l.Log = log.Level(l.Level)
	return l
}

// currently we don't use this function, because Level already defined at struct Logger
func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Error(ctx context.Context, msg string, opts ...interface{}) {
	ze := l.Log.Error()
	if ctx != nil {
		ze = ze.Ctx(ctx)
	}
	ze.Msgf(msg, opts...)
}

func (l Logger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	ze := l.Log.Warn()
	if ctx != nil {
		ze = ze.Ctx(ctx)
	}
	ze.Msgf(msg, opts...)
}

func (l Logger) Info(ctx context.Context, msg string, opts ...interface{}) {
	ze := l.Log.Info()
	if ctx != nil {
		ze = ze.Ctx(ctx)
	}
	ze.Msgf(msg, opts...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	if l.Level >= zerolog.Disabled {
		return
	}

	zl := log.Ctx(ctx)
	var ze *zerolog.Event = zl.WithLevel(l.Level)

	if err != nil {
		ze = zl.Err(err)
	}

	sql, rows := f()
	ze.Str("sql", sql).
		Int64("rows", rows).
		Str("time_elapsed", time.Since(begin).String()).
		Msg("database query")

	return
}
