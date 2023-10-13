package utils

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDeserializeWebsocketMessage(t *testing.T) {
	jsonData := []byte(`{"name":"qq","message":"bvb","HEADERS":{"HX-Request":"true","HX-Trigger":"guestbook-form","HX-Trigger-Name":null,"HX-Target":"guestbook-form","HX-Current-URL":"http://127.0.0.1:8080/"}}`)

	var message WebsocketMessage
	message.Data = make(map[string]string)

	var intermediate map[string]interface{}

	if err := json.Unmarshal(jsonData, &intermediate); err != nil {
		t.Fatalf("Error while deserializing: %v", err)
	}

	for k, v := range intermediate {
		if k == "HEADERS" {
			headersJSON, _ := json.Marshal(v)
			json.Unmarshal(headersJSON, &message.Headers)
		} else {
			if strValue, ok := v.(string); ok {
				message.Data[k] = strValue
			} else if v == nil {
				message.Data[k] = ""
			}
		}
	}

	expectedHeaders := WebsocketHeaders{
		HxRequest:     true,
		HxTrigger:     "guestbook-form",
		HxTarget:      "guestbook-form",
		HxTriggerName: "",
	}

	if !reflect.DeepEqual(message.Headers, expectedHeaders) {
		t.Errorf("Headers mismatch, expected: %+v, got: %+v", expectedHeaders, message.Headers)
	}

	expectedData := map[string]string{
		"name":    "qq",
		"message": "bvb",
	}

	if !reflect.DeepEqual(message.Data, expectedData) {
		t.Errorf("Data mismatch, expected: %+v, got: %+v", expectedData, message.Data)
	}
}
