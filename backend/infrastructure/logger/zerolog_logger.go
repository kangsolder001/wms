package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"wms/pkg/logger"
)

type zerologLogger struct {
	log zerolog.Logger
}

func New(level, output string) logger.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	var writer = os.Stdout
	if output == "file" {
		f, err := os.OpenFile("wms.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			f = os.Stdout
		}
		writer = f
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zl := zerolog.New(writer).
		Level(lvl).
		With().
		Timestamp().
		CallerWithSkipFrameCount(2).
		Str("service", "wms").
		Logger()

	return &zerologLogger{log: zl}
}

func (l *zerologLogger) Debug(msg string, keysAndValues ...any) {
	l.log.Debug().Fields(toMap(keysAndValues)).Msg(msg)
}

func (l *zerologLogger) Info(msg string, keysAndValues ...any) {
	l.log.Info().Fields(toMap(keysAndValues)).Msg(msg)
}

func (l *zerologLogger) Warn(msg string, keysAndValues ...any) {
	l.log.Warn().Fields(toMap(keysAndValues)).Msg(msg)
}

func (l *zerologLogger) Error(msg string, keysAndValues ...any) {
	l.log.Error().Fields(toMap(keysAndValues)).Msg(msg)
}

func (l *zerologLogger) Fatal(msg string, keysAndValues ...any) {
	l.log.Fatal().Fields(toMap(keysAndValues)).Msg(msg)
}

func (l *zerologLogger) With(keysAndValues ...any) logger.Logger {
	return &zerologLogger{log: l.log.With().Fields(toMap(keysAndValues)).Logger()}
}

func (l *zerologLogger) WithField(key string, value any) logger.Logger {
	return &zerologLogger{log: l.log.With().Interface(key, value).Logger()}
}

func (l *zerologLogger) WithFields(fields map[string]any) logger.Logger {
	ctx := l.log.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &zerologLogger{log: ctx.Logger()}
}

func (l *zerologLogger) Debugf(template string, args ...any) {
	l.log.Debug().Msgf(template, args...)
}

func (l *zerologLogger) Infof(template string, args ...any) {
	l.log.Info().Msgf(template, args...)
}

func (l *zerologLogger) Warnf(template string, args ...any) {
	l.log.Warn().Msgf(template, args...)
}

func (l *zerologLogger) Errorf(template string, args ...any) {
	l.log.Error().Msgf(template, args...)
}

func (l *zerologLogger) Fatalf(template string, args ...any) {
	l.log.Fatal().Msgf(template, args...)
}

func toMap(keysAndValues []any) map[string]any {
	if len(keysAndValues) == 0 {
		return nil
	}
	m := make(map[string]any, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("non-string-key-%d", i)
		}
		m[key] = keysAndValues[i+1]
	}
	return m
}
