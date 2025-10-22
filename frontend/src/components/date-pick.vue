<template>
    <v-menu v-model="menu" :close-on-content-click="false" transition="scale-transition" offset-y max-width="290px"
        min-width="auto">
        <template #activator="{ props }">
            <v-text-field v-model="formatValue" :label="label" readonly v-bind="props" variant="outlined"
                density="compact" hide-details="true" :rules="[
                    v => !!v || 'Trường này là bắt buộc'
                ]" />
        </template>

        <v-date-picker v-model="localValue" @update:model-value="menu = false" />
    </v-menu>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { formatVietnamDate } from "../utils/formatters";

const props = defineProps({
    modelValue: String,
    label: String,
});

const emit = defineEmits(["update:modelValue"]);

const menu = ref(false);

const localValue = ref(props.modelValue ? new Date(props.modelValue) : null);

const formatValue = computed(() => {
    return localValue.value ? formatVietnamDate(localValue.value) : "";
});

watch(localValue, (val) => {
    emit("update:modelValue", val ? formatVietnamDate(val) : "");
});

watch(
    () => props.modelValue,
    (val) => {
        if (val) {
            // parse từ chuỗi dd/MM/yyyy về Date object
            const [day, month, year] = val.split("/");
            localValue.value = new Date(year, month - 1, day);
        } else {
            localValue.value = null;
        }
    }
);
</script>
