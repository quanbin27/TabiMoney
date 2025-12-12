<template>
    <v-dialog v-model="props.modelValue" max-width="400">
        <v-card>
            <v-card-title>Thêm đóng góp</v-card-title>
            <v-card-text>
                <v-form ref="contributionFormRef" validate-on="submit">
                    <v-text-field v-model="contributionForm.amount" label="Số tiền" type="number"
                        :rules="[props.rules.required, props.rules.positive]" prefix="₫" required></v-text-field>

                    <v-textarea v-model="contributionForm.note" label="Ghi chú (Tùy chọn)" rows="2"></v-textarea>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn @click="closeDialog">Hủy</v-btn>
                <v-btn color="primary" @click="saveContribution" :loading="savingContribution">
                    Thêm đóng góp
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>
<script setup>
import { goalAPI } from '@/services/api';
import { useAppStore } from '@/stores/app';
import { ref } from 'vue';

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    rules: {
        type: Object,
        default: null,
    },
    selectedGoal: {
        type: Object,
        default: null
    },
    loadGoals: {
        type: Function,
        required: true
    }
});

const emit = defineEmits(["update:modelValue"]);

const contributionFormRef = ref()
const contributionForm = ref({
    amount: null,
    note: ''
})
const contributionFormValid = ref(false)
const savingContribution = ref(false)

const closeDialog = () => {
    emit('update:modelValue', false)
}

const saveContribution = async () => {
    if (contributionFormRef.value) {
        const { valid } = await contributionFormRef.value.validate()
        if (!valid) return
    }
    savingContribution.value = true
    try {
        await goalAPI.addContribution(props.selectedGoal.id, {
            amount: Number(contributionForm.value.amount),
            note: contributionForm.value.note?.trim() || undefined,
        })
        await props.loadGoals()
        closeDialog()
        useAppStore().showSuccess('Thêm đóng góp thành công')
    } catch (error) {
        console.error('Failed to add contribution:', error)
        useAppStore().showError('Thêm đóng góp thất bại')
    } finally {
        contributionFormRef.value?.reset()
        savingContribution.value = false
    }
}
</script>