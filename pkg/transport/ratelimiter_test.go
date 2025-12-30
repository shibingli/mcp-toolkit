package transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRateLimiter(t *testing.T) {
	logger := zap.NewNop()
	rl := NewRateLimiter(10, time.Second, logger)

	assert.NotNil(t, rl)
	assert.Equal(t, 10, rl.maxRequests)
	assert.Equal(t, time.Second, rl.window)
	assert.NotNil(t, rl.requests)
}

func TestRateLimiter_Allow(t *testing.T) {
	logger := zap.NewNop()

	t.Run("allows requests within limit", func(t *testing.T) {
		rl := NewRateLimiter(5, time.Second, logger)
		clientID := "test-client-1"

		// 前5个请求应该被允许 / First 5 requests should be allowed
		for i := 0; i < 5; i++ {
			allowed := rl.Allow(clientID)
			assert.True(t, allowed, "request %d should be allowed", i+1)
		}

		// 第6个请求应该被拒绝 / 6th request should be denied
		allowed := rl.Allow(clientID)
		assert.False(t, allowed, "6th request should be denied")
	})

	t.Run("allows requests after window expires", func(t *testing.T) {
		rl := NewRateLimiter(2, 100*time.Millisecond, logger)
		clientID := "test-client-2"

		// 前2个请求应该被允许 / First 2 requests should be allowed
		assert.True(t, rl.Allow(clientID))
		assert.True(t, rl.Allow(clientID))

		// 第3个请求应该被拒绝 / 3rd request should be denied
		assert.False(t, rl.Allow(clientID))

		// 等待窗口过期 / Wait for window to expire
		time.Sleep(150 * time.Millisecond)

		// 现在应该允许新请求 / Now should allow new requests
		assert.True(t, rl.Allow(clientID))
	})

	t.Run("handles multiple clients independently", func(t *testing.T) {
		rl := NewRateLimiter(2, time.Second, logger)

		// 客户端1的请求 / Client 1 requests
		assert.True(t, rl.Allow("client-1"))
		assert.True(t, rl.Allow("client-1"))
		assert.False(t, rl.Allow("client-1"))

		// 客户端2应该有独立的限制 / Client 2 should have independent limit
		assert.True(t, rl.Allow("client-2"))
		assert.True(t, rl.Allow("client-2"))
		assert.False(t, rl.Allow("client-2"))
	})
}

func TestRateLimiter_Reset(t *testing.T) {
	logger := zap.NewNop()
	rl := NewRateLimiter(2, time.Second, logger)
	clientID := "test-client"

	// 用完限额 / Use up quota
	assert.True(t, rl.Allow(clientID))
	assert.True(t, rl.Allow(clientID))
	assert.False(t, rl.Allow(clientID))

	// 重置后应该可以再次请求 / After reset should allow requests again
	rl.Reset(clientID)
	assert.True(t, rl.Allow(clientID))
}

func TestRateLimiter_ResetAll(t *testing.T) {
	logger := zap.NewNop()
	rl := NewRateLimiter(1, time.Second, logger)

	// 多个客户端用完限额 / Multiple clients use up quota
	assert.True(t, rl.Allow("client-1"))
	assert.False(t, rl.Allow("client-1"))
	assert.True(t, rl.Allow("client-2"))
	assert.False(t, rl.Allow("client-2"))

	// 重置所有后应该都可以请求 / After reset all should allow requests
	rl.ResetAll()
	assert.True(t, rl.Allow("client-1"))
	assert.True(t, rl.Allow("client-2"))
}

func TestRateLimiter_GetStats(t *testing.T) {
	logger := zap.NewNop()
	rl := NewRateLimiter(5, time.Second, logger)

	// 发送一些请求 / Send some requests
	rl.Allow("client-1")
	rl.Allow("client-1")
	rl.Allow("client-2")

	stats := rl.GetStats()
	assert.Equal(t, 2, stats["client-1"])
	assert.Equal(t, 1, stats["client-2"])
}

func TestRateLimiter_Cleanup(t *testing.T) {
	logger := zap.NewNop()
	// 使用很短的窗口进行测试 / Use very short window for testing
	rl := NewRateLimiter(5, 50*time.Millisecond, logger)

	// 添加一些请求 / Add some requests
	rl.Allow("client-1")
	rl.Allow("client-2")

	// 验证客户端存在 / Verify clients exist
	stats := rl.GetStats()
	assert.Len(t, stats, 2)

	// 等待清理运行 / Wait for cleanup to run
	time.Sleep(200 * time.Millisecond)

	// 验证过期记录已被清理 / Verify expired records cleaned up
	rl.mu.RLock()
	clientCount := len(rl.requests)
	rl.mu.RUnlock()
	assert.Equal(t, 0, clientCount, "expired client records should be cleaned up")
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	logger := zap.NewNop()
	rl := NewRateLimiter(100, time.Second, logger)

	// 并发访问测试 / Concurrent access test
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				rl.Allow("concurrent-client")
			}
		}(i)
	}

	// 等待所有goroutine完成 / Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证没有panic / Verify no panic occurred
	stats := rl.GetStats()
	assert.NotNil(t, stats)
}
