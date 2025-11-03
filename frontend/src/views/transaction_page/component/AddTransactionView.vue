<template>
  <v-container class="py-8" style="max-width: 720px">
    <v-dialog v-model="props.modelValue" persistent max-width="800">
      <v-card>
        <v-card-title class="text-h6">Add Transaction</v-card-title>
        <v-card-text>
          <v-form ref="formRef" @submit.prevent="handleSubmit">
            <v-row>
              <v-col cols="12" md="6">
                <v-select :items="typeItems" item-title="label" item-value="value" v-model="form.transaction_type"
                  label="Type" :rules="[v => !!v || 'Type is required']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model.number="form.amount" type="number" label="Amount"
                  :rules="[v => !!v || 'Amount is required', v => v > 0 || 'Amount must be > 0']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-select :items="props.categoryItems" :item-title="formatCategoryTitle" item-value="id"
                  v-model="form.category_id" label="Category" :rules="[v => !!v || 'Category is required']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="form.transaction_date" type="date" label="Date"
                  :rules="[v => !!v || 'Date is required']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="form.transaction_time" type="time" label="Time (Optional)" />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="form.description" label="Description (Optional)" />
              </v-col>
              <v-col cols="12">
                <v-text-field v-model="form.location" label="Location (Optional)" />
              </v-col>
              <v-col cols="12" class="d-flex justify-end">
                <v-btn variant="outlined" color="secondary" class="mr-2" :loading="aiLoading"
                  @click="suggestCategory">AI
                  Suggest Category</v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="closeDialog">Close</v-btn>
          <v-btn @click="handleSubmit" variant="tonal" color="primary" :loading="loading">Create</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- // dialog ai create category -->
    <v-dialog v-model="createDialog" max-width="500">
      <v-card>
        <v-card-title>Create category from AI suggestion</v-card-title>
        <v-card-text>
          <v-text-field v-model="newCategory.name" label="Category Name" required />
          <v-text-field v-model="newCategory.description" label="Description" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="createDialog = false">Cancel</v-btn>
          <v-btn color="primary" @click="createCategoryFromSuggestion">Use this category</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
<script setup>
import { aiAPI, categoryAPI, transactionAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { reactive, ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  categoryItems: {
    type: Array,
    require: true,
    default: []
  }
})

const emit = defineEmits(['update:modelValue'])

const app = useAppStore()
const loading = ref(false)
const formRef = ref()
const aiLoading = ref(false)
const createDialog = ref(false)
const newCategory = reactive({ name: '', description: '' })
function formatCategoryTitle(c) {
  if (!c) return ''
  if (c.name_en) return `${c.name} (${c.name_en})`
  return c.name
}

const typeItems = [
  { label: 'Expense', value: 'expense' },
  { label: 'Income', value: 'income' },
  { label: 'Transfer', value: 'transfer' },
]

const form = reactive({
  category_id: 1,
  amount: 0,
  description: '',
  transaction_type: 'expense',
  transaction_date: new Date().toISOString().slice(0, 10),
  transaction_time: '',
  location: '',
})

const closeDialog = () => {
  emit('update:modelValue', false)
}

async function handleSubmit() {
  if (!formRef.value) return
  const result = await formRef.value.validate()
  if (!result?.valid) return
  loading.value = true
  try {
    const payload = { ...form }

    // Handle transaction_time - keep as HH:MM format or remove if empty
    if (!payload.transaction_time || !payload.transaction_time.trim()) {
      // Remove transaction_time if empty to avoid parsing errors
      delete payload.transaction_time
    }

    // Ensure required fields are present
    if (!payload.category_id) {
      app.showError('Please select a category')
      return
    }
    if (!payload.amount || payload.amount <= 0) {
      app.showError('Please enter a valid amount')
      return
    }
    if (!payload.transaction_type) {
      app.showError('Please select transaction type')
      return
    }
    if (!payload.transaction_date) {
      app.showError('Please select transaction date')
      return
    }

    await transactionAPI.createTransaction(payload)
    app.showSuccess('Created')
    closeDialog()
  } catch (e) {
    app.showError(e?.message || 'Create failed')
  } finally {
    loading.value = false
  }
}

async function suggestCategory() {
  if (!form.description || !form.amount) {
    app.showWarning('Enter description and amount first')
    return
  }
  aiLoading.value = true
  try {
    const { data } = await aiAPI.suggestCategory({
      description: form.description,
      amount: form.amount,
      location: form.location || undefined,
      tags: [],
      existing_categories: props.categoryItems.map(cat => ({
        name: cat.name,
        name_en: cat.name_en,
        description: cat.description
      }))
    })
    const suggestions = data?.suggestions || []
    if (suggestions.length) {
      const top = suggestions[0]

      // AI already determined if category exists
      if (top.is_existing) {
        // Find the existing category by name
        const found = props.categoryItems.value.find((c) => {
          const existingName = (c.name || '').toLowerCase().trim()
          const existingNameEn = (c.name_en || '').toLowerCase().trim()
          const topName = (top.category_name || '').toLowerCase().trim()
          return existingName === topName || existingNameEn === topName
        })

        if (found) {
          form.category_id = found.id
          app.showSuccess(`AI suggested existing category: ${found.name}`)
        } else {
          app.showWarning('AI suggested existing category but not found in list')
        }
      } else {
        // AI suggests new category
        newCategory.name = top.category_name || ''
        newCategory.description = top.reason || `Created from AI suggestion`
        createDialog.value = true
        app.showInfo(`AI suggests new category: "${top.category_name}"`)
      }
    } else {
      app.showInfo('No suggestion')
    }
  } catch (e) {
    app.showError(e?.message || 'AI suggestion failed')
  } finally {
    aiLoading.value = false
  }
}

async function createCategoryFromSuggestion() {
  try {
    const { data } = await categoryAPI.createCategory({
      name: newCategory.name,
      description: newCategory.description || undefined,
    })
    props.categoryItems.value.push(data)
    form.category_id = data.id
    createDialog.value = false
    app.showSuccess('Category created')
  } catch (e) {
    app.showError(e?.message || 'Create category failed')
  }
}
</script>
