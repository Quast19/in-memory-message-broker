package broker

import "sync"

type Broker struct {
	mu     sync.Mutex
	topics map[string]chan string
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]chan string),
	}
}

func (b *Broker) getOrCreateTopic(topic string) chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch, exists := b.topics[topic]
	if !exists {
		ch = make(chan string, 100)
		b.topics[topic] = ch
	}
	return ch
}

func (b *Broker) Publish(topic, message string) bool {
	ch := b.getOrCreateTopic(topic)
	select {
	case ch <- message:
		return true
	default:
		return false
	}
}

func (b *Broker) Consume(topic string) (string, bool) {
	ch := b.getOrCreateTopic(topic)

	select {
	case msg := <-ch:
		return msg, true
	default:
		return "", false
	}
}

func (b *Broker) Count(topic string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.topics[topic])
}
