package transport

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// RateLimiter 请求频率限制器 / Request rate limiter
type RateLimiter struct {
	requests map[string]*clientRequests // 客户端请求记录 / Client request records
	mu       sync.RWMutex               // 保护requests / Protect requests
	logger   *zap.Logger
	// 配置 / Configuration
	maxRequests int           // 时间窗口内最大请求数 / Max requests in time window
	window      time.Duration // 时间窗口 / Time window
	cleanupTick time.Duration // 清理间隔 / Cleanup interval
}

// clientRequests 客户端请求记录 / Client request records
type clientRequests struct {
	timestamps []time.Time // 请求时间戳列表 / Request timestamp list
	mu         sync.Mutex  // 保护timestamps / Protect timestamps
}

// NewRateLimiter 创建频率限制器 / Create rate limiter
// maxRequests: 时间窗口内允许的最大请求数 / Max requests allowed in time window
// window: 时间窗口大小 / Time window size
func NewRateLimiter(maxRequests int, window time.Duration, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string]*clientRequests),
		logger:      logger,
		maxRequests: maxRequests,
		window:      window,
		cleanupTick: window * 2, // 清理间隔为窗口的2倍 / Cleanup interval is 2x window
	}

	// 启动清理goroutine / Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow 检查是否允许请求 / Check if request is allowed
// clientID: 客户端标识（通常是IP地址） / Client identifier (usually IP address)
// 返回true表示允许，false表示超过限制 / Returns true if allowed, false if rate limited
func (rl *RateLimiter) Allow(clientID string) bool {
	now := time.Now()

	rl.mu.Lock()
	cr, exists := rl.requests[clientID]
	if !exists {
		cr = &clientRequests{
			timestamps: make([]time.Time, 0, rl.maxRequests),
		}
		rl.requests[clientID] = cr
	}
	rl.mu.Unlock()

	cr.mu.Lock()
	defer cr.mu.Unlock()

	// 移除过期的时间戳 / Remove expired timestamps
	cutoff := now.Add(-rl.window)
	validTimestamps := make([]time.Time, 0, len(cr.timestamps))
	for _, ts := range cr.timestamps {
		if ts.After(cutoff) {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	cr.timestamps = validTimestamps

	// 检查是否超过限制 / Check if limit exceeded
	if len(cr.timestamps) >= rl.maxRequests {
		rl.logger.Warn("rate limit exceeded",
			zap.String("client_id", clientID),
			zap.Int("requests", len(cr.timestamps)),
			zap.Int("max_requests", rl.maxRequests),
			zap.Duration("window", rl.window))
		return false
	}

	// 记录新请求 / Record new request
	cr.timestamps = append(cr.timestamps, now)
	return true
}

// cleanup 定期清理过期的客户端记录 / Periodically cleanup expired client records
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupTick)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window * 2) // 清理2倍窗口之前的记录 / Cleanup records older than 2x window

		for clientID, cr := range rl.requests {
			cr.mu.Lock()
			// 如果所有时间戳都过期，删除该客户端记录 / If all timestamps expired, delete client record
			if len(cr.timestamps) == 0 || cr.timestamps[len(cr.timestamps)-1].Before(cutoff) {
				delete(rl.requests, clientID)
			}
			cr.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// GetStats 获取统计信息 / Get statistics
func (rl *RateLimiter) GetStats() map[string]int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	stats := make(map[string]int)
	for clientID, cr := range rl.requests {
		cr.mu.Lock()
		stats[clientID] = len(cr.timestamps)
		cr.mu.Unlock()
	}
	return stats
}

// Reset 重置指定客户端的限制 / Reset rate limit for specific client
func (rl *RateLimiter) Reset(clientID string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, clientID)
	rl.logger.Info("rate limit reset for client", zap.String("client_id", clientID))
}

// ResetAll 重置所有客户端的限制 / Reset rate limit for all clients
func (rl *RateLimiter) ResetAll() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.requests = make(map[string]*clientRequests)
	rl.logger.Info("rate limit reset for all clients")
}
