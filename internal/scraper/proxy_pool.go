package scraper

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// ProxyPool 代理池管理
type ProxyPool struct {
	proxies     []string
	index       int
	mu          sync.RWMutex
	failedProxies map[string]time.Time
	healthCheck func(string) bool
}

// NewProxyPool 创建代理池
func NewProxyPool(proxies []string) *ProxyPool {
	return &ProxyPool{
		proxies:      proxies,
		index:        0,
		failedProxies: make(map[string]time.Time),
	}
}

// SetHealthCheck 设置健康检查函数
func (p *ProxyPool) SetHealthCheck(check func(string) bool) {
	p.healthCheck = check
}

// GetNext 获取下一个可用代理
func (p *ProxyPool) GetNext() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 清理过期失败代理 (30分钟内失败的代理重新尝试)
	threshold := time.Now().Add(-30 * time.Minute)
	for proxy, failedTime := range p.failedProxies {
		if failedTime.Before(threshold) {
			delete(p.failedProxies, proxy)
		}
	}

	// 轮换查找可用代理
	start := p.index
	for i := 0; i < len(p.proxies); i++ {
		idx := (start + i) % len(p.proxies)
		proxy := p.proxies[idx]

		// 检查是否在失败列表中
		if _, isFailed := p.failedProxies[proxy]; isFailed {
			continue
		}

		// 健康检查
		if p.healthCheck != nil && !p.healthCheck(proxy) {
			p.failedProxies[proxy] = time.Now()
			continue
		}

		p.index = (idx + 1) % len(p.proxies)
		return proxy
	}

	// 所有代理都失败，返回空
	return ""
}

// GetAll 获取所有代理
func (p *ProxyPool) GetAll() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]string, 0, len(p.proxies))
	for _, proxy := range p.proxies {
		if _, isFailed := p.failedProxies[proxy]; !isFailed {
			result = append(result, proxy)
		}
	}
	return result
}

// MarkFailed 标记代理失败
func (p *ProxyPool) MarkFailed(proxy string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failedProxies[proxy] = time.Now()
}

// MarkSuccess 标记代理成功
func (p *ProxyPool) MarkSuccess(proxy string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.failedProxies, proxy)
}

// Add 添加代理
func (p *ProxyPool) Add(proxies ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.proxies = append(p.proxies, proxies...)
}

// Remove 移除代理
func (p *ProxyPool) Remove(proxy string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, pxy := range p.proxies {
		if pxy == proxy {
			p.proxies = append(p.proxies[:i], p.proxies[i+1:]...)
			break
		}
	}
	delete(p.failedProxies, proxy)
}

// Size 获取代理池大小
func (p *ProxyPool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	count := 0
	for _, proxy := range p.proxies {
		if _, isFailed := p.failedProxies[proxy]; !isFailed {
			count++
		}
	}
	return count
}

// FailedCount 获取失败代理数量
func (p *ProxyPool) FailedCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.failedProxies)
}

// Reset 重置代理池
func (p *ProxyPool) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failedProxies = make(map[string]time.Time)
	p.index = 0
}

// DefaultHealthCheck 默认健康检查
func DefaultHealthCheck(proxy string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 解析代理地址
	proxyURL := fmt.Sprintf("http://%s", proxy)
	transport := &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) {
			return url.Parse(proxyURL)
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 5*time.Second)
		},
	}

	client := &http.Client{Transport: transport}

	// 发送 HEAD 请求到 Google 检测连通性
	req, _ := http.NewRequestWithContext(ctx, "HEAD", "https://www.google.com", nil)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500
}

// RoundRobinProxyHandler 轮换代理的 HTTP Handler
func RoundRobinProxyHandler(pool *ProxyPool, target *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := pool.GetNext()
		if proxy == "" {
			http.Error(w, "No available proxy", http.StatusServiceUnavailable)
			return
		}

		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Header.Set("X-Forwarded-For", getIP(proxy))
		}

		reverseProxy := &httputil.ReverseProxy{
			Director: director,
		}

		reverseProxy.ServeHTTP(w, r)
	})
}

func getIP(proxy string) string {
	host, _, err := net.SplitHostPort(proxy)
	if err != nil {
		return proxy
	}
	return host
}
