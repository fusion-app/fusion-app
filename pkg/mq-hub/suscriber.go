package mqhub

// SubscriberInterface defines how to receive Event
type SubscriberInterface interface {
	Init(config interface{}) error
	SubscribeTo(topic string) (<-chan interface{}, error)
	ListResourceEvents(query interface{}) ([]interface{}, error)
	ListAppEvents(query interface{}) ([]interface{}, error)
}
