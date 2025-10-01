import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../services/api'
import type { User, AuthResponse, LoginRequest, RegisterRequest } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  
  // State
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const initialized = ref(false)
  const loading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)
  const isGuest = computed(() => !isAuthenticated.value)

  // Actions
  const initialize = async () => {
    if (initialized.value) return

    try {
      // Try to get tokens from localStorage
      const storedAccessToken = localStorage.getItem('access_token')
      const storedRefreshToken = localStorage.getItem('refresh_token')

      if (storedAccessToken && storedRefreshToken) {
        accessToken.value = storedAccessToken
        refreshToken.value = storedRefreshToken

        // Set default authorization header
        api.defaults.headers.common['Authorization'] = `Bearer ${storedAccessToken}`

        // Try to get user profile
        try {
          const response = await api.get('/auth/profile')
          user.value = response.data
        } catch (error) {
          // If profile fetch fails, try to refresh token
          await refreshAccessToken()
        }
      }
    } catch (error) {
      console.error('Auth initialization failed:', error)
      // Clear invalid tokens
      await logout()
    } finally {
      initialized.value = true
    }
  }

  const login = async (credentials: LoginRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await api.post<AuthResponse>('/auth/login', credentials)
      const { user: userData, access_token, refresh_token } = response.data

      // Store tokens
      accessToken.value = access_token
      refreshToken.value = refresh_token
      user.value = userData

      // Save to localStorage
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)

      // Set default authorization header
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`

      // Redirect to dashboard
      router.push('/dashboard')
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Login failed')
    } finally {
      loading.value = false
    }
  }

  const register = async (userData: RegisterRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await api.post<AuthResponse>('/auth/register', userData)
      const { user: newUser, access_token, refresh_token } = response.data

      // Store tokens
      accessToken.value = access_token
      refreshToken.value = refresh_token
      user.value = newUser

      // Save to localStorage
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)

      // Set default authorization header
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`

      // Redirect to dashboard
      router.push('/dashboard')
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Registration failed')
    } finally {
      loading.value = false
    }
  }

  const logout = async (): Promise<void> => {
    try {
      if (accessToken.value) {
        await api.post('/auth/logout')
      }
    } catch (error) {
      console.error('Logout request failed:', error)
    } finally {
      // Clear state
      user.value = null
      accessToken.value = null
      refreshToken.value = null

      // Clear localStorage
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')

      // Clear authorization header
      delete api.defaults.headers.common['Authorization']

      // Redirect to login
      router.push('/auth/login')
    }
  }

  const refreshAccessToken = async (): Promise<void> => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await api.post<AuthResponse>('/auth/refresh', {
        refresh_token: refreshToken.value,
      })

      const { access_token, refresh_token: newRefreshToken } = response.data

      // Update tokens
      accessToken.value = access_token
      refreshToken.value = newRefreshToken

      // Save to localStorage
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', newRefreshToken)

      // Update authorization header
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`
    } catch (error) {
      // If refresh fails, logout
      await logout()
      throw new Error('Session expired. Please login again.')
    }
  }

  const updateProfile = async (profileData: Partial<User>): Promise<void> => {
    try {
      const response = await api.put('/auth/profile', profileData)
      user.value = response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Profile update failed')
    }
  }

  const changePassword = async (passwordData: { current_password: string; new_password: string }): Promise<void> => {
    try {
      await api.post('/auth/change-password', passwordData)
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Password change failed')
    }
  }

  return {
    // State
    user,
    accessToken,
    refreshToken,
    initialized,
    loading,
    
    // Getters
    isAuthenticated,
    isGuest,
    
    // Actions
    initialize,
    login,
    register,
    logout,
    refreshAccessToken,
    updateProfile,
    changePassword,
  }
})
