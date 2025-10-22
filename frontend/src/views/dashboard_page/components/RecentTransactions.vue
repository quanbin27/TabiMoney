<template>
    <v-row class="mt-6" v-if="recentTransactions != null">
        <v-col cols="12">
            <v-card>
                <v-card-title>
                    <v-icon left>mdi-receipt</v-icon>
                    Recent Transactions
                </v-card-title>
                <v-card-text>
                    <v-data-table :headers="transactionHeaders" :items="recentTransactions" class=" elevation-0"
                        hide-default-footer>
                        <template v-slot:item.amount="{ item }">
                            <span :class="getAmountClass(item.transaction_type)">
                                {{ formatCurrency(item.amount) }}
                            </span>
                        </template>

                        <template v-slot:item.transaction_date="{ item }">
                            {{ formatDate(item.transaction_date) }}
                        </template>

                        <template v-slot:item.category="{ item }">
                            <v-chip :color="item.category?.color || 'primary'" size="small" variant="tonal">
                                <v-icon start size="small">{{ item.category?.icon || 'mdi-tag' }}</v-icon>
                                {{ item.category?.name || 'Unknown' }}
                            </v-chip>
                        </template>

                        <template v-slot:item.actions="{ item }">
                            <v-btn icon="mdi-pencil" size="small" variant="text" @click="editTransaction(item)" />
                            <v-btn icon="mdi-delete" size="small" variant="text" color="error"
                                @click="deleteTransaction(item.id)" />
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import { transactionAPI } from '../../../services/api'
import { useAppStore } from '../../../stores/app'
import { formatCurrency, formatDate } from '../../../utils/formatters'
const appStore = useAppStore()

const props = defineProps({
    recentTransactions: {
        type: Array,
        required: true
    },
    loadRecentTransactions: {
        type: Function,
        required: true
    }
})

const emit = defineEmits(['openEdit'])

const transactionHeaders = [
    { title: 'Date', key: 'transaction_date', sortable: true },
    { title: 'Description', key: 'description', sortable: true },
    { title: 'Category', key: 'category', sortable: false },
    { title: 'Amount', key: 'amount', sortable: true },
    { title: 'Actions', key: 'actions', sortable: false },
]

const getAmountClass = (type) => {
    return type === 'income' ? 'text-success' : 'text-error'
}

const editTransaction = (item) => {
    emit('openEdit', item)
}

const deleteTransaction = async (id) => {
    const confirmed = await appStore.confirm({
        title: 'Xác nhận xóa',
        text: 'Bạn có chắc muốn xóa mục này?',
        icon: 'warning',
        confirmButtonText: 'Xóa',
        cancelButtonText: 'Hủy',
    })

    if (confirmed) {
        try {
            await transactionAPI.deleteTransaction(id)
            appStore.showSuccess('Transaction deleted successfully')
            props.loadRecentTransactions()
        } catch (error) {
            appStore.showError('Failed to delete transaction')
        }
    }
}

</script>

<style></style>