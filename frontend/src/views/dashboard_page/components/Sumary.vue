<template>
    <v-row class="mb-6" v-if="analytics != null">
        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" color="success" dark>
                <v-card-title class="text-h6">Tổng thu nhập</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ formatCurrency(analytics?.total_income || 0) }}
                    </div>
                    <div class="text-caption opacity-90">
                        Tháng này
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" color="error" dark>
                <v-card-title class="text-h6">Tổng chi tiêu</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ formatCurrency(analytics?.total_expense || 0) }}
                    </div>
                    <div class="text-caption opacity-90">
                        Tháng này
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" :color="netAmountColor" dark>
                <v-card-title class="text-h6">Số dư ròng</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ formatCurrency(analytics?.net_amount || 0) }}
                    </div>
                    <div class="text-caption opacity-90">
                        {{ netAmountLabel }}
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" color="info" dark>
                <v-card-title class="text-h6">Số giao dịch</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ analytics?.transaction_count || 0 }}
                    </div>
                    <div class="text-caption opacity-90">
                        Tháng này
                    </div>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import { formatCurrency } from "@/utils/formatters"
import { computed } from "vue"

const props = defineProps({
    analytics: {
        type: Object,
        required: false,
        default: null
    }
})

const netAmountColor = computed(() => {
    if (!props.analytics.value) return 'info'
    return props.analytics.value.net_amount >= 0 ? 'success' : 'error'
})

const netAmountLabel = computed(() => {
    if (!props.analytics.value) return ''
    return props.analytics.value.net_amount >= 0 ? 'Đang tiết kiệm' : 'Vượt chi'
})

</script>

<style>
.v-card {
    height: 100%;
}
</style>