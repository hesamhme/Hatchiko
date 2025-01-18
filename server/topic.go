package server

import (
	"QueraMQ/queue"
	"github.com/google/uuid"
)

type Topic struct {
	Name        string
	MQ          queue.IMessageQueue
	Subscribers map[string]chan *queue.Message
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name:        name,
		MQ:          queue.NewMessageQueue(),
		Subscribers: make(map[string]chan *queue.Message),
	}
}

func (t *Topic) GetMessageQueue() *queue.MessageQueue {
	return t.MQ.(*queue.MessageQueue)
}

func (t *Topic) Subscribe(clientID string) chan *queue.Message {
	msgChan := make(chan *queue.Message)
	t.Subscribers[clientID] = msgChan
	return msgChan
}

func (t *Topic) Unsubscribe(clientID string) {
	delete(t.Subscribers, clientID)
}

func (t *Topic) Publish(content string, priority int) {
	message := &queue.Message{
		ID:       uuid.New(),
		Content:  content,
		Priority: priority,
	}
	t.MQ.Push(message)
	for _, subscriber := range t.Subscribers {
		subscriber <- message
	}
}
