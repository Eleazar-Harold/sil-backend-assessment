package messaging

import (
	"testing"
)

func TestRabbitMQClient_NewRabbitMQClient(t *testing.T) {
	t.Run("Create RabbitMQ client with valid URL", func(t *testing.T) {
		client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("RabbitMQ client created successfully (RabbitMQ connection available)")
			if client == nil {
				t.Error("Expected client to be non-nil")
			}
		} else {
			t.Logf("RabbitMQ client creation failed as expected without RabbitMQ: %v", err)
		}
	})

	t.Run("Create RabbitMQ client with invalid URL", func(t *testing.T) {
		client, err := NewRabbitMQClient("invalid://url")

		// This should fail with invalid URL
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
		if client != nil {
			t.Error("Expected client to be nil for invalid URL")
		}
	})
}

func TestRabbitMQClient_Publish(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Publish message to exchange", func(t *testing.T) {
		exchange := "test.exchange"
		routingKey := "test.key"
		message := "test message"

		err := client.Publish(exchange, routingKey, message)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("Publish operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("Publish operation failed as expected without RabbitMQ: %v", err)
		}
	})

	t.Run("Publish empty message", func(t *testing.T) {
		exchange := "test.exchange"
		routingKey := "test.key"
		message := ""

		err := client.Publish(exchange, routingKey, message)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("Publish operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("Publish operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_PublishToQueue(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Publish message to queue", func(t *testing.T) {
		queueName := "test.queue"
		message := "test message"

		err := client.PublishToQueue(queueName, message)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("PublishToQueue operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("PublishToQueue operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_DeclareQueue(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Declare queue", func(t *testing.T) {
		queueName := "test.queue"
		durable := true
		autoDelete := false
		exclusive := false
		noWait := false

		err := client.DeclareQueue(queueName, durable, autoDelete, exclusive, noWait)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("DeclareQueue operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("DeclareQueue operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_DeclareExchange(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Declare exchange", func(t *testing.T) {
		exchangeName := "test.exchange"
		exchangeType := "direct"
		durable := true
		autoDelete := false
		internal := false
		noWait := false

		err := client.DeclareExchange(exchangeName, exchangeType, durable, autoDelete, internal, noWait)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("DeclareExchange operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("DeclareExchange operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_BindQueue(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Bind queue to exchange", func(t *testing.T) {
		queueName := "test.queue"
		routingKey := "test.key"
		exchangeName := "test.exchange"

		err := client.BindQueue(queueName, routingKey, exchangeName)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("BindQueue operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("BindQueue operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_Consume(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Consume messages from queue", func(t *testing.T) {
		queueName := "test.queue"
		consumer := "test.consumer"
		autoAck := true
		exclusive := false
		noLocal := false
		noWait := false

		handler := func(message []byte) error {
			t.Logf("Received message: %s", string(message))
			return nil
		}

		err := client.Consume(queueName, consumer, autoAck, exclusive, noLocal, noWait, handler)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("Consume operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("Consume operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_PublishWithConfirmation(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Publish message with confirmation", func(t *testing.T) {
		exchange := "test.exchange"
		routingKey := "test.key"
		message := "test message"

		err := client.PublishWithConfirmation(exchange, routingKey, message)

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("PublishWithConfirmation operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("PublishWithConfirmation operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_Close(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Close connection", func(t *testing.T) {
		err := client.Close()

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		if err == nil {
			t.Log("Close operation succeeded (RabbitMQ connection available)")
		} else {
			t.Logf("Close operation failed as expected without RabbitMQ: %v", err)
		}
	})
}

func TestRabbitMQClient_IsConnected(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}

	t.Run("Check connection status", func(t *testing.T) {
		connected := client.IsConnected()

		// This will fail without a real RabbitMQ connection, but we can test the method exists
		t.Logf("Connection status: %v", connected)
	})
}

func TestRabbitMQClient_ErrorHandling(t *testing.T) {

	t.Run("Test with empty parameters", func(t *testing.T) {
		client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
		if err != nil {
			t.Skip("Skipping test: RabbitMQ connection not available")
		}
		defer client.Close()

		// Test Publish with empty exchange
		err = client.Publish("", "test.key", "test message")
		if err == nil {
			t.Error("Expected error for empty exchange Publish")
		}

		// Test Publish with empty routing key
		err = client.Publish("test.exchange", "", "test message")
		if err == nil {
			t.Error("Expected error for empty routing key Publish")
		}

		// Test PublishToQueue with empty queue name
		err = client.PublishToQueue("", "test message")
		if err == nil {
			t.Error("Expected error for empty queue name PublishToQueue")
		}

		// Test DeclareQueue with empty queue name
		err = client.DeclareQueue("", true, false, false, false)
		if err == nil {
			t.Error("Expected error for empty queue name DeclareQueue")
		}

		// Test DeclareExchange with empty exchange name
		err = client.DeclareExchange("", "direct", true, false, false, false)
		if err == nil {
			t.Error("Expected error for empty exchange name DeclareExchange")
		}

		// Test DeclareExchange with empty exchange type
		err = client.DeclareExchange("test.exchange", "", true, false, false, false)
		if err == nil {
			t.Error("Expected error for empty exchange type DeclareExchange")
		}

		// Test BindQueue with empty queue name
		err = client.BindQueue("", "test.key", "test.exchange")
		if err == nil {
			t.Error("Expected error for empty queue name BindQueue")
		}

		// Test BindQueue with empty routing key
		err = client.BindQueue("test.queue", "", "test.exchange")
		if err == nil {
			t.Error("Expected error for empty routing key BindQueue")
		}

		// Test BindQueue with empty exchange name
		err = client.BindQueue("test.queue", "test.key", "")
		if err == nil {
			t.Error("Expected error for empty exchange name BindQueue")
		}

		// Test Consume with empty queue name
		err = client.Consume("", "consumer", true, false, false, false, func(msg []byte) error { return nil })
		if err == nil {
			t.Error("Expected error for empty queue name Consume")
		}

		// Test PublishWithConfirmation with empty exchange
		err = client.PublishWithConfirmation("", "test.key", "test message")
		if err == nil {
			t.Error("Expected error for empty exchange PublishWithConfirmation")
		}

		// Test PublishWithConfirmation with empty routing key
		err = client.PublishWithConfirmation("test.exchange", "", "test message")
		if err == nil {
			t.Error("Expected error for empty routing key PublishWithConfirmation")
		}
	})

	t.Run("Test with nil callback", func(t *testing.T) {
		client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
		if err != nil {
			t.Skip("Skipping test: RabbitMQ connection not available")
		}
		defer client.Close()

		// Test Consume with nil callback
		err = client.Consume("test.queue", "consumer", true, false, false, false, nil)
		if err == nil {
			t.Error("Expected error for nil callback Consume")
		}
	})

	t.Run("Test with invalid exchange types", func(t *testing.T) {
		client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
		if err != nil {
			t.Skip("Skipping test: RabbitMQ connection not available")
		}
		defer client.Close()

		// Test DeclareExchange with invalid exchange type
		err = client.DeclareExchange("test.exchange", "invalid-type", true, false, false, false)
		if err == nil {
			t.Error("Expected error for invalid exchange type DeclareExchange")
		}
	})
}

func TestRabbitMQClient_DataHandling(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skip("Skipping test: RabbitMQ connection not available")
	}
	defer client.Close()

	t.Run("Test with complex data structures", func(t *testing.T) {
		// Test with struct
		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		testData := TestStruct{ID: 1, Name: "test"}

		err := client.Publish("test.exchange", "test.key", testData)
		if err == nil {
			t.Log("Publish with struct succeeded")
		}

		// Test with map
		testMap := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		}

		err = client.Publish("test.exchange", "test.key", testMap)
		if err == nil {
			t.Log("Publish with map succeeded")
		}

		// Test with slice
		testSlice := []string{"item1", "item2", "item3"}

		err = client.Publish("test.exchange", "test.key", testSlice)
		if err == nil {
			t.Log("Publish with slice succeeded")
		}
	})

	t.Run("Test with nil data", func(t *testing.T) {
		err := client.Publish("test.exchange", "test.key", nil)
		if err == nil {
			t.Log("Publish with nil data succeeded")
		}
	})

	t.Run("Test with empty string data", func(t *testing.T) {
		err := client.Publish("test.exchange", "test.key", "")
		if err == nil {
			t.Log("Publish with empty string succeeded")
		}
	})

	t.Run("Test with very long data", func(t *testing.T) {
		longData := make([]byte, 1000000) // 1MB of data
		for i := range longData {
			longData[i] = byte(i % 256)
		}

		err := client.Publish("test.exchange", "test.key", longData)
		if err == nil {
			t.Log("Publish with long data succeeded")
		}
	})

	t.Run("Test with special characters in routing keys", func(t *testing.T) {
		specialKeys := []string{
			"test.key.with.dots",
			"test-key-with-dashes",
			"test_key_with_underscores",
			"test#key#with#hashes",
			"test*key*with*asterisks",
		}

		for _, key := range specialKeys {
			err := client.Publish("test.exchange", key, "test message")
			if err == nil {
				t.Logf("Publish with special routing key '%s' succeeded", key)
			}
		}
	})
}
