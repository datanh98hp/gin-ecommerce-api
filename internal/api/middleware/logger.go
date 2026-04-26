package middleware

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/dat19/gin-ecommerce-api/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var accessLogger *zap.Logger

func init() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("./logs", 0755); err != nil {
		panic(err)
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/requests.log",
		MaxSize:    10,   // megabytes
		MaxBackups: 5,    // number of old files to retain
		MaxAge:     15,   // days to retain old files
		Compress:   true, // whether to compress old files
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.InfoLevel,
	)

	accessLogger = zap.New(core)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// Inject query tracker into context
		queries := make([]string, 0)
		ctx := context.WithValue(c.Request.Context(), database.QueriesContextKey, &queries)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		statusCode := c.Writer.Status()
		header := c.Request.Header
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		// Log structured data
		accessLogger.Info("request",
			zap.String("method", method),
			zap.Any("header", header),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("body", string(body)),
			zap.Strings("queries", queries),
			zap.String("errors", c.Errors.String()),
			zap.Time("request_time", start),
		)
	}
}
