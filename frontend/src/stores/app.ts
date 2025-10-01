import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // State
  const loading = ref(false)
  const snackbar = ref({
    show: false,
    message: '',
    color: 'info',
    timeout: 4000,
  })

  // Actions
  const setLoading = (value: boolean) => {
    loading.value = value
  }

  const showSnackbar = (message: string, color: string = 'info', timeout: number = 4000) => {
    snackbar.value = {
      show: true,
      message,
      color,
      timeout,
    }
  }

  const hideSnackbar = () => {
    snackbar.value.show = false
  }

  const showSuccess = (message: string) => {
    showSnackbar(message, 'success')
  }

  const showError = (message: string) => {
    showSnackbar(message, 'error')
  }

  const showWarning = (message: string) => {
    showSnackbar(message, 'warning')
  }

  const showInfo = (message: string) => {
    showSnackbar(message, 'info')
  }

  return {
    // State
    loading,
    snackbar,
    
    // Actions
    setLoading,
    showSnackbar,
    hideSnackbar,
    showSuccess,
    showError,
    showWarning,
    showInfo,
  }
})
