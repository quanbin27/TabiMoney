<template>
  <v-container class="py-8" style="max-width: 600px">
    <div class="d-flex align-center mb-6">
      <v-btn
        icon="mdi-arrow-left"
        variant="text"
        @click="$router.back()"
        class="mr-4"
      ></v-btn>
      <h1 class="text-h4">Add Category</h1>
    </div>

    <v-card>
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
          Create Category
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { categoryAPI } from '../services/api'
import { useAppStore } from '../stores/app'

const router = useRouter()
const app = useAppStore()

const loading = ref(false)
const formRef = ref()

const form = reactive({
  name: '',
  name_en: '',
  description: '',
})

async function handleSubmit() {
  if (!formRef.value) return
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  try {
    await categoryAPI.createCategory({
      name: form.name,
      name_en: form.name_en || undefined,
      description: form.description || undefined,
    })
    
    app.showSuccess('Category created successfully!')
    router.push({ name: 'Categories' })
  } catch (e: any) {
    app.showError(e?.message || 'Failed to create category')
  } finally {
    loading.value = false
  }
}
</script>
