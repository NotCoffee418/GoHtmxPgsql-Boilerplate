package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

type WebsocketHeaders struct {
	HxRequest     bool   `json:"Hx-Request,string"`
	HxTrigger     string `json:"Hx-Trigger"`
	HxTarget      string `json:"Hx-Target"`
	HxTriggerName string `json:"Hx-Trigger-Name"`
}

type WebsocketMessage struct {
	Data    map[string]string `json:"-"`
	Headers WebsocketHeaders  `json:"HEADERS"`
}

func DeserializeHtmxWebsocketMessage(jsonStr string) (*WebsocketMessage, error) {
	var msg WebsocketMessage
	msg.Data = make(map[string]string) // Initialize the Data map

	var intermediate map[string]interface{}

	if err := json.Unmarshal([]byte(jsonStr), &intermediate); err != nil {
		return nil, errors.New("failed to unmarshal json")
	}

	// Manually map relevant fields
	for k, v := range intermediate {
		if k == "HEADERS" {
			headersJSON, _ := json.Marshal(v)
			json.Unmarshal(headersJSON, &msg.Headers)
		} else {
			msg.Data[k] = fmt.Sprintf("%v", v)
		}
	}
	return &msg, nil
}
