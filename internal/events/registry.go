package events

type EventHandler func(event *Event) error

type EventRegistry interface {
	Register(eventType EventType, handler EventHandler)
	GetHandlers(eventType EventType) []EventHandler
}

type EventRegistryImpl struct {
	handlers map[EventType][]EventHandler
}

func NewEventRegistry() EventRegistry {
	return &EventRegistryImpl{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (r *EventRegistryImpl) Register(eventType EventType, handler EventHandler) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

func (r *EventRegistryImpl) GetHandlers(eventType EventType) []EventHandler {
	return r.handlers[eventType]
}
