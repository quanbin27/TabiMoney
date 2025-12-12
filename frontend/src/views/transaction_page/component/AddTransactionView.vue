<template>
  <v-container class="py-8" style="max-width: 720px">
    <v-dialog v-model="props.modelValue" persistent max-width="800">
      <v-card>
        <v-card-title class="text-h6">Thêm giao dịch</v-card-title>
        <v-card-text>
          <v-form ref="formRef" @submit.prevent="handleSubmit">
            <v-row>
              <v-col cols="12" md="6">
                <v-select :items="typeItems" item-title="label" item-value="value" v-model="form.transaction_type"
                  label="Loại" :rules="[v => !!v || 'Loại là bắt buộc']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model.number="form.amount" type="number" label="Số tiền"
                  :rules="[v => !!v || 'Số tiền là bắt buộc', v => v > 0 || 'Số tiền phải > 0']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-select :items="props.categoryItems" :item-title="formatCategoryTitle" item-value="id"
                  v-model="form.category_id" label="Danh mục" :rules="[v => !!v || 'Danh mục là bắt buộc']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="form.transaction_date" type="date" label="Ngày"
                  :rules="[v => !!v || 'Ngày là bắt buộc']" required />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="form.transaction_time" type="time" label="Giờ (Tùy chọn)" />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="form.description" label="Mô tả (Tùy chọn)" />
              </v-col>
              <v-col cols="12">
                <v-text-field v-model="form.location" label="Địa điểm (Tùy chọn)" />
              </v-col>
              <v-col cols="12" class="d-flex justify-end">
                <v-btn variant="outlined" color="secondary" class="mr-2" :loading="aiLoading"
                  @click="suggestCategory">AI Gợi ý danh mục</v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="closeDialog">Đóng</v-btn>
          <v-btn @click="handleSubmit" variant="tonal" color="primary" :loading="loading">Tạo</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- // dialog ai create category -->
    <v-dialog v-model="createDialog" max-width="500">
      <v-card>
        <v-card-title>Tạo danh mục từ gợi ý AI</v-card-title>
        <v-card-text>
          <v-text-field v-model="newCategory.name" label="Tên danh mục" required />
          <v-text-field v-model="newCategory.description" label="Mô tả" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="createDialog = false">Hủy</v-btn>
          <v-btn color="primary" @click="createCategoryFromSuggestion">Sử dụng danh mục này</v-btn>
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

const emit = defineEmits(['update:modelValue', 'category-created'])

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
  { label: 'Chi tiêu', value: 'expense' },
  { label: 'Thu nhập', value: 'income' },
  { label: 'Chuyển khoản', value: 'transfer' },
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
      app.showError('Vui lòng chọn danh mục')
      return
    }
    if (!payload.amount || payload.amount <= 0) {
      app.showError('Vui lòng nhập số tiền hợp lệ')
      return
    }
    if (!payload.transaction_type) {
      app.showError('Vui lòng chọn loại giao dịch')
      return
    }
    if (!payload.transaction_date) {
      app.showError('Vui lòng chọn ngày')
      return
    }

    await transactionAPI.createTransaction(payload)
    app.showSuccess('Tạo thành công')
    closeDialog()
  } catch (e) {
    app.showError(e?.message || 'Tạo thất bại')
  } finally {
    loading.value = false
  }
}

async function suggestCategory() {
  if (!form.description || !form.amount) {
    app.showWarning('Vui lòng nhập mô tả và số tiền trước')
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
        const found = props.categoryItems.find((c) => {
          const existingName = (c.name || '').toLowerCase().trim()
          const existingNameEn = (c.name_en || '').toLowerCase().trim()
          const topName = (top.category_name || '').toLowerCase().trim()
          return existingName === topName || existingNameEn === topName
        })

        if (found) {
          form.category_id = found.id
          app.showSuccess(`AI gợi ý danh mục hiện có: ${found.name}`)
        } else {
          app.showWarning('AI gợi ý danh mục hiện có nhưng không tìm thấy trong danh sách')
        }
      } else {
        // AI suggests new category
        newCategory.name = top.category_name || ''
        newCategory.description = top.reason || `Được tạo từ gợi ý AI`
        createDialog.value = true
        app.showInfo(`AI gợi ý danh mục mới: "${top.category_name}"`)
      }
    } else {
      app.showInfo('Không có gợi ý')
    }
  } catch (e) {
    app.showError(e?.message || 'Gợi ý AI thất bại')
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
    // Emit event to parent to update categoryItems
    emit('category-created', data)
    // Also update local props if it's an array (for immediate UI update)
    if (Array.isArray(props.categoryItems)) {
      props.categoryItems.push(data)
    }
    form.category_id = data.id
    createDialog.value = false
    app.showSuccess('Tạo danh mục thành công')
  } catch (e) {
    app.showError(e?.message || 'Tạo danh mục thất bại')
  }
}
</script>
