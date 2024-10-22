package mongodb

import (
	"testing"

	"github.com/rs/zerolog"
)

// Init function executed before start tests to avoid verbose logs.
func init() {
	// Disables logger for unit testing.
	zerolog.SetGlobalLevel(zerolog.Disabled)

}

// Tests valid and invalid publication of json message to the broker.
// func TestInsertParkingJsonMessage(t *testing.T) {

// 	mongoDB := &MongoDBClient{DBWrapper: &MongoMock{}}

// 	type testPublish struct {
// 		Name          string
// 		JsonBody      []byte
// 		ExpectedError error
// 	}

// 	testCases := []testPublish{
// 		{
// 			Name:          "TestValidJson",
// 			JsonBody:      []byte(`{"name": "testValidJson"}`),
// 			ExpectedError: nil,
// 		},
// 		{
// 			Name:          "TestInvalidJson",
// 			JsonBody:      []byte(`{"name": "testInvalidJson"`),
// 			ExpectedError: fmt.Errorf("unmarshal json in ParkingMeterInfo error: unexpected end of JSON input"),
// 		},
// 	}

// 	for _, testCase := range testCases {

// 		t.Run(testCase.Name, func(t *testing.T) {
// 			err := mongoDB.InsertParkingJson(testCase.JsonBody)

// 			// Check if the expected error is obtained
// 			if err == nil && testCase.ExpectedError != nil {
// 				t.Fatalf("Expected an error, but got none.")
// 			} else if err != nil && testCase.ExpectedError == nil {
// 				t.Fatalf("Error not expected, but one given.: %s", err.Error())
// 			} else if err != nil && testCase.ExpectedError != nil && err.Error() != testCase.ExpectedError.Error() {
// 				t.Fatalf("Got a different error than expected. Expected: %s, Got: %s", testCase.ExpectedError.Error(), err.Error())
// 			}
// 		})

// 	}

// }

// Tests that connection func has the expected behavior.
func TestConnectionNotFail(t *testing.T) {

	mongoDB := &MongoDBClient{DBWrapper: &MongoMock{}}

	if err := mongoDB.ConnectMongoClient(); err != nil {
		t.Error("Expected none error but one given", err)
	}

	if mongoDB.Client == nil {
		t.Error("Connection shouldn't be a null pointer.")
	}

	if mongoDB.Collection == nil {
		t.Error("Channel shouldn't be a null pointer.")
	}
}

// Tests that close connection func has the expected behavior.
func TestCloseConnectionNotFail(t *testing.T) {

	mongoDB := &MongoDBClient{DBWrapper: &MongoMock{}}

	if err := mongoDB.DisconnectMongoClient(); err != nil {
		t.Error("Expected none error but one given", err)
	}

	if mongoDB.Client != nil {
		t.Error("Client should be a null pointer after Disconnect client.")
	}

	if mongoDB.Collection != nil {
		t.Error("Collection should be a null pointer.")
	}
}
