<template>
  <v-container class="py-8" style="max-width: 800px">
    <div class="d-flex align-center mb-6">
      <v-btn icon="mdi-arrow-left" variant="text" @click="router.back()" class="mr-4" />
      <h1 class="text-h4">Edit Goal</h1>
      <v-spacer />
      <v-btn color="error" variant="text" :loading="deleting" @click="onDelete">
        <v-icon class="mr-1">mdi-delete</v-icon> Delete
      </v-btn>
    </div>

    <v-card>
      <v-card-text class="pa-6">
        <v-form ref="formRef" @submit.prevent="onSubmit">
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field v-model="form.title" label="Title" :rules="[rules.required]" required />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="form.goal_type"
                label="Goal Type"
                :items="goalTypes"
                item-title="title"
                item-value="value"
                :rules="[rules.required]"
                required
              />
            </v-col>
            <v-col cols="12">
              <v-textarea v-model="form.description" label="Description" rows="3" />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.target_amount"
                label="Target Amount"
                type="number"
                :rules="[rules.required, rules.nonNegative]"
                prefix="₫"
                required
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.current_amount"
                label="Current Amount"
                type="number"
                :rules="[rules.nonNegative]"
                prefix="₫"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.lazy="form.target_date"
                label="Target Date"
                type="date"
                :rules="[rules.required]"
                required
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="form.priority"
                label="Priority"
                :items="priorities"
                item-title="title"
                item-value="value"
                :rules="[rules.required]"
                required
              />
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>
      <v-card-actions class="pa-6 pt-0">
        <v-spacer />
        <v-btn variant="text" @click="router.back()">Cancel</v-btn>
        <v-btn color="primary" :loading="saving" @click="onSubmit">Save Changes</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
  
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { goalAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'

const route = useRoute()
const router = useRouter()
const app = useAppStore()

const id = Number(route.params.id)
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const formRef = ref()

const form = reactive({
  title: '',
  description: '',
  target_amount: null as number | null,
  current_amount: 0 as number | null,
  target_date: '' as string,
  goal_type: null as string | null,
  priority: null as string | null,
})

const goalTypes = [
  { title: 'Savings', value: 'savings' },
  { title: 'Debt Payment', value: 'debt_payment' },
  { title: 'Investment', value: 'investment' },
  { title: 'Purchase', value: 'purchase' },
  { title: 'Other', value: 'other' },
]

const priorities = [
  { title: 'Low', value: 'low' },
  { title: 'Medium', value: 'medium' },
  { title: 'High', value: 'high' },
  { title: 'Urgent', value: 'urgent' },
]

const rules = {
  required: (v: any) => {
    if (v === null || v === undefined) return 'This field is required'
    if (typeof v === 'number') return true
    if (typeof v === 'string') return v.trim().length > 0 || 'This field is required'
    return !!v || 'This field is required'
  },
  nonNegative: (v: any) => {
    if (v === null || v === undefined) return true
    return Number(v) >= 0 || 'Value must be ≥ 0'
  },
}

async function loadGoal() {
  try {
    loading.value = true
    const res = await goalAPI.getGoal(id)
    const g = res.data?.data ?? res.data
    form.title = g.title
    form.description = g.description || ''
    form.target_amount = g.target_amount
    form.current_amount = g.current_amount
    form.goal_type = g.goal_type
    form.priority = g.priority
    // Convert RFC3339 to yyyy-mm-dd
    form.target_date = g.target_date ? new Date(g.target_date).toISOString().slice(0, 10) : ''
  } catch (e: any) {
    app.showError(e?.message || 'Failed to load goal')
  } finally {
    loading.value = false
  }
}

async function onSubmit() {
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) return
  }
  saving.value = true
  try {
    await goalAPI.updateGoal(id, {
      title: (form.title || '').trim(),
      description: form.description?.trim() || '',
      target_amount: Number(form.target_amount),
      current_amount: Number(form.current_amount ?? 0),
      target_date: form.target_date ? `${form.target_date}T23:59:59Z` : null,
      goal_type: form.goal_type,
      priority: form.priority,
    })
    app.showSuccess('Goal updated')
    router.push({ name: 'Goals' })
  } catch (e: any) {
    app.showError(e?.message || 'Failed to update goal')
  } finally {
    saving.value = false
  }
}

async function onDelete() {
  if (!confirm('Delete this goal?')) return
  deleting.value = true
  try {
    await goalAPI.deleteGoal(id)
    app.showSuccess('Goal deleted')
    router.push({ name: 'Goals' })
  } catch (e: any) {
    app.showError(e?.message || 'Failed to delete goal')
  } finally {
    deleting.value = false
  }
}

onMounted(loadGoal)
</script>

<style scoped>
</style>
