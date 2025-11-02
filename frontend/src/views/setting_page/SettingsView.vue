<template>
  <v-container class="py-8">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">Cài đặt</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="8">
        <!-- Notification Preferences -->
        <v-card class="mb-4">
          <v-card-title>
            <v-icon left>mdi-bell-outline</v-icon>
            Cài đặt thông báo
          </v-card-title>

          <v-card-text>
            <v-expansion-panels variant="accordion">
              <!-- Channel Preferences -->
              <v-expansion-panel title="Kênh thông báo">
                <v-expansion-panel-text>
                  <v-switch
                    v-model="notificationPrefs.email_enabled"
                    label="Email"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.telegram_enabled"
                    label="Telegram"
                    color="primary"
                    hide-details
                    :disabled="telegramStatus !== 'connected'"
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.in_app_enabled"
                    label="Trong ứng dụng"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-alert v-if="telegramStatus !== 'connected'" type="info" variant="tonal" density="compact" class="mt-2">
                    Kích hoạt Telegram ở phần dưới để nhận thông báo qua Telegram
                  </v-alert>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <!-- Feature Preferences -->
              <v-expansion-panel title="Loại thông báo">
                <v-expansion-panel-text>
                  <v-switch
                    v-model="notificationPrefs.budget_alerts"
                    label="Cảnh báo ngân sách"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.goal_alerts"
                    label="Cảnh báo mục tiêu"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.ai_alerts"
                    label="Cảnh báo AI"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.transaction_alerts"
                    label="Cảnh báo giao dịch"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.analytics_alerts"
                    label="Báo cáo phân tích"
                    color="primary"
                    hide-details
                  ></v-switch>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <!-- Priority Preferences -->
              <v-expansion-panel title="Mức độ ưu tiên">
                <v-expansion-panel-text>
                  <v-switch
                    v-model="notificationPrefs.urgent_notifications"
                    label="Khẩn cấp"
                    color="red"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.high_notifications"
                    label="Cao"
                    color="orange"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.medium_notifications"
                    label="Trung bình"
                    color="yellow"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.low_notifications"
                    label="Thấp"
                    color="grey"
                    hide-details
                  ></v-switch>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <!-- Frequency Preferences -->
              <v-expansion-panel title="Tần suất">
                <v-expansion-panel-text>
                  <v-switch
                    v-model="notificationPrefs.daily_digest"
                    label="Báo cáo hàng ngày"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.weekly_digest"
                    label="Báo cáo hàng tuần"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.monthly_digest"
                    label="Báo cáo hàng tháng"
                    color="primary"
                    hide-details
                  ></v-switch>
                  <v-switch
                    v-model="notificationPrefs.real_time_alerts"
                    label="Cảnh báo thời gian thực"
                    color="primary"
                    hide-details
                  ></v-switch>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>

            <v-divider class="my-4"></v-divider>

            <div class="d-flex justify-space-between align-center">
              <v-btn color="success" @click="testNotification" :loading="testingNotification">
                <v-icon left>mdi-bell-ring</v-icon>
                Gửi thông báo test
              </v-btn>
              <div>
                <v-btn color="secondary" @click="resetPreferences" :loading="resetting">
                  Khôi phục mặc định
                </v-btn>
                <v-btn color="primary" @click="savePreferences" :loading="saving" class="ml-2">
                  Lưu cài đặt
                </v-btn>
              </div>
            </div>
          </v-card-text>
        </v-card>

        <!-- Telegram Integration -->
        <v-card>
          <v-card-title>
            <v-icon left>mdi-telegram</v-icon>
            Tích hợp Telegram Bot
          </v-card-title>

          <v-card-text>
            <p class="text-body-1 mb-4">
              Liên kết tài khoản TabiMoney với Telegram Bot để sử dụng các tính năng AI và dashboard trên Telegram.
            </p>

            <v-alert v-if="telegramStatus === 'connected'" type="success" class="mb-4">
              <v-icon left>mdi-check-circle</v-icon>
              Tài khoản đã được liên kết với Telegram Bot
            </v-alert>

            <v-alert v-else-if="telegramStatus === 'disconnected'" type="info" class="mb-4">
              <v-icon left>mdi-information</v-icon>
              Tài khoản chưa được liên kết với Telegram Bot
            </v-alert>

            <div v-if="telegramStatus === 'disconnected'">
              <v-btn color="primary" @click="generateLinkCode" :loading="generatingCode" class="mb-4">
                <v-icon left>mdi-link</v-icon>
                Tạo mã liên kết
              </v-btn>

              <v-card v-if="linkCode" class="mt-4" color="primary" variant="outlined">
                <v-card-text>
                  <h3 class="text-h6 mb-2">Mã liên kết của bạn:</h3>
                  <v-text-field :value="linkCode" readonly variant="outlined" append-icon="mdi-content-copy"
                    @click:append="copyToClipboard" class="mb-2" />
                  <p class="text-caption">
                    ⏰ Mã này có hiệu lực trong {{ linkCodeExpiry }} phút
                  </p>
                  <p class="text-body-2">
                    <strong>Hướng dẫn:</strong><br>
                    1. Sao chép mã liên kết ở trên<br>
                    2. Mở Telegram và tìm bot @TabiMoneyBot<br>
                    3. Gửi lệnh /link<br>
                    4. Dán mã liên kết vào bot<br>
                    5. Hoàn tất liên kết!
                  </p>
                </v-card-text>
              </v-card>
            </div>

            <div v-else-if="telegramStatus === 'connected'">
              <v-btn color="error" @click="disconnectTelegram" :loading="disconnecting">
                <v-icon left>mdi-link-off</v-icon>
                Hủy liên kết
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-robot</v-icon>
            Tính năng Telegram Bot
          </v-card-title>

          <v-card-text>
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-chart-line</v-icon>
                </template>
                <v-list-item-title>Dashboard tài chính</v-list-item-title>
                <v-list-item-subtitle>Xem tổng quan chi tiêu</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-chat</v-icon>
                </template>
                <v-list-item-title>Chat với AI</v-list-item-title>
                <v-list-item-subtitle>Phân tích và tư vấn tài chính</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-bell</v-icon>
                </template>
                <v-list-item-title>Thông báo</v-list-item-title>
                <v-list-item-subtitle>Nhận cảnh báo chi tiêu</v-list-item-subtitle>
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
import { useAuthStore } from '@/stores/auth'
import { api } from '@/services/api'
import { notificationPreferencesAPI } from '@/services/api'

const authStore = useAuthStore()

// Reactive data
const telegramStatus = ref('disconnected') // 'connected', 'disconnected', 'loading'
const linkCode = ref('')
const generatingCode = ref(false)
const disconnecting = ref(false)
const linkCodeExpiry = ref(10)

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
const testingNotification = ref(false)

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

      // Auto refresh status after code generation
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
  try {
    await navigator.clipboard.writeText(linkCode.value)
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

const testNotification = async () => {
  try {
    testingNotification.value = true
    await notificationPreferencesAPI.testNotification('in_app')
    showSnackbar('Thông báo test đã được gửi! Kiểm tra hộp thông báo ở góc trên.', 'success')
    // Reload notifications after a short delay to see the new notification
    setTimeout(() => {
      // Trigger notification reload in App.vue by dispatching event or window reload notifications
      window.dispatchEvent(new Event('notification:refresh'))
    }, 500)
  } catch (error) {
    console.error('Error sending test notification:', error)
    showSnackbar('Không thể gửi thông báo test. Vui lòng thử lại.', 'error')
  } finally {
    testingNotification.value = false
  }
}

// Lifecycle
onMounted(() => {
  checkTelegramStatus()
  loadNotificationPreferences()
})
</script>