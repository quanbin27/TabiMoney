<template>
  <div class="text-center mt-10">
    <h1>404 - Page Not Found</h1>
    <p>Trang bạn truy cập không tồn tại. Sẽ quay về trang chủ sau {{ countdown }} giây...</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const countdown = ref(3)

onMounted(() => {
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value === 0) {
      clearInterval(timer)
      router.push('/') // Quay về trang chủ
    }
  }, 1000)
})
</script>
