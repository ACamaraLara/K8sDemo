package logger

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// This struct implements IO Writer interface to set logQueue as the output
// logs from zerolog lib.
type LoggerOutput struct {
	LogQueue chan []byte
}

// Implement Write function IO Writer interface. When calling ZeroLog
// login functions, this function will receive the message and send it
// to the corresponding Queue.
func (o *LoggerOutput) Write(p []byte) (n int, err error) {

	// Make a copy of the slice because slices are passed by reference.
	// This allows 'p' to be used by ZeroLog before the message is
	// extracted from LogQueue.
	pCpy := make([]byte, len(p))
	copy(pCpy, p)

	o.LogQueue <- pCpy

	return len(pCpy), nil
}

// Contains the logger configuration.  This is meant to be a black box for clients of the library,
// that only need to "AddFlagsParams" in order to allow its configuration.
type LoggerConfig struct {
	LogLevel string
}

// Sets the required command-line parameters for the logger.
func AddFlagsParams(cfg *LoggerConfig) {
	flag.StringVar(&cfg.LogLevel, "logLevel", "INFO", "Log level register.")
}

// InitServiceLogger initializes Logger service.
// @param logLevel to filter logger output.
// @param logOutput I/O writer to handle multilevel writer.
// It is used to store the logs and send them to the routine
// that will publish them in Loki.
func InitServiceLogger(cfg LoggerConfig, serviceName string, logOutput *LoggerOutput) error {

	if logOutput == nil {
		return fmt.Errorf("a Logger Output object must be provided")
	}

	// Inits logger global instance to use it all around the project.
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Str("service", serviceName).Logger()

	// Create a multi level writer to allow output logs in console and
	// in the declared byte array.
	writer := zerolog.MultiLevelWriter(os.Stdout, logOutput)

	// Establish multiple level writer as the logger writer.
	log.Logger = log.Output(writer)

	// Set log level.
	switch strings.ToLower(cfg.LogLevel) {
	case zerolog.LevelInfoValue:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case zerolog.LevelDebugValue:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case zerolog.LevelWarnValue:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case zerolog.LevelErrorValue:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case zerolog.LevelFatalValue:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case zerolog.LevelPanicValue:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("Empty loglevel provided. Using INFO level by default.")
	default:
		panic("Wrong loglevel given: " + cfg.LogLevel)
	}

	return nil
}
