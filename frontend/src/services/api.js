import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'

export const api = axios.create({
  baseURL: (import.meta)?.env?.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const authStore = useAuthStore()
    const appStore = useAppStore()

    if (error.response?.status === 401) {
      try {
        await authStore.refreshAccessToken()
        return api.request(error.config)
      } catch (refreshError) {
        await authStore.logout()
        appStore.showError('Session expired. Please login again.')
        return Promise.reject(refreshError)
      }
    }

    if (error.response?.status === 403) {
      appStore.showError('You do not have permission to perform this action.')
    }

    if (error.response?.status === 404) {
      appStore.showError('Resource not found.')
    }

    if (error.response?.status >= 500) {
      appStore.showError('Server error. Please try again later.')
    }

    if (!error.response) {
      appStore.showError('Network error. Please check your connection.')
    }

    return Promise.reject(error)
  }
)

export const authAPI = {
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
  logout: () => api.post('/auth/logout'),
  refresh: (refreshToken) => api.post('/auth/refresh', { refresh_token: refreshToken }),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (profileData) => api.put('/auth/profile', profileData),
  changePassword: (passwordData) => api.post('/auth/change-password', passwordData),
  getIncome: () => api.get('/auth/income'),
  setIncome: (monthly_income) => api.put('/auth/income', { monthly_income }),
}

export const transactionAPI = {
  getTransactions: (params) => api.get('/transactions', { params }),
  getTransaction: (id) => api.get(`/transactions/${id}`),
  createTransaction: (data) => api.post('/transactions', data),
  updateTransaction: (id, data) => api.put(`/transactions/${id}`, data),
  deleteTransaction: (id) => api.delete(`/transactions/${id}`),
}

export const analyticsAPI = {
  getDashboard: (params) => api.get('/analytics/dashboard', { params }),
  getCategorySpending: (params) => api.get('/analytics/category-spending', { params }),
  getSpendingPatterns: (params) => api.get('/analytics/spending-patterns', { params }),
  getAnomalies: (params) => api.get('/analytics/anomalies', { params }),
  getPredictions: (params) => api.get('/analytics/predictions', { params }),
}

export const notificationsAPI = {
  list: (unreadOnly = false) => api.get('/notifications', { params: unreadOnly ? { unread: true } : {} }),
  markRead: (id) => api.post(`/notifications/${id}/read`),
}

export const categoryAPI = {
  getCategories: () => api.get('/categories'),
  getCategory: (id) => api.get(`/categories/${id}`),
  createCategory: (data) => api.post('/categories', data),
  updateCategory: (id, data) => api.put(`/categories/${id}`, data),
  deleteCategory: (id) => api.delete(`/categories/${id}`),
}

export const goalAPI = {
  getGoals: () => api.get('/goals'),
  getGoal: (id) => api.get(`/goals/${id}`),
  createGoal: (data) => api.post('/goals', data),
  updateGoal: (id, data) => api.put(`/goals/${id}`, data),
  deleteGoal: (id) => api.delete(`/goals/${id}`),
  addContribution: (id, data) => api.post(`/goals/${id}/contribute`, data),
}

export const budgetAPI = {
  getBudgets: () => api.get('/budgets'),
  getBudget: (id) => api.get(`/budgets/${id}`),
  createBudget: (data) => api.post('/budgets', data),
  updateBudget: (id, data) => api.put(`/budgets/${id}`, data),
  deleteBudget: (id) => api.delete(`/budgets/${id}`),
}

const aiServiceApi = axios.create({
  baseURL: import.meta.env.VITE_AI_SERVICE_URL || 'http://localhost:8001',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

export const aiAPI = {
  processChat: (data) => aiServiceApi.post('/api/v1/chat/process', data),
  processNLU: (data) => aiServiceApi.post('/api/v1/nlu/process', data),
  predictExpenses: (data) => api.post('/ai/predict-expenses', data),
  detectAnomalies: (data) => api.post('/ai/detect-anomalies', data),
  suggestCategory: (data) => api.post('/ai/suggest-category', data),
  analyzeSpending: (data) => api.post('/ai/analyze-spending', data),
  analyzeGoal: (data) => api.post('/ai/analyze-goal', data),
}


