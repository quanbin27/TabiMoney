<template>
  <v-container>
    <v-overlay :model-value="loading" class="align-center justify-center">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </v-overlay>

    <v-row v-if="!loading && user">
      <v-col cols="12" xl="4">
        <InfoView :user="user" />
      </v-col>

      <v-col cols="12" xl="8">
        <StatsView :user="user" :statistics="statistics" />
      </v-col>
    </v-row>

    <!-- Edit profile basic info -->
    <v-row v-if="!loading && user" class="mt-4">
      <v-col cols="12" xl="8" offset-xl="4">
        <v-card class="pa-4 border-md rounded-xl">
          <v-card-title class="text-h6 pb-2">
            Cập nhật hồ sơ
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="editForm.first_name"
                  label="Tên"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="editForm.last_name"
                  label="Họ"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="editForm.username"
                  label="Tên đăng nhập"
                  variant="outlined"
                  density="comfortable"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="editForm.phone"
                  label="Số điện thoại"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
            </v-row>

            <div class="d-flex justify-end mt-2">
              <v-btn color="primary" :loading="savingProfile" @click="saveProfile">
                Lưu thay đổi
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
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
import { ref, onMounted, computed, watch } from 'vue'
import { InfoView, StatsView } from './components'
import { authAPI, analyticsAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const appStore = useAppStore()
const authStore = useAuthStore()

const loading = ref(true)
const error = ref(null)
const statistics = ref(null)

// Edit profile form
const editForm = ref({
  first_name: '',
  last_name: '',
  username: '',
  phone: '',
})
const savingProfile = ref(false)

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

// Sync edit form when user data changes
watch(user, (val) => {
  if (!val) return
  editForm.value = {
    first_name: val.first_name || '',
    last_name: val.last_name || '',
    username: val.username || '',
    phone: val.phone || '',
  }
}, { immediate: true })

const saveProfile = async () => {
  if (!editForm.value.username) {
    appStore.showWarning('Vui lòng nhập tên đăng nhập')
    return
  }

  try {
    savingProfile.value = true
    await authStore.updateProfile({
      first_name: editForm.value.first_name || null,
      last_name: editForm.value.last_name || null,
      username: editForm.value.username,
      phone: editForm.value.phone || null,
    })
    appStore.showSuccess('Cập nhật hồ sơ thành công!')
  } catch (e) {
    console.error('Failed to update profile:', e)
    const msg = e?.response?.data?.message || 'Không thể cập nhật hồ sơ. Vui lòng thử lại.'
    appStore.showError(msg)
  } finally {
    savingProfile.value = false
  }
}

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