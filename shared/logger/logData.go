package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// Structure to store Log information.
type LogData struct {
	Time        time.Time `json:"time"`
	ServiceName string    `json:"service"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	// This attribute is an interface because its used to that fields that
	// are not pressent in LogData struct. Whith this, you can add any type of field
	// to a log from ZeroLog library.
	// (example: log.Warn().Str("extraKey", "extraValue").Msg("Test warning with extra field.") )
	ExtraFields map[string]interface{} `json:"extra_attributes"`
}

// UnmarshalJSON gets byte array from the output of ZeroLog library and
// deserialize its information in a LogData object.
func (ld *LogData) UnmarshalJSON(data []byte) error {

	// Define Alias structure to unmarshall attributes.
	type Alias LogData

	// Declare aux structure to Unmarshal known and unknown fields.
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ld),
	}

	// Unmarshal known fields.
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	// Check if any mandatory field is empty (That means library bad behavior)
	if ld.Level == "" || ld.Message == "" || ld.ServiceName == "" {
		return fmt.Errorf("LogData struct cannot have empty fields")
	}

	// Unmarshall all fields in 'ExtraFields' map.
	err = json.Unmarshal(data, &aux.Alias.ExtraFields)
	if err != nil {
		return err
	}

	// Delete known fields from ExtraFields map to leave only unknown fields.
	delete(aux.Alias.ExtraFields, "time")
	delete(aux.Alias.ExtraFields, "service")
	delete(aux.Alias.ExtraFields, "level")
	delete(aux.Alias.ExtraFields, "message")

	return nil
}
