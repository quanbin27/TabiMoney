<template>
  <v-container class="py-8">
    <div class="d-flex justify-space-between align-center mb-6">
      <h1 class="text-h4">Categories</h1>
      <v-btn color="primary" prepend-icon="mdi-plus" @click="AddNew">
        Add Category
      </v-btn>
    </div>
    <TableBudgetView :categories="categories" :totalCategories="totalCategories" :onDelete="deleteCategory"
      @on-edit="(item) => handleOnEdit(item)" />
  </v-container>
  <AddCategoryView v-model="isShowAddDialog" />
  <EditCategoryView v-model="isShowEditDialog" :item="selectedItem" />
</template>

<script setup>
import { categoryAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { onMounted, ref } from 'vue'
import { AddCategoryView, EditCategoryView, TableBudgetView } from './component'

const app = useAppStore()

const loading = ref(false)
const categories = ref([])
const filteredCategories = ref([])
const totalCategories = ref(0)
const isShowAddDialog = ref(false);
const isShowEditDialog = ref(false);
const selectedItem = ref(null);

async function loadCategories(options = {}) {
  loading.value = true
  try {
    const response = await categoryAPI.getCategories()
    categories.value = response.data || []
    filteredCategories.value = categories.value
    totalCategories.value = categories.value.length
  } catch (e) {
    app.showError(e?.message || 'Failed to load categories')
  } finally {
    loading.value = false
  }
}

const AddNew = () => {
  isShowAddDialog.value = true
}

const handleOnEdit = (item) => {
  isShowEditDialog.value = true
  selectedItem.value = item
}

function editCategory(category) {
}

async function deleteCategory(category) {
  if (category.is_system) {
    app.showWarning('Cannot delete system categories')
    return
  }

  if (confirm(`Are you sure you want to delete "${category.name}"?`)) {
    try {
      await categoryAPI.deleteCategory(category.id)
      app.showSuccess('Category deleted')
      await loadCategories()
    } catch (e) {
      app.showError(e?.message || 'Failed to delete category')
    }
  }
}

onMounted(() => {
  loadCategories()
})
</script>
