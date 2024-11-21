package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	Conf LoggerConfig
	Buf  *LogBuffer
}

// InitServiceLogger initializes Logger service.
// @param logLevel to filter logger output.
// @param logOutput I/O writer to handle multilevel writer.
// It is used to store the logs and send them to the routine
// that will publish them in Loki.
func InitServiceLogger(cfg LoggerConfig) *Logger {

	// Create channel where zerolog will enqueue logs
	logBuf := &LogBuffer{LogQueue: make(chan []byte, cfg.BuffSize)}

	// Inits logger global instance to use it all around the project.
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create a multi level writer to allow output logs in console and
	// in the declared byte array.
	writer := zerolog.MultiLevelWriter(os.Stdout, logBuf)

	// Establish multiple level writer as the logger writer.
	log.Logger = log.Output(writer)

	SetLogLevel(cfg.LogLevel)

	return &Logger{Conf: cfg, Buf: logBuf}
}
