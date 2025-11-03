<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <h3>Thống kê</h3>
        <v-divider class="my-2"></v-divider>
        <div v-if="statistics" class="stats-container">
            <v-row>
                <v-col cols="12" md="6" class="mb-4">
                    <div class="stat-item">
                        <v-icon color="success" class="mr-2">mdi-arrow-up</v-icon>
                        <div>
                            <p class="text-grey text-caption mb-1">Tổng Thu Nhập</p>
                            <p class="text-h6 font-weight-bold text-success">
                                {{ formatCurrency(statistics.totalIncome, user?.preferences?.currency || 'VND') }}
                            </p>
                        </div>
                    </div>
                </v-col>
                <v-col cols="12" md="6" class="mb-4">
                    <div class="stat-item">
                        <v-icon color="error" class="mr-2">mdi-arrow-down</v-icon>
                        <div>
                            <p class="text-grey text-caption mb-1">Tổng Chi Tiêu</p>
                            <p class="text-h6 font-weight-bold text-error">
                                {{ formatCurrency(statistics.totalExpense, user?.preferences?.currency || 'VND') }}
                            </p>
                        </div>
                    </div>
                </v-col>
                <v-col cols="12" md="6" class="mb-4">
                    <div class="stat-item">
                        <v-icon color="primary" class="mr-2">mdi-cash-multiple</v-icon>
                        <div>
                            <p class="text-grey text-caption mb-1">Số Giao Dịch</p>
                            <p class="text-h6 font-weight-bold">
                                {{ statistics.totalTransactions.toLocaleString() }}
                            </p>
                        </div>
                    </div>
                </v-col>
                <v-col cols="12" md="6" class="mb-4">
                    <div class="stat-item">
                        <v-icon color="info" class="mr-2">mdi-tag-multiple</v-icon>
                        <div>
                            <p class="text-grey text-caption mb-1">Danh Mục Đã Dùng</p>
                            <p class="text-h6 font-weight-bold">
                                {{ statistics.categoriesUsed }}
                            </p>
                        </div>
                    </div>
                </v-col>
            </v-row>
            <v-divider class="my-3"></v-divider>
            <div v-if="statistics.netAmount !== undefined" class="text-center">
                <p class="text-grey text-caption mb-1">Số Dư</p>
                <p class="text-h5 font-weight-bold" :class="statistics.netAmount >= 0 ? 'text-success' : 'text-error'">
                    {{ formatCurrency(statistics.netAmount, user?.preferences?.currency || 'VND') }}
                </p>
            </div>
        </div>
        <div v-else class="text-center py-8">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
            <p class="text-grey mt-4">Đang tải thống kê...</p>
        </div>
    </v-card>
</template>

<script setup>
import { formatCurrency } from '@/utils/formatters'

const props = defineProps({
    user: {
        type: Object,
        default: null
    },
    statistics: {
        type: Object,
        default: null
    }
})
</script>

<style scoped>
.stat-item {
    display: flex;
    align-items: center;
    padding: 12px;
    background-color: rgba(0, 0, 0, 0.02);
    border-radius: 8px;
}
</style>