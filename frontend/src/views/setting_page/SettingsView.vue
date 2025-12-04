<template>
  <v-container class="py-4" fluid>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">Cài đặt</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="8">
        <!-- Notification Preferences -->
        <v-card class="mb-3">
          <v-card-title class="text-h6 py-3">
            <v-icon left size="small">mdi-bell-outline</v-icon>
            Cài đặt thông báo
          </v-card-title>

          <v-card-text class="pt-2 pb-3">
            <v-expansion-panels variant="accordion" density="compact">
              <!-- Channel Preferences -->
              <v-expansion-panel title="Kênh thông báo">
                <v-expansion-panel-text class="py-2">
                  <v-switch v-model="notificationPrefs.email_enabled" label="Email" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.telegram_enabled" label="Telegram" color="primary" hide-details
                    density="compact" :disabled="telegramStatus !== 'connected'"></v-switch>
                  <v-switch v-model="notificationPrefs.in_app_enabled" label="Trong ứng dụng" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-alert v-if="telegramStatus !== 'connected'" type="info" variant="tonal" density="compact"
                    class="mt-2 text-caption">
                    Kích hoạt Telegram ở phần dưới
                  </v-alert>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <!-- Feature Preferences -->
              <v-expansion-panel title="Loại thông báo">
                <v-expansion-panel-text class="py-2">
                  <v-switch v-model="notificationPrefs.budget_alerts" label="Cảnh báo ngân sách" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.goal_alerts" label="Cảnh báo mục tiêu" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.ai_alerts" label="Cảnh báo AI" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.transaction_alerts" label="Cảnh báo giao dịch" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.analytics_alerts" label="Báo cáo phân tích" color="primary"
                    hide-details density="compact"></v-switch>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <!-- Priority Preferences -->
              <v-expansion-panel title="Mức độ ưu tiên">
                <v-expansion-panel-text class="py-2">
                  <v-switch v-model="notificationPrefs.urgent_notifications" label="Khẩn cấp" color="red"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.high_notifications" label="Cao" color="orange"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.medium_notifications" label="Trung bình" color="yellow"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.low_notifications" label="Thấp" color="grey"
                    hide-details density="compact"></v-switch>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <!-- Frequency Preferences -->
              <v-expansion-panel title="Tần suất">
                <v-expansion-panel-text class="py-2">
                  <v-switch v-model="notificationPrefs.daily_digest" label="Báo cáo hàng ngày" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.weekly_digest" label="Báo cáo hàng tuần" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.monthly_digest" label="Báo cáo hàng tháng" color="primary"
                    hide-details density="compact"></v-switch>
                  <v-switch v-model="notificationPrefs.real_time_alerts" label="Cảnh báo thời gian thực" color="primary"
                    hide-details density="compact"></v-switch>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>

            <v-divider class="my-3"></v-divider>

            <div class="d-flex justify-end ga-2">
              <v-btn color="secondary" size="small" @click="resetPreferences" :loading="resetting">
                Mặc định
              </v-btn>
              <v-btn color="primary" size="small" @click="savePreferences" :loading="saving">
                Lưu
              </v-btn>
            </div>
          </v-card-text>
        </v-card>

        <!-- Transaction Settings -->
        <v-card class="mb-3">
          <v-card-title class="text-h6 py-3">
            <v-icon left size="small">mdi-cash-multiple</v-icon>
            Cài đặt giao dịch
          </v-card-title>

          <v-card-text class="pt-2 pb-3">
            <v-text-field
              v-model.number="largeTransactionThreshold"
              label="Ngưỡng giao dịch lớn (VND)"
              type="number"
              hint="Cảnh báo khi giao dịch chi tiêu vượt ngưỡng"
              persistent-hint
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-alert-circle"
              :rules="[rules.required, rules.min]"
              class="mb-3"
            >
              <template v-slot:append>
                <v-btn
                  icon="mdi-refresh"
                  variant="text"
                  size="small"
                  @click="resetThreshold"
                  :disabled="savingThreshold"
                  title="Mặc định: 1,000,000 VND"
                ></v-btn>
              </template>
            </v-text-field>

            <v-alert type="info" variant="tonal" density="compact" class="mb-3 text-caption">
              <strong>Lưu ý:</strong> Chỉ áp dụng cho giao dịch <strong>chi tiêu</strong>, không áp dụng cho thu nhập.
            </v-alert>

            <div class="d-flex justify-end">
              <v-btn
                color="primary"
                size="small"
                @click="saveThreshold"
                :loading="savingThreshold"
              >
                <v-icon left size="small">mdi-content-save</v-icon>
                Lưu
              </v-btn>
            </div>
          </v-card-text>
        </v-card>

        <!-- Telegram Integration -->
        <v-card>
          <v-card-title class="text-h6 py-3">
            <v-icon size="24" color="blue" class="mr-2">
              <svg viewBox="0 0 24 24" width="24" height="24">
                <path fill="currentColor"
                  d="M9.04 16.54L9.25 12.97L17.77 5.52C18.12 5.22 17.73 5.09 17.26 5.33L6.75 10.57L3.26 9.47C2.5 9.24 2.49 8.71 3.42 8.35L20.42 1.7C21.09 1.45 21.67 1.9 21.45 2.84L18.6 16.27C18.44 17.01 17.98 17.21 17.36 16.87L13.83 14.29L12.13 15.92C11.98 16.07 11.86 16.19 11.6 16.19L9.04 16.54Z" />
              </svg>
            </v-icon>
            Tích hợp Telegram Bot
          </v-card-title>

          <v-card-text class="pt-2 pb-3">
            <p class="text-body-2 mb-3 text-medium-emphasis">
              Liên kết tài khoản với Telegram Bot để sử dụng AI và dashboard.
            </p>

            <v-alert v-if="telegramStatus === 'connected'" type="success" density="compact" class="mb-3">
              Đã liên kết với Telegram Bot
            </v-alert>

            <v-alert v-else-if="telegramStatus === 'disconnected'" type="info" density="compact" class="mb-3">
              Chưa liên kết với Telegram Bot
            </v-alert>

            <div v-if="telegramStatus === 'disconnected'">
              <v-btn color="primary" size="small" @click="generateLinkCode" :loading="generatingCode" class="mb-3">
                <v-icon left size="small">mdi-link</v-icon>
                Tạo mã liên kết
              </v-btn>

              <v-card v-if="linkCode" class="mt-3" color="primary" variant="outlined" density="compact">
                <v-card-text class="py-3">
                  <h3 class="text-subtitle-1 mb-2">Mã liên kết:</h3>
                  <v-text-field :value="linkCode" readonly variant="outlined" density="compact" 
                    append-inner-icon="mdi-content-copy"
                    @click:append-inner="copyToClipboard" class="mb-2" hide-details />
                  <p class="text-caption mb-2">
                    ⏰ Hiệu lực: {{ linkCodeExpiry }} phút
                  </p>
                  <p class="text-caption">
                    <strong>Hướng dẫn:</strong> Sao chép mã → Mở Telegram → Tìm @TabiMoneyBot → Gửi /link → Dán mã
                  </p>
                </v-card-text>
              </v-card>
            </div>

            <div v-else-if="telegramStatus === 'connected'">
              <v-btn color="error" size="small" @click="disconnectTelegram" :loading="disconnecting">
                <v-icon left size="small">mdi-link-off</v-icon>
                Hủy liên kết
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="text-h6 py-3">
            <v-icon left size="small">mdi-robot</v-icon>
            Tính năng Telegram Bot
          </v-card-title>

          <v-card-text class="pt-2 pb-3">
            <v-list density="compact">
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon size="small">mdi-chart-line</v-icon>
                </template>
                <v-list-item-title class="text-body-2">Dashboard tài chính</v-list-item-title>
                <v-list-item-subtitle class="text-caption">Xem tổng quan chi tiêu</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon size="small">mdi-chat</v-icon>
                </template>
                <v-list-item-title class="text-body-2">Chat với AI</v-list-item-title>
                <v-list-item-subtitle class="text-caption">Phân tích và tư vấn</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon size="small">mdi-bell</v-icon>
                </template>
                <v-list-item-title class="text-body-2">Thông báo</v-list-item-title>
                <v-list-item-subtitle class="text-caption">Nhận cảnh báo chi tiêu</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" :timeout="snackbar.timeout">
      {{ snackbar.message }}
      <template v-slot:actions>
        <v-btn color="white" variant="text" @click="snackbar.show = false">
          Đóng
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api, authAPI } from '@/services/api'
import { notificationPreferencesAPI } from '@/services/api'
import { useAuthStore } from '../../stores/auth'
const authStore = useAuthStore()

// Reactive data
const telegramStatus = ref('disconnected') // 'connected', 'disconnected', 'loading'
const linkCode = ref('')
const generatingCode = ref(false)
const disconnecting = ref(false)
const linkCodeExpiry = ref(10)

// Transaction settings
const largeTransactionThreshold = ref(null) // null means use default
const savingThreshold = ref(false)
const rules = {
  required: (v) => v !== null && v !== undefined && v !== '' || 'Vui lòng nhập ngưỡng',
  min: (v) => v === null || v === undefined || v >= 0 || 'Ngưỡng phải lớn hơn hoặc bằng 0'
}

// Notification preferences
const notificationPrefs = ref({
  email_enabled: true,
  telegram_enabled: true,
  in_app_enabled: true,
  push_enabled: false,
  budget_alerts: true,
  goal_alerts: true,
  ai_alerts: true,
  transaction_alerts: true,
  analytics_alerts: true,
  urgent_notifications: true,
  high_notifications: true,
  medium_notifications: true,
  low_notifications: false,
  daily_digest: false,
  weekly_digest: true,
  monthly_digest: true,
  real_time_alerts: true,
  quiet_hours_start: '22:00',
  quiet_hours_end: '08:00',
  timezone: 'Asia/Ho_Chi_Minh'
})
const saving = ref(false)
const resetting = ref(false)

const snackbar = ref({
  show: false,
  message: '',
  color: 'success',
  timeout: 3000
})

// Methods
const generateLinkCode = async () => {
  try {
    generatingCode.value = true

    const response = await api.post('/auth/telegram/generate-link-code')

    if (response.data.success) {
      linkCode.value = response.data.link_code
      linkCodeExpiry.value = response.data.expiry_minutes

      showSnackbar('Mã liên kết đã được tạo thành công!', 'success')

      setTimeout(() => {
        checkTelegramStatus()
      }, 1000)
    } else {
      showSnackbar('Không thể tạo mã liên kết. Vui lòng thử lại.', 'error')
    }
  } catch (error) {
    console.error('Error generating link code:', error)
    showSnackbar('Có lỗi xảy ra khi tạo mã liên kết.', 'error')
  } finally {
    generatingCode.value = false
  }
}

const checkTelegramStatus = async () => {
  try {
    telegramStatus.value = 'loading'

    const response = await api.get('/auth/telegram/status')

    if (response.data.connected) {
      telegramStatus.value = 'connected'
    } else {
      telegramStatus.value = 'disconnected'
    }
  } catch (error) {
    console.error('Error checking Telegram status:', error)
    telegramStatus.value = 'disconnected'
  }
}

const disconnectTelegram = async () => {
  try {
    disconnecting.value = true

    const response = await api.post('/auth/telegram/disconnect')

    if (response.data.success) {
      telegramStatus.value = 'disconnected'
      linkCode.value = ''
      showSnackbar('Đã hủy liên kết Telegram thành công!', 'success')
    } else {
      showSnackbar('Không thể hủy liên kết. Vui lòng thử lại.', 'error')
    }
  } catch (error) {
    console.error('Error disconnecting Telegram:', error)
    showSnackbar('Có lỗi xảy ra khi hủy liên kết.', 'error')
  } finally {
    disconnecting.value = false
  }
}

const copyToClipboard = async () => {
  if (!linkCode.value) {
    showSnackbar('Chưa có mã liên kết để sao chép.', 'warning')
    return
  }

  try {
    const canUseNavigatorClipboard = typeof navigator !== 'undefined'
      && navigator.clipboard
      && typeof navigator.clipboard.writeText === 'function'

    if (canUseNavigatorClipboard) {
      await navigator.clipboard.writeText(linkCode.value)
    } else if (typeof document !== 'undefined') {
      const textarea = document.createElement('textarea')
      textarea.value = linkCode.value
      textarea.setAttribute('readonly', '')
      textarea.style.position = 'absolute'
      textarea.style.left = '-9999px'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    } else {
      throw new Error('Clipboard API is not available')
    }

    showSnackbar('Đã sao chép mã liên kết!', 'success')
  } catch (error) {
    console.error('Error copying to clipboard:', error)
    showSnackbar('Không thể sao chép mã liên kết.', 'error')
  }
}

const showSnackbar = (message, color = 'success') => {
  snackbar.value = {
    show: true,
    message,
    color,
    timeout: 3000
  }
}

// Notification preferences methods
const loadNotificationPreferences = async () => {
  try {
    const response = await notificationPreferencesAPI.getPreferences()
    if (response.data?.data) {
      notificationPrefs.value = { ...notificationPrefs.value, ...response.data.data }
    }
  } catch (error) {
    console.error('Error loading notification preferences:', error)
  }
}

const savePreferences = async () => {
  try {
    saving.value = true
    await notificationPreferencesAPI.updatePreferences(notificationPrefs.value)
    showSnackbar('Cài đặt thông báo đã được lưu thành công!', 'success')
  } catch (error) {
    console.error('Error saving notification preferences:', error)
    showSnackbar('Không thể lưu cài đặt thông báo. Vui lòng thử lại.', 'error')
  } finally {
    saving.value = false
  }
}

const resetPreferences = async () => {
  try {
    resetting.value = true
    await notificationPreferencesAPI.resetToDefaults()
    await loadNotificationPreferences()
    showSnackbar('Cài đặt đã được khôi phục về mặc định!', 'success')
  } catch (error) {
    console.error('Error resetting notification preferences:', error)
    showSnackbar('Không thể khôi phục cài đặt. Vui lòng thử lại.', 'error')
  } finally {
    resetting.value = false
  }
}

// Transaction threshold methods
const loadThreshold = async () => {
  try {
    const response = await api.get('/auth/profile')
    if (response.data?.profile?.large_transaction_threshold) {
      largeTransactionThreshold.value = response.data.profile.large_transaction_threshold
    } else {
      largeTransactionThreshold.value = null
    }
  } catch (error) {
    console.error('Error loading threshold:', error)
  }
}

const saveThreshold = async () => {
  try {
    savingThreshold.value = true
    const threshold = largeTransactionThreshold.value === null || largeTransactionThreshold.value === '' 
      ? null 
      : parseFloat(largeTransactionThreshold.value)
    
    await authAPI.setLargeTransactionThreshold(threshold)
    showSnackbar(
      threshold === null 
        ? 'Đã khôi phục ngưỡng về mặc định (1,000,000 VND)' 
        : `Đã lưu ngưỡng giao dịch lớn: ${threshold.toLocaleString('vi-VN')} VND`,
      'success'
    )
  } catch (error) {
    console.error('Error saving threshold:', error)
    showSnackbar('Không thể lưu ngưỡng. Vui lòng thử lại.', 'error')
  } finally {
    savingThreshold.value = false
  }
}

const resetThreshold = () => {
  largeTransactionThreshold.value = null
  saveThreshold()
}

// Lifecycle
onMounted(() => {
  checkTelegramStatus()
  loadNotificationPreferences()
  loadThreshold()
})
</script>