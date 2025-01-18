package test

import (
	"QueraMQ/server"
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBasicFunctionality(t *testing.T) {
	serverAddr := "127.0.0.1:8080"
	srv := server.NewServer(serverAddr)
	go func() {
		if err := srv.Run(); err != nil {
			t.Errorf("Server failed to run: %v", err)
		}
	}()
	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	go runConsumer(t, serverAddr, "test_topic", &wg)
	time.Sleep(100 * time.Millisecond)

	wg.Add(1)
	go runProducer(t, serverAddr, "test_topic", "Hello, World!", 1, &wg)

	wg.Wait()

	// Check queue
	_, exists := srv.GetTopic("test_topic")
	require.True(t, exists, "Failed to get topic")

	conn, err := net.Dial("tcp", serverAddr)
	require.NoError(t, err, "Producer failed to connect")
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	message := map[string]interface{}{"action": "shutdown"}
	encoder.Encode(message)
}

func runConsumer(t *testing.T, address, topic string, wg *sync.WaitGroup) {

	consumerConn, err := net.Dial("tcp", address)
	require.NoError(t, err, "Consumer failed to connect")
	defer consumerConn.Close()

	consumerEncoder := json.NewEncoder(consumerConn)
	consumerDecoder := json.NewDecoder(consumerConn)

	// Send subscribe request
	subscribeRequest := map[string]interface{}{
		"action": "subscribe",
		"topic":  topic,
	}
	err = consumerEncoder.Encode(subscribeRequest)
	require.NoError(t, err, "Failed to send subscribe request")

	// Read subscribe response
	var subscribeResponse map[string]interface{}
	err = consumerDecoder.Decode(&subscribeResponse)
	require.NoError(t, err, "Failed to read subscribe response")
	require.Equal(t, "ok", subscribeResponse["status"], "Subscribe failed")

	wg.Done()

	// Wait for message delivery
	var deliverMessage map[string]interface{}
	err = consumerDecoder.Decode(&deliverMessage)
	require.NoError(t, err, "Failed to read delivered message")
	require.Equal(t, "deliver", deliverMessage["action"], "Expected action 'deliver'")

	_, ok := deliverMessage["message"].(map[string]interface{})
	require.True(t, ok, "Invalid message format")

	message := map[string]interface{}{"action": "close_connection"}
	consumerEncoder.Encode(message)
}

func runProducer(t *testing.T, address, topic, content string, priority int, wg *sync.WaitGroup) {
	defer wg.Done()

	producerConn, err := net.Dial("tcp", address)
	require.NoError(t, err, "Producer failed to connect")
	defer producerConn.Close()

	producerEncoder := json.NewEncoder(producerConn)
	producerDecoder := json.NewDecoder(producerConn)

	// Send publish request
	publishRequest := map[string]interface{}{
		"action": "publish",
		"message": map[string]interface{}{
			"topic":    topic,
			"content":  content,
			"priority": priority,
		},
	}
	err = producerEncoder.Encode(publishRequest)
	require.NoError(t, err, "Failed to send publish request")

	// Read publish response
	var publishResponse map[string]interface{}
	err = producerDecoder.Decode(&publishResponse)
	require.NoError(t, err, "Failed to read publish response")
	require.Equal(t, "ok", publishResponse["status"], "Publish failed")

	message := map[string]interface{}{"action": "close_connection"}
	producerEncoder.Encode(message)
}
