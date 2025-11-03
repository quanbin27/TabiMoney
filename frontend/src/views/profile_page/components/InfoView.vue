<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <v-row>
            <v-col cols="12" md="4" class="d-flex justify-center">
                <v-avatar size="100">
                    <v-img v-if="props.user?.avatar" :src="props.user.avatar" alt="Avatar"></v-img>
                    <v-icon v-else size="60" color="grey">mdi-account</v-icon>
                </v-avatar>
            </v-col>
            <v-col cols="12" md="8">
                <v-row>
                    <v-col cols="8" md="9">
                        <h3>{{ props.user?.name || 'User' }}</h3>
                        <p class="text-body-2">{{ props.user?.email || '' }}</p>
                        <p v-if="props.user?.phone" class="text-body-2">{{ props.user.phone }}</p>
                        <p v-if="props.user?.username" class="text-body-2 text-grey">@{{ props.user.username }}</p>
                    </v-col>
                    <v-col cols="4" md="3" class="d-flex justify-end align-center">
                        <v-chip v-if="props.user?.is_verified" color="success" size="small" prepend-icon="mdi-check-circle">
                            Verified
                        </v-chip>
                    </v-col>
                </v-row>
            </v-col>
        </v-row>
        <v-divider class="my-2"></v-divider>
        <div v-if="props.user?.phone" class="mb-2">
            <p><strong>Số điện thoại:</strong> {{ props.user.phone }}</p>
        </div>
        <div v-if="props.user?.created_at" class="mb-2">
            <p><strong>Thành viên từ:</strong> {{ formatDate(props.user.created_at) }}</p>
        </div>
        <div v-if="props.user?.last_login_at">
            <p><strong>Đăng nhập lần cuối:</strong> {{ formatDate(props.user.last_login_at) }}</p>
        </div>
        <div v-if="!props.user?.phone && !props.user?.created_at">
            <p class="text-grey">Thông tin bổ sung chưa được cập nhật</p>
        </div>
    </v-card>
</template>

<script setup>
import { formatVietnamDate } from '@/utils/formatters'

const props = defineProps(['user'])

const formatDate = (dateString) => {
    if (!dateString) return ''
    return formatVietnamDate(dateString, false)
}
</script>

<style></style>