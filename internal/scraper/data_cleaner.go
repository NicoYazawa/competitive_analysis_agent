package scraper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// DataCleaner 数据清洗标准化
type DataCleaner struct{}

// NewDataCleaner 创建数据清洗器
func NewDataCleaner() *DataCleaner {
	return &DataCleaner{}
}

// CompetitorData 竞品原始数据
type CompetitorData struct {
	Name              string `json:"name"`
	Platform          string `json:"platform"`
	PlatformProductID string `json:"platform_product_id"`
	Price            string `json:"price"`
	Currency         string `json:"currency"`
	Rating           string `json:"rating"`
	ReviewCount      string `json:"review_count"`
	SellerRating     string `json:"seller_rating"`
	SellerReviewCount string `json:"seller_review_count"`
	SourceURL        string `json:"source_url"`
}

// CleanCompetitorData 清洗竞品数据
func (d *DataCleaner) CleanCompetitorData(raw *CompetitorData) (*CleanedCompetitorData, error) {
	cleaned := &CleanedCompetitorData{
		SourceURL: strings.TrimSpace(raw.SourceURL),
		Currency:  "USD",
	}

	// 清洗名称
	cleaned.Name = d.CleanText(raw.Name)
	if cleaned.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// 清洗平台
	cleaned.Platform = d.CleanPlatform(raw.Platform)
	if cleaned.Platform == "" {
		return nil, fmt.Errorf("platform is required")
	}

	// 清洗平台产品ID
	cleaned.PlatformProductID = d.CleanText(raw.PlatformProductID)

	// 解析价格
	price, err := ParsePrice(raw.Price)
	if err != nil {
		cleaned.Price = 0
	} else {
		cleaned.Price = price
	}

	// 解析货币
	if raw.Currency != "" {
		cleaned.Currency = d.CleanCurrency(raw.Currency)
	}

	// 解析评分
	rating, err := ParseRating(raw.Rating)
	if err != nil {
		cleaned.Rating = 0
	} else {
		cleaned.Rating = rating
	}

	// 解析评论数
	reviews, err := ParseReviewCount(raw.ReviewCount)
	if err != nil {
		cleaned.ReviewCount = 0
	} else {
		cleaned.ReviewCount = reviews
	}

	// 解析卖家评分
	sellerRating, err := ParseRating(raw.SellerRating)
	if err != nil {
		cleaned.SellerRating = 0
	} else {
		cleaned.SellerRating = sellerRating
	}

	// 解析卖家评论数
	sellerReviews, err := ParseReviewCount(raw.SellerReviewCount)
	if err != nil {
		cleaned.SellerReviewCount = 0
	} else {
		cleaned.SellerReviewCount = sellerReviews
	}

	return cleaned, nil
}

// CleanedCompetitorData 清洗后的竞品数据
type CleanedCompetitorData struct {
	Name              string  `json:"name"`
	Platform          string  `json:"platform"`
	PlatformProductID string  `json:"platform_product_id"`
	Price            float64 `json:"price"`
	Currency         string  `json:"currency"`
	Rating           float64 `json:"rating"`
	ReviewCount      int     `json:"review_count"`
	SellerRating     float64 `json:"seller_rating"`
	SellerReviewCount int    `json:"seller_review_count"`
	SourceURL        string  `json:"source_url"`
}

// CleanText 清洗文本
func (d *DataCleaner) CleanText(text string) string {
	if text == "" {
		return ""
	}

	// 移除多余空白
	text = strings.Join(strings.Fields(text), " ")

	// 移除特殊字符但保留中文、英文、数字和基本标点
	var result strings.Builder
	for _, r := range text {
		if unicode.IsPrint(r) && (unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || unicode.IsPunct(r)) {
			result.WriteRune(r)
		}
	}
	text = result.String()

	return strings.TrimSpace(text)
}

// CleanPlatform 标准化平台名称
func (d *DataCleaner) CleanPlatform(platform string) string {
	platform = strings.ToLower(strings.TrimSpace(platform))

	platforms := map[string]string{
		"amazon":           "amazon",
		"amazon.com":       "amazon",
		"www.amazon.com":   "amazon",
		"aliexpress":       "aliexpress",
		"aliexpress.com":   "aliexpress",
		"www.aliexpress.com": "aliexpress",
		"ebay":             "ebay",
		"www.ebay.com":     "ebay",
		"temu":             "temu",
		"www.temu.com":     "temu",
		"shopify":          "shopify",
		"taobao":           "taobao",
		"tmall":            "tmall",
		"jd":               "jd",
		"jd.com":           "jd",
	}

	if v, ok := platforms[platform]; ok {
		return v
	}

	return platform
}

// CleanCurrency 标准化货币
func (d *DataCleaner) CleanCurrency(currency string) string {
	currency = strings.ToUpper(strings.TrimSpace(currency))

	currencies := map[string]string{
		"USD": "USD",
		"US$": "USD",
		"$":   "USD",
		"EUR": "EUR",
		"€":   "EUR",
		"GBP": "GBP",
		"£":   "GBP",
		"CNY": "CNY",
		"JPY": "JPY",
	}

	if v, ok := currencies[currency]; ok {
		return v
	}

	return currency
}

// CleanPrice 清洗价格
func (d *DataCleaner) CleanPrice(price string) (float64, error) {
	return ParsePrice(price)
}

// CleanRating 清洗评分
func (d *DataCleaner) CleanRating(rating string) (float64, error) {
	return ParseRating(rating)
}

// CleanReviewCount 清洗评论数
func (d *DataCleaner) CleanReviewCount(count string) (int, error) {
	return ParseReviewCount(count)
}

// RemoveHTMLTags 移除 HTML 标签
func (d *DataCleaner) RemoveHTMLTags(html string) string {
	// 简单 HTML 标签移除
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")

	// 解码 HTML 实体
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	return d.CleanText(text)
}

// TruncateString 截断字符串
func (d *DataCleaner) TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// ValidateURL 验证 URL 格式
func (d *DataCleaner) ValidateURL(urlStr string) bool {
	re := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(urlStr)
}

// ExtractDomain 提取域名
func (d *DataCleaner) ExtractDomain(urlStr string) string {
	re := regexp.MustCompile(`https?://([^/]+)`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// NormalizeWhitespace 规范化空白字符
func (d *DataCleaner) NormalizeWhitespace(s string) string {
	// 将所有连续空白替换为单个空格
	space := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(space.ReplaceAllString(s, " "))
}

// JSONToMap JSON 转为 Map
func (d *DataCleaner) JSONToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}
	return result, nil
}

// MapToJSON Map 转为 JSON
func (d *DataCleaner) MapToJSON(data map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}
	return string(jsonBytes), nil
}

// SanitizeForDB 清洗用于数据库存储
func (d *DataCleaner) SanitizeForDB(s string) string {
	// 移除或转义特殊字符防止 SQL 注入 (假设使用参数化查询，此处做额外防护)
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\\", "\\\\")
	return s
}

// ExtractProductID 提取产品ID
func (d *DataCleaner) ExtractProductID(url string, platform string) string {
	switch platform {
	case "amazon":
		// Amazon ASIN 提取 (ASIN 通常是 10 个字符)
		re := regexp.MustCompile(`/dp/([A-Z0-9]+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "aliexpress":
		// AliExpress 产品 ID 提取
		re := regexp.MustCompile(`/item/(\d+)\.html`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "ebay":
		// eBay 产品 ID 提取
		re := regexp.MustCompile(`/itm/(\d+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "temu":
		// Temu 产品 ID 提取
		re := regexp.MustCompile(`/item/([\w-]+)\.html`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "shopify":
		// Shopify 产品 ID 从 URL 提取
		re := regexp.MustCompile(`/products/([\w-]+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
