package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/utils"
)

const lokiPostPath string = "/api/prom/push"

// Returns URL to post in Loki logging service.
func getLokiPostUrl() string {
	lokiBaseUrl := utils.GetEnvironWithDefault("LOKI_URL", "http://loki:3100")
	fmt.Println("Loki URL:", lokiBaseUrl)
	return lokiBaseUrl + lokiPostPath
}

// StartLokiLogPublishRoutine starts routine to send logs coming
// from ZeroLog to LokiDB.
// @param logWritter I/O writer to handle multilevel writer.
// It is used to store the logs and send them to the routine
// that will publish them in Loki.
func (log *Logger) StartLokiLogPublishRoutine() error {

	// Declare url and http client to post logs to Loki db.
	lokiPostURL := getLokiPostUrl()

	lokiClient := http.DefaultClient

	// Infinite loop that listens to a LoqQueue channel for service logs.
	go func() {
		for log := range log.Buf.LogQueue {

			// // Deserialize byte array in a LogData object.
			var logData LogData

			if err := logData.UnmarshalJSON(log); err != nil {
				fmt.Println("Error marshaling log entry to JSON:", err)
				continue
			}

			// Convert logData structure in a byte array with Loki format.
			logLoki, err := getLogInLokiFormat(&logData)
			if err != nil {
				fmt.Println("Error getting LokiLog structure. Error: ", err)
			}

			req, _ := http.NewRequest("POST", lokiPostURL, bytes.NewBuffer(logLoki))

			// Establecer el encabezado "Content-Type" en la solicitud
			req.Header.Set("Content-Type", "application/json")

			lokiClient.Do(req)
		}
	}()
	return nil
}

// Serialiizes a LogData struct in a []byte.
func getLogInLokiFormat(logData *LogData) ([]byte, error) {

	// store LogData attributes in Loki specified format.
	logLoki := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": func() map[string]interface{} {
					stream := map[string]interface{}{
						"severity": logData.Level,
					}

					// Add extra fields as key-value pairs.
					for key, value := range logData.ExtraFields {
						stream[key] = value
					}

					return stream
				}(),
				"values": []string{logData.Message},
			},
		},
	}

	// Serialize new structure in a byte array.
	payload, err := json.Marshal(logLoki)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
