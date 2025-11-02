<template>
  <v-container>
    <v-row>
      <v-col cols="12" xl="4">
        <InfoView :user="user"></InfoView>
      </v-col>

      <v-col cols="12" xl="8">
        <AccountView :user="user"></AccountView>

      </v-col>
    </v-row>

    <v-row class="mt-4">
      <v-col cols="12" xl="4">
        <PreferencesView :user="user"></PreferencesView>

      </v-col>

      <v-col cols="12" xl="8">
        <StatsView :user="user"></StatsView>

      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { AccountView, InfoView, PreferencesView, StatsView } from './components'

const user = {
  "id": 1,
  "name": "Nguyen Van A",
  "email": "nguyenvana@example.com",
  "phone": "+84 912 345 678",
  "dob": "1995-06-15",
  "gender": "Male",
  "address": "123 Nguyễn Trãi, Hà Nội, Việt Nam",
  "avatar": "https://i.pravatar.cc/150?img=3",
  "accounts": [
    {
      "id": 1,
      "type": "Bank Account",
      "name": "Vietcombank",
      "balance": 15000000
    },
    {
      "id": 2,
      "type": "Cash",
      "name": "Ví tiền mặt",
      "balance": 5000000
    },
    {
      "id": 3,
      "type": "E-Wallet",
      "name": "Momo",
      "balance": 2000000
    },
    {
      "id": 4,
      "type": "Credit Card",
      "name": "Techcombank Visa",
      "balance": -3000000
    }
  ],
  "preferences": {
    "currency": "VND",
    "language": "vi",
    "theme": "light",
    "notification": true
  },
  "statistics": {
    "totalIncome": 450000000,
    "totalExpense": 320000000,
    "totalTransactions": 1240,
    "categoriesUsed": 12
  }
}

</script>