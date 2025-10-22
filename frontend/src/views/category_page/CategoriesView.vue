<template>
  <v-container class="py-8">
    <div class="d-flex justify-space-between align-center mb-6">
      <h1 class="text-h4">Categories</h1>
      <v-btn color="primary" prepend-icon="mdi-plus">
        Add Category
      </v-btn>
    </div>

    <v-card>
      <v-data-table :headers="headers" :items="categories" :loading="loading" :items-per-page="pagination.limit"
        :total-items="totalCategories" @update:options="loadCategories" class="elevation-1">
        <template v-slot:top>
          <v-toolbar flat>
            <v-toolbar-title>Manage Categories</v-toolbar-title>
            <v-divider class="mx-4" inset vertical></v-divider>
            <v-spacer></v-spacer>
            <v-text-field v-model="search" append-icon="mdi-magnify" label="Search categories" single-line hide-details
              @input="filterCategories"></v-text-field>
          </v-toolbar>
        </template>

        <template v-slot:item.is_system="{ item }">
          <v-chip :color="item.is_system ? 'success' : 'primary'" size="small">
            {{ item.is_system ? 'System' : 'Custom' }}
          </v-chip>
        </template>


        <template v-slot:item.actions="{ item }">
          <v-btn v-if="!item.is_system" icon="mdi-pencil" size="small" variant="text"
            @click="editCategory(item)"></v-btn>
          <v-btn v-if="!item.is_system" icon="mdi-delete" size="small" variant="text" color="error"
            @click="deleteCategory(item)"></v-btn>
          <span v-if="item.is_system" class="text-caption text-grey">System</span>
        </template>

        <template v-slot:no-data>
          <v-alert :value="true" color="info" icon="mdi-information">
            No categories found. Create your first category!
          </v-alert>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>

<script setup>
import { categoryAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const app = useAppStore()

const loading = ref(false)
const categories = ref([])
const filteredCategories = ref([])
const totalCategories = ref(0)
const search = ref('')

const pagination = ref({
  page: 1,
  limit: 20,
})

const headers = [
  { title: 'Name', key: 'name', sortable: true },
  { title: 'English Name', key: 'name_en', sortable: true },
  { title: 'Description', key: 'description', sortable: false },
  { title: 'Type', key: 'is_system', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false },
]

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

function filterCategories() {
  if (!search.value.trim()) {
    filteredCategories.value = categories.value
    return
  }

  const query = search.value.toLowerCase()
  filteredCategories.value = categories.value.filter(cat =>
    cat.name?.toLowerCase().includes(query) ||
    cat.name_en?.toLowerCase().includes(query) ||
    cat.description?.toLowerCase().includes(query)
  )
}

function editCategory(category) {
  router.push({ name: 'EditCategory', params: { id: category.id } })
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
