<template>
  <v-container fluid v-if="!loading">
    <Welcome v-model="isShowAddDialog" />
    <Sumary :analytics="analytics" />
    <ChartHeath :analytics="analytics" />
    <RecentTransactions :recentTransactions="recentTransactions" :loadRecentTransactions="loadRecentTransactions"
      @openEdit="handleOpenEditDialog" />
    <AddTransactionView v-model="isShowAddDialog" @update:modelValue="onModelUpdated" :categoryItems="categoryItems" />
    <EditTransactionView v-model="isShowEditDialog" :transaction="selectedTransaction"
      @update:modelValue="onModelUpdated" :categoryItems="categoryItems" />
  </v-container>
</template>

<script setup>
import { analyticsAPI, categoryAPI, transactionAPI } from '@/services/api'
import { onMounted, ref } from 'vue'
import {
  ChartHeath,
  RecentTransactions, Sumary, Welcome
} from "./components"

import { AddTransactionView, EditTransactionView } from '../transaction_page/component/index'

// State
const analytics = ref(null)
const recentTransactions = ref([])
const loading = ref(false)
const isShowAddDialog = ref(false)
const isShowEditDialog = ref(false)
const selectedTransaction = ref(null)
const categoryItems = ref([])

// Methods
const loadDashboardData = async () => {
  loading.value = true
  try {
    const [analyticsResponse, transactionsResponse] = await Promise.all([
      analyticsAPI.getDashboard(),
      transactionAPI.getTransactions({ limit: 10, sort_by: 'transaction_date', sort_order: 'desc' })
    ])

    analytics.value = analyticsResponse.data
    recentTransactions.value = transactionsResponse.data.data
  } catch (error) {
    appStore.showError('Failed to load dashboard data')
  } finally {
    loading.value = false
  }
}

const loadRecentTransactions = async () => {
  try {
    const transactionsResponse = await transactionAPI.getTransactions({ limit: 10, sort_by: 'transaction_date', sort_order: 'desc' })
    recentTransactions.value = transactionsResponse.data.data
  } catch (error) {
    appStore.showError('Failed to load recent transactions')
  }
}

const handleOpenEditDialog = (transaction) => {
  selectedTransaction.value = transaction
  isShowEditDialog.value = true
}

const onModelUpdated = (value) => {
  isShowAddDialog.value = value
  loadRecentTransactions()
}

async function loadCategories() {
  try {
    const { data } = await categoryAPI.getCategories()
    categoryItems.value = Array.isArray(data) ? data : []
  } catch (e) {
    appStore.showError(e?.message || 'Failed to load categories')
  }
}

// Lifecycle
onMounted(() => {
  loadDashboardData()
  loadCategories()
})
</script>

<style scoped>
.v-card {
  border-radius: 12px;
}

.text-success {
  color: rgb(var(--v-theme-success)) !important;
}

.text-error {
  color: rgb(var(--v-theme-error)) !important;
}
</style>
