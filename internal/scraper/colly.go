package scraper

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// CollyScraper 基于 colly 的 HTML 爬虫
type CollyScraper struct {
	collector *colly.Collector
	proxies   []string
	index     int
}

// NewCollyScraper 创建 CollyScraper 实例
func NewCollyScraper() *CollyScraper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(3),
	)

	return &CollyScraper{
		collector: c,
		proxies:   []string{},
		index:     0,
	}
}

// SetProxies 设置代理池
func (c *CollyScraper) SetProxies(proxies []string) {
	c.proxies = proxies
}

// getNextProxy 轮换获取代理
func (c *CollyScraper) getNextProxy() string {
	if len(c.proxies) == 0 {
		return ""
	}
	proxy := c.proxies[c.index]
	c.index = (c.index + 1) % len(c.proxies)
	return proxy
}

// ScrapResult 爬取结果
type ScrapResult struct {
	URL     string
	Title   string
	Content string
	Price   string
	Rating  string
	Reviews string
	Seller  string
	RawData map[string]string
	Error   error
}

// ScrapePage 爬取单个页面
func (c *CollyScraper) ScrapePage(ctx context.Context, url string) (*ScrapResult, error) {
	result := &ScrapResult{
		URL:     url,
		RawData: make(map[string]string),
	}

	// 创建上下文控制的 collector
	collector := c.collector.Clone()

	// 设置代理轮换
	if proxy := c.getNextProxy(); proxy != "" {
		collector.SetProxy(proxy)
	}

	var scrapeErr error

	collector.OnHTML("title", func(e *colly.HTMLElement) {
		result.Title = e.Text
	})

	collector.OnHTML("span.a-price-whole", func(e *colly.HTMLElement) {
		result.Price = e.Text
	})

	collector.OnHTML("span.a-icon-alt", func(e *colly.HTMLElement) {
		if result.Rating == "" {
			result.Rating = e.Text
		}
	})

	collector.OnHTML("span.a-size-base", func(e *colly.HTMLElement) {
		if result.Reviews == "" {
			result.Reviews = e.Text
		}
	})

	collector.OnHTML("#sellerName", func(e *colly.HTMLElement) {
		result.Seller = e.Text
	})

	collector.OnHTML("div#productDescription", func(e *colly.HTMLElement) {
		result.Content = e.Text
	})

	collector.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("colly scrape error: %w", err)
	})

	err := collector.Visit(url)
	if err != nil {
		return result, fmt.Errorf("visit error: %w", err)
	}

	collector.Wait()

	if scrapeErr != nil {
		result.Error = scrapeErr
	}

	return result, scrapeErr
}

// ScrapeMultiple 批量爬取
func (c *CollyScraper) ScrapeMultiple(ctx context.Context, urls []string) []*ScrapResult {
	results := make([]*ScrapResult, 0, len(urls))

	for _, url := range urls {
		result, err := c.ScrapePage(ctx, url)
		if err != nil {
			result.Error = err
		}
		results = append(results, result)
	}

	return results
}

// ScrapeWithCallback 带回调的爬取
func (c *CollyScraper) ScrapeWithCallback(ctx context.Context, url string, cb func(*ScrapResult) error) error {
	result, err := c.ScrapePage(ctx, url)
	if err != nil {
		return err
	}
	return cb(result)
}

// Close 关闭爬虫
func (c *CollyScraper) Close() {
	c.collector.Wait()
}

// ParsePrice 解析价格字符串
func ParsePrice(priceStr string) (float64, error) {
	// 移除货币符号和逗号
	cleaned := strings.ReplaceAll(priceStr, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	var price float64
	_, err := fmt.Sscanf(cleaned, "%f", &price)
	if err != nil {
		return 0, fmt.Errorf("parse price error: %w", err)
	}

	return price, nil
}

// ParseRating 解析评分
func ParseRating(ratingStr string) (float64, error) {
	// 尝试解析 "4.5 out of 5" 格式
	var whole, tenth int
	rest := ""
	_, err := fmt.Sscanf(ratingStr, "%d.%d%s", &whole, &tenth, &rest)
	if err == nil && whole > 0 && tenth >= 0 {
		// 检查分数部分是否合理 (tenth 应该 < 10)
		if tenth < 10 {
			return float64(whole) + float64(tenth)/10, nil
		}
	}

	// 尝试直接解析浮点数
	var rating float64
	_, err = fmt.Sscanf(ratingStr, "%f", &rating)
	if err == nil && rating > 0 {
		return rating, nil
	}

	// 提取第一个合理的浮点数
	cleaned := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		return ' '
	}, ratingStr)

	parts := strings.Fields(cleaned)
	for _, part := range parts {
		if parsed, err := fmt.Sscanf(part, "%f", &rating); err == nil && parsed == 1 && rating > 0 && rating <= 5 {
			return rating, nil
		}
	}

	return 0, fmt.Errorf("parse rating error: could not parse rating from %q", ratingStr)
}

// ParseReviewCount 解析评论数
func ParseReviewCount(reviewStr string) (int, error) {
	// 移除逗号和非数字字符
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, reviewStr)

	var count int
	_, err := fmt.Sscanf(cleaned, "%d", &count)
	if err != nil {
		return 0, fmt.Errorf("parse review count error: %w", err)
	}

	return count, nil
}
