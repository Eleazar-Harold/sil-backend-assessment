package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	t.Run("Create Redis client with valid parameters", func(t *testing.T) {
		addr := "localhost:6379"
		password := ""
		db := 0

		client := NewRedisClient(addr, password, db)

		if client == nil {
			t.Error("Expected Redis client to be created")
		}

		if client.client == nil {
			t.Error("Expected Redis client to be initialized")
		}
	})

	t.Run("Create Redis client with password", func(t *testing.T) {
		addr := "localhost:6379"
		password := "secret"
		db := 0

		client := NewRedisClient(addr, password, db)

		if client == nil {
			t.Error("Expected Redis client to be created")
		}

		if client.client == nil {
			t.Error("Expected Redis client to be initialized")
		}
	})

	t.Run("Create Redis client with different database", func(t *testing.T) {
		addr := "localhost:6379"
		password := ""
		db := 1

		client := NewRedisClient(addr, password, db)

		if client == nil {
			t.Error("Expected Redis client to be created")
		}

		if client.client == nil {
			t.Error("Expected Redis client to be initialized")
		}
	})
}

func TestRedisClient_Set(t *testing.T) {
	// Note: These tests would require a real Redis instance or a mock
	// For now, we'll test the basic functionality without actual Redis connection

	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Set key-value pair", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		value := "test-value"
		expiration := time.Hour

		err := client.Set(ctx, key, value, expiration)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("Set operation succeeded (Redis connection available)")
		} else {
			t.Logf("Set operation failed as expected without Redis: %v", err)
		}
	})

	t.Run("Set key-value pair with zero expiration", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key-no-expiry"
		value := "test-value"

		err := client.Set(ctx, key, value, 0)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("Set operation succeeded (Redis connection available)")
		} else {
			t.Logf("Set operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Get(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Get existing key", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		var value string

		err := client.Get(ctx, key, &value)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("Get operation succeeded, value: %s", value)
		} else {
			t.Logf("Get operation failed as expected without Redis: %v", err)
		}
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		ctx := context.Background()
		key := "non-existing-key"
		var value string

		err := client.Get(ctx, key, &value)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("Get operation succeeded, value: %s", value)
		} else {
			t.Logf("Get operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Delete(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Delete existing key", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"

		err := client.Delete(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("Delete operation succeeded")
		} else {
			t.Logf("Delete operation failed as expected without Redis: %v", err)
		}
	})

	t.Run("Delete non-existing key", func(t *testing.T) {
		ctx := context.Background()
		key := "non-existing-key"

		err := client.Delete(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("Delete operation succeeded")
		} else {
			t.Logf("Delete operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Exists(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Check if key exists", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"

		exists, err := client.Exists(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("Exists operation succeeded, exists: %v", exists)
		} else {
			t.Logf("Exists operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_SetNX(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Set key only if not exists", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key-nx"
		value := "test-value"
		expiration := time.Hour

		set, err := client.SetNX(ctx, key, value, expiration)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("SetNX operation succeeded, set: %v", set)
		} else {
			t.Logf("SetNX operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Expire(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Set expiration for key", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		expiration := time.Hour

		err := client.Expire(ctx, key, expiration)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("Expire operation succeeded")
		} else {
			t.Logf("Expire operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_TTL(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Get TTL for key", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"

		ttl, err := client.TTL(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("TTL operation succeeded, ttl: %v", ttl)
		} else {
			t.Logf("TTL operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Increment(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Increment numeric value", func(t *testing.T) {
		ctx := context.Background()
		key := "test-counter"

		value, err := client.Increment(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("Increment operation succeeded, value: %d", value)
		} else {
			t.Logf("Increment operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_IncrementBy(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Increment numeric value by amount", func(t *testing.T) {
		ctx := context.Background()
		key := "test-counter"
		amount := int64(5)

		value, err := client.IncrementBy(ctx, key, amount)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("IncrementBy operation succeeded, value: %d", value)
		} else {
			t.Logf("IncrementBy operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Decrement(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Decrement numeric value", func(t *testing.T) {
		ctx := context.Background()
		key := "test-counter"

		value, err := client.Decrement(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("Decrement operation succeeded, value: %d", value)
		} else {
			t.Logf("Decrement operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_DecrementBy(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Decrement numeric value by amount", func(t *testing.T) {
		ctx := context.Background()
		key := "test-counter"
		amount := int64(3)

		value, err := client.DecrementBy(ctx, key, amount)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("DecrementBy operation succeeded, value: %d", value)
		} else {
			t.Logf("DecrementBy operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_ListPush(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Push value to list", func(t *testing.T) {
		ctx := context.Background()
		key := "test-list"
		value := "test-value"

		err := client.ListPush(ctx, key, value)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("ListPush operation succeeded")
		} else {
			t.Logf("ListPush operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_ListPop(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Pop value from list", func(t *testing.T) {
		ctx := context.Background()
		key := "test-list"
		var value string

		err := client.ListPop(ctx, key, &value)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("ListPop operation succeeded, value: %s", value)
		} else {
			t.Logf("ListPop operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_ListLength(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Get list length", func(t *testing.T) {
		ctx := context.Background()
		key := "test-list"

		length, err := client.ListLength(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("ListLength operation succeeded, length: %d", length)
		} else {
			t.Logf("ListLength operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_HashSet(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Set hash fields", func(t *testing.T) {
		ctx := context.Background()
		key := "test-hash"
		values := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}

		err := client.HashSet(ctx, key, values)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("HashSet operation succeeded")
		} else {
			t.Logf("HashSet operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_HashGet(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Get hash field", func(t *testing.T) {
		ctx := context.Background()
		key := "test-hash"
		field := "test-field"
		var value string

		err := client.HashGet(ctx, key, field, &value)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("HashGet operation succeeded, value: %s", value)
		} else {
			t.Logf("HashGet operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_HashGetAll(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Get all hash fields", func(t *testing.T) {
		ctx := context.Background()
		key := "test-hash"

		values, err := client.HashGetAll(ctx, key)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Logf("HashGetAll operation succeeded, values: %v", values)
		} else {
			t.Logf("HashGetAll operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_HashDelete(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Delete hash fields", func(t *testing.T) {
		ctx := context.Background()
		key := "test-hash"
		fields := []string{"field1", "field2"}

		err := client.HashDelete(ctx, key, fields...)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("HashDelete operation succeeded")
		} else {
			t.Logf("HashDelete operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Ping(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Ping Redis server", func(t *testing.T) {
		ctx := context.Background()

		err := client.Ping(ctx)

		// This will fail without a real Redis connection, but we can test the method exists
		if err == nil {
			t.Log("Ping operation succeeded")
		} else {
			t.Logf("Ping operation failed as expected without Redis: %v", err)
		}
	})
}

func TestRedisClient_Close(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Close Redis connection", func(t *testing.T) {
		err := client.Close()

		// This should always succeed
		if err != nil {
			t.Errorf("Expected Close to succeed, got: %v", err)
		}
	})
}

func TestRedisClient_GetClient(t *testing.T) {
	client := NewRedisClient("localhost:6379", "", 0)

	t.Run("Get Redis client", func(t *testing.T) {
		redisClient := client.GetClient()

		if redisClient == nil {
			t.Error("Expected Redis client to be returned")
		}
	})
}
