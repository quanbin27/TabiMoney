<template>
    <v-row v-if="props.goals.length > 0">
        <v-col v-for="goal in props.goals" :key="goal.id" cols="12" md="6" lg="4">
            <v-card class="goal-card" :class="getGoalStatusClass(goal)" elevation="2" hover>
                <v-card-title class="d-flex align-center justify-space-between">
                    <div class="d-flex align-center">
                        <v-icon :color="getGoalStatusColor(goal)" class="mr-2" size="large">
                            {{ getGoalIcon(goal) }}
                        </v-icon>
                        <div>
                            <h3 class="text-h6">{{ goal.title }}</h3>
                            <p class="text-caption mb-0">{{ goal.description }}</p>
                        </div>
                    </div>
                    <v-menu>
                        <template v-slot:activator="{ props }">
                            <v-btn icon="mdi-dots-vertical" variant="text" v-bind="props"></v-btn>
                        </template>
                        <v-list>
                            <v-list-item @click="openEditDialog(goal)">
                                <template v-slot:prepend>
                                    <v-icon>mdi-pencil</v-icon>
                                </template>
                                <v-list-item-title>Sửa</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="deleteGoal(goal.id)">
                                <template v-slot:prepend>
                                    <v-icon>mdi-delete</v-icon>
                                </template>
                                <v-list-item-title>Xóa</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-card-title>

                <v-card-text>
                    <!-- Progress Bar -->
                    <div class="mb-4">
                        <div class="d-flex justify-space-between mb-2">
                            <span class="text-body-2">Tiến độ</span>
                            <span class="text-body-2 font-weight-bold">
                                {{ getProgressPercentage(goal) }}%
                            </span>
                        </div>
                        <v-progress-linear :model-value="getProgressPercentage(goal)" :color="getGoalStatusColor(goal)"
                            height="8" rounded></v-progress-linear>
                    </div>

                    <!-- Goal Details -->
                    <v-row class="text-center">
                        <v-col cols="6">
                            <div class="text-h6 text-primary">
                                {{ formatCurrency(goal.current_amount) }}
                            </div>
                            <div class="text-caption text-grey">Hiện tại</div>
                        </v-col>
                        <v-col cols="6">
                            <div class="text-h6 text-grey">
                                {{ formatCurrency(goal.target_amount) }}
                            </div>
                            <div class="text-caption text-grey">Mục tiêu</div>
                        </v-col>
                    </v-row>

                    <!-- Goal Status -->
                    <v-chip :color="getGoalStatusColor(goal)" size="small" class="mt-3">
                        {{ getGoalStatusText(goal) }}
                    </v-chip>

                    <!-- Deadline -->
                    <div class="mt-3 text-caption text-grey">
                        <v-icon size="small" class="mr-1">mdi-calendar</v-icon>
                        Hạn chót: {{ formatDate(goal.target_date) }}
                    </div>
                </v-card-text>

                <v-card-actions>
                    <v-btn color="primary" variant="text" @click="openDetails(goal)">
                        Xem chi tiết
                    </v-btn>
                    <v-spacer></v-spacer>
                    <v-btn color="success" variant="text" @click="addContribution(goal)">
                        Thêm đóng góp
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import { formatCurrency, formatDate } from '@/utils/formatters';

const props = defineProps({
    goals: Array,
});

const emit = defineEmits(["open-details", "open-contribution", "open-edit", "on-delete"]);

const openDetails = (goal) => {
    emit('open-details', goal)
}

const addContribution = (goal) => {
    emit('open-contribution', goal)
}

const openEditDialog = (goal) => {
    emit('open-edit', goal)
}

const deleteGoal = (goalID) => {
    emit('on-delete', goalID)
}

const getGoalIcon = (goal) => {
    const status = getGoalStatus(goal)
    const iconMap = {
        completed: 'mdi-check-circle',
        overdue: 'mdi-alert-circle',
        on_track: 'mdi-trending-up',
        in_progress: 'mdi-progress-clock',
        not_started: 'mdi-target'
    }
    return iconMap[status] || 'mdi-target'
}

const getProgressPercentage = (goal) => {
    if (goal.target_amount === 0) return 0
    return Math.min((goal.current_amount / goal.target_amount) * 100, 100).toFixed(2)
}

const getGoalStatus = (goal) => {
    const progress = Number(getProgressPercentage(goal))
    const now = new Date()
    const targetDate = goal.target_date ? new Date(goal.target_date) : null

    // Nếu backend đã đánh dấu đạt mục tiêu thì ưu tiên trạng thái này
    if (goal.is_achieved) return 'completed'

    // Nếu chưa có ngày mục tiêu hoặc ngày mục tiêu nằm trong quá khứ mà chưa đạt -> quá hạn
    if (targetDate && targetDate < now && progress < 100) return 'overdue'

    // Nếu chưa có đóng góp nào
    if (!goal.current_amount || goal.current_amount === 0) return 'not_started'

    // Đã có tiền hiện tại nhưng chưa đạt
    if (progress < 50) return 'in_progress'
    if (progress < 100) return 'on_track'

    return 'in_progress'
}

const getGoalStatusText = (goal) => {
    const status = getGoalStatus(goal)
    const statusMap = {
        completed: 'Hoàn thành',
        overdue: 'Quá hạn',
        on_track: 'Đúng tiến độ',
        in_progress: 'Đang thực hiện',
        not_started: 'Chưa bắt đầu'
    }
    return statusMap[status] || 'Không xác định'
}

const getGoalStatusColor = (goal) => {
    const status = getGoalStatus(goal)
    const colorMap = {
        completed: 'success',
        overdue: 'error',
        on_track: 'primary',
        in_progress: 'warning',
        not_started: 'grey'
    }
    return colorMap[status] || 'grey'
}

const getGoalStatusClass = (goal) => {
    const status = getGoalStatus(goal)
    return `goal-${status}`
}
</script>