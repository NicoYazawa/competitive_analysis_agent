package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProxyPool(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	pool := NewProxyPool(proxies)

	assert.Equal(t, 2, pool.Size())
	assert.Equal(t, 0, pool.FailedCount())
}

func TestProxyPool_GetNext(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080", "proxy3:8080"}
	pool := NewProxyPool(proxies)

	// 设置健康检查始终返回 true
	pool.SetHealthCheck(func(s string) bool { return true })

	// 测试轮换
	first := pool.GetNext()
	assert.Equal(t, "proxy1:8080", first)

	second := pool.GetNext()
	assert.Equal(t, "proxy2:8080", second)

	third := pool.GetNext()
	assert.Equal(t, "proxy3:8080", third)

	// 循环回第一个
	fourth := pool.GetNext()
	assert.Equal(t, "proxy1:8080", fourth)
}

func TestProxyPool_GetNext_AllFailed(t *testing.T) {
	proxies := []string{"proxy1:8080"}
	pool := NewProxyPool(proxies)

	// 健康检查始终失败
	pool.SetHealthCheck(func(s string) bool { return false })

	result := pool.GetNext()
	assert.Equal(t, "", result) // 无可用代理
}

func TestProxyPool_MarkFailed(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	pool := NewProxyPool(proxies)

	pool.MarkFailed("proxy1:8080")

	assert.Equal(t, 1, pool.Size())
	assert.Equal(t, 1, pool.FailedCount())

	// 获取应该返回 proxy2
	result := pool.GetNext()
	assert.Equal(t, "proxy2:8080", result)
}

func TestProxyPool_MarkSuccess(t *testing.T) {
	proxies := []string{"proxy1:8080"}
	pool := NewProxyPool(proxies)

	pool.MarkFailed("proxy1:8080")
	assert.Equal(t, 1, pool.FailedCount())

	pool.MarkSuccess("proxy1:8080")
	assert.Equal(t, 0, pool.FailedCount())
	assert.Equal(t, 1, pool.Size())
}

func TestProxyPool_Add(t *testing.T) {
	pool := NewProxyPool([]string{"proxy1:8080"})
	pool.Add("proxy2:8080", "proxy3:8080")

	assert.Equal(t, 3, pool.Size())
}

func TestProxyPool_Remove(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080", "proxy3:8080"}
	pool := NewProxyPool(proxies)

	pool.Remove("proxy2:8080")

	assert.Equal(t, 2, pool.Size())
	// GetNext should return one of the remaining proxies
	next := pool.GetNext()
	assert.NotEqual(t, "", next)
}

func TestProxyPool_Size(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	pool := NewProxyPool(proxies)

	assert.Equal(t, 2, pool.Size())

	pool.MarkFailed("proxy1:8080")
	assert.Equal(t, 1, pool.Size())
}

func TestProxyPool_FailedCount(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	pool := NewProxyPool(proxies)

	assert.Equal(t, 0, pool.FailedCount())

	pool.MarkFailed("proxy1:8080")
	assert.Equal(t, 1, pool.FailedCount())

	pool.MarkFailed("proxy2:8080")
	assert.Equal(t, 2, pool.FailedCount())
}

func TestProxyPool_Reset(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	pool := NewProxyPool(proxies)

	pool.MarkFailed("proxy1:8080")
	pool.MarkFailed("proxy2:8080")

	pool.Reset()

	assert.Equal(t, 2, pool.Size())
	assert.Equal(t, 0, pool.FailedCount())
}

func TestProxyPool_GetAll(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080", "proxy3:8080"}
	pool := NewProxyPool(proxies)

	pool.MarkFailed("proxy2:8080")

	all := pool.GetAll()
	assert.Len(t, all, 2)
	assert.Contains(t, all, "proxy1:8080")
	assert.Contains(t, all, "proxy3:8080")
	assert.NotContains(t, all, "proxy2:8080")
}

func TestProxyPool_GetNext_AfterRemoval(t *testing.T) {
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	pool := NewProxyPool(proxies)

	pool.SetHealthCheck(func(s string) bool { return true })

	// 获取第一个
	first := pool.GetNext()
	assert.Equal(t, "proxy1:8080", first)

	// 移除第一个
	pool.Remove("proxy1:8080")

	// 再次获取应该得到 proxy2
	second := pool.GetNext()
	assert.Equal(t, "proxy2:8080", second)
}

func TestDefaultHealthCheck_InvalidProxy(t *testing.T) {
	// 无效代理格式不应该 panic
	result := DefaultHealthCheck("invalid-proxy-format")
	assert.False(t, result)
}

func TestGetIP(t *testing.T) {
	tests := []struct {
		proxy    string
		expected string
	}{
		{"192.168.1.1:8080", "192.168.1.1"},
		{"proxy.example.com:8080", "proxy.example.com"},
		{"192.168.1.1", "192.168.1.1"}, // 无端口
	}

	for _, tt := range tests {
		t.Run(tt.proxy, func(t *testing.T) {
			result := getIP(tt.proxy)
			assert.Equal(t, tt.expected, result)
		})
	}
}
