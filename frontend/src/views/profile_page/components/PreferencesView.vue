<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <v-row>
            <v-col cols="12">
                <h3>Tùy chọn</h3>
            </v-col>
        </v-row>
        <v-divider class="my-2"></v-divider>
        <div v-if="user?.preferences">
            <p class="mb-2"><strong>Tiền tệ:</strong> {{ user.preferences.currency || 'VND' }}</p>
            <p class="mb-2"><strong>Ngôn ngữ:</strong> {{ getLanguageLabel(user.preferences.language) }}</p>
            <p class="mb-2"><strong>Múi giờ:</strong> {{ user.preferences.timezone || 'Asia/Ho_Chi_Minh' }}</p>
            <p class="mb-2" v-if="user.preferences.monthly_income">
                <strong>Thu nhập hàng tháng:</strong> 
                {{ formatCurrency(user.preferences.monthly_income, user.preferences.currency) }}
            </p>
            <p class="mb-2">
                <strong>Thông báo:</strong> 
                <v-chip 
                    :color="user.preferences.notification ? 'success' : 'grey'" 
                    size="small"
                    class="ml-2"
                >
                    {{ user.preferences.notification ? 'Bật' : 'Tắt' }}
                </v-chip>
            </p>
        </div>
        <div v-else class="text-grey">
            <p>Chưa có tùy chọn được thiết lập</p>
        </div>
    </v-card>
</template>

<script setup>
import { formatCurrency } from '@/utils/formatters'

const props = defineProps(['user'])

const getLanguageLabel = (lang) => {
    const languages = {
        'vi': 'Tiếng Việt',
        'en': 'English',
        'zh': '中文'
    }
    return languages[lang] || lang || 'Tiếng Việt'
}
</script>

<style></style>