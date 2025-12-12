<template>
    <v-dialog v-model="props.modelValue" max-width="600">
        <v-card>
            <v-card-title>Chi tiết mục tiêu</v-card-title>
            <v-card-text>
                <v-list density="compact">
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Tên mục tiêu</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.title }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Mô tả</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.description || '—' }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Số tiền mục tiêu</v-list-item-title>
                        <v-list-item-subtitle>
                            {{ formatCurrency(props.detailsGoal?.target_amount || 0) }}
                        </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Số tiền hiện tại</v-list-item-title>
                        <v-list-item-subtitle>
                            {{ formatCurrency(props.detailsGoal?.current_amount || 0) }}
                        </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Ngày hoàn thành dự kiến</v-list-item-title>
                        <v-list-item-subtitle>{{ formatDate(props.detailsGoal?.target_date) }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Loại</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.goal_type }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Mức độ ưu tiên</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.priority }}</v-list-item-subtitle>
                    </v-list-item>
                </v-list>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn variant="text" @click="closeDialog">Đóng</v-btn>
                <v-btn color="primary" variant="text" @click="goEdit(detailsGoal)">Sửa</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>
<script setup>
import { formatCurrency, formatDate } from '@/utils/formatters';

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    detailsGoal: Object,
});

const emit = defineEmits(["update:modelValue", "open-edit"]);

const closeDialog = () => {
    emit('update:modelValue', false)
}

const goEdit = (goal) => {
    emit('open-edit', goal)
}
</script>