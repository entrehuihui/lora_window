package testhandler

import "github.com/brocaar/lora-app-server/internal/handler"

// TestHandler implements a Handler for testing.
type TestHandler struct {
	SendDataUpChan               chan handler.DataUpPayload
	SendJoinNotificationChan     chan handler.JoinNotification
	SendACKNotificationChan      chan handler.ACKNotification
	SendErrorNotificationChan    chan handler.ErrorNotification
	DataDownPayloadChan          chan handler.DataDownPayload
	SendStatusNotificationChan   chan handler.StatusNotification
	SendLocationNotificationChan chan handler.LocationNotification
}

// NewTestHandler returns a TestHandler.
func NewTestHandler() *TestHandler {
	return &TestHandler{
		SendDataUpChan:               make(chan handler.DataUpPayload, 100),
		SendJoinNotificationChan:     make(chan handler.JoinNotification, 100),
		SendACKNotificationChan:      make(chan handler.ACKNotification, 100),
		SendErrorNotificationChan:    make(chan handler.ErrorNotification, 100),
		DataDownPayloadChan:          make(chan handler.DataDownPayload, 100),
		SendStatusNotificationChan:   make(chan handler.StatusNotification, 100),
		SendLocationNotificationChan: make(chan handler.LocationNotification, 100),
	}
}

// Close method.
func (t *TestHandler) Close() error {
	return nil
}

// SendDataUp method.
func (t *TestHandler) SendDataUp(payload handler.DataUpPayload) error {
	t.SendDataUpChan <- payload
	return nil
}

// SendJoinNotification Method.
func (t *TestHandler) SendJoinNotification(payload handler.JoinNotification) error {
	t.SendJoinNotificationChan <- payload
	return nil
}

// SendACKNotification method.
func (t *TestHandler) SendACKNotification(payload handler.ACKNotification) error {
	t.SendACKNotificationChan <- payload
	return nil
}

// SendErrorNotification method.
func (t *TestHandler) SendErrorNotification(payload handler.ErrorNotification) error {
	t.SendErrorNotificationChan <- payload
	return nil
}

// DataDownChan method.
func (t *TestHandler) DataDownChan() chan handler.DataDownPayload {
	return t.DataDownPayloadChan
}

// SendStatusNotification method.
func (t *TestHandler) SendStatusNotification(payload handler.StatusNotification) error {
	t.SendStatusNotificationChan <- payload
	return nil
}

// SendLocationNotification method.
func (t *TestHandler) SendLocationNotification(payload handler.LocationNotification) error {
	t.SendLocationNotificationChan <- payload
	return nil
}
