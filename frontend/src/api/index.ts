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

export const healthAPI = {
  check: () => axios.get('/api/health')
}

export const productsAPI = {
  list: () => apiClient.get('/products'),
  get: (id: string) => apiClient.get(`/products/${id}`)
}

export const competitorsAPI = {
  list: () => apiClient.get('/competitors'),
  get: (id: string) => apiClient.get(`/competitors/${id}`)
}

export const trendsAPI = {
  list: () => apiClient.get('/trends')
}

export const strategyAPI = {
  getPricing: () => apiClient.get('/strategy/pricing')
}
