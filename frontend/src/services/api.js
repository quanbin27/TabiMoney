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
        appStore.showError('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
        return Promise.reject(refreshError)
      }
    }

    if (error.response?.status === 403) {
      appStore.showError('Bạn không có quyền thực hiện thao tác này.')
    }

    if (error.response?.status === 404) {
      appStore.showError('Không tìm thấy tài nguyên.')
    }

    if (error.response?.status >= 500) {
      appStore.showError('Lỗi máy chủ. Vui lòng thử lại sau.')
    }

    if (!error.response) {
      appStore.showError('Lỗi mạng. Vui lòng kiểm tra kết nối của bạn.')
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
  setLargeTransactionThreshold: (threshold) => api.put('/auth/large-transaction-threshold', { threshold }),
}

export const transactionAPI = {
  getTransactions: (params) => api.get('/transactions', { params }),
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

export const notificationPreferencesAPI = {
  getPreferences: () => api.get('/notification-preferences'),
  updatePreferences: (data) => api.put('/notification-preferences', data),
  getSummary: () => api.get('/notification-preferences/summary'),
  resetToDefaults: () => api.post('/notification-preferences/reset'),
  getEnabledChannels: () => api.get('/notification-preferences/channels'),
  testNotification: (channel = 'in_app') => api.post('/notification-preferences/test', null, { params: { channel } }),
}

export const categoryAPI = {
  getCategories: () => api.get('/categories'),
  createCategory: (data) => api.post('/categories', data),
  updateCategory: (id, data) => api.put(`/categories/${id}`, data),
  deleteCategory: (id) => api.delete(`/categories/${id}`),
}

export const goalAPI = {
  getGoals: () => api.get('/goals'),
  createGoal: (data) => api.post('/goals', data),
  updateGoal: (id, data) => api.put(`/goals/${id}`, data),
  deleteGoal: (id) => api.delete(`/goals/${id}`),
  addContribution: (id, data) => api.post(`/goals/${id}/contribute`, data),
}

export const budgetAPI = {
  getBudgets: () => api.get('/budgets'),
  createBudget: (data) => api.post('/budgets', data),
  updateBudget: (id, data) => api.put(`/budgets/${id}`, data),
  deleteBudget: (id) => api.delete(`/budgets/${id}`),
  getInsights: () => api.get('/budgets/insights'),
  getAutoSuggestions: () => api.get('/budgets/auto/suggestions'),
  createFromSuggestions: (payload) => api.post('/budgets/auto/create', payload),
}

// AI Service API - Use relative path for production, localhost for development
const getAIServiceURL = () => {
  const url = import.meta.env.VITE_AI_SERVICE_URL || '/ai-service'
  // If it's already a full URL (starts with http), use it as is
  // Otherwise, it's a relative path and axios will use current origin
  return url.startsWith('http') ? url : url
}

const aiServiceApi = axios.create({
  baseURL: getAIServiceURL(),
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

export const aiAPI = {
  processChat: (data) => aiServiceApi.post('/api/v1/chat/process', data),
  suggestCategory: (data) => api.post('/ai/suggest-category', data),
}


