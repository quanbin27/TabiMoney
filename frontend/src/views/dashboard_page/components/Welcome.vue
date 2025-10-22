<template>
    <v-row class="mb-6">
        <v-col cols="12">
            <v-card class="pa-6" color="primary" dark>
                <v-row align="center">
                    <v-col cols="12" md="8">
                        <h1 class="text-h4 font-weight-bold mb-2">
                            Welcome back, {{ authStore.user?.first_name || authStore.user?.username }}!
                        </h1>
                        <p class="text-h6 opacity-90">
                            Here's your financial overview for {{ currentMonth }}
                        </p>
                    </v-col>
                    <v-col cols="12" md="4" class="text-right">
                        <v-btn color="white" variant="outlined" size="large" @click="openAddTransactionDialog">
                            <v-icon left>mdi-plus</v-icon>
                            Add Transaction
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import { computed } from "vue"
import { useAuthStore } from '../../../stores/auth'
const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
})
const emit = defineEmits(['update:modelValue'])
const authStore = useAuthStore()

const openAddTransactionDialog = () => {
    emit('update:modelValue', true)
}

const currentMonth = computed(() => {
    return new Date().toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
})
</script>

<style></style>