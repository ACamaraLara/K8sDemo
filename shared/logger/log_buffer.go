package logger

// This struct implements IO Writer interface to set logQueue as the output
// logs from zerolog lib.
type LogBuffer struct {
	LogQueue chan []byte
}

// Implement Write function IO Writer interface. When calling ZeroLog
// login functions, this function will receive the message and send it
// to the corresponding Queue.
func (o *LogBuffer) Write(p []byte) (n int, err error) {

	// Make a copy of the slice because slices are passed by reference.
	// This allows 'p' to be used by ZeroLog before the message is
	// extracted from LogQueue.
	pCpy := make([]byte, len(p))
	copy(pCpy, p)

	o.LogQueue <- pCpy

	return len(pCpy), nil
}
