<template>
  <v-container>
    <v-overlay :model-value="loading" class="align-center justify-center">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </v-overlay>

    <v-row v-if="!loading && user">
      <v-col cols="12" xl="4">
        <InfoView :user="user"></InfoView>
      </v-col>

      <v-col cols="12" xl="8">
        <AccountView :user="user"></AccountView>
      </v-col>
    </v-row>

    <v-row v-if="!loading && user" class="mt-4">
      <v-col cols="12" xl="4">
        <PreferencesView :user="user"></PreferencesView>
      </v-col>

      <v-col cols="12" xl="8">
        <StatsView :user="user" :statistics="statistics"></StatsView>
      </v-col>
    </v-row>

    <v-row v-if="!loading && error">
      <v-col cols="12">
        <v-alert type="error" :text="error"></v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { AccountView, InfoView, PreferencesView, StatsView } from './components'
import { authAPI, analyticsAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const appStore = useAppStore()
const authStore = useAuthStore()

const loading = ref(true)
const error = ref(null)
const statistics = ref(null)

// Transform backend user data to component-expected format
const user = computed(() => {
  if (!authStore.user) return null

  const backendUser = authStore.user
  const profile = backendUser.profile || {}

  // Parse notification settings if it's a JSON string
  let notificationSettings = {}
  if (profile.notification_settings) {
    try {
      notificationSettings = typeof profile.notification_settings === 'string' 
        ? JSON.parse(profile.notification_settings) 
        : profile.notification_settings
    } catch (e) {
      notificationSettings = {}
    }
  }

  return {
    id: backendUser.id,
    name: [backendUser.first_name, backendUser.last_name].filter(Boolean).join(' ') || backendUser.username || 'User',
    email: backendUser.email,
    username: backendUser.username,
    first_name: backendUser.first_name,
    last_name: backendUser.last_name,
    phone: backendUser.phone || '',
    avatar: backendUser.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(backendUser.email)}&background=random`,
    avatar_url: backendUser.avatar_url,
    is_verified: backendUser.is_verified,
    created_at: backendUser.created_at,
    last_login_at: backendUser.last_login_at,
    // Accounts - not available in backend, using empty array
    accounts: [],
    // Preferences from profile
    preferences: {
      currency: profile.currency || 'VND',
      language: profile.language || 'vi',
      timezone: profile.timezone || 'Asia/Ho_Chi_Minh',
      monthly_income: profile.monthly_income || 0,
      notification: notificationSettings.enabled !== false, // Default to true
      notification_settings: notificationSettings,
      ai_settings: profile.ai_settings ? (typeof profile.ai_settings === 'string' ? JSON.parse(profile.ai_settings) : profile.ai_settings) : {}
    },
    // Raw profile data for updates
    profile: profile
  }
})

const fetchProfile = async () => {
  try {
    loading.value = true
    error.value = null
    
    // Fetch user profile if not already in store
    if (!authStore.user) {
      const response = await authAPI.getProfile()
      authStore.user = response.data
    }

    // Fetch statistics from analytics
    try {
      const analyticsResponse = await analyticsAPI.getDashboard()
      const analyticsData = analyticsResponse.data
      statistics.value = {
        totalIncome: analyticsData.total_income || 0,
        totalExpense: analyticsData.total_expense || 0,
        totalTransactions: analyticsData.transaction_count || 0,
        categoriesUsed: analyticsData.top_categories?.length || 0,
        netAmount: analyticsData.net_amount || 0
      }
    } catch (analyticsError) {
      console.warn('Failed to fetch analytics:', analyticsError)
      // Set default statistics if analytics fails
      statistics.value = {
        totalIncome: 0,
        totalExpense: 0,
        totalTransactions: 0,
        categoriesUsed: 0,
        netAmount: 0
      }
    }
  } catch (err) {
    console.error('Failed to fetch profile:', err)
    error.value = err?.response?.data?.message || 'Không thể tải thông tin profile. Vui lòng thử lại.'
    appStore.showError(error.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>