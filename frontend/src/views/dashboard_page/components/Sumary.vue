<template>
    <v-row class="mb-6" v-if="analytics != null">
        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" color="success" dark>
                <v-card-title class="text-h6">Total Income</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ formatCurrency(analytics?.total_income || 0) }}
                    </div>
                    <div class="text-caption opacity-90">
                        This month
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" color="error" dark>
                <v-card-title class="text-h6">Total Expenses</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ formatCurrency(analytics?.total_expense || 0) }}
                    </div>
                    <div class="text-caption opacity-90">
                        This month
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="3">
            <v-card class="pa-4" :color="netAmountColor" dark>
                <v-card-title class="text-h6">Net Amount</v-card-title>
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
                <v-card-title class="text-h6">Transactions</v-card-title>
                <v-card-text>
                    <div class="text-h4 font-weight-bold">
                        {{ analytics?.transaction_count || 0 }}
                    </div>
                    <div class="text-caption opacity-90">
                        This month
                    </div>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import { computed } from "vue"
import { formatCurrency } from "../../../utils/formatters"

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
    return props.analytics.value.net_amount >= 0 ? 'Saved' : 'Overspent'
})

</script>

<style>
.v-card {
    height: 100%;
}
</style>