package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQClient handles RabbitMQ messaging operations
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish publishes a message to an exchange
func (c *RabbitMQClient) Publish(exchange, routingKey string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Timestamp:   time.Now(),
		},
	)
}

// PublishToQueue publishes a message directly to a queue
func (c *RabbitMQClient) PublishToQueue(queueName string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.channel.Publish(
		"",        // exchange (empty for direct queue)
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Timestamp:   time.Now(),
		},
	)
}

// DeclareQueue declares a queue
func (c *RabbitMQClient) DeclareQueue(queueName string, durable, autoDelete, exclusive, noWait bool) error {
	_, err := c.channel.QueueDeclare(
		queueName,  // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		nil,        // arguments
	)
	return err
}

// DeclareExchange declares an exchange
func (c *RabbitMQClient) DeclareExchange(exchangeName, exchangeType string, durable, autoDelete, internal, noWait bool) error {
	return c.channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		durable,      // durable
		autoDelete,   // auto-deleted
		internal,     // internal
		noWait,       // no-wait
		nil,          // arguments
	)
}

// BindQueue binds a queue to an exchange
func (c *RabbitMQClient) BindQueue(queueName, routingKey, exchangeName string) error {
	return c.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

// Consume starts consuming messages from a queue
func (c *RabbitMQClient) Consume(queueName, consumer string, autoAck, exclusive, noLocal, noWait bool, handler func([]byte) error) error {
	msgs, err := c.channel.Consume(
		queueName, // queue
		consumer,  // consumer
		autoAck,   // auto-ack
		exclusive, // exclusive
		noLocal,   // no-local
		noWait,    // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error processing message: %v", err)
				if !autoAck {
					msg.Nack(false, true) // Reject and requeue
				}
			} else if !autoAck {
				msg.Ack(false) // Acknowledge message
			}
		}
	}()

	return nil
}

// PublishWithConfirmation publishes a message and waits for confirmation
func (c *RabbitMQClient) PublishWithConfirmation(exchange, routingKey string, data interface{}) error {
	// Enable publisher confirms
	if err := c.channel.Confirm(false); err != nil {
		return fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = c.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Wait for confirmation using NotifyPublish
	confirms := c.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	select {
	case conf := <-confirms:
		if !conf.Ack {
			return fmt.Errorf("message not confirmed by broker")
		}
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for broker confirmation")
	}

	return nil
}

// Close closes the RabbitMQ connection
func (c *RabbitMQClient) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
			return err
		}
	}
	return nil
}

// IsConnected returns true if the client is connected
func (c *RabbitMQClient) IsConnected() bool {
	return c.conn != nil && !c.conn.IsClosed()
}
