<template>
  <v-container fluid class="pa-4">
    <!-- Header -->
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card color="primary" dark>
          <v-card-title class="d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-icon class="mr-3" size="large">mdi-target</v-icon>
              <div>
                <h2 class="text-h5">Mục tiêu tài chính</h2>
                <p class="text-subtitle-1 mb-0 opacity-90">
                  Đặt và theo dõi các mục tiêu tài chính của bạn
                </p>
              </div>
            </div>
            <v-btn color="white" variant="outlined" @click="openCreateDialog" prepend-icon="mdi-plus">
              Mục tiêu mới
            </v-btn>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <!-- Goals Grid -->
    <GirdGoalView v-if="goals.length > 0" :goals='goals' @open-details="(goal) => handleOpenDetail(goal)"
      @open-contribution="(goal) => handleContribution(goal)" @open-edit="(goal) => handleOpenEdit(goal)"
      @on-delete="(goalID) => deleteGoal(goalID)" />

    <!-- Empty State -->
    <v-row v-else>
      <v-col cols="12">
        <v-card class="text-center pa-8">
          <v-icon size="64" color="grey-lighten-1">mdi-target</v-icon>
          <h3 class="text-h6 mt-4 text-grey">Chưa có mục tiêu nào</h3>
          <p class="text-body-2 text-grey mb-4">
            Bắt đầu bằng cách tạo mục tiêu tài chính đầu tiên của bạn
          </p>
          <v-btn color="primary" @click="openCreateDialog" prepend-icon="mdi-plus">
            Tạo mục tiêu đầu tiên
          </v-btn>
        </v-card>
      </v-col>
    </v-row>

    <!-- Create/Edit Goal Dialog -->
    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title>
          {{ isEditing ? 'Chỉnh sửa mục tiêu' : 'Tạo mục tiêu mới' }}
        </v-card-title>
        <v-card-text>
          <v-form ref="formRef" validate-on="submit">
            <v-text-field v-model="form.title" label="Tên mục tiêu" :rules="[rules.required]" required
              clearable></v-text-field>

            <v-textarea v-model="form.description" label="Mô tả" rows="3"></v-textarea>

            <v-text-field v-model="form.target_amount" label="Số tiền mục tiêu" type="number"
              :rules="[rules.required, rules.positive]" prefix="₫" required clearable></v-text-field>

            <v-text-field v-model="form.current_amount" label="Số tiền hiện tại" type="number" :rules="[rules.positive]"
              prefix="₫" clearable></v-text-field>

            <v-text-field v-model.lazy="form.target_date" label="Ngày hoàn thành dự kiến" type="date"
              :rules="[rules.required, rules.futureDate]" required></v-text-field>

            <v-select v-model="form.goal_type" label="Loại mục tiêu" :items="goalTypes" item-title="title"
              item-value="value" :rules="[rules.required]" required></v-select>

            <v-select v-model="form.priority" label="Mức độ ưu tiên" :items="priorities" item-title="title" item-value="value"
              :rules="[rules.required]" required></v-select>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeDialog">Hủy</v-btn>
          <v-btn type="submit" color="primary" @click="saveGoal" :loading="saving">
            {{ isEditing ? 'Cập nhật' : 'Tạo mới' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Add Contribution Dialog -->
    <AddContributionView v-model="contributionDialog" :rules="rules" :selectedGoal="selectedGoal"
      :loadGoals="loadGoals" />
  </v-container>
  <!-- Goal Details Dialog (read-only) placed outside conditional blocks to avoid v-else adjacency issues -->
  <DetailGoalView v-model="detailsDialog" :detailsGoal='detailsGoal' @open-edit="(goal) => handleOpenEdit(goal)" />
</template>

<script setup>
import { goalAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { onMounted, ref } from 'vue'
import { AddContributionView, DetailGoalView, GirdGoalView } from './component/index'
// Reactive data
const goals = ref([])
const dialog = ref(false)
const contributionDialog = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const formValid = ref(false)
const selectedGoal = ref(null)
const detailsDialog = ref(false)
const detailsGoal = ref(null)

// Form data
const formRef = ref()
const form = ref({
  title: '',
  description: '',
  target_amount: null,
  current_amount: 0,
  target_date: '',
  goal_type: null,
  priority: null,
})

// Options
const goalTypes = [
  { title: 'Tiết kiệm', value: 'savings' },
  { title: 'Trả nợ', value: 'debt_payment' },
  { title: 'Đầu tư', value: 'investment' },
  { title: 'Mua sắm', value: 'purchase' },
  { title: 'Khác', value: 'other' }
]

const priorities = [
  { title: 'Thấp', value: 'low' },
  { title: 'Trung bình', value: 'medium' },
  { title: 'Cao', value: 'high' },
  { title: 'Khẩn cấp', value: 'urgent' }
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
    return value >= 0 || 'Giá trị phải lớn hơn hoặc bằng 0'
  },
  futureDate: (value) => {
    if (!value) return true
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    return new Date(value) > today || 'Ngày phải nằm trong tương lai'
  }
}

// Methods
const loadGoals = async () => {
  try {
    const response = await goalAPI.getGoals()
    goals.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load goals:', error)
    useAppStore().showError('Không thể tải danh sách mục tiêu')
  }
}

const handleOpenDetail = (goal) => {
  detailsDialog.value = true
  detailsGoal.value = goal
}

const openCreateDialog = () => {
  isEditing.value = false
  form.value = {
    title: '',
    description: '',
    target_amount: null,
    current_amount: 0,
    target_date: '',
    goal_type: null,
    priority: null
  }
  dialog.value = true
}

const handleOpenEdit = (goal) => {
  isEditing.value = true
  form.value = {
    title: goal.title,
    description: goal.description,
    target_amount: goal.target_amount,
    current_amount: goal.current_amount,
    target_date: goal.target_date ? new Date(goal.target_date).toISOString().slice(0, 10) : '',
    goal_type: goal.goal_type,
    priority: goal.priority
  }
  selectedGoal.value = goal
  dialog.value = true
}

const closeDialog = () => {
  dialog.value = false
  selectedGoal.value = null
}

const saveGoal = async () => {
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) return
  }
  saving.value = true
  try {
    const basePayload = {
      title: (form.value.title || '').trim(),
      description: form.value.description?.trim() || '',
      target_amount: form.value.target_amount != null ? Number(form.value.target_amount) : null,
      target_date: form.value.target_date ? `${form.value.target_date}T23:59:59Z` : null,
      goal_type: form.value.goal_type,
      priority: form.value.priority,
    }

    if (isEditing.value) {
      const updatePayload = {
        ...basePayload,
        current_amount: form.value.current_amount != null ? Number(form.value.current_amount) : 0,
      }
      await goalAPI.updateGoal(selectedGoal.value.id, updatePayload)
    } else {
      await goalAPI.createGoal(basePayload)
    }
    await loadGoals()
    closeDialog()
    useAppStore().showSuccess(
      isEditing.value ? 'Cập nhật mục tiêu thành công' : 'Tạo mục tiêu thành công'
    )
  } catch (error) {
    console.error('Failed to save goal:', error)
    useAppStore().showError('Không thể lưu mục tiêu')
  } finally {
    saving.value = false
  }
}

const deleteGoal = async (goalId) => {
  if (confirm('Bạn có chắc chắn muốn xoá mục tiêu này?')) {
    try {
      await goalAPI.deleteGoal(goalId)
      await loadGoals()
      useAppStore().showSuccess('Đã xoá mục tiêu thành công')
    } catch (error) {
      console.error('Failed to delete goal:', error)
      useAppStore().showError('Không thể xoá mục tiêu')
    }
  }
}

const handleContribution = (goal) => {
  selectedGoal.value = goal
  contributionDialog.value = true
}

// Lifecycle
onMounted(() => {
  loadGoals()
})
</script>

<style scoped>
.goal-card {
  transition: all 0.3s ease;
}

.goal-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.goal-completed {
  border-left: 4px solid #4caf50;
}

.goal-overdue {
  border-left: 4px solid #f44336;
}

.goal-on_track {
  border-left: 4px solid #2196f3;
}

.goal-in_progress {
  border-left: 4px solid #ff9800;
}

.goal-not_started {
  border-left: 4px solid #9e9e9e;
}
</style>
