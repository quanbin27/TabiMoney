<template>
  <v-container class="py-8">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">Cài đặt</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="8">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-telegram</v-icon>
            Tích hợp Telegram Bot
          </v-card-title>
          
          <v-card-text>
            <p class="text-body-1 mb-4">
              Liên kết tài khoản TabiMoney với Telegram Bot để sử dụng các tính năng AI và dashboard trên Telegram.
            </p>

            <v-alert
              v-if="telegramStatus === 'connected'"
              type="success"
              class="mb-4"
            >
              <v-icon left>mdi-check-circle</v-icon>
              Tài khoản đã được liên kết với Telegram Bot
            </v-alert>

            <v-alert
              v-else-if="telegramStatus === 'disconnected'"
              type="info"
              class="mb-4"
            >
              <v-icon left>mdi-information</v-icon>
              Tài khoản chưa được liên kết với Telegram Bot
            </v-alert>

            <div v-if="telegramStatus === 'disconnected'">
              <v-btn
                color="primary"
                @click="generateLinkCode"
                :loading="generatingCode"
                class="mb-4"
              >
                <v-icon left>mdi-link</v-icon>
                Tạo mã liên kết
              </v-btn>

              <v-card v-if="linkCode" class="mt-4" color="primary" variant="outlined">
                <v-card-text>
                  <h3 class="text-h6 mb-2">Mã liên kết của bạn:</h3>
                  <v-text-field
                    :value="linkCode"
                    readonly
                    variant="outlined"
                    append-icon="mdi-content-copy"
                    @click:append="copyToClipboard"
                    class="mb-2"
                  />
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
              <v-btn
                color="error"
                @click="disconnectTelegram"
                :loading="disconnecting"
              >
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
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="snackbar.timeout"
    >
      {{ snackbar.message }}
      <template v-slot:actions>
        <v-btn
          color="white"
          variant="text"
          @click="snackbar.show = false"
        >
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

const authStore = useAuthStore()

// Reactive data
const telegramStatus = ref('disconnected') // 'connected', 'disconnected', 'loading'
const linkCode = ref('')
const generatingCode = ref(false)
const disconnecting = ref(false)
const linkCodeExpiry = ref(10)

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

// Lifecycle
onMounted(() => {
  checkTelegramStatus()
})
</script>
