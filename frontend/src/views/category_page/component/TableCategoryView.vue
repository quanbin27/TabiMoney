<template>
    <v-card v-if="props.categories">
        <v-data-table :headers="headers" :items="props.categories" :items-per-page="pagination.limit"
            :total-items="props.totalCategories" class="elevation-1">
            <template v-slot:top>
                <v-toolbar flat>
                    <v-toolbar-title>Quản lý danh mục</v-toolbar-title>
                    <v-divider class="mx-4" inset vertical></v-divider>
                    <v-spacer></v-spacer>
                    <v-text-field v-model="search" append-icon="mdi-magnify" label="Tìm kiếm danh mục" single-line
                        hide-details @input="filterCategories"></v-text-field>
                </v-toolbar>
            </template>

            <template v-slot:item.is_system="{ item }">
                <v-chip :color="item.is_system ? 'success' : 'primary'" size="small">
                    {{ item.is_system ? 'Hệ thống' : 'Tự tạo' }}
                </v-chip>
            </template>

            <template v-slot:item.actions="{ item }">
                <v-btn v-if="!item.is_system" icon="mdi-pencil" size="small" variant="text"
                    @click="onEdit(item)"></v-btn>
                <v-btn v-if="!item.is_system" icon="mdi-delete" size="small" variant="text" color="error"
                    @click="deleteCategory(item)"></v-btn>
                <span v-if="item.is_system" class="text-caption text-grey">Hệ thống</span>
            </template>

            <template v-slot:no-data>
                <v-alert :value="true" color="info" icon="mdi-information">
                    Chưa có danh mục nào. Hãy tạo danh mục đầu tiên!
                </v-alert>
            </template>
        </v-data-table>
    </v-card>
</template>
<script setup>
import { ref } from 'vue'
import { useAppStore } from '../../../stores/app'
import { categoryAPI } from '../../../services/api'

const props = defineProps({
    categories: {
        type: Array,
        require: true,
    },
    totalCategories: {
        type: Number,
        require: true
    },
    load: {
        type: Function,
        require: true
    }
})

const app = useAppStore()

const emit = defineEmits(['on-edit'])

const pagination = ref({
    page: 1,
    limit: 20,
})

const search = ref('')

const headers = [
    { title: 'Tên', key: 'name', sortable: true },
    { title: 'Tên tiếng Anh', key: 'name_en', sortable: true },
    { title: 'Mô tả', key: 'description', sortable: false },
    { title: 'Loại', key: 'is_system', sortable: true },
    { title: 'Thao tác', key: 'actions', sortable: false },
]

const onEdit = (item) => {
    emit('on-edit', item)
}

function filterCategories() {
    if (!search.value.trim()) {
        filteredCategories.value = categories.value
        return
    }

    const query = search.value.toLowerCase()
    filteredCategories.value = categories.value.filter(cat =>
        cat.name?.toLowerCase().includes(query) ||
        cat.name_en?.toLowerCase().includes(query) ||
        cat.description?.toLowerCase().includes(query)
    )
}

async function deleteCategory(category) {
    if (category.is_system) {
        app.showWarning('Không thể xoá danh mục hệ thống')
        return
    }
    const confirmed = await app.confirm({
        title: 'Xác nhận xóa',
        text: 'Bạn có chắc muốn xóa mục này?',
        icon: 'warning',
        confirmButtonText: 'Xóa',
        cancelButtonText: 'Hủy',
    })

    if (confirmed) {
        try {
            await categoryAPI.deleteCategory(category.id)
            app.showSuccess('Đã xoá danh mục')
            await props.load()
        } catch (e) {
            app.showError(e?.message || 'Không thể xoá danh mục')
        }
    }
}
</script>