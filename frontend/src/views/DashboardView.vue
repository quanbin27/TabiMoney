<template>
  <v-container fluid>
    <!-- Welcome Section -->
    <v-row class="mb-6">
      <v-col cols="12">
        <v-card class="pa-6" color="primary" dark>
          <v-row align="center">
            <v-col cols="12" md="8">
              <h1 class="text-h4 font-weight-bold mb-2">
                Welcome back, {{ authStore.user?.first_name || authStore.user?.username }}!
              </h1>
              <p class="text-h6 opacity-90">
                Here's your financial overview for {{ currentMonth }}
              </p>
            </v-col>
            <v-col cols="12" md="4" class="text-right">
              <v-btn
                color="white"
                variant="outlined"
                size="large"
                @click="$router.push('/transactions/add')"
              >
                <v-icon left>mdi-plus</v-icon>
                Add Transaction
              </v-btn>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Financial Summary Cards -->
    <v-row class="mb-6">
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4" color="success" dark>
          <v-card-title class="text-h6">Total Income</v-card-title>
          <v-card-text>
            <div class="text-h4 font-weight-bold">
              {{ formatCurrency(analytics?.total_income || 0) }}
            </div>
            <div class="text-caption opacity-90">
              This month
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4" color="error" dark>
          <v-card-title class="text-h6">Total Expenses</v-card-title>
          <v-card-text>
            <div class="text-h4 font-weight-bold">
              {{ formatCurrency(analytics?.total_expense || 0) }}
            </div>
            <div class="text-caption opacity-90">
              This month
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4" :color="netAmountColor" dark>
          <v-card-title class="text-h6">Net Amount</v-card-title>
          <v-card-text>
            <div class="text-h4 font-weight-bold">
              {{ formatCurrency(analytics?.net_amount || 0) }}
            </div>
            <div class="text-caption opacity-90">
              {{ netAmountLabel }}
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" sm="6" md="3">
        <v-card class="pa-4" color="info" dark>
          <v-card-title class="text-h6">Transactions</v-card-title>
          <v-card-text>
            <div class="text-h4 font-weight-bold">
              {{ analytics?.transaction_count || 0 }}
            </div>
            <div class="text-caption opacity-90">
              This month
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Charts and Analytics -->
    <v-row>
      <!-- Spending by Category Chart -->
      <v-col cols="12" md="8">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-chart-pie</v-icon>
            Spending by Category
          </v-card-title>
          <v-card-text>
            <div v-if="loading" class="text-center pa-8">
              <v-progress-circular indeterminate color="primary" />
            </div>
            <div v-else-if="categoryChartData">
              <DoughnutChart
                :data="categoryChartData"
                :options="chartOptions"
              />
            </div>
            <div v-else class="text-center pa-8 text-medium-emphasis">
              No data available
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <!-- Financial Health -->
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-heart-pulse</v-icon>
            Financial Health
          </v-card-title>
          <v-card-text>
            <div v-if="analytics?.financial_health" class="text-center">
              <v-progress-circular
                :model-value="analytics.financial_health.score"
                :color="healthColor"
                size="120"
                width="12"
                class="mb-4"
              >
                <span class="text-h4 font-weight-bold">
                  {{ Math.round(analytics.financial_health.score) }}
                </span>
              </v-progress-circular>
              
              <div class="text-h6 font-weight-bold mb-2">
                {{ analytics.financial_health.level.toUpperCase() }}
              </div>
              
              <div class="text-body-2 text-medium-emphasis mb-4">
                Savings Rate: {{ analytics.financial_health.savings_rate.toFixed(1) }}%
              </div>
              
              <v-list density="compact">
                <v-list-item
                  v-for="recommendation in analytics.financial_health.recommendations"
                  :key="recommendation"
                  class="px-0"
                >
                  <template v-slot:prepend>
                    <v-icon size="small" color="primary">mdi-lightbulb</v-icon>
                  </template>
                  <v-list-item-title class="text-caption">
                    {{ recommendation }}
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </div>
            <div v-else class="text-center pa-4 text-medium-emphasis">
              No health data available
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Recent Transactions -->
    <v-row class="mt-6">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-receipt</v-icon>
            Recent Transactions
          </v-card-title>
          <v-card-text>
            <v-data-table
              :headers="transactionHeaders"
              :items="recentTransactions"
              :loading="loading"
              class="elevation-0"
            >
              <template v-slot:item.amount="{ item }">
                <span :class="getAmountClass(item.transaction_type)">
                  {{ formatCurrency(item.amount) }}
                </span>
              </template>
              
              <template v-slot:item.transaction_date="{ item }">
                {{ formatDate(item.transaction_date) }}
              </template>
              
              <template v-slot:item.category="{ item }">
                <v-chip
                  :color="item.category?.color || 'primary'"
                  size="small"
                  variant="tonal"
                >
                  <v-icon start size="small">{{ item.category?.icon || 'mdi-tag' }}</v-icon>
                  {{ item.category?.name || 'Unknown' }}
                </v-chip>
              </template>
              
              <template v-slot:item.actions="{ item }">
                <v-btn
                  icon="mdi-pencil"
                  size="small"
                  variant="text"
                  @click="editTransaction(item.id)"
                />
                <v-btn
                  icon="mdi-delete"
                  size="small"
                  variant="text"
                  color="error"
                  @click="deleteTransaction(item.id)"
                />
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import { analyticsAPI, transactionAPI } from '../services/api'
import { formatCurrency, formatDate } from '../utils/formatters'
import DoughnutChart from '../components/charts/DoughnutChart.vue'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// State
const analytics = ref(null)
const recentTransactions = ref([])
const loading = ref(false)

// Computed
const currentMonth = computed(() => {
  return new Date().toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
})

const netAmountColor = computed(() => {
  if (!analytics.value) return 'info'
  return analytics.value.net_amount >= 0 ? 'success' : 'error'
})

const netAmountLabel = computed(() => {
  if (!analytics.value) return ''
  return analytics.value.net_amount >= 0 ? 'Saved' : 'Overspent'
})

const healthColor = computed(() => {
  if (!analytics.value?.financial_health) return 'info'
  const score = analytics.value.financial_health.score
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'error'
})

const categoryChartData = computed(() => {
  if (!analytics.value?.category_breakdown) return null
  
  return {
    labels: analytics.value.category_breakdown.map(c => c.category_name),
    datasets: [{
      data: analytics.value.category_breakdown.map(c => c.amount),
      backgroundColor: [
        '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF',
        '#FF9F40', '#FF6384', '#C9CBCF', '#4BC0C0', '#FF6384'
      ],
      borderWidth: 2,
      borderColor: '#fff'
    }]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom' },
  },
}

const transactionHeaders = [
  { title: 'Date', key: 'transaction_date', sortable: true },
  { title: 'Description', key: 'description', sortable: true },
  { title: 'Category', key: 'category', sortable: false },
  { title: 'Amount', key: 'amount', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false },
]

// Methods
const loadDashboardData = async () => {
  loading.value = true
  try {
    const [analyticsResponse, transactionsResponse] = await Promise.all([
      analyticsAPI.getDashboard(),
      transactionAPI.getTransactions({ limit: 10, sort_by: 'transaction_date', sort_order: 'desc' })
    ])
    
    analytics.value = analyticsResponse.data
    recentTransactions.value = transactionsResponse.data.data
  } catch (error) {
    appStore.showError('Failed to load dashboard data')
  } finally {
    loading.value = false
  }
}

const getAmountClass = (type) => {
  return type === 'income' ? 'text-success' : 'text-error'
}

const editTransaction = (id) => {
  router.push(`/transactions/${id}/edit`)
}

const deleteTransaction = async (id) => {
  if (confirm('Are you sure you want to delete this transaction?')) {
    try {
      await transactionAPI.deleteTransaction(id)
      appStore.showSuccess('Transaction deleted successfully')
      loadDashboardData()
    } catch (error) {
      appStore.showError('Failed to delete transaction')
    }
  }
}

// Lifecycle
onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.v-card {
  border-radius: 12px;
}

.text-success {
  color: rgb(var(--v-theme-success)) !important;
}

.text-error {
  color: rgb(var(--v-theme-error)) !important;
}
</style>
