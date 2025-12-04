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
                <h2 class="text-h5">Quản lý ngân sách</h2>
                <p class="text-subtitle-1 mb-0 opacity-90">
                  Theo dõi và quản lý ngân sách hàng tháng của bạn
                </p>
              </div>
            </div>
            <v-btn color="white" variant="outlined" @click="openCreateDialog" prepend-icon="mdi-plus">
              Ngân sách mới
            </v-btn>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <!-- Auto Budget Suggestion Banner -->
    <v-row class="mb-4">
      <v-col cols="12">
        <v-alert type="info" variant="tonal" border="start" prominent v-if="suggestBannerVisible">
          <div class="d-flex align-center justify-space-between w-100">
            <div>
              <div class="text-subtitle-1 font-weight-medium">Đề xuất ngân sách tháng này đã sẵn sàng</div>
              <div class="text-body-2 opacity-80">Tạo nhanh ngân sách theo gợi ý chỉ với 1 lần nhấn</div>
            </div>
            <div>
              <v-btn color="primary" @click="openSuggestDialog" prepend-icon="mdi-magic-staff">Xem gợi ý</v-btn>
            </div>
          </div>
        </v-alert>
      </v-col>
    </v-row>

    <!-- Budget Overview Cards -->
    <v-row class="mb-4" v-if="budgets.length > 0">
      <v-col cols="12" md="3">
        <v-card color="success" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-success">Tổng ngân sách</div>
            <div class="text-h4">{{ formatCurrency(totalBudget) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="error" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-error">Tổng đã chi</div>
            <div class="text-h4">{{ formatCurrency(totalSpent) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="primary" variant="tonal">
          <v-card-text>
            <div class="text-h6 text-primary">Còn lại</div>
            <div class="text-h4">{{ formatCurrency(totalRemaining) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card :color="budgetHealthColor" variant="tonal">
          <v-card-text>
            <div class="text-h6">Sức khỏe ngân sách</div>
            <div class="text-h4">{{ budgetHealthPercentage }}%</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="info" variant="tonal">
          <v-card-text>
            <div class="text-h6">Safe to Spend (Ngày)</div>
            <div class="text-h5">{{ formatCurrency(insights?.safe_to_spend_daily || 0) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card color="info" variant="tonal">
          <v-card-text>
            <div class="text-h6">Safe to Spend (Tuần)</div>
            <div class="text-h5">{{ formatCurrency(insights?.safe_to_spend_weekly || 0) }}</div>
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
            Ngân sách hàng tháng
          </v-card-title>
          <v-card-text>
            <v-data-table :headers="headers" :items="budgets" :loading="loading" class="elevation-1">
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
                <span class="font-weight-bold" :class="getRemainingAmountClass(item.remaining_amount)">
                  {{ formatCurrency(item.remaining_amount) }}
                </span>
              </template>

              <template v-slot:item.progress="{ item }">
                <div class="d-flex align-center">
                  <v-progress-linear :model-value="getProgressPercentage(item)" :color="getProgressColor(item)"
                    height="20" rounded class="mr-2"></v-progress-linear>
                  <span class="text-caption">{{ getProgressPercentage(item) }}%</span>
                </div>
              </template>

              <template v-slot:item.status="{ item }">
                <v-chip :color="getStatusColor(item)" size="small">
                  {{ getStatusText(item) }}
                </v-chip>
              </template>

              <template v-slot:item.actions="{ item }">
                <v-btn icon="mdi-pencil" size="small" variant="text" @click="openEditDialog(item)"></v-btn>
                <v-btn icon="mdi-delete" size="small" variant="text" color="error"
                  @click="deleteBudget(item.id)"></v-btn>
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
          <h3 class="text-h6 mt-4 text-grey">Chưa có ngân sách nào</h3>
          <p class="text-body-2 text-grey mb-4">
            Tạo ngân sách đầu tiên để bắt đầu theo dõi chi tiêu
          </p>
          <v-btn color="primary" @click="openCreateDialog" prepend-icon="mdi-plus">
            Tạo ngân sách đầu tiên
          </v-btn>
        </v-card>
      </v-col>
    </v-row>

    <!-- Create/Edit Budget Dialog -->
    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title>
          {{ isEditing ? 'Chỉnh sửa ngân sách' : 'Tạo ngân sách mới' }}
        </v-card-title>
        <v-card-text>
          <v-form ref="formRef" validate-on="submit">
            <v-select v-model="form.category_id" label="Danh mục" :items="categories" item-title="name" item-value="id"
              :rules="[rules.required]" required></v-select>

            <v-text-field v-model="form.name" label="Tên ngân sách" :rules="[rules.required]" required
              clearable></v-text-field>

            <v-text-field
              v-model="form.amount"
              label="Số tiền ngân sách"
              type="number"
              :rules="[rules.required, rules.positive]"
              prefix="₫"
              required
              clearable
            ></v-text-field>

            <v-text-field v-model="form.start_date" label="Từ ngày" type="date" :rules="[rules.required]"
              required></v-text-field>

            <v-text-field v-model="form.end_date" label="Đến ngày" type="date" :rules="[rules.required]"
              required></v-text-field>

            <v-text-field v-model="form.alert_threshold" label="Ngưỡng cảnh báo (%)" type="number"
              :rules="[rules.required, rules.positive]" suffix="%" min="0" max="100"></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeDialog">Hủy</v-btn>
          <v-btn type="submit" color="primary" @click="saveBudget" :loading="saving">
            {{ isEditing ? 'Cập nhật' : 'Tạo mới' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Auto Budget Suggestions Dialog -->
    <v-dialog v-model="suggestDialog" max-width="720">
      <v-card>
        <v-card-title>
          Gợi ý ngân sách tháng này
        </v-card-title>
        <v-card-text>
          <div class="mb-4 d-flex align-center justify-space-between">
            <div class="text-subtitle-2">Tổng đề xuất: {{ formatCurrency(totalSuggested) }}</div>
            <div class="text-caption">Kỳ: {{ suggestionPeriodLabel }}</div>
          </div>
          <v-table density="comfortable">
            <thead>
              <tr>
                <th>Tên</th>
                <th style="width: 180px">Danh mục</th>
                <th style="width: 180px">Số tiền</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in suggestItems" :key="idx">
                <td>{{ item.name }}</td>
                <td>
                  <v-select v-model="item.category_id" :items="categories" item-title="name" item-value="id" clearable
                    density="compact" />
                </td>
                <td>
                  <v-text-field v-model.number="item.suggested_amount" type="number" min="0" density="compact"
                    prefix="₫" />
                </td>
              </tr>
            </tbody>
          </v-table>
          <div class="text-caption mt-2" v-if="suggestNotes.length">
            <ul>
              <li v-for="(n, i) in suggestNotes" :key="i">{{ n }}</li>
            </ul>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="suggestDialog = false">Đóng</v-btn>
          <v-btn color="primary" :loading="creatingFromSuggest" @click="createFromSuggestions" prepend-icon="mdi-check">
            Tạo tất cả
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
            Ngân sách vs Đã chi
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
            Chi tiêu theo danh mục
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
import { budgetAPI, categoryAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { formatCurrency } from '@/utils/formatters'
import Chart from 'chart.js/auto'
import { computed, nextTick, onMounted, ref } from 'vue'

// Reactive data
const budgets = ref([])
const insights = ref(null)
const suggestDialog = ref(false)
const suggestItems = ref([])
const suggestNotes = ref([])
const suggestionPeriodLabel = ref('')
const creatingFromSuggest = ref(false)
const suggestBannerVisible = ref(true)
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
  { title: 'Tên', key: 'name', sortable: true },
  { title: 'Danh mục', key: 'category_name', sortable: true },
  { title: 'Ngân sách', key: 'amount', sortable: true },
  { title: 'Đã chi', key: 'spent_amount', sortable: true },
  { title: 'Còn lại', key: 'remaining_amount', sortable: true },
  { title: 'Tiến độ', key: 'progress', sortable: false },
  { title: 'Trạng thái', key: 'status', sortable: true },
  { title: 'Thao tác', key: 'actions', sortable: false }
]

// Validation rules
const rules = {
  required: (value) => {
    if (value === null || value === undefined) return 'Trường này là bắt buộc'
    if (typeof value === 'number') return true // allow 0 as valid
    if (typeof value === 'string') return value.trim().length > 0 || 'Trường này là bắt buộc'
    return !!value || 'Trường này là bắt buộc'
  },
  positive: (value) => {
    if (value === null || value === undefined) return true
    return Number(value) > 0 || 'Giá trị phải lớn hơn 0'
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
    useAppStore().showError('Không thể tải danh sách ngân sách')
  } finally {
    loading.value = false
  }
}

const loadInsights = async () => {
  try {
    const response = await budgetAPI.getInsights()
    insights.value = response.data?.data || null
  } catch (e) {
    console.error('Failed to load insights', e)
  }
}

const openSuggestDialog = async () => {
  try {
    const res = await budgetAPI.getAutoSuggestions()
    const data = res.data?.data
    suggestItems.value = (data?.suggestions || []).map(s => ({
      category_id: s.category_id ?? null,
      name: s.name,
      suggested_amount: Number(s.suggested_amount || 0),
    }))
    suggestNotes.value = data?.notes || []
    suggestionPeriodLabel.value = `${new Date(data.start_date).toLocaleDateString()} - ${new Date(data.end_date).toLocaleDateString()}`
    suggestDialog.value = true
  } catch (e) {
    console.error('Failed to get suggestions', e)
    useAppStore().showError('Không thể tải gợi ý ngân sách')
  }
}

const totalSuggested = computed(() => suggestItems.value.reduce((s, i) => s + (Number(i.suggested_amount) || 0), 0))

const createFromSuggestions = async () => {
  creatingFromSuggest.value = true
  try {
    // default to current month period
    const now = new Date()
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
    const payload = {
      period: 'monthly',
      start_date: startOfMonth.toISOString(),
      end_date: new Date(endOfMonth.getFullYear(), endOfMonth.getMonth(), endOfMonth.getDate(), 23, 59, 59).toISOString(),
      alert_threshold: 80,
      budgets: suggestItems.value.map(i => ({
        category_id: i.category_id ?? null,
        name: i.name,
        suggested_amount: Number(i.suggested_amount || 0),
      })),
    }
    await budgetAPI.createFromSuggestions(payload)
    suggestDialog.value = false
    suggestBannerVisible.value = false
    await loadBudgets()
    await loadInsights()
    useAppStore().showSuccess('Đã tạo ngân sách từ gợi ý')
  } catch (e) {
    console.error('Create from suggestions failed', e)
    useAppStore().showError('Không thể tạo ngân sách từ gợi ý')
  } finally {
    creatingFromSuggest.value = false
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
    // Tự động suy ra chu kỳ từ độ dài khoảng thời gian để tránh phải bắt user chọn thêm một field
    let period = 'monthly'
    if (form.value.start_date && form.value.end_date) {
      const start = new Date(form.value.start_date)
      const end = new Date(form.value.end_date)
      const diffDays = Math.max(1, Math.round((end - start) / (1000 * 60 * 60 * 24)) + 1)
      if (diffDays <= 10) {
        period = 'weekly'
      } else if (diffDays >= 330) {
        period = 'yearly'
      } else {
        period = 'monthly'
      }
    }

    // Normalize payload for backend (RFC3339 dates, numeric fields)
    const payload = {
      category_id: form.value.category_id,
      name: (form.value.name || '').trim(),
      amount: form.value.amount != null ? Number(form.value.amount) : null,
      period,
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
      isEditing.value ? 'Cập nhật ngân sách thành công' : 'Tạo ngân sách thành công'
    )
  } catch (error) {
    console.error('Failed to save budget:', error)
    const backendMessage = error?.response?.data?.message || ''

    // Hiển thị message rõ ràng hơn khi trùng ngân sách active cùng danh mục & khoảng thời gian
    if (
      backendMessage.includes('already an active budget for this category') ||
      backendMessage.includes('another active budget for this category')
    ) {
      useAppStore().showError(
        'Đã tồn tại ngân sách đang hoạt động cho danh mục này trong khoảng thời gian bạn chọn. Vui lòng chỉnh lại thời gian hoặc tắt ngân sách cũ trước khi tạo mới.'
      )
    } else if (backendMessage) {
      useAppStore().showError(backendMessage)
    } else {
      useAppStore().showError('Không thể lưu ngân sách')
    }
  } finally {
    saving.value = false
  }
}

const deleteBudget = async (budgetId) => {
  if (confirm('Bạn có chắc chắn muốn xoá ngân sách này?')) {
    try {
      await budgetAPI.deleteBudget(budgetId)
      await loadBudgets()
      useAppStore().showSuccess('Đã xoá ngân sách thành công')
    } catch (error) {
      console.error('Failed to delete budget:', error)
      useAppStore().showError('Không thể xoá ngân sách')
    }
  }
}

// Helper methods
const getProgressPercentage = (budget) => {
  if (budget.amount === 0) return 0
  return Math.min((budget.spent_amount / budget.amount) * 100, 100).toFixed(2)
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
        labels: ['Ngân sách', 'Đã chi'],
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
          label: 'Ngân sách',
          data: budgets.value.map(b => b.amount),
          backgroundColor: '#2196F3'
        }, {
          label: 'Đã chi',
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
  await loadInsights()
  await createCharts()
})
</script>

<style scoped>
.v-progress-linear {
  border-radius: 4px;
}
</style>
