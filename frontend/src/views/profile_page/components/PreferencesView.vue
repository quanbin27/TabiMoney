<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <v-row>
            <v-col cols="9">
                <h3>Preferences</h3>
            </v-col>
            <v-col cols="3" class="d-flex justify-end">
                <v-btn icon="mdi-cog"></v-btn>
            </v-col>
        </v-row>
        <v-divider class="my-2"></v-divider>
        <p><strong>Currency:</strong> {{ user.preferences.currency }}</p>
        <p><strong>Language:</strong> {{ user.preferences.language }}</p>
        <p><strong>Theme:</strong> {{ user.preferences.theme }}</p>
        <p><strong>Notifications:</strong> {{ user.preferences.notification ? 'Enabled' : 'Disabled' }}</p>
    </v-card>
</template>

<script setup>
const props = defineProps(['user'])
</script>

<style></style>