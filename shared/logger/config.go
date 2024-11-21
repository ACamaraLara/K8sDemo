package logger

import (
	"flag"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Contains the logger configuration.  This is meant to be a black box for clients of the library,
// that only need to "AddFlagsParams" in order to allow its configuration.
type LoggerConfig struct {
	LogLevel string
	BuffSize int
}

// Sets the required command-line parameters for the logger.
func AddFlagsParams(cfg *LoggerConfig) {
	flag.StringVar(&cfg.LogLevel, "logLevel", "INFO", "Log level register.")
	flag.IntVar(&cfg.BuffSize, "logBufSize", 1000, "Log queue size")
}

func SetLogLevel(logLevel string) {
	// Set log level.
	switch strings.ToLower(logLevel) {
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
	case "disabled": // No macro for "disabled".
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("Empty loglevel provided. Using INFO level by default.")
	default:
		panic("Wrong loglevel given: " + logLevel)
	}

	log.Info().Msgf("Log level set at %s", logLevel)
}
