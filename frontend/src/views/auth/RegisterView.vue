<template>
  <v-container fluid class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center pa-8">
            <v-icon size="48" color="primary" class="mb-4">mdi-account-plus</v-icon>
            <h2 class="text-h4 font-weight-bold">Create your account</h2>
            <p class="text-subtitle-1 text-medium-emphasis mt-2">
              Join TabiMoney to manage your finances with AI
            </p>
          </v-card-title>

          <v-card-text class="pa-8">
            <v-form @submit.prevent="handleRegister" ref="formRef">
              <v-text-field
                v-model="form.email"
                label="Email"
                type="email"
                :rules="emailRules"
                :error-messages="errors.email"
                required
                class="mb-3"
              />

              <v-text-field
                v-model="form.username"
                label="Username"
                :rules="[(v: string) => !!v || 'Username is required']"
                :error-messages="errors.username"
                required
                class="mb-3"
              />

              <v-text-field
                v-model="form.password"
                label="Password"
                type="password"
                :rules="passwordRules"
                :error-messages="errors.password"
                required
                class="mb-3"
              />

              <v-text-field v-model="form.first_name" label="First name" class="mb-3" />
              <v-text-field v-model="form.last_name" label="Last name" class="mb-3" />
              <v-text-field v-model="form.phone" label="Phone" class="mb-6" />

              <v-btn
                type="submit"
                color="primary"
                size="large"
                block
                :loading="loading"
                :disabled="!isFormValid"
              >
                Create account
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-actions class="pa-8 pt-0">
            <v-row>
              <v-col cols="12" class="text-center">
                <p class="text-body-2">
                  Already have an account?
                  <router-link to="/auth/login" class="text-primary text-decoration-none">
                    Login here
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

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useAppStore } from '../../stores/app'

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
  email: [] as string[],
  username: [] as string[],
  password: [] as string[],
})

const emailRules = [
  (v: string) => !!v || 'Email is required',
  (v: string) => /.+@.+\..+/.test(v) || 'Email must be valid',
]

const passwordRules = [
  (v: string) => !!v || 'Password is required',
  (v: string) => v.length >= 6 || 'Password must be at least 6 characters',
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

    appStore.showSuccess('Registration successful!')
  } catch (error: any) {
    const message = error?.message || ''
    if (message.includes('email')) {
      errors.email = [message]
    } else if (message.includes('username')) {
      errors.username = [message]
    } else if (message.includes('password')) {
      errors.password = [message]
    } else {
      appStore.showError(message || 'Registration failed')
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
