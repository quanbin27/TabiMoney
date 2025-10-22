<template>
  <v-container fluid class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center pa-8">
            <v-icon size="48" color="primary" class="mb-4">mdi-wallet</v-icon>
            <h2 class="text-h4 font-weight-bold">TabiMoney</h2>
            <p class="text-subtitle-1 text-medium-emphasis mt-2">
              AI-Powered Personal Finance Management
            </p>
          </v-card-title>

          <v-card-text class="pa-8">
            <v-form @submit.prevent="handleLogin" ref="formRef">
              <v-text-field v-model="form.email" label="Email" type="email" variant="outlined" :rules="emailRules"
                :error-messages="errors.email" required class="mb-4" />

              <v-text-field v-model="form.password" label="Password" type="password" variant="outlined"
                :rules="passwordRules" :error-messages="errors.password" required class="mb-4" />

              <v-checkbox v-model="form.remember" label="Remember me" color="primary" class="mb-4" />

              <v-btn type="submit" color="primary" size="large" block :loading="loading" :disabled="!isFormValid">
                Login
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-actions class="pa-8 pt-0">
            <v-row>
              <v-col cols="12" class="text-center">
                <p class="text-body-2">
                  Don't have an account?
                  <router-link to="/auth/register" class="text-primary text-decoration-none">
                    Sign up here
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
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// Form data
const form = reactive({
  email: '',
  password: '',
  remember: false,
})

// Form validation
const emailRules = [
  (v) => !!v || 'Email is required',
  (v) => /.+@.+\..+/.test(v) || 'Email must be valid',
]

const passwordRules = [
  (v) => !!v || 'Password is required',
  (v) => v.length >= 6 || 'Password must be at least 6 characters',
]

// Form state
const formRef = ref()
const loading = ref(false)
const errors = reactive({
  email: [],
  password: [],
})

// Computed
const isFormValid = computed(() => {
  return form.email && form.password && form.email.includes('@') && form.password.length >= 6
})

// Methods
const handleLogin = async () => {
  if (!formRef.value) return
  const result = await formRef.value.validate()
  if (!result?.valid) return

  loading.value = true
  errors.email = []
  errors.password = []

  try {
    await authStore.login({
      email: form.email,
      password: form.password,
    })

    appStore.showSuccess('Login successful!')
  } catch (error) {
    if (error.message.includes('email')) {
      errors.email = [error.message]
    } else if (error.message.includes('password')) {
      errors.password = [error.message]
    } else {
      appStore.showError(error.message)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.v-card {
  border-radius: 16px;
}

.v-text-field {
  margin-bottom: 16px;
}
</style>
