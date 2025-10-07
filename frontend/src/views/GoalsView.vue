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
                <h2 class="text-h5">Financial Goals</h2>
                <p class="text-subtitle-1 mb-0 opacity-90">
                  Set and track your financial objectives
                </p>
              </div>
            </div>
            <v-btn
              color="white"
              variant="outlined"
              @click="openCreateDialog"
              prepend-icon="mdi-plus"
            >
              New Goal
            </v-btn>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <!-- Goals Grid -->
    <v-row v-if="goals.length > 0">
      <v-col
        v-for="goal in goals"
        :key="goal.id"
        cols="12"
        md="6"
        lg="4"
      >
        <v-card
          class="goal-card"
          :class="getGoalStatusClass(goal)"
          elevation="2"
          hover
        >
          <v-card-title class="d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-icon
                :color="getGoalStatusColor(goal)"
                class="mr-2"
                size="large"
              >
                {{ getGoalIcon(goal) }}
              </v-icon>
              <div>
                <h3 class="text-h6">{{ goal.name }}</h3>
                <p class="text-caption mb-0">{{ goal.description }}</p>
              </div>
            </div>
            <v-menu>
              <template v-slot:activator="{ props }">
                <v-btn
                  icon="mdi-dots-vertical"
                  variant="text"
                  v-bind="props"
                ></v-btn>
              </template>
              <v-list>
                <v-list-item @click="openEditDialog(goal)">
                  <template v-slot:prepend>
                    <v-icon>mdi-pencil</v-icon>
                  </template>
                  <v-list-item-title>Edit</v-list-item-title>
                </v-list-item>
                <v-list-item @click="deleteGoal(goal.id)">
                  <template v-slot:prepend>
                    <v-icon>mdi-delete</v-icon>
                  </template>
                  <v-list-item-title>Delete</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-card-title>

          <v-card-text>
            <!-- Progress Bar -->
            <div class="mb-4">
              <div class="d-flex justify-space-between mb-2">
                <span class="text-body-2">Progress</span>
                <span class="text-body-2 font-weight-bold">
                  {{ getProgressPercentage(goal) }}%
                </span>
              </div>
              <v-progress-linear
                :model-value="getProgressPercentage(goal)"
                :color="getGoalStatusColor(goal)"
                height="8"
                rounded
              ></v-progress-linear>
            </div>

            <!-- Goal Details -->
            <v-row class="text-center">
              <v-col cols="6">
                <div class="text-h6 text-primary">
                  {{ formatCurrency(goal.current_amount) }}
                </div>
                <div class="text-caption text-grey">Current</div>
              </v-col>
              <v-col cols="6">
                <div class="text-h6 text-grey">
                  {{ formatCurrency(goal.target_amount) }}
                </div>
                <div class="text-caption text-grey">Target</div>
              </v-col>
            </v-row>

            <!-- Goal Status -->
            <v-chip
              :color="getGoalStatusColor(goal)"
              size="small"
              class="mt-3"
            >
              {{ getGoalStatusText(goal) }}
            </v-chip>

            <!-- Deadline -->
            <div class="mt-3 text-caption text-grey">
              <v-icon size="small" class="mr-1">mdi-calendar</v-icon>
              Deadline: {{ formatDate(goal.target_date) }}
            </div>
          </v-card-text>

          <v-card-actions>
            <v-btn
              color="primary"
              variant="text"
              @click="viewGoalDetails(goal)"
            >
              View Details
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
              color="success"
              variant="text"
              @click="addContribution(goal)"
            >
              Add Contribution
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- Empty State -->
    <v-row v-else>
      <v-col cols="12">
        <v-card class="text-center pa-8">
          <v-icon size="64" color="grey-lighten-1">mdi-target</v-icon>
          <h3 class="text-h6 mt-4 text-grey">No Goals Yet</h3>
          <p class="text-body-2 text-grey mb-4">
            Start by creating your first financial goal
          </p>
          <v-btn
            color="primary"
            @click="openCreateDialog"
            prepend-icon="mdi-plus"
          >
            Create Your First Goal
          </v-btn>
        </v-card>
      </v-col>
    </v-row>

    <!-- Create/Edit Goal Dialog -->
    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title>
          {{ isEditing ? 'Edit Goal' : 'Create New Goal' }}
        </v-card-title>
        <v-card-text>
          <v-form ref="form" v-model="formValid">
            <v-text-field
              v-model="form.name"
              label="Goal Name"
              :rules="[rules.required]"
              required
            ></v-text-field>
            
            <v-textarea
              v-model="form.description"
              label="Description"
              rows="3"
            ></v-textarea>
            
            <v-text-field
              v-model="form.target_amount"
              label="Target Amount"
              type="number"
              :rules="[rules.required, rules.positive]"
              prefix="₫"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="form.current_amount"
              label="Current Amount"
              type="number"
              :rules="[rules.positive]"
              prefix="₫"
            ></v-text-field>
            
            <v-text-field
              v-model="form.target_date"
              label="Target Date"
              type="date"
              :rules="[rules.required, rules.futureDate]"
              required
            ></v-text-field>
            
            <v-select
              v-model="form.goal_type"
              label="Goal Type"
              :items="goalTypes"
              :rules="[rules.required]"
              required
            ></v-select>
            
            <v-select
              v-model="form.priority"
              label="Priority"
              :items="priorities"
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
            @click="saveGoal"
            :disabled="!formValid"
            :loading="saving"
          >
            {{ isEditing ? 'Update' : 'Create' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Add Contribution Dialog -->
    <v-dialog v-model="contributionDialog" max-width="400">
      <v-card>
        <v-card-title>Add Contribution</v-card-title>
        <v-card-text>
          <v-form ref="contributionForm" v-model="contributionFormValid">
            <v-text-field
              v-model="contributionForm.amount"
              label="Amount"
              type="number"
              :rules="[rules.required, rules.positive]"
              prefix="₫"
              required
            ></v-text-field>
            
            <v-textarea
              v-model="contributionForm.note"
              label="Note (Optional)"
              rows="2"
            ></v-textarea>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeContributionDialog">Cancel</v-btn>
          <v-btn
            color="primary"
            @click="saveContribution"
            :disabled="!contributionFormValid"
            :loading="savingContribution"
          >
            Add Contribution
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { goalAPI } from '@/services/api'
import { formatCurrency, formatDate } from '@/utils/formatters'
import { useAppStore } from '@/stores/app'

// Reactive data
const goals = ref<any[]>([])
const dialog = ref(false)
const contributionDialog = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const savingContribution = ref(false)
const formValid = ref(false)
const contributionFormValid = ref(false)
const selectedGoal = ref<any>(null)

// Form data
const form = ref({
  name: '',
  description: '',
  target_amount: 0,
  current_amount: 0,
  target_date: '',
  goal_type: '',
  priority: ''
})

const contributionForm = ref({
  amount: 0,
  note: ''
})

// Options
const goalTypes = [
  { title: 'Emergency Fund', value: 'emergency_fund' },
  { title: 'Vacation', value: 'vacation' },
  { title: 'Home Purchase', value: 'home_purchase' },
  { title: 'Car Purchase', value: 'car_purchase' },
  { title: 'Education', value: 'education' },
  { title: 'Retirement', value: 'retirement' },
  { title: 'Investment', value: 'investment' },
  { title: 'Other', value: 'other' }
]

const priorities = [
  { title: 'Low', value: 'low' },
  { title: 'Medium', value: 'medium' },
  { title: 'High', value: 'high' },
  { title: 'Critical', value: 'critical' }
]

// Validation rules
const rules = {
  required: (value: any) => !!value || 'This field is required',
  positive: (value: number) => value >= 0 || 'Value must be positive',
  futureDate: (value: string) => {
    if (!value) return true
    return new Date(value) > new Date() || 'Date must be in the future'
  }
}

// Methods
const loadGoals = async () => {
  try {
    const response = await goalAPI.getGoals()
    goals.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load goals:', error)
    useAppStore().showError('Failed to load goals')
  }
}

const openCreateDialog = () => {
  isEditing.value = false
  form.value = {
    name: '',
    description: '',
    target_amount: 0,
    current_amount: 0,
    target_date: '',
    goal_type: '',
    priority: ''
  }
  dialog.value = true
}

const openEditDialog = (goal: any) => {
  isEditing.value = true
  form.value = {
    name: goal.name,
    description: goal.description,
    target_amount: goal.target_amount,
    current_amount: goal.current_amount,
    target_date: goal.target_date,
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
  saving.value = true
  try {
    if (isEditing.value) {
      await goalAPI.updateGoal(selectedGoal.value.id, form.value)
    } else {
      await goalAPI.createGoal(form.value)
    }
    await loadGoals()
    closeDialog()
    useAppStore().showSuccess(
      isEditing.value ? 'Goal updated successfully' : 'Goal created successfully'
    )
  } catch (error) {
    console.error('Failed to save goal:', error)
    useAppStore().showError('Failed to save goal')
  } finally {
    saving.value = false
  }
}

const deleteGoal = async (goalId: number) => {
  if (confirm('Are you sure you want to delete this goal?')) {
    try {
      await goalAPI.deleteGoal(goalId)
      await loadGoals()
      useAppStore().showSuccess('Goal deleted successfully')
    } catch (error) {
      console.error('Failed to delete goal:', error)
      useAppStore().showError('Failed to delete goal')
    }
  }
}

const addContribution = (goal: any) => {
  selectedGoal.value = goal
  contributionForm.value = {
    amount: 0,
    note: ''
  }
  contributionDialog.value = true
}

const closeContributionDialog = () => {
  contributionDialog.value = false
  selectedGoal.value = null
}

const saveContribution = async () => {
  savingContribution.value = true
  try {
    // This would typically be a separate API endpoint
    // For now, we'll update the goal's current amount
    const updatedGoal = {
      ...selectedGoal.value,
      current_amount: selectedGoal.value.current_amount + contributionForm.value.amount
    }
    await goalAPI.updateGoal(selectedGoal.value.id, updatedGoal)
    await loadGoals()
    closeContributionDialog()
    useAppStore().showSuccess('Contribution added successfully')
  } catch (error) {
    console.error('Failed to add contribution:', error)
    useAppStore().showError('Failed to add contribution')
  } finally {
    savingContribution.value = false
  }
}

const viewGoalDetails = (goal: any) => {
  // Navigate to goal details page or show detailed modal
  console.log('View goal details:', goal)
}

// Helper methods
const getProgressPercentage = (goal: any) => {
  if (goal.target_amount === 0) return 0
  return Math.min((goal.current_amount / goal.target_amount) * 100, 100)
}

const getGoalStatus = (goal: any) => {
  const progress = getProgressPercentage(goal)
  const now = new Date()
  const targetDate = new Date(goal.target_date)
  
  if (progress >= 100) return 'completed'
  if (targetDate < now) return 'overdue'
  if (progress >= 75) return 'on_track'
  if (progress >= 50) return 'in_progress'
  return 'not_started'
}

const getGoalStatusText = (goal: any) => {
  const status = getGoalStatus(goal)
  const statusMap = {
    completed: 'Completed',
    overdue: 'Overdue',
    on_track: 'On Track',
    in_progress: 'In Progress',
    not_started: 'Not Started'
  }
  return statusMap[status] || 'Unknown'
}

const getGoalStatusColor = (goal: any) => {
  const status = getGoalStatus(goal)
  const colorMap = {
    completed: 'success',
    overdue: 'error',
    on_track: 'primary',
    in_progress: 'warning',
    not_started: 'grey'
  }
  return colorMap[status] || 'grey'
}

const getGoalStatusClass = (goal: any) => {
  const status = getGoalStatus(goal)
  return `goal-${status}`
}

const getGoalIcon = (goal: any) => {
  const status = getGoalStatus(goal)
  const iconMap = {
    completed: 'mdi-check-circle',
    overdue: 'mdi-alert-circle',
    on_track: 'mdi-trending-up',
    in_progress: 'mdi-progress-clock',
    not_started: 'mdi-target'
  }
  return iconMap[status] || 'mdi-target'
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
