package logger

import (
	"calendly/lib/web"
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var entry *logrus.Entry

const (
	DEBUG   = 0
	INFO    = 1
	WARNING = 2
	ERROR   = 3
	FATAL   = 4
)

func Init(mode int) {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Out = os.Stdout

	switch mode {
	case DEBUG:
		log.Level = logrus.DebugLevel
	case INFO:
		log.Level = logrus.InfoLevel
	case WARNING:
		log.Level = logrus.WarnLevel
	case ERROR:
		log.Level = logrus.ErrorLevel
	case FATAL:
		log.Level = logrus.FatalLevel
	}

	entry = logrus.NewEntry(log)
}

func Info(ctx context.Context, msg string, data map[string]any) {
	log(ctx, logrus.InfoLevel, msg, data)
}

func Debug(ctx context.Context, msg string, data map[string]any) {
	log(ctx, logrus.DebugLevel, msg, data)
}

func Warn(ctx context.Context, msg string, data map[string]any) {
	log(ctx, logrus.WarnLevel, msg, data)
}

func Error(ctx context.Context, msg string, data map[string]any) {
	log(ctx, logrus.ErrorLevel, msg, data)
}

func log(ctx context.Context, logLevel logrus.Level, msg string, data map[string]any) {
	entry.WithFields(buildFieldsFromContext(ctx, data)).Log(logLevel, msg)
}

func buildFieldsFromContext(ctx context.Context, data map[string]any) logrus.Fields {
	fields := logrus.Fields(map[string]any{"payload": data})
	if requestID := web.GetRequestIDFromContext(ctx); requestID != "" {
		fields["request_id"] = requestID
	}
	return fields
}
