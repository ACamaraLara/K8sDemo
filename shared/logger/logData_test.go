package logger

import (
	"testing"
	"time"
)

// Test deserialization of a normal log without extra fields.
func TestUnmarshalSimpleJsonNotFail(t *testing.T) {

	testLog := []byte("{\"level\":\"info\",\"service\":\"TestService\"" +
		",\"time\":\"2023-06-01T11:32:55+02:00\",\"message\":\"TestMessage.\"}\n")

	var testLogData LogData

	if err := testLogData.UnmarshalJSON(testLog); err != nil {
		t.Error("Unexpected error unmarshaling JSON -> ", err)
	}

	if testLogData.Level != "info" {
		t.Errorf("Expected log level = %s, but got = %s -> ",
			"info", testLogData.Level)
	}

	if testLogData.Message != "TestMessage." {
		t.Errorf("Expected message = %s, but got = %s -> ",
			"TestMessage.", testLogData.Message)
	}

	timeTest, _ := time.Parse("2006-01-02T15:04:05-07:00", "2023-06-01T11:32:55+02:00")

	if testLogData.Time != timeTest {
		t.Errorf("Expected log level = %s, but got = %s -> ",
			timeTest.String(), testLogData.Time.String())
	}

	if len(testLogData.ExtraFields) != 0 {
		t.Errorf("Expected no extra fields but %d was given",
			len(testLogData.ExtraFields))
	}

}

// Test deserialization of a log with extra fields.
func TestUnmarshalExtraJsonNotFail(t *testing.T) {

	testLog := []byte("{\"level\":\"info\",\"service\":\"TestService\"" +
		",\"testKey\":\"testValue\",\"testKey2\":\"testValue2\",\"message\":\"TestMessage.\"}\n")

	var testLogData LogData

	if err := testLogData.UnmarshalJSON(testLog); err != nil {
		t.Error("Expected none error but one was given -> ", err)
	}

	if len(testLogData.ExtraFields) != 2 {
		t.Fatalf("Expected one extra field but %d was given",
			len(testLogData.ExtraFields))
	}

	valueVar, present := testLogData.ExtraFields["testKey"]

	if !present || valueVar != "testValue" {
		t.Fatalf("Obtained value (%s) is different than expected (%s)",
			valueVar, "testValue")
	}

	valueVar2, present := testLogData.ExtraFields["testKey2"]

	if !present || valueVar2 != "testValue2" {
		t.Fatalf("Obtained value (%s) is different than expected (%s)",
			valueVar, "testValue2")
	}
}

// Tests if a mandatory field is empty.
func TestUnmarshallEmptyLogLevelFail(t *testing.T) {
	testLog := []byte("{\"service\":\"TestService\"" +
		",\"time\":\"2023-06-01T11:32:55+02:00\",\"message\":\"TestMessage.\"}\n")

	var testLogData LogData

	if err := testLogData.UnmarshalJSON(testLog); err.Error() != "LogData struct cannot have empty fields" {
		t.Error("Obtained different error than expected. Obtained = ", err)
	}
}

// Tests if function fails with a bad formated Json.
func TestUnmarshallBadFormatedLogFail(t *testing.T) {
	testLog := []byte("This is not a Json")

	var testLogData LogData

	if err := testLogData.UnmarshalJSON(testLog); err.Error() != "invalid character 'T' looking for beginning of value" {
		t.Error("Obtained different error than expected. Obtained = ", err)
	}
}
