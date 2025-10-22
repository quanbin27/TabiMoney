import { defineStore } from 'pinia'
import { ref } from 'vue'
import Swal from 'sweetalert2'

export const useAppStore = defineStore('app', () => {
  const loading = ref(false)

  const setLoading = (value) => {
    loading.value = value
  }

  const showAlert = (message, icon = 'info', timer = 3000) => {
    Swal.fire({
      text: message,
      icon,
      timer,
      showConfirmButton: false,
      toast: true,
      position: 'top-end',
      timerProgressBar: true,
    })
  }

  const showSuccess = (message) => showAlert(message, 'success')
  const showError = (message) => showAlert(message, 'error')
  const showWarning = (message) => showAlert(message, 'warning')
  const showInfo = (message) => showAlert(message, 'info')

  const confirm = async ({
    title = 'Bạn có chắc không?',
    text = '',
    icon = 'warning',
    confirmButtonText = 'Đồng ý',
    cancelButtonText = 'Hủy',
  } = {}) => {
    const result = await Swal.fire({
      title,
      text,
      icon,
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText,
      cancelButtonText,
    })
    return result.isConfirmed
  }

  return {
    loading,
    setLoading,
    showAlert,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    confirm,
  }
})
