import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()

  const user = ref(null)
  const accessToken = ref(null)
  const refreshToken = ref(null)
  const initialized = ref(false)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)
  const isGuest = computed(() => !isAuthenticated.value)

  const initialize = async () => {
    if (initialized.value) return
    try {
      const storedAccessToken = localStorage.getItem('access_token')
      const storedRefreshToken = localStorage.getItem('refresh_token')
      if (storedAccessToken && storedRefreshToken) {
        accessToken.value = storedAccessToken
        refreshToken.value = storedRefreshToken
        api.defaults.headers.common['Authorization'] = `Bearer ${storedAccessToken}`
        try {
          const response = await api.get('/auth/profile')
          user.value = response.data
        } catch (error) {
          await refreshAccessToken()
        }
      }
    } catch (error) {
      console.error('Auth initialization failed:', error)
      await logout()
    } finally {
      initialized.value = true
    }
  }

  const login = async (credentials) => {
    loading.value = true
    try {
      const response = await api.post('/auth/login', credentials)
      const { user: userData, access_token, refresh_token } = response.data
      accessToken.value = access_token
      refreshToken.value = refresh_token
      user.value = userData
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`
      router.push('/dashboard')
    } catch (error) {
      const message = error?.response?.data?.message || 'Đăng nhập thất bại'
      throw new Error(message)
    } finally {
      loading.value = false
    }
  }

  const register = async (userData) => {
    loading.value = true
    try {
      const response = await api.post('/auth/register', userData)
      const { user: newUser, access_token, refresh_token } = response.data
      accessToken.value = access_token
      refreshToken.value = refresh_token
      user.value = newUser
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`
      router.push('/dashboard')
    } catch (error) {
      const message = error?.response?.data?.message || 'Đăng ký thất bại'
      throw new Error(message)
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      if (accessToken.value) {
        await api.post('/auth/logout')
      }
    } catch (error) {
      console.error('Logout request failed:', error)
    } finally {
      user.value = null
      accessToken.value = null
      refreshToken.value = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      delete api.defaults.headers.common['Authorization']
      router.push('/auth/login')
    }
  }

  const refreshAccessToken = async () => {
    if (!refreshToken.value) throw new Error('No refresh token available')
    try {
      const response = await api.post('/auth/refresh', { refresh_token: refreshToken.value })
      const { access_token, refresh_token: newRefreshToken } = response.data
      accessToken.value = access_token
      refreshToken.value = newRefreshToken
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', newRefreshToken)
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`
    } catch (error) {
      await logout()
      throw new Error('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
    }
  }

  const updateProfile = async (profileData) => {
    // Backend requires a valid URL for AvatarURL, so always send a proper avatar_url
    let avatarUrl = user.value?.avatar_url || ''
    if (!avatarUrl) {
      const baseName = user.value?.email || user.value?.username || 'User'
      avatarUrl = `https://ui-avatars.com/api/?name=${encodeURIComponent(baseName)}&background=random`
    }

    const payload = {
      ...profileData,
      avatar_url: avatarUrl,
    }

    const response = await api.put('/auth/profile', payload)
    user.value = response.data
  }

  const changePassword = async (passwordData) => {
    await api.post('/auth/change-password', passwordData)
  }

  return {
    user,
    accessToken,
    refreshToken,
    initialized,
    loading,
    isAuthenticated,
    isGuest,
    initialize,
    login,
    register,
    logout,
    refreshAccessToken,
    updateProfile,
    changePassword,
  }
})


