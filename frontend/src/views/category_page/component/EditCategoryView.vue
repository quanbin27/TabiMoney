<template>
  <v-container class="py-8" style="max-width: 720px">
    <v-dialog v-model="props.modelValue" persistent max-width="800">
      <v-card>
        <v-card-title>Edit</v-card-title>

        <v-card-text class="pa-6">
          <v-form ref="formRef" @submit.prevent="handleSubmit">
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="form.name" label="Category Name" :rules="[v => !!v || 'Name is required']"
                  required />
              </v-col>

              <v-col cols="12">
                <v-text-field v-model="form.name_en" label="English Name (Optional)" />
              </v-col>

              <v-col cols="12">
                <v-textarea v-model="form.description" label="Description (Optional)" rows="3" />
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>

        <v-card-actions class="pa-6 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="closeDialog">
            Cancel
          </v-btn>
          <v-btn color="primary" :loading="loading" @click="handleSubmit">
            Update Category
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { categoryAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { onMounted, reactive, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  item: {
    type: Object,
    require: true
  }
})

const emit = defineEmits(['update:modelValue'])

const form = ref({
  name: '',
  name_en: '',
  description: '',
})

watch(
  () => props.item,
  (newVal) => {
    if (newVal) {
      form.value = newVal
    }
  }
)

const app = useAppStore()

const loading = ref(false)
const formRef = ref()
const category = ref(null)

const closeDialog = () => {
  emit('update:modelValue', false)
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

</script>
