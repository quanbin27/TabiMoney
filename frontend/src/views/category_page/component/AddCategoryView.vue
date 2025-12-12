<template>
  <v-container class="py-8" style="max-width: 720px">
    <v-dialog v-model="props.modelValue" persistent max-width="800">
      <v-card>
        <v-card-text class="pa-6">
          <v-form ref="formRef" @submit.prevent="handleSubmit">
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="form.name" label="Tên danh mục" :rules="[v => !!v || 'Tên là bắt buộc']"
                  required />
              </v-col>

              <v-col cols="12">
                <v-text-field v-model="form.name_en" label="Tên tiếng Anh (Tùy chọn)" />
              </v-col>

              <v-col cols="12">
                <v-textarea v-model="form.description" label="Mô tả (Tùy chọn)" rows="3" />
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>

        <v-card-actions class="pa-6 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="closeDialog">
            Hủy
          </v-btn>
          <v-btn color="primary" :loading="loading" @click="handleSubmit">
            Tạo danh mục
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { categoryAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { reactive, ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  load: {
    type: Function,
    require: true
  }
})

const emit = defineEmits(['update:modelValue'])

const app = useAppStore()

const loading = ref(false)
const formRef = ref()

const form = reactive({
  name: '',
  name_en: '',
  description: '',
})

const closeDialog = () => {
  formRef.value?.reset()
  emit('update:modelValue', false)
}

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

    app.showSuccess('Tạo danh mục thành công!')
    props.load()
    closeDialog()
  } catch (e) {
    app.showError(e?.message || 'Tạo danh mục thất bại')
  } finally {
    loading.value = false
  }
}
</script>
