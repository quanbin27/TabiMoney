<template>
    <v-card class="pa-4 border-md rounded-xl w-100 h-100">
        <v-row>
            <v-col cols="12" md="4" class="d-flex justify-center">
                <v-avatar size="100">
                    <v-img :src="props.user.avatar" alt="Avatar"></v-img>
                </v-avatar>
            </v-col>
            <v-col cols="12" md="8">
                <v-row>
                    <v-col cols="8" md="9">
                        <h3>{{ props.user.name }}</h3>
                        <p class="text-body-2">{{ props.user.email }}</p>
                        <p class="text-body-2">{{ props.user.phone }}</p>
                    </v-col>
                    <v-col cols="4" md="3" class="d-flex justify-end align-center">
                        <v-btn icon="mdi-cog"></v-btn>
                    </v-col>
                </v-row>

            </v-col>
        </v-row>
        <v-divider class="my-2"></v-divider>
        <p><strong>Date of Birth:</strong> {{ props.user.dob }}</p>
        <p><strong>Gender:</strong> {{ props.user.gender }}</p>
        <p><strong>Address:</strong> {{ props.user.address }}</p>
    </v-card>
</template>

<script setup>
const props = defineProps(['user'])
</script>

<style></style>