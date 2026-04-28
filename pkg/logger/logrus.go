package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger é a interface para logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// LogrusLogger implementa a interface Logger usando logrus
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogrusLogger cria uma nova instância do logger
func NewLogrusLogger(level string) *LogrusLogger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	return &LogrusLogger{logger: log}
}

// Info registra uma mensagem de informação
func (l *LogrusLogger) Info(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Info(msg)
}

// Error registra uma mensagem de erro
func (l *LogrusLogger) Error(msg string, err error, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).WithError(err).Error(msg)
}

// Warn registra uma mensagem de aviso
func (l *LogrusLogger) Warn(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Warn(msg)
}

// Debug registra uma mensagem de debug
func (l *LogrusLogger) Debug(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Debug(msg)
}

func parseFields(fields ...interface{}) logrus.Fields {
	result := logrus.Fields{}
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			result[fields[i].(string)] = fields[i+1]
		}
	}
	return result
}
