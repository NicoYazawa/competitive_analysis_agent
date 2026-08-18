package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// BrowserEngine 浏览器引擎接口 (用于测试 Mock)
type BrowserEngine interface {
	Page() (*rod.Page, error)
	Close()
}

// RodScraper 基于 rod 的 JS 动态渲染爬虫
type RodScraper struct {
	browser  *rod.Browser
	pagePool []*rod.Page
	poolSize int
	proxies  []string
	index    int
}

// NewRodScraper 创建 RodScraper 实例
func NewRodScraper() (*RodScraper, error) {
	browser := rod.New()

	err := browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	return &RodScraper{
		browser:  browser,
		pagePool: make([]*rod.Page, 0),
		poolSize: 5,
		proxies:  []string{},
		index:    0,
	}, nil
}

// NewRodScraperWithBrowser 使用已有的 browser 创建 (便于测试)
func NewRodScraperWithBrowser(browser *rod.Browser) *RodScraper {
	return &RodScraper{
		browser:  browser,
		pagePool: make([]*rod.Page, 0),
		poolSize: 5,
		proxies:  []string{},
		index:    0,
	}
}

// SetPoolSize 设置页面池大小
func (r *RodScraper) SetPoolSize(size int) {
	r.poolSize = size
}

// SetProxies 设置代理池
func (r *RodScraper) SetProxies(proxies []string) {
	r.proxies = proxies
}

// getNextProxy 轮换获取代理
func (r *RodScraper) getNextProxy() string {
	if len(r.proxies) == 0 {
		return ""
	}
	proxy := r.proxies[r.index]
	r.index = (r.index + 1) % len(r.proxies)
	return proxy
}

// getPage 获取或创建页面
func (r *RodScraper) getPage(ctx context.Context) (*rod.Page, error) {
	if len(r.pagePool) > 0 {
		page := r.pagePool[len(r.pagePool)-1]
		r.pagePool = r.pagePool[:len(r.pagePool)-1]
		return page, nil
	}

	page, err := r.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return page, nil
}

// releasePage 归还页面到池中
func (r *RodScraper) releasePage(page *rod.Page) {
	if len(r.pagePool) < r.poolSize {
		r.pagePool = append(r.pagePool, page)
	} else {
		page.Close()
	}
}

// RodScrapResult Rod 爬取结果
type RodScrapResult struct {
	URL     string
	Title   string
	Price   string
	Rating  string
	Reviews string
	Seller  string
	Content string
	RawData map[string]string
	Error   error
	RawHTML string
}

// ScrapePage 爬取单个 JS 渲染页面
func (r *RodScraper) ScrapePage(ctx context.Context, urlStr string) (*RodScrapResult, error) {
	result := &RodScrapResult{
		URL:    urlStr,
		RawData: make(map[string]string),
	}

	page, err := r.getPage(ctx)
	if err != nil {
		result.Error = err
		return result, err
	}
	defer r.releasePage(page)

	// 导航到目标页面
	err = page.Timeout(defaultTimeout).Navigate(urlStr)
	if err != nil {
		result.Error = fmt.Errorf("navigation error: %w", err)
		return result, result.Error
	}

	// 等待页面加载完成
	err = page.WaitLoad()
	if err != nil {
		result.Error = fmt.Errorf("wait load error: %w", err)
		return result, result.Error
	}

	// 提取标题
	if titleEl := page.MustElement("title"); titleEl != nil {
		result.Title, _ = titleEl.Text()
	}

	// 提取价格
	if priceEl := page.MustElement("span.a-price-whole"); priceEl != nil {
		result.Price, _ = priceEl.Text()
	}

	// 提取评分
	if ratingEl := page.MustElement("span.a-icon-alt"); ratingEl != nil {
		result.Rating, _ = ratingEl.Text()
	}

	// 提取评论数
	if reviewEl := page.MustElement("span.a-size-base"); reviewEl != nil {
		result.Reviews, _ = reviewEl.Text()
	}

	// 提取卖家信息
	if sellerEl := page.MustElement("#sellerName"); sellerEl != nil {
		result.Seller, _ = sellerEl.Text()
	}

	// 提取商品描述
	if descEl := page.MustElement("#productDescription"); descEl != nil {
		result.Content, _ = descEl.Text()
	}

	// 获取完整 HTML
	result.RawHTML, _ = page.HTML()

	return result, nil
}

// ScrapeMultiple 批量爬取
func (r *RodScraper) ScrapeMultiple(ctx context.Context, urls []string) []*RodScrapResult {
	results := make([]*RodScrapResult, 0, len(urls))

	for _, urlStr := range urls {
		result, err := r.ScrapePage(ctx, urlStr)
		if err != nil {
			result.Error = err
		}
		results = append(results, result)
	}

	return results
}

// ScrapeWithCallback 带回调的爬取
func (r *RodScraper) ScrapeWithCallback(ctx context.Context, urlStr string, cb func(*RodScrapResult) error) error {
	result, err := r.ScrapePage(ctx, urlStr)
	if err != nil {
		return err
	}
	return cb(result)
}

// ExecuteScript 在页面上执行 JavaScript
func (r *RodScraper) ExecuteScript(page *rod.Page, script string) (interface{}, error) {
	return page.Evaluate(&rod.EvalOptions{JS: script})
}

// GetPageHTML 获取页面完整 HTML
func (r *RodScraper) GetPageHTML(page *rod.Page) (string, error) {
	return page.HTML()
}

// Close 关闭爬虫
func (r *RodScraper) Close() {
	for _, page := range r.pagePool {
		page.Close()
	}
	if r.browser != nil {
		r.browser.Close()
	}
}

// Helper functions

// ExtractJSONFromPage 从页面提取 JSON 数据
func ExtractJSONFromPage(script string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(script), &result)
	if err != nil {
		return nil, fmt.Errorf("extract json error: %w", err)
	}
	return result, nil
}

// CleanText 清理文本
func CleanText(text string) string {
	// 移除多余空白
	text = strings.Join(strings.Fields(text), " ")
	// 移除特殊字符
	text = strings.TrimSpace(text)
	return text
}

// defaultTimeout 默认超时时间
const defaultTimeout = 30_000_000_000 // 30 seconds in nanoseconds
