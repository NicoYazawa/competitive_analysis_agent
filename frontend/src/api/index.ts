import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

apiClient.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export default apiClient

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------
export const healthAPI = {
  check: () => axios.get('/api/health')
}

// ---------------------------------------------------------------------------
// Competitors
// ---------------------------------------------------------------------------
export const competitorsAPI = {
  list: (platform?: string) => {
    const params = platform ? { platform } : {}
    return apiClient.get('/competitors/', { params })
  },
  get: (id: string) => apiClient.get(`/competitors/${id}`),
  search: (q: string, platform?: string) => {
    const params: Record<string, string> = { q }
    if (platform) params.platform = platform
    return apiClient.get('/competitors/search', { params })
  },
  priceChanges: (threshold?: number) => {
    const params = threshold ? { threshold } : {}
    return apiClient.get('/competitors/price-changes', { params })
  },
  getPriceHistory: (id: string, days = 30) =>
    apiClient.get(`/competitors/${id}/price-history`, { params: { days } })
}

// ---------------------------------------------------------------------------
// Trends
// ---------------------------------------------------------------------------
export const trendOverviewAPI = {
  get: () => apiClient.get('/trends/overview')
}

export const trendsAPI = {
  analyze: (query: string) => apiClient.get('/trends/analyze', { params: { query } }),
  getAnalysis: (query: string) => apiClient.get('/trends/analysis', { params: { query } }),
  listCategories: () => apiClient.get('/trends/categories'),
  getCategory: (category: string) => apiClient.get(`/trends/categories/${category}`),
  getDemandSignals: (platform?: string) => {
    const params = platform ? { platform } : {}
    return apiClient.get('/trends/demand-signals', { params })
  },
  getTrending: (limit = 10) => apiClient.get('/trends/trending', { params: { limit } })
}

// ---------------------------------------------------------------------------
// Strategy / Pricing
// ---------------------------------------------------------------------------
export const strategyAPI = {
  getPricing: () => apiClient.get('/strategy/pricing/'),
  getCompetitorPricing: (productId: string) =>
    apiClient.get(`/strategy/pricing/competitors/${productId}`),
  getPriceElasticity: (productId: string) =>
    apiClient.get(`/strategy/pricing/elasticity/${productId}`),
  getPricingScenario: (payload: object) =>
    apiClient.post('/strategy/pricing/scenario', payload),
  getPricingHistory: (productId: string) =>
    apiClient.get(`/strategy/pricing/history/${productId}`),
  getDynamicPricing: (productId: string) =>
    apiClient.get(`/strategy/pricing/dynamic/${productId}`)
}

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------
export const productsAPI = {
  getRecommendations: () => apiClient.get('/products/recommendations'),
  compare: (ids: string[]) => apiClient.get('/products/compare', { params: { ids: ids.join(',') } }),
  getTrends: () => apiClient.get('/products/trends'),
  getById: (id: string) => apiClient.get(`/products/${id}`),
  getAnalysis: (id: string) => apiClient.get(`/products/${id}/analysis`),
  getCategoryInsights: (category: string) =>
    apiClient.get(`/products/categories/${category}/insights`)
}
