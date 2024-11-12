package logger

import (
	"github.com/rs/zerolog"
	"os"
)

type Logger struct {
	dbLogger     zerolog.Logger
	serverLogger zerolog.Logger
}

func NewLogger() *Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05", NoColor: false}

	dbLogger := zerolog.New(writer).With().
		Str("service", "database").
		Logger()

	serverLogger := zerolog.New(writer).With().
		Str("service", "server").
		Logger()

	return &Logger{
		dbLogger:     dbLogger,
		serverLogger: serverLogger,
	}
}

func (l *Logger) LogDBError(err error, place string) {
	l.dbLogger.Error().Err(err).Str("place:", place).Msg("DB query failed")
}

func (l *Logger) LogDBQuery(query string) {
	l.dbLogger.Info().Str("query", query).Msg("DB query executed")
}

func (l *Logger) LogServerRequest(method string, url string, status int) {
	l.serverLogger.Info().Str("method", method).Str("url", url).Int("status", status).Msg("HTTP request handled")
}

func (l *Logger) LogServerError(err error, place string, url string) {
	l.serverLogger.Error().Err(err).Str("place", place).Str("url", url).Msg("HTTP request failed")
}

func (l *Logger) LogServerAccess(method string, url string, status int, userID string, requestID string) {
	l.serverLogger.Info().
		Str("method", method).
		Str("url", url).
		Int("status", status).
		Str("user_id", userID).
		Str("request_id", requestID).
		Msg("Access log")
}
