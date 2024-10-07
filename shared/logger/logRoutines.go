package logger

import (
	"fmt"
)

// Returns URL to post in Loki logging service.
// func getLokiPostUrl() (string, error) {
// 	lokiBaseUrl := config.GetEnvironWithDefault("LOKI_URL", "http://localhost:3100")
// 	fmt.Println("Loki URL:", lokiBaseUrl)
// 	return lokiBaseUrl + "/loki/api/v1/push", nil
// }

// StartLokiLogPublishRoutine starts routine to send logs coming
// from ZeroLog to LokiDB.
// @param logWritter I/O writer to handle multilevel writer.
// It is used to store the logs and send them to the routine
// that will publish them in Loki.
func StartLokiLogPublishRoutine(logWriter *LoggerOutput) error {

	// Declare url and http client to post logs to Loki db.
	// lokiPostURL, err := getLokiPostUrl()
	// if err != nil {
	// 	return fmt.Errorf("Error obtaining loki post URL " + err.Error())
	// }

	// lokiClient := http.DefaultClient

	// Infinite loop that listens to a LoqQueue channel for service logs.
	go func() {
		for log := range logWriter.LogQueue {

			// // Deserialize byte array in a LogData object.
			// var logData LogData
			fmt.Printf("log: %v\n", log)
			// if err := logData.UnmarshalJSON(log); err != nil {
			// 	fmt.Println("Error marshaling log entry to JSON:", err)
			// 	continue
			// }

			// // Convert logData structure in a byte array with Loki format.
			// logLoki, err := getLogInLokiFormat(&logData)
			// if err != nil {
			// 	fmt.Println("Error getting LokiLog structure. Error: ", err)
			// }

			// req, err := http.NewRequest("POST", lokiPostURL, bytes.NewBuffer(logLoki))

			// // Make HTTP request to Loki post URL.
			// if err != nil {
			// 	fmt.Println("Error creating Loki POST request. Error: ", err)
			// 	continue
			// }

			// // Establecer el encabezado "Content-Type" en la solicitud
			// req.Header.Set("Content-Type", "application/json")

			// resp, err := lokiClient.Do(req)

			// if err != nil {
			// 	fmt.Println("Error sending POST request to Loki server. Error: ", err)
			// 	continue
			// }

			// if resp.StatusCode != http.StatusNoContent {
			// 	fmt.Println("Unexpected status code from loki, status code: ", resp.StatusCode)
			// }

			// resp.Body.Close()

		}
	}()
	return nil
}

// Serialiizes a LogData struct in a []byte.
// func getLogInLokiFormat(logData *LogData) ([]byte, error) {

// 	// store LogData attributes in Loki specified format.
// 	logLoki := map[string]interface{}{
// 		"streams": []map[string]interface{}{
// 			{
// 				"stream": func() map[string]interface{} {
// 					stream := map[string]interface{}{
// 						"severity": logData.Level,
// 						"job":      logData.ServiceName,
// 					}

// 					// Add extra fields as key-value pairs.
// 					for key, value := range logData.ExtraFields {
// 						stream[key] = value
// 					}

// 					return stream
// 				}(),
// 				"values": [][]string{
// 					{strconv.FormatInt(time.Now().UnixNano(), 10), logData.Message},
// 				},
// 			},
// 		},
// 	}

// 	// Serialize new structure in a byte array.
// 	payload, err := json.Marshal(logLoki)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return payload, nil
// }
