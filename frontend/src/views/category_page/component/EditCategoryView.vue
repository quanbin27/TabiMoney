<template>
  <v-container class="py-8" style="max-width: 600px">
    <div class="d-flex align-center mb-6">
      <v-btn
        icon="mdi-arrow-left"
        variant="text"
        @click="$router.back()"
        class="mr-4"
      ></v-btn>
      <h1 class="text-h4">Edit Category</h1>
    </div>

    <v-card v-if="category">
      <v-card-text class="pa-6">
        <v-form ref="formRef" @submit.prevent="handleSubmit">
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="form.name"
                label="Category Name"
                :rules="[v => !!v || 'Name is required']"
                required
              />
            </v-col>
            
            <v-col cols="12">
              <v-text-field
                v-model="form.name_en"
                label="English Name (Optional)"
              />
            </v-col>
            
            <v-col cols="12">
              <v-textarea
                v-model="form.description"
                label="Description (Optional)"
                rows="3"
              />
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>
      
      <v-card-actions class="pa-6 pt-0">
        <v-spacer />
        <v-btn
          variant="text"
          @click="$router.back()"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :loading="loading"
          @click="handleSubmit"
        >
          Update Category
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card v-else-if="loading">
      <v-card-text class="pa-6 text-center">
        <v-progress-circular indeterminate></v-progress-circular>
        <p class="mt-4">Loading category...</p>
      </v-card-text>
    </v-card>

    <v-card v-else>
      <v-card-text class="pa-6 text-center">
        <v-icon size="64" color="error">mdi-alert-circle</v-icon>
        <h3 class="mt-4">Category not found</h3>
        <p class="text-grey">The category you're looking for doesn't exist.</p>
        <v-btn color="primary" @click="$router.push({ name: 'Categories' })">
          Back to Categories
        </v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { categoryAPI } from '../services/api'
import { useAppStore } from '../stores/app'

const router = useRouter()
const route = useRoute()
const app = useAppStore()

const loading = ref(false)
const formRef = ref()
const category = ref(null)

const form = reactive({
  name: '',
  name_en: '',
  description: '',
})

async function loadCategory() {
  const categoryId = String(route.params.id || '')
  if (!categoryId) {
    router.push({ name: 'Categories' })
    return
  }

  loading.value = true
  try {
    const response = await categoryAPI.getCategory(parseInt(categoryId))
    category.value = response.data
    
    // Populate form
    form.name = category.value.name || ''
    form.name_en = category.value.name_en || ''
    form.description = category.value.description || ''
  } catch (e) {
    app.showError(e?.message || 'Failed to load category')
    router.push({ name: 'Categories' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!formRef.value || !category.value) return
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  try {
    await categoryAPI.updateCategory(category.value.id, {
      name: form.name,
      name_en: form.name_en || undefined,
      description: form.description || undefined,
    })
    
    app.showSuccess('Category updated successfully!')
    router.push({ name: 'Categories' })
  } catch (e) {
    app.showError(e?.message || 'Failed to update category')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategory()
})
</script>
