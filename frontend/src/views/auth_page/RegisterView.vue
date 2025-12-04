<template>
  <v-container fluid class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center pa-8">
            <v-icon size="48" color="primary" class="mb-4">mdi-account-plus</v-icon>
            <h2 class="text-h4 font-weight-bold">Tạo tài khoản mới</h2>
            <p class="text-subtitle-1 text-medium-emphasis mt-2">
              Tham gia TabiMoney để quản lý tài chính với AI
            </p>
          </v-card-title>

          <v-card-text class="pa-8">
            <v-form @submit.prevent="handleRegister" ref="formRef">
              <v-text-field v-model="form.email" label="Email" type="email" :rules="emailRules"
                :error-messages="errors.email" required class="mb-3" />

              <v-text-field
                v-model="form.username"
                label="Tên đăng nhập"
                :rules="[(v) => !!v || 'Vui lòng nhập tên đăng nhập']"
                :error-messages="errors.username" required class="mb-3" />

              <v-text-field v-model="form.password" label="Mật khẩu" type="password" :rules="passwordRules"
                :error-messages="errors.password" required class="mb-3" />

              <v-text-field v-model="form.first_name" label="Tên" class="mb-3" />
              <v-text-field v-model="form.last_name" label="Họ" class="mb-3" />
              <v-text-field v-model="form.phone" label="Số điện thoại" class="mb-6" />

              <v-btn type="submit" color="primary" size="large" block :loading="loading" :disabled="!isFormValid">
                Tạo tài khoản
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-actions class="pa-8 pt-0">
            <v-row>
              <v-col cols="12" class="text-center">
                <p class="text-body-2">
                  Đã có tài khoản?
                  <router-link to="/auth/login" class="text-primary text-decoration-none">
                    Đăng nhập tại đây
                  </router-link>
                </p>
              </v-col>
            </v-row>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { computed, reactive, ref } from 'vue'

const authStore = useAuthStore()
const appStore = useAppStore()

const formRef = ref()
const loading = ref(false)

const form = reactive({
  email: '',
  username: '',
  password: '',
  first_name: '',
  last_name: '',
  phone: '',
})

const errors = reactive({
  email: [],
  username: [],
  password: [],
})

const emailRules = [
  (v) => !!v || 'Vui lòng nhập email',
  (v) => /.+@.+\..+/.test(v) || 'Email không hợp lệ',
]

const passwordRules = [
  (v) => !!v || 'Vui lòng nhập mật khẩu',
  (v) => v.length >= 6 || 'Mật khẩu phải có ít nhất 6 ký tự',
]

const isFormValid = computed(() => {
  return (
    form.email && /.+@.+\..+/.test(form.email) &&
    form.username &&
    form.password && form.password.length >= 6
  )
})

const handleRegister = async () => {
  if (!formRef.value?.validate()) return

  loading.value = true
  errors.email = []
  errors.username = []
  errors.password = []

  try {
    await authStore.register({
      email: form.email,
      username: form.username,
      password: form.password,
      first_name: form.first_name || undefined,
      last_name: form.last_name || undefined,
      phone: form.phone || undefined,
    })

    appStore.showSuccess('Đăng ký tài khoản thành công!')
  } catch (error) {
    const message = error?.message || ''
    if (message.includes('email')) {
      errors.email = [message]
    } else if (message.includes('username')) {
      errors.username = [message]
    } else if (message.includes('password')) {
      errors.password = [message]
    } else {
      appStore.showError(message || 'Đăng ký thất bại')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
  background: linear-gradient(135deg, #43cea2 0%, #185a9d 100%);
}

.v-card {
  border-radius: 16px;
}
</style>
