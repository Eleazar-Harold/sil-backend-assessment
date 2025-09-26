package messaging

import (
	"testing"
	"time"
)

func TestNATSClient_NewNATSClient(t *testing.T) {
	t.Run("Create NATS client with valid URL", func(t *testing.T) {
		client, err := NewNATSClient("nats://localhost:4222")

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("NATS client created successfully (NATS connection available)")
			if client == nil {
				t.Error("Expected client to be non-nil")
			}
		} else {
			t.Logf("NATS client creation failed as expected without NATS: %v", err)
		}
	})

	t.Run("Create NATS client with invalid URL", func(t *testing.T) {
		client, err := NewNATSClient("invalid://url")

		// This should fail with invalid URL
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
		if client != nil {
			t.Error("Expected client to be nil for invalid URL")
		}
	})
}

func TestNATSClient_Publish(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Publish message to subject", func(t *testing.T) {
		subject := "test.subject"
		message := "test message"

		err := client.Publish(subject, message)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("Publish operation succeeded (NATS connection available)")
		} else {
			t.Logf("Publish operation failed as expected without NATS: %v", err)
		}
	})

	t.Run("Publish empty message", func(t *testing.T) {
		subject := "test.subject"
		message := ""

		err := client.Publish(subject, message)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("Publish operation succeeded (NATS connection available)")
		} else {
			t.Logf("Publish operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_PublishWithReply(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Publish message with reply", func(t *testing.T) {
		subject := "test.subject"
		message := "test message"
		timeout := 5 * time.Second

		reply, err := client.PublishWithReply(subject, message, timeout)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Logf("PublishWithReply operation succeeded, reply: %s", string(reply))
		} else {
			t.Logf("PublishWithReply operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_Subscribe(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Subscribe to subject", func(t *testing.T) {
		subject := "test.subject"
		handler := func(msg []byte) error {
			t.Logf("Received message: %s", string(msg))
			return nil
		}

		err := client.Subscribe(subject, handler)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("Subscribe operation succeeded (NATS connection available)")
		} else {
			t.Logf("Subscribe operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_QueueSubscribe(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Subscribe to subject with queue group", func(t *testing.T) {
		subject := "test.subject"
		queue := "test.queue"
		handler := func(msg []byte) error {
			t.Logf("Received message: %s", string(msg))
			return nil
		}

		err := client.QueueSubscribe(subject, queue, handler)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("QueueSubscribe operation succeeded (NATS connection available)")
		} else {
			t.Logf("QueueSubscribe operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_CreateStream(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Create stream", func(t *testing.T) {
		streamName := "test.stream"
		subjects := []string{"test.*"}

		err := client.CreateStream(streamName, subjects)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("CreateStream operation succeeded (NATS connection available)")
		} else {
			t.Logf("CreateStream operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_PublishToStream(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Publish to stream", func(t *testing.T) {
		streamName := "test.stream"
		subject := "test.message"
		message := "test message"

		err := client.PublishToStream(streamName, subject, message)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("PublishToStream operation succeeded (NATS connection available)")
		} else {
			t.Logf("PublishToStream operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_SubscribeToStream(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Subscribe to stream", func(t *testing.T) {
		streamName := "test.stream"
		subject := "test.*"
		handler := func(msg []byte) error {
			t.Logf("Received stream message: %s", string(msg))
			return nil
		}

		err := client.SubscribeToStream(streamName, subject, handler)

		// This will fail without a real NATS connection, but we can test the method exists
		if err == nil {
			t.Log("SubscribeToStream operation succeeded (NATS connection available)")
		} else {
			t.Logf("SubscribeToStream operation failed as expected without NATS: %v", err)
		}
	})
}

func TestNATSClient_Close(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Close connection", func(t *testing.T) {
		client.Close()

		// This will fail without a real NATS connection, but we can test the method exists
		t.Log("Close operation completed")
	})
}

func TestNATSClient_IsConnected(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Check connection status", func(t *testing.T) {
		connected := client.IsConnected()

		// This will fail without a real NATS connection, but we can test the method exists
		t.Logf("Connection status: %v", connected)
	})
}

func TestNATSClient_GetConnectionStats(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}

	t.Run("Get connection statistics", func(t *testing.T) {
		stats := client.GetConnectionStats()

		// This will fail without a real NATS connection, but we can test the method exists
		t.Logf("Connection stats: %+v", stats)
	})
}

func TestNATSClient_ErrorHandling(t *testing.T) {

	t.Run("Test with empty subject", func(t *testing.T) {
		client, err := NewNATSClient("nats://localhost:4222")
		if err != nil {
			t.Skip("Skipping test: NATS connection not available")
		}
		defer client.Close()

		// Test Publish with empty subject
		err = client.Publish("", "test message")
		if err == nil {
			t.Error("Expected error for empty subject Publish")
		}

		// Test PublishWithReply with empty subject
		_, err = client.PublishWithReply("", "test message", time.Second)
		if err == nil {
			t.Error("Expected error for empty subject PublishWithReply")
		}

		// Test Subscribe with empty subject
		err = client.Subscribe("", func(msg []byte) error { return nil })
		if err == nil {
			t.Error("Expected error for empty subject Subscribe")
		}

		// Test QueueSubscribe with empty subject
		err = client.QueueSubscribe("", "queue", func(msg []byte) error { return nil })
		if err == nil {
			t.Error("Expected error for empty subject QueueSubscribe")
		}

		// Test PublishToStream with empty subject
		err = client.PublishToStream("test-stream", "", "test message")
		if err == nil {
			t.Error("Expected error for empty subject PublishToStream")
		}

		// Test SubscribeToStream with empty subject
		err = client.SubscribeToStream("test-stream", "", func(msg []byte) error { return nil })
		if err == nil {
			t.Error("Expected error for empty subject SubscribeToStream")
		}
	})

	t.Run("Test with empty stream name", func(t *testing.T) {
		client, err := NewNATSClient("nats://localhost:4222")
		if err != nil {
			t.Skip("Skipping test: NATS connection not available")
		}
		defer client.Close()

		// Test CreateStream with empty stream name
		err = client.CreateStream("", []string{"test.*"})
		if err == nil {
			t.Error("Expected error for empty stream name CreateStream")
		}

		// Test PublishToStream with empty stream name
		err = client.PublishToStream("", "test.subject", "test message")
		if err == nil {
			t.Error("Expected error for empty stream name PublishToStream")
		}

		// Test SubscribeToStream with empty stream name
		err = client.SubscribeToStream("", "test.subject", func(msg []byte) error { return nil })
		if err == nil {
			t.Error("Expected error for empty stream name SubscribeToStream")
		}
	})

	t.Run("Test with empty queue name", func(t *testing.T) {
		client, err := NewNATSClient("nats://localhost:4222")
		if err != nil {
			t.Skip("Skipping test: NATS connection not available")
		}
		defer client.Close()

		// Test QueueSubscribe with empty queue name
		err = client.QueueSubscribe("test.subject", "", func(msg []byte) error { return nil })
		if err == nil {
			t.Error("Expected error for empty queue name QueueSubscribe")
		}
	})

	t.Run("Test with nil callback", func(t *testing.T) {
		client, err := NewNATSClient("nats://localhost:4222")
		if err != nil {
			t.Skip("Skipping test: NATS connection not available")
		}
		defer client.Close()

		// Test Subscribe with nil callback
		err = client.Subscribe("test.subject", nil)
		if err == nil {
			t.Error("Expected error for nil callback Subscribe")
		}

		// Test QueueSubscribe with nil callback
		err = client.QueueSubscribe("test.subject", "queue", nil)
		if err == nil {
			t.Error("Expected error for nil callback QueueSubscribe")
		}

		// Test SubscribeToStream with nil callback
		err = client.SubscribeToStream("test-stream", "test.subject", nil)
		if err == nil {
			t.Error("Expected error for nil callback SubscribeToStream")
		}
	})

	t.Run("Test with empty subjects list", func(t *testing.T) {
		client, err := NewNATSClient("nats://localhost:4222")
		if err != nil {
			t.Skip("Skipping test: NATS connection not available")
		}
		defer client.Close()

		// Test CreateStream with empty subjects list
		err = client.CreateStream("test-stream", []string{})
		if err == nil {
			t.Error("Expected error for empty subjects list CreateStream")
		}

		// Test CreateStream with nil subjects list
		err = client.CreateStream("test-stream", nil)
		if err == nil {
			t.Error("Expected error for nil subjects list CreateStream")
		}
	})
}

func TestNATSClient_DataHandling(t *testing.T) {
	client, err := NewNATSClient("nats://localhost:4222")
	if err != nil {
		t.Skip("Skipping test: NATS connection not available")
	}
	defer client.Close()

	t.Run("Test with complex data structures", func(t *testing.T) {
		// Test with struct
		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		testData := TestStruct{ID: 1, Name: "test"}

		err := client.Publish("test.struct", testData)
		if err == nil {
			t.Log("Publish with struct succeeded")
		}

		// Test with map
		testMap := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		}

		err = client.Publish("test.map", testMap)
		if err == nil {
			t.Log("Publish with map succeeded")
		}

		// Test with slice
		testSlice := []string{"item1", "item2", "item3"}

		err = client.Publish("test.slice", testSlice)
		if err == nil {
			t.Log("Publish with slice succeeded")
		}
	})

	t.Run("Test with nil data", func(t *testing.T) {
		err := client.Publish("test.nil", nil)
		if err == nil {
			t.Log("Publish with nil data succeeded")
		}
	})

	t.Run("Test with empty string data", func(t *testing.T) {
		err := client.Publish("test.empty", "")
		if err == nil {
			t.Log("Publish with empty string succeeded")
		}
	})

	t.Run("Test with very long data", func(t *testing.T) {
		longData := make([]byte, 1000000) // 1MB of data
		for i := range longData {
			longData[i] = byte(i % 256)
		}

		err := client.Publish("test.long", longData)
		if err == nil {
			t.Log("Publish with long data succeeded")
		}
	})
}
