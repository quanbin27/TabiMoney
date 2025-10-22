<template>
    <v-container>
        <v-text-field class="mb-4" v-model="search" label="Search" prepend-inner-icon="mdi-magnify" variant="outlined"
            hide-details single-line></v-text-field>
        <v-data-table :headers="headers" :items="props.items" :loading="props.loading" :items-per-page="10"
            :search="search" loading-text="Loading transactions...">
            <template #item.amount="{ item }">
                {{ formatCurrency(item.amount) }}
            </template>
            <template #item.transaction_date="{ item }">
                {{ new Date(item.transaction_date).toLocaleDateString() }}
            </template>

            <template v-slot:item.actions="{ item }">
                <v-btn icon="mdi-pencil" size="small" variant="text" @click="editTransaction(item)" />
                <v-btn icon="mdi-delete" size="small" variant="text" color="error"
                    @click="deleteTransaction(item.id)" />
            </template>
        </v-data-table>
    </v-container>
</template>

<script setup>
import { transactionAPI } from '@/services/api'
import { useAppStore } from '@/stores/app'
import { formatCurrency } from '@/utils/formatters'
import { ref } from 'vue'

const props = defineProps({
    items: {
        type: Array,
        required: true
    },
    loading: {
        type: Boolean,
        required: true
    },
    load: {
        type: Function,
        required: true
    }
})

const emit = defineEmits(['openEdit'])

const search = ref(null)

const appStore = useAppStore()

const headers = [
    { title: 'Date', key: 'transaction_date', sortable: true },
    { title: 'Type', key: 'transaction_type' },
    { title: 'Category', key: 'category.name' },
    { title: 'Description', key: 'description' },
    { title: 'Amount', key: 'amount' },
    { title: 'Actions', key: 'actions' },
]

const editTransaction = (transaction) => {
    emit('openEdit', transaction);
};

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
            props.load()
        } catch (error) {
            appStore.showError('Failed to delete transaction')
        }
    }
}

</script>