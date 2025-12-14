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
    <v-row class="ai-layout" align="stretch">
      <!-- Conversation Panel -->
      <v-col cols="12" class="chat-column">
        <v-card class="chat-shell" elevation="3">
          <div class="chat-header d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-avatar size="44" color="primary" class="me-3">
                <v-icon size="28" color="white">mdi-robot</v-icon>
              </v-avatar>
              <div>
                <div class="text-h6 mb-1">Trợ lý TabiMoney AI</div>
                <div class="text-body-2 text-medium-emphasis">
                  Luôn sẵn sàng hỗ trợ tài chính của bạn
                </div>
              </div>
            </div>
            <v-chip color="success" variant="tonal" prepend-icon="mdi-circle-slice-8">
              Trực tuyến
            </v-chip>
          </div>

          <v-divider />

          <!-- Conversation Body -->
          <div class="chat-body" ref="chatMessages">
            <div v-if="messages.length === 0" class="empty-state text-center">
              <v-icon size="80" color="primary">mdi-message-processing</v-icon>
              <h3 class="text-h6 mt-4 mb-2">Bắt đầu cuộc trò chuyện đầu tiên</h3>
              <p class="text-body-2 text-medium-emphasis">
                Đặt câu hỏi về chi tiêu, kế hoạch tiết kiệm hoặc bất cứ điều gì về tài chính cá nhân.
              </p>
              <v-btn color="primary" class="mt-4" @click="sendQuickMessage(quickSuggestions[0]?.text)">
                Gợi ý câu hỏi
              </v-btn>
            </div>

            <div v-for="(message, index) in messages" :key="index" class="message"
              :class="message.isUser ? 'user-message' : 'ai-message'">
              <div class="message-meta">
                <div class="d-flex align-center">
                  <v-icon size="18" :color="message.isUser ? 'primary' : 'secondary'" class="me-1">
                    {{ message.isUser ? 'mdi-account' : 'mdi-robot' }}
                  </v-icon>
                  <span>{{ message.isUser ? 'Bạn' : 'Tabi AI' }}</span>
                  <v-chip
                    v-if="!message.isUser && message.intent && message.intent !== 'error' && message.intent !== 'general'"
                    size="x-small"
                    variant="tonal"
                    color="primary"
                    class="ms-2 intent-badge"
                  >
                    {{ formatIntent(message.intent) }}
                  </v-chip>
                </div>
                <span class="text-caption text-medium-emphasis">{{ formatTime(message.timestamp) }}</span>
              </div>

              <div class="message-bubble" :class="{ 'user-bubble': message.isUser }">
                <div class="text-body-1" v-html="formatMessage(message.content)"></div>
                
                <!-- AI Suggestions - Context-aware suggestions từ AI -->
                <div v-if="!message.isUser && message.suggestions && message.suggestions.length > 0" class="message-suggestions">
                  <div class="text-caption mb-2 opacity-75 d-flex align-center">
                    <v-icon size="14" class="me-1">mdi-lightbulb-outline</v-icon>
                    Gợi ý tiếp theo:
                  </div>
                  <div class="suggestion-chips">
                    <v-chip
                      v-for="(suggestion, idx) in message.suggestions.slice(0, 4)"
                      :key="idx"
                      size="small"
                      variant="outlined"
                      class="suggestion-chip"
                      @click="sendQuickMessage(suggestion)"
                      prepend-icon="mdi-arrow-right"
                    >
                      {{ suggestion }}
                    </v-chip>
                  </div>
                </div>
              </div>
            </div>

            <!-- Typing Indicator -->
            <div v-if="isTyping" class="message ai-message">
              <div class="message-meta">
                <div class="d-flex align-center">
                  <v-icon size="18" color="secondary" class="me-1">mdi-robot</v-icon>
                  <span>Tabi AI</span>
                </div>
                <span class="text-caption text-medium-emphasis">Đang soạn...</span>
              </div>
              <div class="message-bubble">
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </div>

          <v-divider />

          <!-- Composer -->
          <div class="chat-composer">
            <v-form @submit.prevent="sendMessage" class="w-100">
              <div class="d-flex composer-inner">
                <v-textarea v-model="inputMessage" placeholder="Hỏi tôi về tài chính của bạn..." minlength="1"
                  rows="1" max-rows="4" auto-grow variant="outlined" class="flex-grow-1 me-3 composer-textarea"
                  :disabled="isTyping" @keydown.enter.exact.prevent="sendMessage"
                  @keydown.enter.shift.exact="inputMessage += '\n'" hide-details></v-textarea>
                <v-btn type="submit" color="primary" class="composer-send" size="large" :loading="isTyping"
                  :disabled="!inputMessage.trim()">
                  <v-icon>mdi-send</v-icon>
                </v-btn>
              </div>
            </v-form>
          </div>
        </v-card>
      </v-col>

    </v-row>

    <!-- Insight Section -->
    <v-row class="insight-row mt-6" dense>
      <v-col cols="12" md="6">
        <v-card class="insight-card suggestion-card" variant="flat">
          <v-card-title class="d-flex align-center justify-space-between">
            <div>
              <div class="text-subtitle-1">Câu hỏi gợi ý</div>
              <div class="text-caption text-medium-emphasis">
                Những câu lệnh AI xử lý tốt nhất
              </div>
            </div>
            <v-icon color="primary">mdi-comment-text-multiple</v-icon>
          </v-card-title>
          <v-divider />
          <div class="suggestion-grid">
            <v-sheet v-for="(suggestion, index) in quickSuggestions" :key="index" class="suggestion-pill"
              color="primary" variant="tonal" @click="sendQuickMessage(suggestion.text)" rounded="lg">
              <v-icon size="18" class="me-2">{{ suggestion.icon || 'mdi-lightbulb' }}</v-icon>
              <div class="text-body-2 flex-grow-1">
                {{ suggestion.text }}
              </div>
              <v-icon size="18">mdi-arrow-top-right</v-icon>
            </v-sheet>
          </div>
        </v-card>
      </v-col>

      <v-col cols="12" md="6">
        <v-card class="insight-card transactions-card" variant="flat">
          <v-card-title class="d-flex align-center justify-space-between">
            <div>
              <div class="text-subtitle-1">Giao dịch gần đây</div>
              <div class="text-caption text-medium-emphasis">
                Hữu ích khi hỏi về ngân sách và dòng tiền
              </div>
            </div>
            <v-icon color="primary">mdi-chart-donut</v-icon>
          </v-card-title>
          <v-divider />
          <div class="transaction-list" v-if="recentTransactions.length">
            <div v-for="transaction in recentTransactions.slice(0, 4)" :key="transaction.id"
              class="transaction-row">
              <div class="transaction-icon" :class="transaction.transaction_type === 'income' ? 'is-income' : 'is-expense'">
                <v-icon size="18">
                  {{ transaction.transaction_type === 'income' ? 'mdi-arrow-up-bold' : 'mdi-arrow-down-bold' }}
                </v-icon>
              </div>
              <div class="transaction-info">
                <div class="transaction-title">
                  {{ transaction.description || 'Không có mô tả' }}
                </div>
                <div class="transaction-date text-caption">
                  {{ formatDate(transaction.transaction_date) }}
                </div>
              </div>
              <div class="transaction-amount"
                :class="transaction.transaction_type === 'income' ? 'text-success' : 'text-error'">
                {{ transaction.transaction_type === 'income' ? '+' : '-' }}{{ formatCurrency(transaction.amount) }}
              </div>
            </div>
          </div>
          <div v-else class="transaction-empty text-center text-medium-emphasis py-6">
            Chưa có giao dịch nào
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { aiAPI, transactionAPI } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { formatCurrency, formatDate } from '@/utils/formatters'
import { nextTick, onMounted, ref, watch } from 'vue'

// Reactive data
const messages = ref([])

const inputMessage = ref('')
const isTyping = ref(false)
const chatMessages = ref(null)

// Quick suggestions - General suggestions, không thay đổi theo context
// Message suggestions sẽ được AI tạo context-aware và hiển thị dưới mỗi message
const quickSuggestions = ref([
  { text: 'Tôi đã chi bao nhiêu tiền tháng này?', icon: 'mdi-cash' },
  { text: 'Danh mục nào tôi chi nhiều nhất?', icon: 'mdi-chart-pie' },
  { text: 'Tình hình ngân sách của tôi thế nào?', icon: 'mdi-wallet' },
  { text: 'Tiến độ mục tiêu của tôi?', icon: 'mdi-target' },
  { text: 'Gợi ý tiết kiệm cho tôi', icon: 'mdi-lightbulb' },
  { text: 'Dự đoán chi tiêu tháng tới', icon: 'mdi-trending-up' }
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

    // Add AI response with suggestions
    const aiResponse = response.data.response || 'Xin lỗi, tôi không thể trả lời câu hỏi này.'
    const suggestions = response.data.suggestions || []
    
    messages.value.push({
      content: aiResponse,
      isUser: false,
      timestamp: new Date(),
      suggestions: suggestions,
      intent: response.data.intent,
      entities: response.data.entities || []
    })

    // Không cập nhật quick suggestions từ AI suggestions
    // Quick suggestions nên giữ nguyên (general suggestions)
    // Message suggestions đã context-aware và hiển thị dưới mỗi message

  } catch (error) {
    console.error('Chat error:', error)
    const errorMessage = error.response?.data?.detail || error.message || 'Xin lỗi, có lỗi xảy ra khi xử lý tin nhắn của bạn.'
    messages.value.push({
      content: errorMessage,
      isUser: false,
      timestamp: new Date(),
      suggestions: ['Thử lại', 'Xem số dư', 'Thêm giao dịch mới']
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
    const container = chatMessages.value?.$el ?? chatMessages.value
    if (container) {
      container.scrollTop = container.scrollHeight
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

const formatIntent = (intent) => {
  const intentMap = {
    'add_transaction': 'Thêm giao dịch',
    'query_balance': 'Truy vấn số dư',
    'analyze_data': 'Phân tích dữ liệu',
    'budget_management': 'Quản lý ngân sách',
    'goal_tracking': 'Theo dõi mục tiêu',
    'smart_recommendations': 'Gợi ý thông minh',
    'expense_forecasting': 'Dự đoán chi tiêu',
    'general': 'Chung',
    'error': 'Lỗi'
  }
  return intentMap[intent] || intent
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

// Auto-scroll whenever messages change or typing indicator appears
watch(messages, () => {
  scrollToBottom()
})

watch(isTyping, (typing) => {
  if (typing) {
    scrollToBottom()
  }
})
</script>

<style scoped>
.ai-layout {
  gap: 24px;
}

.chat-column {
  display: flex;
}

.chat-shell {
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: 720px;
  height: calc(100vh - 150px);
}

.chat-header {
  padding: 24px;
}

.chat-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  background: linear-gradient(180deg, rgba(248, 249, 250, 0.8), rgba(248, 249, 250, 0.2));
}

.empty-state {
  padding: 80px 24px;
}

.message {
  max-width: 90%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.message.user-message {
  margin-left: auto;
}

.message-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: rgba(0, 0, 0, 0.6);
}

.message.user-message .message-meta {
  flex-direction: row-reverse;
  text-align: right;
}

.message-bubble {
  padding: 16px 20px;
  border-radius: 20px;
  background: white;
  box-shadow: 0 6px 24px rgba(15, 23, 42, 0.08);
  line-height: 1.5;
}

.message-bubble.user-bubble,
.message.user-message .message-bubble {
  background: #1e88e5;
  color: white;
}

.chat-composer {
  padding: 18px 24px;
  background: white;
}

.composer-inner {
  gap: 12px;
}

.composer-textarea :deep(textarea) {
  font-size: 16px;
  line-height: 1.4;
}

.composer-send {
  min-width: 64px;
}

.typing-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: currentColor;
  opacity: 0.7;
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

.insight-row {
  margin-top: 32px;
}

.insight-card {
  width: 100%;
  border-radius: 18px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
  background: #ffffff;
}

.suggestion-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
  padding: 16px 20px 24px;
}

.suggestion-pill {
  display: flex;
  align-items: center;
  padding: 12px 14px;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.suggestion-pill:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 18px rgba(30, 136, 229, 0.18);
}

.transaction-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
}

.transaction-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: #f9fafc;
}

.transaction-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.transaction-icon.is-income {
  background: rgba(46, 125, 50, 0.12);
  color: #2e7d32;
}

.transaction-icon.is-expense {
  background: rgba(229, 57, 53, 0.12);
  color: #d32f2f;
}

.transaction-info {
  flex: 1;
  min-width: 0;
}

.transaction-title {
  font-weight: 500;
  margin-bottom: 2px;
}

.transaction-amount {
  font-weight: 600;
}

.message-suggestions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.message-bubble.user-bubble .message-suggestions {
  border-top-color: rgba(255, 255, 255, 0.3);
}

.suggestion-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.suggestion-chip {
  cursor: pointer;
  transition: all 0.2s ease;
}

.suggestion-chip:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.intent-badge {
  font-size: 10px;
  height: 18px;
  padding: 0 6px;
}

@media (max-width: 1280px) {
  .chat-shell {
    min-height: 580px;
    height: auto;
  }

  .message {
    max-width: 100%;
  }

  .suggestion-grid {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .insight-row {
    margin-top: 24px;
  }
}
</style>
