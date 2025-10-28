<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <h3>Statistics</h3>
        <v-divider class="my-2"></v-divider>
        <p><strong>Total Income:</strong> {{ user.statistics.totalIncome.toLocaleString() }} VND</p>
        <p><strong>Total Expense:</strong> {{ user.statistics.totalExpense.toLocaleString() }} VND</p>
        <p><strong>Transactions:</strong> {{ user.statistics.totalTransactions }}</p>
        <p><strong>Categories Used:</strong> {{ user.statistics.categoriesUsed }}</p>
    </v-card>
</template>

<script setup>
const props = defineProps(['user'])
</script>

<style></style>