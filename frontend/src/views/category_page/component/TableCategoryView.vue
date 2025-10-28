<template>
    <v-card v-if="props.categories">
        <v-data-table :headers="headers" :items="props.categories" :items-per-page="pagination.limit"
            :total-items="props.totalCategories" class="elevation-1">
            <template v-slot:top>
                <v-toolbar flat>
                    <v-toolbar-title>Manage Categories</v-toolbar-title>
                    <v-divider class="mx-4" inset vertical></v-divider>
                    <v-spacer></v-spacer>
                    <v-text-field v-model="search" append-icon="mdi-magnify" label="Search categories" single-line
                        hide-details @input="filterCategories"></v-text-field>
                </v-toolbar>
            </template>

            <template v-slot:item.is_system="{ item }">
                <v-chip :color="item.is_system ? 'success' : 'primary'" size="small">
                    {{ item.is_system ? 'System' : 'Custom' }}
                </v-chip>
            </template>

            <template v-slot:item.actions="{ item }">
                <v-btn v-if="!item.is_system" icon="mdi-pencil" size="small" variant="text"
                    @click="onEdit(item)"></v-btn>
                <v-btn v-if="!item.is_system" icon="mdi-delete" size="small" variant="text" color="error"
                    @click="props.onDelete(item)"></v-btn>
                <span v-if="item.is_system" class="text-caption text-grey">System</span>
            </template>

            <template v-slot:no-data>
                <v-alert :value="true" color="info" icon="mdi-information">
                    No categories found. Create your first category!
                </v-alert>
            </template>
        </v-data-table>
    </v-card>
</template>
<script setup>
import { ref } from 'vue'

const props = defineProps({
    categories: {
        type: Array,
        require: true,
    },
    totalCategories: {
        type: Number,
        require: true
    },
    onDelete: {
        type: Function,
        require: true
    }
})

const emit = defineEmits(['on-edit'])

const pagination = ref({
    page: 1,
    limit: 20,
})

const search = ref('')

const headers = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'English Name', key: 'name_en', sortable: true },
    { title: 'Description', key: 'description', sortable: false },
    { title: 'Type', key: 'is_system', sortable: true },
    { title: 'Actions', key: 'actions', sortable: false },
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
</script>