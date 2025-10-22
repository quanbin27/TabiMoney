<template>
    <!-- Filter type -->
    <v-row class="mb-2 ">
        <v-col cols="12" md="3">
            <v-select v-model="filterType" :items="timeFilters" item-title="label" item-value="value"
                hide-details="true" label="Loại thời gian" variant="outlined" />
        </v-col>

        <!-- Theo ngày -->
        <v-col cols="12" md="8" v-if="filterType === 'day'" class="d-flex ga-6 ">
            <datePick v-model="dateRange.start" label="Từ ngày" />
            <datePick v-model="dateRange.end" label="Đến ngày" />
        </v-col>

        <!-- Theo tuần -->
        <v-col cols="12" md="8" v-if="filterType === 'week'" class="d-flex ga-6 ">
            <v-select v-model="selectedYear" :items="years" label="Chọn năm" variant="outlined" hide-details="true" />
            <v-select v-model="selectedMonth" :items="months" item-title="label" item-value="value" label="Chọn tháng"
                variant="outlined" hide-details="true" />
        </v-col>

        <!-- Theo tháng -->
        <v-col cols="12" md="8" v-if="filterType === 'month'">
            <v-select v-model="selectedYear" :items="years" label="Chọn năm" variant="outlined" hide-details="true" />
        </v-col>
        <!-- Theo năm -->
        <v-col cols="12" md="8" v-if="filterType === 'year'" class="d-flex ga-6 ">
            <v-select v-model="yearRange.start" :items="years" label="Từ năm" variant="outlined" hide-details="true" />
            <v-select v-model="yearRange.end" :items="years" label="Đến năm" variant="outlined" hide-details="true" />
        </v-col>

        <v-col cols="12" md="1" class="d-flex justify-center align-center ">
            <v-btn @click="handleFilter" color="primary" block>Lọc</v-btn>
        </v-col>
    </v-row>
</template>

<script setup>
import datePick from "@/components/date-pick.vue";
import { ref } from "vue";

const timeFilters = [
    { label: "Theo ngày", value: "day" },
    { label: "Theo tuần", value: "week" },
    { label: "Theo tháng", value: "month" },
    { label: "Theo năm", value: "year" },
];

const filterType = ref("day");

// Ngày
const dateRange = ref({ start: null, end: null });

// Tuần + tháng
const selectedYear = ref(new Date().getFullYear());
const selectedMonth = ref(new Date().getMonth() + 1);

// Năm (từ - đến)
const yearRange = ref({ start: 2020, end: new Date().getFullYear() });

const years = Array.from({ length: 15 }, (_, i) => 2011 + i); // 2015 → 2029
const months = [
    { label: "Tháng 1", value: 1 },
    { label: "Tháng 2", value: 2 },
    { label: "Tháng 3", value: 3 },
    { label: "Tháng 4", value: 4 },
    { label: "Tháng 5", value: 5 },
    { label: "Tháng 6", value: 6 },
    { label: "Tháng 7", value: 7 },
    { label: "Tháng 8", value: 8 },
    { label: "Tháng 9", value: 9 },
    { label: "Tháng 10", value: 10 },
    { label: "Tháng 11", value: 11 },
    { label: "Tháng 12", value: 12 },
];


const handleFilter = () => {
    console.log("filterType", filterType.value)
    console.log("dateRange", dateRange.value)
    console.log("selectedYear", selectedYear.value)
    console.log("selectedMonth", selectedMonth.value)
    console.log("yearRange", yearRange.value)
}
</script>
