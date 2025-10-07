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
                  {{ item.category_name }}
                </div>
              </template>

              <template v-slot:item.budget_amount="{ item }">
                <span class="font-weight-bold">{{ formatCurrency(item.budget_amount) }}</span>
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
          <v-form ref="form" v-model="formValid">
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
              v-model="form.budget_amount"
              label="Budget Amount"
              type="number"
              :rules="[rules.required, rules.positive]"
              prefix="₫"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="form.month"
              label="Month"
              type="month"
              :rules="[rules.required]"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="form.description"
              label="Description (Optional)"
            ></v-text-field>
            
            <v-select
              v-model="form.budget_type"
              label="Budget Type"
              :items="budgetTypes"
              :rules="[rules.required]"
              required
            ></v-select>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeDialog">Cancel</v-btn>
          <v-btn
            color="primary"
            @click="saveBudget"
            :disabled="!formValid"
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

<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import { budgetAPI, categoryAPI } from '@/services/api'
import { formatCurrency } from '@/utils/formatters'
import { useAppStore } from '@/stores/app'
import Chart from 'chart.js/auto'

// Reactive data
const budgets = ref<any[]>([])
const categories = ref<any[]>([])
const dialog = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const loading = ref(false)
const formValid = ref(false)
const selectedBudget = ref<any>(null)
const budgetChart = ref<HTMLCanvasElement>()
const categoryChart = ref<HTMLCanvasElement>()

// Form data
const form = ref({
  category_id: null,
  budget_amount: 0,
  month: '',
  description: '',
  budget_type: ''
})

// Table headers
const headers = [
  { title: 'Category', key: 'category_name', sortable: true },
  { title: 'Budget Amount', key: 'budget_amount', sortable: true },
  { title: 'Spent Amount', key: 'spent_amount', sortable: true },
  { title: 'Remaining', key: 'remaining_amount', sortable: true },
  { title: 'Progress', key: 'progress', sortable: false },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false }
]

// Options
const budgetTypes = [
  { title: 'Monthly', value: 'monthly' },
  { title: 'Weekly', value: 'weekly' },
  { title: 'Yearly', value: 'yearly' },
  { title: 'Custom', value: 'custom' }
]

// Validation rules
const rules = {
  required: (value: any) => !!value || 'This field is required',
  positive: (value: number) => value > 0 || 'Value must be positive'
}

// Computed properties
const totalBudget = computed(() => {
  return budgets.value.reduce((sum, budget) => sum + budget.budget_amount, 0)
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
    categories.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const openCreateDialog = () => {
  isEditing.value = false
  form.value = {
    category_id: null,
    budget_amount: 0,
    month: new Date().toISOString().slice(0, 7), // Current month
    description: '',
    budget_type: 'monthly'
  }
  dialog.value = true
}

const openEditDialog = (budget: any) => {
  isEditing.value = true
  form.value = {
    category_id: budget.category_id,
    budget_amount: budget.budget_amount,
    month: budget.month,
    description: budget.description || '',
    budget_type: budget.budget_type
  }
  selectedBudget.value = budget
  dialog.value = true
}

const closeDialog = () => {
  dialog.value = false
  selectedBudget.value = null
}

const saveBudget = async () => {
  saving.value = true
  try {
    if (isEditing.value) {
      await budgetAPI.updateBudget(selectedBudget.value.id, form.value)
    } else {
      await budgetAPI.createBudget(form.value)
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

const deleteBudget = async (budgetId: number) => {
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
const getProgressPercentage = (budget: any) => {
  if (budget.budget_amount === 0) return 0
  return Math.min((budget.spent_amount / budget.budget_amount) * 100, 100)
}

const getProgressColor = (budget: any) => {
  const percentage = getProgressPercentage(budget)
  if (percentage >= 100) return 'error'
  if (percentage >= 80) return 'warning'
  return 'success'
}

const getStatusText = (budget: any) => {
  const percentage = getProgressPercentage(budget)
  if (percentage >= 100) return 'Over Budget'
  if (percentage >= 80) return 'Near Limit'
  return 'On Track'
}

const getStatusColor = (budget: any) => {
  const percentage = getProgressPercentage(budget)
  if (percentage >= 100) return 'error'
  if (percentage >= 80) return 'warning'
  return 'success'
}

const getRemainingAmountClass = (amount: number) => {
  return amount < 0 ? 'text-error' : 'text-success'
}

const getCategoryColor = (categoryId: number) => {
  const colors = ['primary', 'secondary', 'success', 'warning', 'error', 'info']
  return colors[categoryId % colors.length]
}

const getCategoryIcon = (categoryId: number) => {
  const icons = ['mdi-home', 'mdi-car', 'mdi-food', 'mdi-shopping', 'mdi-gamepad', 'mdi-medical-bag']
  return icons[categoryId % icons.length]
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
        labels: budgets.value.map(b => b.category_name),
        datasets: [{
          label: 'Budget',
          data: budgets.value.map(b => b.budget_amount),
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
