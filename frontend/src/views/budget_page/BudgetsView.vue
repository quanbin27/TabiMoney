<template>
  <v-container fluid class="pa-4">
    <!-- Header -->
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card color="primary" dark>
          <v-card-title class="d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-icon class="mr-3" size="large">mdi-wallet</v-icon>
              <div>
                <h2 class="text-h5">Budget Management</h2>
                <p class="text-subtitle-1 mb-0 opacity-90">
                  Track and manage your monthly budgets
                </p>
              </div>
            </div>
            <v-btn
              color="white"
              variant="outlined"
              @click="openCreateDialog"
              prepend-icon="mdi-plus"
            >
              New Budget
            </v-btn>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <!-- Budget Overview Cards -->
    <v-row class="mb-4" v-if="budgets.length > 0">
      <v-col cols="12" md="3">
        <v-card color="success" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-success">Total Budget</div>
            <div class="text-h4">{{ formatCurrency(totalBudget) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="error" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-error">Total Spent</div>
            <div class="text-h4">{{ formatCurrency(totalSpent) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="primary" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-primary">Remaining</div>
            <div class="text-h4">{{ formatCurrency(totalRemaining) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card :color="budgetHealthColor" variant="tonal">
          <v-card-text>
            <div class="text-h6">Budget Health</div>
            <div class="text-h4">{{ budgetHealthPercentage }}%</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Budgets List -->
    <v-row v-if="budgets.length > 0">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-chart-pie</v-icon>
            Monthly Budgets
          </v-card-title>
          <v-card-text>
            <v-data-table
              :headers="headers"
              :items="budgets"
              :loading="loading"
              class="elevation-1"
            >
              <template v-slot:item.category_name="{ item }">
                <div class="d-flex align-center">
                  <v-icon class="mr-2" :color="getCategoryColor(item.category_id)">
                    {{ getCategoryIcon(item.category_id) }}
                  </v-icon>
                  {{ getCategoryName(item.category_id) }}
                </div>
              </template>

              <template v-slot:item.amount="{ item }">
                <span class="font-weight-bold">{{ formatCurrency(item.amount) }}</span>
              </template>

              <template v-slot:item.spent_amount="{ item }">
                <span class="font-weight-bold">{{ formatCurrency(item.spent_amount) }}</span>
              </template>

              <template v-slot:item.remaining_amount="{ item }">
                <span 
                  class="font-weight-bold"
                  :class="getRemainingAmountClass(item.remaining_amount)"
                >
                  {{ formatCurrency(item.remaining_amount) }}
                </span>
              </template>

              <template v-slot:item.progress="{ item }">
                <div class="d-flex align-center">
                  <v-progress-linear
                    :model-value="getProgressPercentage(item)"
                    :color="getProgressColor(item)"
                    height="20"
                    rounded
                    class="mr-2"
                  ></v-progress-linear>
                  <span class="text-caption">{{ getProgressPercentage(item) }}%</span>
                </div>
              </template>

              <template v-slot:item.status="{ item }">
                <v-chip
                  :color="getStatusColor(item)"
                  size="small"
                >
                  {{ getStatusText(item) }}
                </v-chip>
              </template>

              <template v-slot:item.actions="{ item }">
                <v-btn
                  icon="mdi-pencil"
                  size="small"
                  variant="text"
                  @click="openEditDialog(item)"
                ></v-btn>
                <v-btn
                  icon="mdi-delete"
                  size="small"
                  variant="text"
                  color="error"
                  @click="deleteBudget(item.id)"
                ></v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Empty State -->
    <v-row v-else>
      <v-col cols="12">
        <v-card class="text-center pa-8">
          <v-icon size="64" color="grey-lighten-1">mdi-wallet</v-icon>
          <h3 class="text-h6 mt-4 text-grey">No Budgets Set</h3>
          <p class="text-body-2 text-grey mb-4">
            Create your first budget to start tracking your spending
          </p>
          <v-btn
            color="primary"
            @click="openCreateDialog"
            prepend-icon="mdi-plus"
          >
            Create Your First Budget
          </v-btn>
        </v-card>
      </v-col>
    </v-row>

    <!-- Create/Edit Budget Dialog -->
    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title>
          {{ isEditing ? 'Edit Budget' : 'Create New Budget' }}
        </v-card-title>
        <v-card-text>
          <v-form ref="formRef" validate-on="submit">
            <v-select
              v-model="form.category_id"
              label="Category"
              :items="categories"
              item-title="name"
              item-value="id"
              :rules="[rules.required]"
              required
            ></v-select>
            
            <v-text-field
              v-model="form.name"
              label="Budget Name"
              :rules="[rules.required]"
              required
              clearable
            ></v-text-field>
            
            <v-text-field
              v-model="form.amount"
              label="Budget Amount"
              type="number"
              :rules="[rules.required, rules.positive]"
              prefix="₫"
              required
              clearable
            ></v-text-field>
            
            <v-select
              v-model="form.period"
              label="Period"
              :items="budgetTypes"
              item-title="title"
              item-value="value"
              :rules="[rules.required]"
              required
            ></v-select>
            
            <v-text-field
              v-model="form.start_date"
              label="Start Date"
              type="date"
              :rules="[rules.required]"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="form.end_date"
              label="End Date"
              type="date"
              :rules="[rules.required]"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="form.alert_threshold"
              label="Alert Threshold (%)"
              type="number"
              :rules="[rules.required, rules.positive]"
              suffix="%"
              min="0"
              max="100"
            ></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeDialog">Cancel</v-btn>
          <v-btn
            type="submit"
            color="primary"
            @click="saveBudget"
            :loading="saving"
          >
            {{ isEditing ? 'Update' : 'Create' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Budget Analysis Chart -->
    <v-row v-if="budgets.length > 0" class="mt-4">
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-chart-donut</v-icon>
            Budget vs Spent
          </v-card-title>
          <v-card-text>
            <canvas ref="budgetChart" height="300"></canvas>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-chart-bar</v-icon>
            Category Spending
          </v-card-title>
          <v-card-text>
            <canvas ref="categoryChart" height="300"></canvas>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { budgetAPI, categoryAPI } from '@/services/api'
import { formatCurrency } from '@/utils/formatters'
import { useAppStore } from '@/stores/app'
import Chart from 'chart.js/auto'

// Reactive data
const budgets = ref([])
const categories = ref([])
const dialog = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const loading = ref(false)
const formValid = ref(false)
const selectedBudget = ref(null)
const budgetChart = ref(null)
const categoryChart = ref(null)

// Form data
const formRef = ref()
const form = ref({
  category_id: null,
  name: '',
  amount: null,
  period: 'monthly',
  start_date: '',
  end_date: '',
  alert_threshold: 80
})

// Table headers
const headers = [
  { title: 'Name', key: 'name', sortable: true },
  { title: 'Category', key: 'category_name', sortable: true },
  { title: 'Budget Amount', key: 'amount', sortable: true },
  { title: 'Spent Amount', key: 'spent_amount', sortable: true },
  { title: 'Remaining', key: 'remaining_amount', sortable: true },
  { title: 'Progress', key: 'progress', sortable: false },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false }
]

// Options
const budgetTypes = [
  { title: 'Weekly', value: 'weekly' },
  { title: 'Monthly', value: 'monthly' },
  { title: 'Yearly', value: 'yearly' }
]

// Validation rules
const rules = {
  required: (value) => {
    if (value === null || value === undefined) return 'This field is required'
    if (typeof value === 'number') return true // allow 0 as valid
    if (typeof value === 'string') return value.trim().length > 0 || 'This field is required'
    return !!value || 'This field is required'
  },
  positive: (value) => {
    if (value === null || value === undefined) return true
    return Number(value) > 0 || 'Value must be positive'
  }
}

// Computed properties
const totalBudget = computed(() => {
  return budgets.value.reduce((sum, budget) => sum + budget.amount, 0)
})

const totalSpent = computed(() => {
  return budgets.value.reduce((sum, budget) => sum + budget.spent_amount, 0)
})

const totalRemaining = computed(() => {
  return totalBudget.value - totalSpent.value
})

const budgetHealthPercentage = computed(() => {
  if (totalBudget.value === 0) return 0
  return Math.round((totalRemaining.value / totalBudget.value) * 100)
})

const budgetHealthColor = computed(() => {
  const percentage = budgetHealthPercentage.value
  if (percentage >= 80) return 'success'
  if (percentage >= 60) return 'primary'
  if (percentage >= 40) return 'warning'
  return 'error'
})

// Methods
const loadBudgets = async () => {
  loading.value = true
  try {
    const response = await budgetAPI.getBudgets()
    budgets.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load budgets:', error)
    useAppStore().showError('Failed to load budgets')
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const response = await categoryAPI.getCategories()
    const raw = response.data
    // Handle both wrapped ({ data: [...] }) and unwrapped ([]) responses
    categories.value = Array.isArray(raw) ? raw : (raw?.data ?? [])
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const openCreateDialog = () => {
  isEditing.value = false
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  
  form.value = {
    category_id: null,
    name: '',
    amount: null,
    period: 'monthly',
    start_date: startOfMonth.toISOString().split('T')[0],
    end_date: endOfMonth.toISOString().split('T')[0],
    alert_threshold: 80
  }
  dialog.value = true
}

const openEditDialog = (budget) => {
  isEditing.value = true
  form.value = {
    category_id: budget.category_id,
    name: budget.name,
    amount: budget.amount,
    period: budget.period,
    start_date: budget.start_date ? new Date(budget.start_date).toISOString().split('T')[0] : '',
    end_date: budget.end_date ? new Date(budget.end_date).toISOString().split('T')[0] : '',
    alert_threshold: budget.alert_threshold
  }
  selectedBudget.value = budget
  dialog.value = true
}

const closeDialog = () => {
  dialog.value = false
  selectedBudget.value = null
}

const saveBudget = async () => {
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) return
  }
  saving.value = true
  try {
    // Normalize payload for backend (RFC3339 dates, numeric fields)
    const payload = {
      category_id: form.value.category_id,
      name: (form.value.name || '').trim(),
      amount: form.value.amount != null ? Number(form.value.amount) : null,
      period: form.value.period,
      start_date: form.value.start_date ? `${form.value.start_date}T00:00:00Z` : null,
      end_date: form.value.end_date ? `${form.value.end_date}T23:59:59Z` : null,
      alert_threshold: form.value.alert_threshold != null ? Number(form.value.alert_threshold) : 80,
    }

    if (isEditing.value) {
      await budgetAPI.updateBudget(selectedBudget.value.id, payload)
    } else {
      await budgetAPI.createBudget(payload)
    }
    await loadBudgets()
    closeDialog()
    useAppStore().showSuccess(
      isEditing.value ? 'Budget updated successfully' : 'Budget created successfully'
    )
  } catch (error) {
    console.error('Failed to save budget:', error)
    useAppStore().showError('Failed to save budget')
  } finally {
    saving.value = false
  }
}

const deleteBudget = async (budgetId) => {
  if (confirm('Are you sure you want to delete this budget?')) {
    try {
      await budgetAPI.deleteBudget(budgetId)
      await loadBudgets()
      useAppStore().showSuccess('Budget deleted successfully')
    } catch (error) {
      console.error('Failed to delete budget:', error)
      useAppStore().showError('Failed to delete budget')
    }
  }
}

// Helper methods
const getProgressPercentage = (budget) => {
  if (budget.amount === 0) return 0
  return Math.min((budget.spent_amount / budget.amount) * 100, 100)
}

const getProgressColor = (budget) => {
  const percentage = getProgressPercentage(budget)
  if (percentage >= 100) return 'error'
  if (percentage >= 80) return 'warning'
  return 'success'
}

const getStatusText = (budget) => {
  const percentage = getProgressPercentage(budget)
  if (percentage >= 100) return 'Over Budget'
  if (percentage >= 80) return 'Near Limit'
  return 'On Track'
}

const getStatusColor = (budget) => {
  const percentage = getProgressPercentage(budget)
  if (percentage >= 100) return 'error'
  if (percentage >= 80) return 'warning'
  return 'success'
}

const getRemainingAmountClass = (amount) => {
  return amount < 0 ? 'text-error' : 'text-success'
}

const getCategoryColor = (categoryId) => {
  const colors = ['primary', 'secondary', 'success', 'warning', 'error', 'info']
  return colors[categoryId % colors.length]
}

const getCategoryIcon = (categoryId) => {
  const icons = ['mdi-home', 'mdi-car', 'mdi-food', 'mdi-shopping', 'mdi-gamepad', 'mdi-medical-bag']
  return icons[categoryId % icons.length]
}

const getCategoryName = (categoryId) => {
  const cat = categories.value.find((c) => c.id === categoryId)
  return cat?.name ?? '—'
}

const createCharts = async () => {
  await nextTick()
  
  if (budgetChart.value && budgets.value.length > 0) {
    new Chart(budgetChart.value, {
      type: 'doughnut',
      data: {
        labels: ['Budget', 'Spent'],
        datasets: [{
          data: [totalBudget.value, totalSpent.value],
          backgroundColor: ['#4CAF50', '#F44336'],
          borderWidth: 2,
          borderColor: '#fff'
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'bottom'
          }
        }
      }
    })
  }

  if (categoryChart.value && budgets.value.length > 0) {
    new Chart(categoryChart.value, {
      type: 'bar',
      data: {
        labels: budgets.value.map(b => b.name),
        datasets: [{
          label: 'Budget',
          data: budgets.value.map(b => b.amount),
          backgroundColor: '#2196F3'
        }, {
          label: 'Spent',
          data: budgets.value.map(b => b.spent_amount),
          backgroundColor: '#F44336'
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: true
          }
        }
      }
    })
  }
}

// Lifecycle
onMounted(async () => {
  await loadCategories()
  await loadBudgets()
  await createCharts()
})
</script>

<style scoped>
.v-progress-linear {
  border-radius: 4px;
}
</style>
