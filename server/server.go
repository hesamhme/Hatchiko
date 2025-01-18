package server

import (
	"encoding/json"
	"net"
	"sync"
)

type Server struct {
	Addr    string
	Topics  map[string]*Topic
	Clients map[string]net.Conn
	mu      sync.Mutex
}

func NewServer(address string) *Server {
	return &Server{
		Addr:    address,
		Topics:  make(map[string]*Topic),
		Clients: make(map[string]net.Conn),
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) Stop() {
	// Graceful shutdown logic
	for _, conn := range s.Clients {
		conn.Close()
	}
}

func (s *Server) GetTopic(topicName string) (*Topic, bool) {
	topic, exists := s.Topics[topicName]
	return topic, exists
}

func (s *Server) GetClientConnections() []net.Conn {
	var conns []net.Conn
	for _, conn := range s.Clients {
		conns = append(conns, conn)
	}
	return conns
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	var request map[string]interface{}

	for {
		if err := decoder.Decode(&request); err != nil {
			return
		}

		action, ok := request["action"].(string)
		if !ok {
			encoder.Encode(map[string]string{"error": "invalid action"})
			continue
		}

		switch action {
		case "publish":
			s.handlePublish(request, encoder)
		case "subscribe":
			s.handleSubscribe(request, encoder, conn)
		case "unsubscribe":
			s.handleUnsubscribe(request, encoder, conn)
		case "shutdown":
			s.handleShutdown()
			return
		case "close_connection":
			conn.Close()
			return
		default:
			encoder.Encode(map[string]string{"error": "unknown action"})
		}
	}
}

func (s *Server) handlePublish(request map[string]interface{}, encoder *json.Encoder) {
	message, ok := request["message"].(map[string]interface{})
	if !ok {
		encoder.Encode(map[string]string{"error": "message is required"})
		return
	}

	topicName, ok := message["topic"].(string)
	if !ok || topicName == "" {
		encoder.Encode(map[string]string{"error": "topic is required"})
		return
	}

	content, ok := message["content"].(string)
	if !ok || content == "" {
		encoder.Encode(map[string]string{"error": "message content is required"})
		return
	}

	priority, ok := message["priority"].(float64)
	if !ok {
		encoder.Encode(map[string]string{"error": "priority is required"})
		return
	}

	topic, exists := s.GetTopic(topicName)
	if !exists {
		topic = NewTopic(topicName)
		s.Topics[topicName] = topic
	}

	topic.Publish(content, int(priority))
	encoder.Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleSubscribe(request map[string]interface{}, encoder *json.Encoder, conn net.Conn) {
	topicName, ok := request["topic"].(string)
	if !ok || topicName == "" {
		encoder.Encode(map[string]string{"error": "topic is required"})
		return
	}

	topic, exists := s.GetTopic(topicName)
	if !exists {
		topic = NewTopic(topicName)
		s.Topics[topicName] = topic
	}

	clientID := conn.RemoteAddr().String()
	msgChan := topic.Subscribe(clientID)

	encoder.Encode(map[string]string{"status": "ok"})

	go func() {
		for message := range msgChan {
			encoder.Encode(map[string]interface{}{
				"action": "deliver",
				"message": map[string]interface{}{
					"message_id": message.ID.String(),
					"topic":      topicName,
					"content":    message.Content,
					"priority":   message.Priority,
				},
			})
		}
	}()
}

func (s *Server) handleUnsubscribe(request map[string]interface{}, encoder *json.Encoder, conn net.Conn) {
	topicName, ok := request["topic"].(string)
	if !ok || topicName == "" {
		encoder.Encode(map[string]string{"error": "topic is required"})
		return
	}

	topic, exists := s.GetTopic(topicName)
	if exists {
		clientID := conn.RemoteAddr().String()
		topic.Unsubscribe(clientID)
		encoder.Encode(map[string]string{"status": "ok"})
	}
}

func (s *Server) handleShutdown() {
	for _, conn := range s.Clients {
		conn.Close()
	}
}
