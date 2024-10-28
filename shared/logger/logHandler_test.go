package logger

import (
	"testing"

	"github.com/rs/zerolog/log"
)

// Checks correct initialization of
func TestInitServiceLoggerNotFail(t *testing.T) {

	logWriter := &LoggerOutput{LogQueue: make(chan []byte, 1000)}

	if err := InitServiceLogger(LoggerConfig{"INFO"}, logWriter); err != nil {
		t.Fatal("Unexpected error initializing logger: ", err)
	}

	if len(logWriter.LogQueue) > 0 {
		t.Fatalf("The number of messages in logs queue should be 0 at init (%d)",
			len(logWriter.LogQueue))
	}

	log.Info().Msg("This is a test log")

	if len(logWriter.LogQueue) != 1 {
		t.Fatalf("Number of messages in logs queue is different than expected (%d)",
			len(logWriter.LogQueue))
	}

}

// Checks correct initialization of
func TestInitServiceLoggerFail(t *testing.T) {

	// Expected error got from file LogHandler.go in line 33.
	if err := InitServiceLogger(LoggerConfig{"INFO"}, nil); err.Error() != "a Logger Output object must be provided" {
		t.Fatalf("Unexpected error message. Expected-> %s. Got -> %s", "a Logger Output object must be provided", err.Error())
	}

}
