<template>
    <v-dialog v-model="props.modelValue" max-width="600">
        <v-card>
            <v-card-title>Goal Details</v-card-title>
            <v-card-text>
                <v-list density="compact">
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Title</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.title }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Description</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.description || '—' }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Target Amount</v-list-item-title>
                        <v-list-item-subtitle>
                            {{ formatCurrency(props.detailsGoal?.target_amount || 0) }}
                        </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Current Amount</v-list-item-title>
                        <v-list-item-subtitle>
                            {{ formatCurrency(props.detailsGoal?.current_amount || 0) }}
                        </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Target Date</v-list-item-title>
                        <v-list-item-subtitle>{{ formatDate(props.detailsGoal?.target_date) }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Type</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.goal_type }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title class="font-weight-bold">Priority</v-list-item-title>
                        <v-list-item-subtitle>{{ props.detailsGoal?.priority }}</v-list-item-subtitle>
                    </v-list-item>
                </v-list>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn variant="text" @click="closeDialog">Close</v-btn>
                <v-btn color="primary" variant="text" @click="goEdit(detailsGoal)">Edit</v-btn>
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