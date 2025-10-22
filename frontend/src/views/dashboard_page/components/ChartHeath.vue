<template>
    <v-row>
        <!-- Spending by Category Chart -->
        <v-col cols="12" md="8" v-if="analytics">
            <v-card>
                <v-card-title>
                    <v-icon left>mdi-chart-pie</v-icon>
                    Spending by Category
                </v-card-title>
                <v-card-text>
                    <div v-if="categoryChartData" class="text-center pa-8">
                        <DoughnutChart :data="categoryChartData" :options="chartOptions" />
                    </div>
                    <div v-else class="text-center pa-8 text-medium-emphasis">
                        No data available
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <!-- Financial Health -->
        <v-col cols="12" md="4">
            <v-card>
                <v-card-title>
                    <v-icon left>mdi-heart-pulse</v-icon>
                    Financial Health
                </v-card-title>
                <v-card-text class="">
                    <div v-if="analytics?.financial_health" class="text-center ">
                        <v-progress-circular :model-value="analytics.financial_health.score" :color="healthColor"
                            size="120" width="12" class="mb-4">
                            <span class="text-h4 font-weight-bold">
                                {{ Math.round(analytics.financial_health.score) }}
                            </span>
                        </v-progress-circular>

                        <div class="text-h6 font-weight-bold mb-2">
                            {{ analytics.financial_health.level.toUpperCase() }}
                        </div>

                        <div class="text-body-2 text-medium-emphasis mb-4">
                            Savings Rate: {{ analytics.financial_health.savings_rate.toFixed(1) }}%
                        </div>

                        <v-list density="compact">
                            <v-list-item v-for="recommendation in analytics.financial_health.recommendations"
                                :key="recommendation" class="px-0">
                                <template v-slot:prepend>
                                    <v-icon size="small" color="primary">mdi-lightbulb</v-icon>
                                </template>
                                <v-list-item-title class="text-caption">
                                    {{ recommendation }}
                                </v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </div>
                    <div v-else class="text-center pa-4 text-medium-emphasis">
                        No health data available
                    </div>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import DoughnutChart from "@/components/DoughnutChart.vue"
import { computed } from "vue"

const props = defineProps({
    analytics: {
        type: Object,
        required: false,
        default: null
    },
})

const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { position: 'bottom' },
    },
}

const categoryChartData = computed(() => {
    if (!props.analytics?.category_breakdown) return null

    return {
        labels: props.analytics?.category_breakdown.map(c => c.category_name),
        datasets: [{
            data: props.analytics?.category_breakdown.map(c => c.amount),
            backgroundColor: [
                '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF',
                '#FF9F40', '#FF6384', '#C9CBCF', '#4BC0C0', '#FF6384'
            ],
            borderWidth: 2,
            borderColor: '#fff'
        }]
    }
})

const healthColor = computed(() => {
    if (!props.analytics?.financial_health) return 'info'
    const score = props.analytics?.financial_health.score
    if (score >= 80) return 'success'
    if (score >= 60) return 'warning'
    return 'error'
})

</script>

<style>
.v-card {
    height: 100%;
}
</style>