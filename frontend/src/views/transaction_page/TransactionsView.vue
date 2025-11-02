<template>
  <v-container>
    <!-- <FilterTransactionView /> -->

    <v-row class="mb-4" align="center" justify="space-between">
      <v-col cols="12" md="6">
        <h1 class="text-h5">Transactions</h1>
      </v-col>
      <v-col cols="12" md="6" class="text-right">
        <v-btn color="primary" @click="isShowAddDialog = true">Add Transaction</v-btn>
      </v-col>
    </v-row>

    <TableTransactionView :items="items" :loading="loading" @openEdit="handleOpenEditDialog" :load="load" />
    <AddTransactionView v-model="isShowAddDialog" @update:modelValue="onModelUpdated" :categoryItems="categoryItems" />
    <EditTransactionView v-model="isShowEditDialog" :transaction="selectedTransaction"
      @update:modelValue="onModelUpdated" :categoryItems="categoryItems" />
  </v-container>
</template>
<script setup>
import { categoryAPI, transactionAPI } from '@/services/api';
import { useAppStore } from '@/stores/app';
import { onMounted, ref } from 'vue';
import { AddTransactionView, EditTransactionView, TableTransactionView } from './component/index';

const isShowAddDialog = ref(false);
const isShowEditDialog = ref(false);
const selectedTransaction = ref(null);
const loading = ref(false);
const items = ref([])
const categoryItems = ref([])

const appStore = useAppStore()

const handleOpenEditDialog = (transaction) => {
  selectedTransaction.value = transaction;
  isShowEditDialog.value = true;
};

const onModelUpdated = (value) => {
  isShowEditDialog.value = value;
  if (!value) {
    load();
  }
};

async function load() {
  loading.value = true
  try {
    const { data } = await transactionAPI.getTransactions({ page: 1, limit: 50 })
    items.value = Array.isArray(data.data) ? data.data : []
  } catch (e) {
    appStore.showError(e?.message || 'Failed to load transactions')
  } finally {
    loading.value = false
  }
}
async function loadCategories() {
  try {
    const { data } = await categoryAPI.getCategories()
    categoryItems.value = Array.isArray(data) ? data : []
  } catch (e) {
    appStore.showError(e?.message || 'Failed to load categories')
  }
}

onMounted(() => {
  load();
  loadCategories()
});
</script>
