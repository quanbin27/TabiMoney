<template>
  <div class="chart-container">
    <canvas ref="chartRef"></canvas>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  DoughnutController,
} from 'chart.js'

ChartJS.register(DoughnutController, ArcElement, Tooltip, Legend)

const props = defineProps({
  data: { type: Object, required: true },
  options: { type: Object, default: () => ({}) },
})

const chartRef = ref(null)
let chartInstance = null

const createChart = () => {
  if (!chartRef.value || !props.data || !props.data.labels?.length) return

  const config = {
    type: 'doughnut',
    data: props.data,
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'bottom',
          labels: {
            padding: 20,
            usePointStyle: true,
          },
        },
        tooltip: {
          callbacks: {
            label(context) {
              const label = context.label || ''
              const value = context.parsed
              const total = context.dataset.data.reduce((a, b) => a + b, 0)
              const percentage = ((value / total) * 100).toFixed(1)
              return `${label}: ${value.toLocaleString()} (${percentage}%)`
            },
          },
        },
      },
      ...props.options,
    },
  }

  chartInstance = new ChartJS(chartRef.value, config)
}

const updateChart = () => {
  if (!props.data || !props.data.labels?.length) return
  if (chartInstance) {
    chartInstance.data = props.data
    chartInstance.update()
  } else {
    createChart()
  }
}

const destroyChart = () => {
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
}

watch(() => props.data, () => nextTick(updateChart), { deep: true })

onMounted(() => nextTick(createChart))
onUnmounted(destroyChart)
</script>


<style scoped>
.chart-container {
  position: relative;
  height: 300px;
  width: 100%;
}
</style>
