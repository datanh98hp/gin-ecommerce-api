package database

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type contextKey string

const QueriesContextKey contextKey = "db_queries"

type GormLogger struct {
	logger.Interface
}

func NewGormLogger(base logger.Interface) *GormLogger {
	return &GormLogger{Interface: base}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()

	// Capture query if key exists in context
	if queries, ok := ctx.Value(QueriesContextKey).(*[]string); ok {
		*queries = append(*queries, sql)
	}

	l.Interface.Trace(ctx, begin, fc, err)
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Interface.Info(ctx, msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Interface.Warn(ctx, msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Interface.Error(ctx, msg, data...)
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{Interface: l.Interface.LogMode(level)}
}
