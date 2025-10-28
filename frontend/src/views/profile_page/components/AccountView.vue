<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <h3>Accounts</h3>
        <v-divider class="my-2"></v-divider>
        <v-list>
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