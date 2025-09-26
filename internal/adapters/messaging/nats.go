package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// NATSClient handles NATS messaging operations
type NATSClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewNATSClient creates a new NATS client
func NewNATSClient(url string) (*NATSClient, error) {
	// Connect to NATS server
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &NATSClient{
		conn: conn,
		js:   js,
	}, nil
}

// Publish publishes a message to a subject
func (c *NATSClient) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.conn.Publish(subject, payload)
}

// PublishWithReply publishes a message and waits for a reply
func (c *NATSClient) PublishWithReply(subject string, data interface{}, timeout time.Duration) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	msg, err := c.conn.Request(subject, payload, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to publish request: %w", err)
	}

	return msg.Data, nil
}

// Subscribe subscribes to a subject
func (c *NATSClient) Subscribe(subject string, handler func([]byte) error) error {
	_, err := c.conn.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	})
	return err
}

// QueueSubscribe subscribes to a subject with queue group
func (c *NATSClient) QueueSubscribe(subject, queue string, handler func([]byte) error) error {
	_, err := c.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	})
	return err
}

// CreateStream creates a JetStream stream
func (c *NATSClient) CreateStream(streamName string, subjects []string) error {
	stream, err := c.js.StreamInfo(streamName)
	if err != nil && err != nats.ErrStreamNotFound {
		return fmt.Errorf("failed to get stream info: %w", err)
	}

	if stream == nil {
		_, err = c.js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream: %w", err)
		}
	}

	return nil
}

// PublishToStream publishes a message to a JetStream stream
func (c *NATSClient) PublishToStream(streamName, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = c.js.PublishAsync(subject, payload)
	return err
}

// SubscribeToStream subscribes to a JetStream stream
func (c *NATSClient) SubscribeToStream(streamName, subject string, handler func([]byte) error) error {
	_, err := c.js.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			log.Printf("Error processing message: %v", err)
		}
		msg.Ack()
	})
	return err
}

// Close closes the NATS connection
func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// IsConnected returns true if the client is connected
func (c *NATSClient) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

// GetConnectionStats returns connection statistics
func (c *NATSClient) GetConnectionStats() nats.Statistics {
	return c.conn.Stats()
}
