package message

import (
	"fmt"
	"testing"
)

type testHandler struct {
	message Message
}

func (h *testHandler) ContextName() string {
	return "TestingContext"
}

func (h *testHandler) MessageName() string {
	return "TestingMessage"
}

func (h *testHandler) Handle(msg Message) error {
	h.message = msg
	return nil
}

type failingTestHandler struct {
	message Message
}

func (h *failingTestHandler) ContextName() string {
	return "TestingContext"
}

func (h *failingTestHandler) MessageName() string {
	return "TestingMessage"
}

func (h *failingTestHandler) Handle(msg Message) error {
	h.message = msg
	return fmt.Errorf("Failed on purpose")
}

type testMessage struct {
	context string
	name    string
}

func (m *testMessage) Context() string {
	return m.context
}

func (m *testMessage) Name() string {
	return m.name
}

func Test_NewSimpleMessageBus_shouldDispatchToHandlerWithNameAndContext(t *testing.T) {
	handler := &testHandler{}
	var message Message = &testMessage{
		context: handler.ContextName(),
		name:    handler.MessageName(),
	}
	messageBus := NewSimpleMessageBus([]Handler{
		handler,
	})

	err := messageBus.Dispatch(message)

	if err != nil {
		t.Errorf("Should not return error")
	}

	if handler.message != message {
		t.Errorf("Message should be equal")
	}
}
func Test_NewSimpleMessageBus_shouldDispatchToMultipleHandlers(t *testing.T) {
	handler := &testHandler{}
	handlerTwo := &testHandler{}
	var message Message = &testMessage{
		context: handler.ContextName(),
		name:    handler.MessageName(),
	}
	messageBus := NewSimpleMessageBus([]Handler{
		handler,
		handlerTwo,
	})

	err := messageBus.Dispatch(message)

	if err != nil {
		t.Errorf("Should not return error")
	}

	if handler.message != message {
		t.Errorf("Message should be equal")
	}

	if handlerTwo.message != message {
		t.Errorf("Message should be equal")
	}
}
func Test_NewSimpleMessageBus_shoulNotDispatchToSecondHandlerIfFirstFails(t *testing.T) {
	handler := &failingTestHandler{}
	handlerTwo := &testHandler{}
	var message Message = &testMessage{
		context: handler.ContextName(),
		name:    handler.MessageName(),
	}
	messageBus := NewSimpleMessageBus([]Handler{
		handler,
		handlerTwo,
	})

	err := messageBus.Dispatch(message)

	if err == nil {
		t.Errorf("Should have returned error")
	}

	if handler.message != message {
		t.Errorf("Message should be equal")
	}

	if handlerTwo.message == message {
		t.Errorf("Message should not be set")
	}
}
func Test_NewSimpleMessageBus_shouldReturnErrorIfNoHandlerInContextFound(t *testing.T) {
	handler := &testHandler{}
	var message Message = &testMessage{
		context: "unhandledcontext",
		name:    handler.MessageName(),
	}
	messageBus := NewSimpleMessageBus([]Handler{
		handler,
	})

	err := messageBus.Dispatch(message)

	if err == nil {
		t.Errorf("Should return error")
	}
}
func Test_NewSimpleMessageBus_shouldReturnErrorIfNoHandlerWithMessageNameFound(t *testing.T) {
	handler := &testHandler{}
	var message Message = &testMessage{
		context: handler.MessageName(),
		name:    "unhandledname",
	}
	messageBus := NewSimpleMessageBus([]Handler{
		handler,
	})

	err := messageBus.Dispatch(message)

	if err == nil {
		t.Errorf("Should return error")
	}
}

func Test_NewSimpleMessageBus_shouldReturnErrorIfNoHandlerFound(t *testing.T) {
	handler := &testHandler{}
	var message Message = &testMessage{
		context: "unhandledcontext",
		name:    "unhandledname",
	}
	messageBus := NewSimpleMessageBus([]Handler{
		handler,
	})

	err := messageBus.Dispatch(message)

	if err == nil {
		t.Errorf("Should return error")
	}
}
