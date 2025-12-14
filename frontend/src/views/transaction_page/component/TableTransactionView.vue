<template>
    <v-container>
        <v-text-field class="mb-4" v-model="search" label="Tìm kiếm" prepend-inner-icon="mdi-magnify" variant="outlined"
            hide-details single-line></v-text-field>
        <v-data-table 
            :headers="headers" 
            :items="props.items" 
            :loading="props.loading" 
            :items-per-page="50"
            :items-per-page-options="[25, 50, 100, -1]"
            :search="search" 
            loading-text="Đang tải danh sách giao dịch..."
            class="elevation-1">
            <template #item.amount="{ item }">
                {{ formatCurrency(item.amount) }}
            </template>
            <template #item.transaction_date="{ item }">
                {{ new Date(item.transaction_date).toLocaleDateString() }}
            </template>
            <template #item.transaction_type="{ item }">
                <v-chip 
                    :color="item.transaction_type === 'income' ? 'success' : 'error'" 
                    size="small"
                    variant="tonal">
                    {{ item.transaction_type === 'income' ? 'Thu nhập' : 'Chi tiêu' }}
                </v-chip>
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
    { title: 'Ngày', key: 'transaction_date', sortable: true },
    { title: 'Loại', key: 'transaction_type' },
    { title: 'Danh mục', key: 'category.name' },
    { title: 'Mô tả', key: 'description' },
    { title: 'Số tiền', key: 'amount' },
    { title: 'Thao tác', key: 'actions' },
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
            appStore.showSuccess('Xoá giao dịch thành công')
            props.load()
        } catch (error) {
            appStore.showError('Không thể xoá giao dịch')
        }
    }
}

</script>