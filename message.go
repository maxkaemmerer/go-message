package message

import "fmt"

type MessageBus interface {
	Dispatch(msg Message) error
}

type simpleMessageBus struct {
	handlerMap map[string][]Handler
}

func (mb *simpleMessageBus) Dispatch(msg Message) error {
	namespace := fmt.Sprintf("%s-%s", msg.Context(), msg.Name())
	handlers, ok := mb.handlerMap[namespace]
	if ok {
		for _, handler := range handlers {
			err := handler.Handle(msg)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("No handler found for message %s", namespace)
	}
	return nil
}

func NewSimpleMessageBus(handlers []Handler) MessageBus {
	handlerMap := make(map[string][]Handler)
	for _, handler := range handlers {
		namespace := fmt.Sprintf("%s-%s", handler.ContextName(), handler.MessageName())
		handlerMap[namespace] = append(handlerMap[namespace], handler)
	}
	return &simpleMessageBus{
		handlerMap: handlerMap,
	}
}

type Message interface {
	Name() string
	Context() string
}

type Handler interface {
	ContextName() string
	MessageName() string
	Handle(msg Message) error
}
