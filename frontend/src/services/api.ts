import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'

// Create axios instance
export const api = axios.create({
  baseURL: (import.meta as any)?.env?.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    
    // Add auth token if available
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const authStore = useAuthStore()
    const appStore = useAppStore()
    
    // Handle 401 errors (unauthorized)
    if (error.response?.status === 401) {
      // Try to refresh token
      try {
        await authStore.refreshAccessToken()
        // Retry the original request
        return api.request(error.config)
      } catch (refreshError) {
        // If refresh fails, logout
        await authStore.logout()
        appStore.showError('Session expired. Please login again.')
        return Promise.reject(refreshError)
      }
    }
    
    // Handle 403 errors (forbidden)
    if (error.response?.status === 403) {
      appStore.showError('You do not have permission to perform this action.')
    }
    
    // Handle 404 errors (not found)
    if (error.response?.status === 404) {
      appStore.showError('Resource not found.')
    }
    
    // Handle 500 errors (server error)
    if (error.response?.status >= 500) {
      appStore.showError('Server error. Please try again later.')
    }
    
    // Handle network errors
    if (!error.response) {
      appStore.showError('Network error. Please check your connection.')
    }
    
    return Promise.reject(error)
  }
)

// API endpoints
export const authAPI = {
  login: (credentials: any) => api.post('/auth/login', credentials),
  register: (userData: any) => api.post('/auth/register', userData),
  logout: () => api.post('/auth/logout'),
  refresh: (refreshToken: string) => api.post('/auth/refresh', { refresh_token: refreshToken }),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (profileData: any) => api.put('/auth/profile', profileData),
  changePassword: (passwordData: any) => api.post('/auth/change-password', passwordData),
  getIncome: () => api.get('/auth/income'),
  setIncome: (monthly_income: number) => api.put('/auth/income', { monthly_income }),
}

export const transactionAPI = {
  getTransactions: (params?: any) => api.get('/transactions', { params }),
  getTransaction: (id: number) => api.get(`/transactions/${id}`),
  createTransaction: (data: any) => api.post('/transactions', data),
  updateTransaction: (id: number, data: any) => api.put(`/transactions/${id}`, data),
  deleteTransaction: (id: number) => api.delete(`/transactions/${id}`),
}

export const analyticsAPI = {
  getDashboard: (params?: any) => api.get('/analytics/dashboard', { params }),
  getCategorySpending: (params?: any) => api.get('/analytics/category-spending', { params }),
  getSpendingPatterns: (params?: any) => api.get('/analytics/spending-patterns', { params }),
  getAnomalies: (params?: any) => api.get('/analytics/anomalies', { params }),
  getPredictions: (params?: any) => api.get('/analytics/predictions', { params }),
}

export const notificationsAPI = {
  list: (unreadOnly = false) => api.get('/notifications', { params: unreadOnly ? { unread: true } : {} }),
  markRead: (id: number) => api.post(`/notifications/${id}/read`),
}

export const categoryAPI = {
  getCategories: () => api.get('/categories'),
  getCategory: (id: number) => api.get(`/categories/${id}`),
  createCategory: (data: any) => api.post('/categories', data),
  updateCategory: (id: number, data: any) => api.put(`/categories/${id}`, data),
  deleteCategory: (id: number) => api.delete(`/categories/${id}`),
}

export const goalAPI = {
  getGoals: () => api.get('/goals'),
  getGoal: (id: number) => api.get(`/goals/${id}`),
  createGoal: (data: any) => api.post('/goals', data),
  updateGoal: (id: number, data: any) => api.put(`/goals/${id}`, data),
  deleteGoal: (id: number) => api.delete(`/goals/${id}`),
  addContribution: (id: number, data: { amount: number; note?: string }) => api.post(`/goals/${id}/contribute`, data),
}

export const budgetAPI = {
  getBudgets: () => api.get('/budgets'),
  getBudget: (id: number) => api.get(`/budgets/${id}`),
  createBudget: (data: any) => api.post('/budgets', data),
  updateBudget: (id: number, data: any) => api.put(`/budgets/${id}`, data),
  deleteBudget: (id: number) => api.delete(`/budgets/${id}`),
}

// AI Service direct client (bypass backend)
const aiServiceApi = axios.create({
  baseURL: import.meta.env.VITE_AI_SERVICE_URL || 'http://localhost:8001',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const aiAPI = {
  // Direct AI service calls
  processChat: (data: any) => aiServiceApi.post('/api/v1/chat/process', data),
  processNLU: (data: any) => aiServiceApi.post('/api/v1/nlu/process', data),
  
  // Backend proxy calls (for other AI features)
  predictExpenses: (data: any) => api.post('/ai/predict-expenses', data),
  detectAnomalies: (data: any) => api.post('/ai/detect-anomalies', data),
  suggestCategory: (data: any) => api.post('/ai/suggest-category', data),
  analyzeSpending: (data: any) => api.post('/ai/analyze-spending', data),
  analyzeGoal: (data: any) => api.post('/ai/analyze-goal', data),
}
