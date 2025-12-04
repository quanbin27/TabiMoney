<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2">mdi-chart-line</v-icon>
            Bảng phân tích tài chính
          </v-card-title>
          <v-card-subtitle>
            Phân tích chi tiêu và gợi ý từ AI
          </v-card-subtitle>
        </v-card>
      </v-col>
    </v-row>

    <!-- Monthly Income Settings -->
    <v-row class="mt-2">
      <v-col cols="12" md="6">
        <v-card variant="outlined">
          <v-card-title>
            <v-icon class="mr-2">mdi-currency-usd</v-icon>
            Thu nhập hàng tháng
          </v-card-title>
          <v-card-text class="d-flex align-center gap-2">
            <v-text-field v-model.number="monthlyIncome" type="number" label="Thu nhập hàng tháng" prefix="₫"
              density="comfortable" hide-details class="mr-2" />
            <v-btn color="primary" :loading="incomeSaving" @click="saveIncome">Lưu</v-btn>
            <div class="ml-4 text-caption" v-if="dashboardData">
              Đã chi tháng này: <strong>{{ percentOfIncome.toFixed(1) }}%</strong>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card variant="outlined">
          <v-card-title>
            <v-icon class="mr-2">mdi-compare</v-icon>
            So sánh tháng với tháng trước
          </v-card-title>
          <v-card-text>
            <div v-if="mom">
              <div class="d-flex justify-space-between">
                <div>
                  <div class="text-caption">Tháng hiện tại</div>
                  <div class="text-body-2">Thu nhập: {{ formatCurrency(mom.current.income) }}</div>
                  <div class="text-body-2">Chi tiêu: {{ formatCurrency(mom.current.expense) }}</div>
                </div>
                <div class="text-right">
                  <div class="text-caption">Tháng trước</div>
                  <div class="text-body-2">Thu nhập: {{ formatCurrency(mom.prev.income) }}</div>
                  <div class="text-body-2">Chi tiêu: {{ formatCurrency(mom.prev.expense) }}</div>
                </div>
              </div>
              <v-divider class="my-2" />
              <div>
                <div>Thay đổi thu nhập: <strong :class="mom.incomeChange >= 0 ? 'text-success' : 'text-error'">{{
                  mom.incomeChange.toFixed(1) }}%</strong></div>
                <div>Thay đổi chi tiêu: <strong :class="mom.expenseChange >= 0 ? 'text-error' : 'text-success'">{{
                  mom.expenseChange.toFixed(1) }}%</strong></div>
              </div>
            </div>
            <div v-else class="text-caption">Chọn tháng để so sánh.</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Date Range Filter -->
    <v-row>
      <v-col cols="12" md="3">
        <v-select v-model="selectedPeriod" :items="periodOptions" label="Khoảng thời gian"
          @update:model-value="loadAnalytics" />
      </v-col>
      <v-col cols="12" md="3">
        <v-text-field v-model="startDate" type="date" label="Từ ngày" @update:model-value="loadAnalytics" />
      </v-col>
      <v-col cols="12" md="3">
        <v-text-field v-model="endDate" type="date" label="Đến ngày" @update:model-value="loadAnalytics" />
      </v-col>
      <v-col cols="12" md="3" class="d-flex align-center">
        <v-btn color="primary" @click="loadAnalytics" :loading="loading">
          <v-icon class="mr-2">mdi-refresh</v-icon>
          Làm mới
        </v-btn>
      </v-col>
    </v-row>

    <!-- Transaction Filters (Interactive) -->
    <v-row class="mt-1">
      <v-col cols="12" md="2">
        <v-select :items="txTypeOptions" v-model="txType" label="Loại giao dịch" clearable />
      </v-col>
      <v-col cols="12" md="2">
        <v-text-field v-model.number="minAmount" type="number" label="Số tiền tối thiểu" />
      </v-col>
      <v-col cols="12" md="2">
        <v-text-field v-model.number="maxAmount" type="number" label="Số tiền tối đa" />
      </v-col>
      <v-col cols="12" md="4">
        <v-text-field v-model="search" label="Tìm theo mô tả/địa điểm" clearable />
      </v-col>
      <v-col cols="12" md="2" class="d-flex align-center">
        <v-btn color="secondary" :loading="txLoading" @click="loadFilteredTransactions">Lọc</v-btn>
      </v-col>
    </v-row>

    <!-- Financial Summary Cards -->
    <v-row v-if="dashboardData">
      <v-col cols="12" md="3">
        <v-card color="success" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-success">Tổng thu nhập</div>
            <div class="text-h4">{{ formatCurrency(dashboardData.total_income) }}</div>
            <div class="text-caption">{{ dashboardData.transaction_count }} giao dịch</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="error" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-error">Tổng chi tiêu</div>
            <div class="text-h4">{{ formatCurrency(dashboardData.total_expense) }}</div>
            <div class="text-caption">{{ dashboardData.transaction_count }} giao dịch</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card :color="dashboardData.net_amount >= 0 ? 'success' : 'error'" variant="tonal">
          <v-card-text>
            <div class="text-h6" :class="dashboardData.net_amount >= 0 ? 'text-success' : 'text-error'">
              Số dư ròng
            </div>
            <div class="text-h4">{{ formatCurrency(dashboardData.net_amount) }}</div>
            <div class="text-caption">Thu nhập - Chi tiêu</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="info" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-info">Tỷ lệ tiết kiệm</div>
            <div class="text-h4">{{ dashboardData.financial_health?.savings_rate?.toFixed(1) || '0.0' }}%</div>
            <div class="text-caption">Sức khỏe tài chính</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Filtered Transactions Summary -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-filter</v-icon>
            Giao dịch sau khi lọc
          </v-card-title>
          <v-card-text class="d-flex align-center justify-space-between">
            <div>
              <div class="text-body-2">Kết quả: <strong>{{ filteredTotal }}</strong></div>
              <div class="text-caption">Hiển thị {{ Math.min(10, filteredItems.length) }} dòng mới nhất trên tổng số
                {{ filteredTotal }}</div>
            </div>
              <div class="text-caption" v-if="filteredItems.length">
              Tổng: <strong>{{formatCurrency(filteredItems.reduce((a, t) => a + Number(t.amount || 0), 0))}}</strong>
            </div>
          </v-card-text>
          <v-divider />
          <v-card-text>
            <v-list density="compact">
            <v-list-item v-for="t in filteredItems.slice(0, 10)" :key="t.id">
                <v-list-item-title>{{ t.description || '(không có mô tả)' }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ t.transaction_type }} · {{ t.category?.name || 'Chưa phân loại' }} · {{ formatCurrency(t.amount) }}
                  · {{ new Date(t.transaction_date).toLocaleDateString() }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Charts Row -->
    <v-row>
      <!-- Category Spending Chart -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-chart-donut</v-icon>
            Chi tiêu theo danh mục
          </v-card-title>
          <v-card-text>
            <div v-if="categorySpendingChartData" style="height: 300px;">
              <DoughnutChart :data="categorySpendingChartData" :options="doughnutChartOptions" />
            </div>
            <div v-else class="text-center pa-8">
              <v-icon size="48" color="grey">mdi-chart-donut</v-icon>
              <div class="text-h6 mt-2">Chưa có dữ liệu chi tiêu</div>
              <div class="text-caption">Thêm một vài giao dịch để xem phân tích</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Spending Trends -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-chart-line</v-icon>
            Xu hướng chi tiêu
          </v-card-title>
          <v-card-text>
            <div v-if="spendingPatterns && spendingPatterns.insights && spendingPatterns.insights.length > 0"
              style="height: 300px;">
              <div class="text-h6 mb-4">Nhận định từ AI</div>
              <v-alert v-for="(insight, index) in spendingPatterns.insights" :key="index" :type="insight.type"
                variant="tonal" class="mb-2">
                {{ insight.description }}
              </v-alert>
            </div>
            <div v-else class="text-center pa-8">
              <v-icon size="48" color="grey">mdi-chart-line</v-icon>
              <div class="text-h6 mt-2">Chưa có dữ liệu xu hướng</div>
              <div class="text-caption">AI đang phân tích thói quen chi tiêu của bạn</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- AI Insights and Anomalies -->
    <v-row>
      <!-- AI Recommendations -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-robot</v-icon>
            Gợi ý từ AI
          </v-card-title>
          <v-card-text>
            <div
              v-if="spendingPatterns && spendingPatterns.recommendations && spendingPatterns.recommendations.length > 0">
              <v-list>
                <v-list-item v-for="(rec, index) in spendingPatterns.recommendations" :key="index"
                  :prepend-icon="getRecommendationIcon(rec.priority)">
                  <v-list-item-title>{{ rec.title }}</v-list-item-title>
                  <v-list-item-subtitle>{{ rec.description }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </div>
            <div v-else class="text-center pa-4">
              <v-icon size="32" color="grey">mdi-robot</v-icon>
              <div class="text-body-2 mt-2">AI đang học thói quen chi tiêu của bạn...</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Anomaly Detection -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-alert-circle</v-icon>
            Chi tiêu bất thường
          </v-card-title>
          <v-card-text>
            <div v-if="anomalies && anomalies.anomalies && anomalies.anomalies.length > 0">
              <v-alert v-for="(anomaly, index) in anomalies.anomalies" :key="index" type="warning" variant="tonal"
                class="mb-2">
                <div class="text-subtitle-2">{{ anomaly.description }}</div>
                <div class="text-caption">
                  Số tiền: {{ formatCurrency(anomaly.amount) }} |
                  Ngày: {{ formatDate(anomaly.date) }}
                </div>
              </v-alert>
            </div>
            <div v-else class="text-center pa-4">
              <v-icon size="32" color="success">mdi-check-circle</v-icon>
              <div class="text-body-2 mt-2">Không phát hiện bất thường</div>
              <div class="text-caption">Chi tiêu của bạn đang khá ổn định</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Predictions -->
    <v-row v-if="predictions && predictions.user_id > 0">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-crystal-ball</v-icon>
            Dự đoán từ AI
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-card-text>
                    <div class="text-h6 text-primary">Dự đoán tháng tới</div>
                    <div class="text-h4">{{ formatCurrency(Number(predictions.predicted_amount ?? 0)) }}</div>
                    <div class="text-caption">Chi tiêu dự kiến</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-card-text>
                    <div class="text-h6 text-info">Mức độ tin cậy</div>
                    <div class="text-h4">{{ (Number(predictions.confidence_score ?? 0) * 100).toFixed(1) }}%</div>
                    <div class="text-caption">Độ chính xác ước tính</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-card-text>
                    <div class="text-h6 text-warning">Xu hướng</div>
                    <div class="text-h4">{{ predictions.trends?.[0]?.trend || 'ổn định' }}</div>
                    <div class="text-caption">Xu hướng chi tiêu</div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- AI Service Status -->
    <v-row v-else-if="predictions && predictions.user_id === 0">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-robot</v-icon>
            Trạng thái dịch vụ AI
          </v-card-title>
          <v-card-text>
            <v-alert type="info" variant="tonal">
              <div class="text-subtitle-2">Dịch vụ AI đang khởi động</div>
              <div class="text-caption">
                Dịch vụ AI đang được khởi chạy. Các tính năng dự đoán và phân tích nâng cao sẽ sẵn sàng sau ít phút,
                đặc biệt là lần khởi động đầu tiên.
              </div>
            </v-alert>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import DoughnutChart from '@/components/DoughnutChart.vue'
import { analyticsAPI, authAPI, transactionAPI } from '@/services/api'
import { formatCurrency, formatDate } from '@/utils/formatters'
import { computed, onMounted, ref } from 'vue'

// Reactive data
const loading = ref(false)
const dashboardData = ref(null)
const categorySpending = ref([])
const spendingPatterns = ref(null)
const anomalies = ref(null)
const predictions = ref(null)
const monthlyIncome = ref(0)
const incomeSaving = ref(false)
const mom = ref(null)

// Date range
const selectedPeriod = ref('last_month')
const startDate = ref('')
const endDate = ref('')

const periodOptions = [
  { title: 'Tháng trước', value: 'last_month' },
  { title: '3 tháng gần đây', value: 'last_3_months' },
  { title: '6 tháng gần đây', value: 'last_6_months' },
  { title: 'Năm qua', value: 'last_year' },
  { title: 'Khoảng tuỳ chọn', value: 'custom' }
]

// Tx filters
const txTypeOptions = [
  { title: 'Thu nhập', value: 'income' },
  { title: 'Chi tiêu', value: 'expense' },
  { title: 'Chuyển khoản', value: 'transfer' },
]
const txType = ref(null)
const minAmount = ref(null)
const maxAmount = ref(null)
const search = ref('')
const txLoading = ref(false)
const filteredItems = ref([])
const filteredTotal = ref(0)

// Chart data (fallback to AI patterns if spending breakdown empty)
const categorySpendingChartData = computed(() => {
  if (categorySpending.value.length > 0) {
    return {
      labels: categorySpending.value.map((item) => item.category_name),
      datasets: [{
        data: categorySpending.value.map((item) => Number(item.amount || 0)),
        backgroundColor: [
          '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF',
          '#FF9F40', '#FF6384', '#C9CBCF', '#4BC0C0', '#FF6384'
        ],
        borderWidth: 1,
        borderColor: '#ffffff'
      }]
    }
  }

  const patterns = (spendingPatterns.value && spendingPatterns.value.patterns) || []
  if (Array.isArray(patterns) && patterns.length > 0) {
    // Use only expense-like categories (exclude obvious income names)
    const expenseLike = patterns.filter((p) => {
      const name = String(p.category_name || '').toLowerCase()
      return !name.includes('thu nhập') && !name.includes('income')
    })
    if (expenseLike.length === 0) return null
    return {
      labels: expenseLike.map((p) => p.category_name),
      datasets: [{
        data: expenseLike.map((p) => Number(p.total_amount || 0)),
        backgroundColor: [
          '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF',
          '#FF9F40', '#FF6384', '#C9CBCF', '#4BC0C0', '#FF6384'
        ],
        borderWidth: 1,
        borderColor: '#ffffff'
      }]
    }
  }

  return null
})

const doughnutChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom'
    }
  }
}

// Methods
const loadAnalytics = async () => {
  loading.value = true
  try {
    const params = getDateRangeParams()

    // Load all analytics data in parallel
    const [dashboard, spending, patterns, anomalyData, predictionData] = await Promise.all([
      analyticsAPI.getDashboard(params),
      analyticsAPI.getCategorySpending(params),
      analyticsAPI.getSpendingPatterns(params),
      analyticsAPI.getAnomalies(params),
      analyticsAPI.getPredictions(params)
    ])

    dashboardData.value = dashboard.data
    categorySpending.value = spending.data
    // Normalize spending patterns to support both string and object formats
    const patternsData = patterns.data || {}
    const normalizedInsights = Array.isArray(patternsData.insights)
      ? patternsData.insights.map((it) => {
        if (typeof it === 'string') {
          return { type: 'info', description: it }
        }
        return {
          type: it.type || 'info',
          description: it.description || String(it)
        }
      })
      : []
    const translateToVi = (text) => {
      if (!text) return text
      const map = {
        'Set up budget alerts for your top spending categories.': 'Thiết lập cảnh báo ngân sách cho các danh mục chi tiêu hàng đầu của bạn.',
        'Consider setting monthly spending limits for discretionary categories.': 'Cân nhắc đặt hạn mức chi tiêu hàng tháng cho các danh mục tùy ý.',
        'Your spending patterns show consistent behavior across categories.': 'Mẫu chi tiêu của bạn cho thấy hành vi ổn định giữa các danh mục.',
        'Consider reviewing your top spending categories for optimization opportunities.': 'Hãy xem lại các danh mục chi tiêu hàng đầu để tối ưu hoá.'
      }
      return map[text] || text
    }
    const normalizedRecs = Array.isArray(patternsData.recommendations)
      ? patternsData.recommendations.map((it) => {
        if (typeof it === 'string') {
          return { title: 'Gợi ý', description: translateToVi(it), priority: 'medium' }
        }
        return {
          title: it.title ? translateToVi(it.title) : 'Gợi ý',
          description: translateToVi(it.description || String(it)),
          priority: it.priority || 'medium'
        }
      })
      : []
    spendingPatterns.value = {
      ...patternsData,
      insights: normalizedInsights,
      recommendations: normalizedRecs,
    }
    anomalies.value = anomalyData.data
    predictions.value = predictionData.data

    // Debug logging
    console.log('Dashboard Data:', dashboardData.value)
    console.log('Category Spending:', categorySpending.value)
  } catch (error) {
    console.error('Failed to load analytics:', error)
  } finally {
    loading.value = false
  }
}

const getDateRangeParams = () => {
  const params = {}

  if (selectedPeriod.value === 'custom') {
    if (startDate.value) params.start_date = startDate.value
    if (endDate.value) params.end_date = endDate.value
  } else {
    const now = new Date()
    const end = new Date(now)
    let start = new Date(now)

    switch (selectedPeriod.value) {
      case 'last_month':
        start.setMonth(start.getMonth() - 1)
        break
      case 'last_3_months':
        start.setMonth(start.getMonth() - 3)
        break
      case 'last_6_months':
        start.setMonth(start.getMonth() - 6)
        break
      case 'last_year':
        start.setFullYear(start.getFullYear() - 1)
        break
    }

    params.start_date = start.toISOString().split('T')[0]
    params.end_date = end.toISOString().split('T')[0]
  }

  return params
}

const getRecommendationIcon = (priority) => {
  switch (priority) {
    case 'high': return 'mdi-alert-circle'
    case 'medium': return 'mdi-information'
    case 'low': return 'mdi-lightbulb'
    default: return 'mdi-information'
  }
}

const loadIncome = async () => {
  try {
    const res = await authAPI.getIncome()
    monthlyIncome.value = Number(res.data?.monthly_income ?? 0)
  } catch (e) {
    monthlyIncome.value = 0
  }
}

const saveIncome = async () => {
  try {
    incomeSaving.value = true
    await authAPI.setIncome(Number(monthlyIncome.value || 0))
    await loadAnalytics()
  } finally {
    incomeSaving.value = false
  }
}

const computeMoM = async () => {
  try {
    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth() + 1
    const prevDate = new Date(year, month - 2, 1)
    const [cur, prev] = await Promise.all([
      analyticsAPI.getDashboard({ year, month }),
      analyticsAPI.getDashboard({ year: prevDate.getFullYear(), month: prevDate.getMonth() + 1 }),
    ])
    const c = cur.data
    const p = prev.data
    const incomeChange = p.total_income ? ((c.total_income - p.total_income) / p.total_income) * 100 : 0
    const expenseChange = p.total_expense ? ((c.total_expense - p.total_expense) / p.total_expense) * 100 : 0
    mom.value = {
      current: { income: c.total_income, expense: c.total_expense },
      prev: { income: p.total_income, expense: p.total_expense },
      incomeChange,
      expenseChange,
    }
  } catch (e) {
    mom.value = null
  }
}

const loadFilteredTransactions = async () => {
  txLoading.value = true
  try {
    const params = {
      page: 1,
      limit: 50,
      start_date: startDate.value,
      end_date: endDate.value,
      search: search.value || undefined,
      min_amount: minAmount.value || undefined,
      max_amount: maxAmount.value || undefined,
      transaction_type: txType.value || undefined,
      sort_by: 'transaction_date',
      sort_order: 'desc',
    }
    const res = await transactionAPI.getTransactions(params)
    filteredItems.value = Array.isArray(res.data?.data) ? res.data.data : res.data?.items || []
    filteredTotal.value = Number(res.data?.total || (res.data?.pagination?.total ?? filteredItems.value.length))
  } catch (e) {
    filteredItems.value = []
    filteredTotal.value = 0
  } finally {
    txLoading.value = false
  }
}

const percentOfIncome = computed(() => {
  const exp = Number(dashboardData.value?.total_expense || 0)
  const inc = Number(monthlyIncome.value || 0)
  if (!inc || inc <= 0) return 0
  return (exp / inc) * 100
})

// Initialize date range
const initializeDateRange = () => {
  const now = new Date()
  const lastMonth = new Date(now)
  lastMonth.setMonth(lastMonth.getMonth() - 1)

  startDate.value = lastMonth.toISOString().split('T')[0]
  endDate.value = now.toISOString().split('T')[0]
}

// Lifecycle
onMounted(() => {
  initializeDateRange()
  loadAnalytics()
  loadIncome()
  computeMoM()
})
</script>
