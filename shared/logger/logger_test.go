package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Checks correct initialization of
func TestInitServiceLoggerNotFail(t *testing.T) {

	logger := InitServiceLogger(LoggerConfig{zerolog.LevelDebugValue, 1000})

	// There is one log just after set log level, that's the reason of spect 1 log at the queue.
	if len(logger.Buf.LogQueue) != 1 {
		t.Fatalf("The number of messages in logs queue should be 1 at init (%d)",
			len(logger.Buf.LogQueue))
	}

	log.Info().Msg("This is a test log")

	if len(logger.Buf.LogQueue) != 2 {
		t.Fatalf("Number of messages in logs queue is different than expected (%d)",
			len(logger.Buf.LogQueue))
	}

}
