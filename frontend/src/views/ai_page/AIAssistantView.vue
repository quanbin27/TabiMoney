<template>
  <v-container fluid class="pa-4">
    <!-- Header -->
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card color="primary" dark>
          <v-card-title class="d-flex align-center">
            <v-icon class="me-3" size="large">mdi-robot</v-icon>
            <div>
              <h2 class="text-h5">AI Assistant</h2>
              <p class="text-subtitle-1 mb-0 opacity-90">
                Hỏi tôi về tài chính cá nhân của bạn
              </p>
            </div>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <!-- Chat Interface -->
    <v-row>
      <v-col cols="12" md="8">
        <v-card class="chat-container" height="600">
          <!-- Chat Messages -->
          <v-card-text class="chat-messages pa-4" ref="chatMessages">
            <div v-if="messages.length === 0" class="text-center pa-8">
              <v-icon size="64" color="grey-lighten-1">mdi-chat-outline</v-icon>
              <h3 class="text-h6 mt-4 text-grey">Chào mừng đến với AI Assistant!</h3>
              <p class="text-body-2 text-grey">Hãy hỏi tôi về tài chính của bạn</p>
            </div>

            <div v-for="(message, index) in messages" :key="index" class="message mb-4"
              :class="message.isUser ? 'user-message' : 'ai-message'">
              <div class="d-flex" :class="message.isUser ? 'justify-end' : 'justify-start'">
                <v-card :color="message.isUser ? 'primary' : 'grey-lighten-4'" :dark="message.isUser"
                  class="message-bubble pa-3" max-width="70%">
                  <div class="d-flex align-start">
                    <v-avatar :color="message.isUser ? 'white' : 'primary'" size="32" class="me-3">
                      <v-icon :color="message.isUser ? 'primary' : 'white'">
                        {{ message.isUser ? 'mdi-account' : 'mdi-robot' }}
                      </v-icon>
                    </v-avatar>
                    <div class="flex-grow-1">
                      <div class="text-body-1" v-html="formatMessage(message.content)"></div>
                      <div class="text-caption mt-1 opacity-70">
                        {{ formatTime(message.timestamp) }}
                      </div>
                    </div>
                  </div>
                </v-card>
              </div>
            </div>

            <!-- Typing Indicator -->
            <div v-if="isTyping" class="ai-message mb-4">
              <div class="d-flex justify-start">
                <v-card color="grey-lighten-4" class="message-bubble pa-3">
                  <div class="d-flex align-center">
                    <v-avatar color="primary" size="32" class="me-3">
                      <v-icon color="white">mdi-robot</v-icon>
                    </v-avatar>
                    <div class="typing-indicator">
                      <span></span>
                      <span></span>
                      <span></span>
                    </div>
                  </div>
                </v-card>
              </div>
            </div>
          </v-card-text>

          <!-- Input Area -->
          <v-divider></v-divider>
          <v-card-text class="pa-4">
            <v-form @submit.prevent="sendMessage">
              <v-row align="center">
                <v-col cols="12">
                  <v-textarea v-model="inputMessage" placeholder="Hỏi tôi về tài chính của bạn..." rows="2" auto-grow
                    variant="outlined" :disabled="isTyping" @keydown.enter.exact.prevent="sendMessage"
                    @keydown.enter.shift.exact="inputMessage += '\n'" hide-details></v-textarea>
                </v-col>
                <v-col cols="auto">
                  <v-btn type="submit" color="primary" :loading="isTyping" :disabled="!inputMessage.trim()"
                    size="large">
                    <v-icon>mdi-send</v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Quick Actions Sidebar -->
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>
            <v-icon class="me-2">mdi-lightning-bolt</v-icon>
            Câu hỏi nhanh
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item v-for="(suggestion, index) in quickSuggestions" :key="index"
                @click="sendQuickMessage(suggestion)" class="mb-2">
                <template v-slot:prepend>
                  <v-icon color="primary">mdi-chat</v-icon>
                </template>
                <v-list-item-title class="text-body-2">
                  {{ suggestion }}
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>

        <!-- Recent Transactions Context -->
        <v-card class="mt-4">
          <v-card-title>
            <v-icon class="me-2">mdi-history</v-icon>
            Giao dịch gần đây
          </v-card-title>
          <v-card-text>
            <div v-if="recentTransactions.length > 0">
              <v-list density="compact">
                <v-list-item v-for="transaction in recentTransactions.slice(0, 3)" :key="transaction.id" class="px-0">
                  <template v-slot:prepend>
                    <v-icon :color="transaction.transaction_type === 'income' ? 'success' : 'error'" size="small">
                      {{ transaction.transaction_type === 'income' ? 'mdi-arrow-up' : 'mdi-arrow-down' }}
                    </v-icon>
                  </template>
                  <v-list-item-title class="text-body-2">
                    {{ transaction.description || 'Không có mô tả' }}
                  </v-list-item-title>
                  <v-list-item-subtitle class="text-caption">
                    {{ formatCurrency(transaction.amount) }} - {{ formatDate(transaction.transaction_date) }}
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </div>
            <div v-else class="text-center pa-4 text-grey">
              <v-icon>mdi-information</v-icon>
              <p class="text-caption mt-2">Chưa có giao dịch nào</p>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { aiAPI, transactionAPI } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { formatCurrency, formatDate } from '@/utils/formatters'
import { nextTick, onMounted, ref } from 'vue'

// Reactive data
const messages = ref([])

const inputMessage = ref('')
const isTyping = ref(false)
const chatMessages = ref(null)

// Quick suggestions
const quickSuggestions = ref([
  'Tôi đã chi bao nhiêu tiền tháng này?',
  'Danh mục nào tôi chi nhiều nhất?',
  'Tôi có thể tiết kiệm được bao nhiêu?',
  'Dự đoán chi tiêu tháng tới',
  'Gợi ý cách quản lý tài chính',
  'Phân tích xu hướng chi tiêu'
])

const recentTransactions = ref([])

// Auth store
const authStore = useAuthStore()

// Methods
const sendMessage = async () => {
  if (!inputMessage.value.trim() || isTyping.value) return

  const userMessage = inputMessage.value.trim()
  inputMessage.value = ''

  // Add user message
  messages.value.push({
    content: userMessage,
    isUser: true,
    timestamp: new Date()
  })

  scrollToBottom()
  isTyping.value = true

  try {
    const response = await aiAPI.processChat({
      message: userMessage,
      user_id: authStore.user?.id || 0 // Use actual user ID from auth store
    })

    // Simulate typing delay
    await new Promise(resolve => setTimeout(resolve, 1000))

    // Add AI response
    messages.value.push({
      content: response.data.response || 'Xin lỗi, tôi không thể trả lời câu hỏi này.',
      isUser: false,
      timestamp: new Date()
    })

  } catch (error) {
    console.error('Chat error:', error)
    messages.value.push({
      content: 'Xin lỗi, có lỗi xảy ra khi xử lý tin nhắn của bạn.',
      isUser: false,
      timestamp: new Date()
    })
  } finally {
    isTyping.value = false
    scrollToBottom()
  }
}

const sendQuickMessage = (suggestion) => {
  inputMessage.value = suggestion
  sendMessage()
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatMessages.value) {
      chatMessages.value.scrollTop = chatMessages.value.scrollHeight
    }
  })
}

const formatMessage = (content) => {
  // Simple formatting for better readability
  return content
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
}

const formatTime = (date) => {
  return date.toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadRecentTransactions = async () => {
  try {
    const response = await transactionAPI.getTransactions({
      limit: 5,
      sort_by: 'transaction_date',
      sort_order: 'desc'
    })
    recentTransactions.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load recent transactions:', error)
  }
}

// Lifecycle
onMounted(() => {
  loadRecentTransactions()

  // Add welcome message
  messages.value.push({
    content: 'Xin chào! Tôi là AI Assistant của TabiMoney. Tôi có thể giúp bạn phân tích tài chính, đưa ra gợi ý và trả lời các câu hỏi về quản lý tiền bạc. Hãy hỏi tôi bất cứ điều gì!',
    isUser: false,
    timestamp: new Date()
  })
})
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  max-height: 500px;
}

.message-bubble {
  border-radius: 18px;
}

.user-message .message-bubble {
  border-bottom-right-radius: 4px;
}

.ai-message .message-bubble {
  border-bottom-left-radius: 4px;
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #666;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {

  0%,
  80%,
  100% {
    transform: scale(0.8);
    opacity: 0.5;
  }

  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.v-list-item {
  border-radius: 8px;
  margin-bottom: 4px;
}

.v-list-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}
</style>
