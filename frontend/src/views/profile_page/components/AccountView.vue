<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <h3>Tài khoản</h3>
        <v-divider class="my-2"></v-divider>
        <div v-if="!user?.accounts || user.accounts.length === 0" class="text-center py-8">
            <v-icon size="64" color="grey-lighten-1">mdi-wallet-outline</v-icon>
            <p class="text-grey mt-4">Chức năng quản lý tài khoản sẽ được cập nhật trong phiên bản tương lai</p>
        </div>
        <v-list v-else>
            <v-list-item v-for="acc in user.accounts" :key="acc.id">
                <v-list-item-title>
                    {{ acc.name }} ({{ acc.type }})
                </v-list-item-title>
                <v-list-item-subtitle>
                    {{ acc.balance.toLocaleString() }} VND
                </v-list-item-subtitle>
            </v-list-item>
        </v-list>
    </v-card>
</template>

<script setup>
const props = defineProps(['user'])
</script>

<style></style>