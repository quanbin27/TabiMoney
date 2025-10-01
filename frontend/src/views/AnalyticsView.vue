<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2">mdi-chart-line</v-icon>
            Financial Analytics Dashboard
          </v-card-title>
          <v-card-subtitle>
            AI-powered insights and spending analysis
          </v-card-subtitle>
        </v-card>
      </v-col>
    </v-row>

    <!-- Date Range Filter -->
    <v-row>
      <v-col cols="12" md="3">
        <v-select
          v-model="selectedPeriod"
          :items="periodOptions"
          label="Time Period"
          @update:model-value="loadAnalytics"
        />
      </v-col>
      <v-col cols="12" md="3">
        <v-text-field
          v-model="startDate"
          type="date"
          label="Start Date"
          @update:model-value="loadAnalytics"
        />
      </v-col>
      <v-col cols="12" md="3">
        <v-text-field
          v-model="endDate"
          type="date"
          label="End Date"
          @update:model-value="loadAnalytics"
        />
      </v-col>
      <v-col cols="12" md="3" class="d-flex align-center">
        <v-btn
          color="primary"
          @click="loadAnalytics"
          :loading="loading"
        >
          <v-icon class="mr-2">mdi-refresh</v-icon>
          Refresh
        </v-btn>
      </v-col>
    </v-row>

    <!-- Financial Summary Cards -->
    <v-row v-if="dashboardData">
      <v-col cols="12" md="3">
        <v-card color="success" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-success">Total Income</div>
            <div class="text-h4">{{ formatCurrency(dashboardData.total_income) }}</div>
            <div class="text-caption">{{ dashboardData.transaction_count }} transactions</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="error" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-error">Total Expenses</div>
            <div class="text-h4">{{ formatCurrency(dashboardData.total_expense) }}</div>
            <div class="text-caption">{{ dashboardData.transaction_count }} transactions</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card :color="dashboardData.net_amount >= 0 ? 'success' : 'error'" variant="tonal">
          <v-card-text>
            <div class="text-h6" :class="dashboardData.net_amount >= 0 ? 'text-success' : 'text-error'">
              Net Amount
            </div>
            <div class="text-h4">{{ formatCurrency(dashboardData.net_amount) }}</div>
            <div class="text-caption">Income - Expenses</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="info" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-info">Savings Rate</div>
            <div class="text-h4">{{ dashboardData.financial_health?.savings_rate?.toFixed(1) || '0.0' }}%</div>
            <div class="text-caption">Financial Health</div>
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
            Spending by Category
          </v-card-title>
          <v-card-text>
            <div v-if="categorySpending.length > 0" style="height: 300px;">
              <DoughnutChart
                :data="categorySpendingChartData"
                :options="doughnutChartOptions"
              />
            </div>
            <div v-else class="text-center pa-8">
              <v-icon size="48" color="grey">mdi-chart-donut</v-icon>
              <div class="text-h6 mt-2">No spending data</div>
              <div class="text-caption">Add some transactions to see analytics</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Spending Trends -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-chart-line</v-icon>
            Spending Trends
          </v-card-title>
          <v-card-text>
            <div v-if="spendingPatterns && spendingPatterns.insights && spendingPatterns.insights.length > 0" style="height: 300px;">
              <div class="text-h6 mb-4">AI Insights</div>
              <v-alert
                v-for="(insight, index) in spendingPatterns.insights"
                :key="index"
                :type="insight.type"
                variant="tonal"
                class="mb-2"
              >
                {{ insight.description }}
              </v-alert>
            </div>
            <div v-else class="text-center pa-8">
              <v-icon size="48" color="grey">mdi-chart-line</v-icon>
              <div class="text-h6 mt-2">No trend data</div>
              <div class="text-caption">AI is analyzing your spending patterns</div>
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
            AI Recommendations
          </v-card-title>
          <v-card-text>
            <div v-if="spendingPatterns && spendingPatterns.recommendations && spendingPatterns.recommendations.length > 0">
              <v-list>
                <v-list-item
                  v-for="(rec, index) in spendingPatterns.recommendations"
                  :key="index"
                  :prepend-icon="getRecommendationIcon(rec.priority)"
                >
                  <v-list-item-title>{{ rec.title }}</v-list-item-title>
                  <v-list-item-subtitle>{{ rec.description }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </div>
            <div v-else class="text-center pa-4">
              <v-icon size="32" color="grey">mdi-robot</v-icon>
              <div class="text-body-2 mt-2">AI is learning your patterns...</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Anomaly Detection -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-alert-circle</v-icon>
            Spending Anomalies
          </v-card-title>
          <v-card-text>
            <div v-if="anomalies && anomalies.anomalies && anomalies.anomalies.length > 0">
              <v-alert
                v-for="(anomaly, index) in anomalies.anomalies"
                :key="index"
                type="warning"
                variant="tonal"
                class="mb-2"
              >
                <div class="text-subtitle-2">{{ anomaly.description }}</div>
                <div class="text-caption">
                  Amount: {{ formatCurrency(anomaly.amount) }} | 
                  Date: {{ formatDate(anomaly.date) }}
                </div>
              </v-alert>
            </div>
            <div v-else class="text-center pa-4">
              <v-icon size="32" color="success">mdi-check-circle</v-icon>
              <div class="text-body-2 mt-2">No anomalies detected</div>
              <div class="text-caption">Your spending looks normal</div>
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
            AI Predictions
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-card-text>
                    <div class="text-h6 text-primary">Next Month Prediction</div>
                    <div class="text-h4">{{ formatCurrency(Number(predictions.predicted_amount ?? 0)) }}</div>
                    <div class="text-caption">Expected expenses</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-card-text>
                    <div class="text-h6 text-info">Confidence Score</div>
                    <div class="text-h4">{{ (Number(predictions.confidence_score ?? 0) * 100).toFixed(1) }}%</div>
                    <div class="text-caption">Prediction accuracy</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-card-text>
                    <div class="text-h6 text-warning">Trend Direction</div>
                    <div class="text-h4">{{ predictions.trends?.[0]?.trend || 'stable' }}</div>
                    <div class="text-caption">Spending trend</div>
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
            AI Service Status
          </v-card-title>
          <v-card-text>
            <v-alert type="info" variant="tonal">
              <div class="text-subtitle-2">AI Service Initializing</div>
              <div class="text-caption">
                The AI service is starting up. Predictions and advanced analytics will be available shortly.
                This may take a few minutes on first startup.
              </div>
            </v-alert>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { analyticsAPI } from '@/services/api'
import { formatCurrency, formatDate } from '@/utils/formatters'
import DoughnutChart from '@/components/charts/DoughnutChart.vue'

// Reactive data
const loading = ref(false)
const dashboardData = ref(null)
const categorySpending = ref([])
const spendingPatterns = ref(null)
const anomalies = ref(null)
const predictions = ref(null)

// Date range
const selectedPeriod = ref('last_month')
const startDate = ref('')
const endDate = ref('')

const periodOptions = [
  { title: 'Last Month', value: 'last_month' },
  { title: 'Last 3 Months', value: 'last_3_months' },
  { title: 'Last 6 Months', value: 'last_6_months' },
  { title: 'Last Year', value: 'last_year' },
  { title: 'Custom Range', value: 'custom' }
]

// Chart data
const categorySpendingChartData = computed(() => {
  if (!categorySpending.value.length) return null
  
  return {
    labels: categorySpending.value.map(item => item.category_name),
    datasets: [{
      data: categorySpending.value.map(item => item.amount),
      backgroundColor: [
        '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF',
        '#FF9F40', '#FF6384', '#C9CBCF', '#4BC0C0', '#FF6384'
      ],
      borderWidth: 1,
      borderColor: '#ffffff'
    }]
  }
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
      ? patternsData.insights.map((it: any) => {
          if (typeof it === 'string') {
            return { type: 'info', description: it }
          }
          return {
            type: it.type || 'info',
            description: it.description || String(it)
          }
        })
      : []
    const translateToVi = (text: string) => {
      if (!text) return text
      const map: Record<string, string> = {
        'Set up budget alerts for your top spending categories.': 'Thiết lập cảnh báo ngân sách cho các danh mục chi tiêu hàng đầu của bạn.',
        'Consider setting monthly spending limits for discretionary categories.': 'Cân nhắc đặt hạn mức chi tiêu hàng tháng cho các danh mục tùy ý.',
        'Your spending patterns show consistent behavior across categories.': 'Mẫu chi tiêu của bạn cho thấy hành vi ổn định giữa các danh mục.',
        'Consider reviewing your top spending categories for optimization opportunities.': 'Hãy xem lại các danh mục chi tiêu hàng đầu để tối ưu hoá.'
      }
      return map[text] || text
    }
    const normalizedRecs = Array.isArray(patternsData.recommendations)
      ? patternsData.recommendations.map((it: any) => {
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
  const params: any = {}
  
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

const getRecommendationIcon = (priority: string) => {
  switch (priority) {
    case 'high': return 'mdi-alert-circle'
    case 'medium': return 'mdi-information'
    case 'low': return 'mdi-lightbulb'
    default: return 'mdi-information'
  }
}

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
})
</script>
