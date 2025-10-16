import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const loading = ref(false)
  const snackbar = ref({
    show: false,
    message: '',
    color: 'info',
    timeout: 4000,
  })

  const setLoading = (value) => {
    loading.value = value
  }

  const showSnackbar = (message, color = 'info', timeout = 4000) => {
    snackbar.value = { show: true, message, color, timeout }
  }

  const hideSnackbar = () => {
    snackbar.value.show = false
  }

  const showSuccess = (message) => { showSnackbar(message, 'success') }
  const showError = (message) => { showSnackbar(message, 'error') }
  const showWarning = (message) => { showSnackbar(message, 'warning') }
  const showInfo = (message) => { showSnackbar(message, 'info') }

  return {
    loading,
    snackbar,
    setLoading,
    showSnackbar,
    hideSnackbar,
    showSuccess,
    showError,
    showWarning,
    showInfo,
  }
})


