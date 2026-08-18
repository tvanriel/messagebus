package messagebus

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync"
)

type MessageBus struct {
	mu       sync.RWMutex
	handlers handlers
}

type handler struct {
	topic    Topic
	callback reflect.Value
	queue    chan []reflect.Value
}

type handlers []*handler

func New() *MessageBus {
	return &MessageBus{}
}

func buildHandlerArgs(args []any) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0, len(args))

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}

func (m *MessageBus) Subscribe(topic Topic, fn any) error {
	err := isValidHandler(fn)
	if err != nil {
		return err
	}

	h := &handler{
		topic:    topic,
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value),
	}

	go m.workQueue(h)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.handlers = append(m.handlers, h)

	return nil
}


var ErrUnsuitableHandler = errors.New("not a suitable handler")

func isValidHandler(fn any) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%w, expected reflect.Func, %s given", ErrUnsuitableHandler, reflect.TypeOf(fn))
	}

	return nil
}

func (m *MessageBus) Unsubscribe(topic Topic, fn any) error {
	err := isValidHandler(fn)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(fn)

	m.mu.Lock()
	defer m.mu.Unlock()

	for i, h := range m.handlers {
		if h.callback == rv && h.topic == topic {
			close(h.queue)

			m.handlers = slices.Delete(m.handlers, i, 1)
		}
	}

	return nil
}

func (m *MessageBus) Close(topic Topic) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, h := range m.handlers {
		if h.topic == topic {
			close(h.queue)

			m.handlers = slices.Delete(m.handlers, i, 1)
		}
	}
}

func (m *MessageBus) Publish(rk RouteKey, args ...any) {
	rArgs := buildHandlerArgs(args)

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, h := range m.handlers {
		if h.topic.Matches(rk) {
			h.queue <- rArgs
		}
	}
}

func (m *MessageBus) NumHandlers() int {
	return len(m.handlers)
}

func (m *MessageBus) workQueue(h *handler) {
	for args := range h.queue {
		h.callback.Call(args)
	}
}
