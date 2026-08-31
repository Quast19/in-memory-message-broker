package broker

import "sync"

type Broker struct {
	mu     sync.Mutex
	topics map[string][]string
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string][]string),
	}
}

func (b *Broker) Publish(topic, message string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.topics[topic] = append(b.topics[topic], message)
}

func (b *Broker) Consume(topic string) (string, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	messages := b.topics[topic]
	if len(messages) == 0 {
		return "", false
	}

	msg := messages[0]
	b.topics[topic] = messages[1:]
	return msg, true
}

func (b *Broker) Count(topic string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.topics[topic])
}
