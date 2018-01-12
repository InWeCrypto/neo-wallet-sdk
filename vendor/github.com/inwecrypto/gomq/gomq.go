package gomq

// Producer mq producer client
type Producer interface {
	Produce(topic string, key []byte, content interface{}) error
}

// Message MQ message
type Message interface {
	Key() []byte
	Topic() string
	Value() []byte
	Offset() int64
}

// Consumer mq consumer client
type Consumer interface {
	Close()
	Messages() <-chan Message
	Errors() <-chan error
	Commit(message Message)
}
