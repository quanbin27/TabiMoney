<template>
  <v-container class="py-8">
    <v-row class="mb-4" align="center" justify="space-between">
      <v-col cols="12" md="6">
        <h1 class="text-h5">Transactions</h1>
      </v-col>
      <v-col cols="12" md="6" class="text-right">
        <v-btn color="primary" :to="{ name: 'AddTransaction' }">Add Transaction</v-btn>
      </v-col>
    </v-row>

    <v-data-table :headers="headers" :items="items" :loading="loading" :items-per-page="10">
      <template #item.amount="{ item }">
        {{ formatCurrency(item.amount) }}
      </template>
      <template #item.transaction_date="{ item }">
        {{ new Date(item.transaction_date).toLocaleDateString() }}
      </template>
    </v-data-table>
  </v-container>
  
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { transactionAPI } from '../services/api'
import { useAppStore } from '../stores/app'

const app = useAppStore()
const headers = [
  { title: 'Date', key: 'transaction_date' },
  { title: 'Type', key: 'transaction_type' },
  { title: 'Category', key: 'category.name' },
  { title: 'Description', key: 'description' },
  { title: 'Amount', key: 'amount' },
]

const items = ref([])
const loading = ref(false)

function formatCurrency(value) {
  return new Intl.NumberFormat(undefined, { style: 'currency', currency: 'VND' }).format(value)
}

async function load() {
  loading.value = true
  try {
    const { data } = await transactionAPI.getTransactions({ page: 1, limit: 50 })
    items.value = Array.isArray(data.data) ? data.data : []
  } catch (e) {
    app.showError(e?.message || 'Failed to load transactions')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
