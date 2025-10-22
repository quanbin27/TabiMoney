<template>
  <v-app>
    <v-app-bar v-if="!isAuthPage" app color="primary" dark elevation="2">
      <v-app-bar-nav-icon @click="drawer = !drawer" />

      <v-toolbar-title>TabiMoney</v-toolbar-title>

      <v-spacer />

      <!-- Notifications -->
      <v-menu v-if="authStore.isAuthenticated" offset-y>
        <template v-slot:activator="{ props }">
          <v-btn icon v-bind="props" @click="loadNotifications">
            <v-badge :content="unreadCount" color="red" v-if="unreadCount > 0" overlap>
              <v-icon>mdi-bell</v-icon>
            </v-badge>
            <v-icon v-else>mdi-bell</v-icon>
          </v-btn>
        </template>
        <v-list style="min-width: 320px; max-width: 420px;">
          <v-list-subheader>Notifications</v-list-subheader>
          <v-list-item v-if="notifications.length === 0">
            <v-list-item-title>No notifications</v-list-item-title>
          </v-list-item>
          <v-list-item v-for="n in notifications" :key="n.id" :class="{ 'opacity-70': n.is_read }">
            <v-list-item-title class="d-flex align-center justify-space-between">
              <span>{{ n.title }}</span>
              <v-chip size="x-small" :color="chipColor(n.notification_type)">{{ n.notification_type }}</v-chip>
            </v-list-item-title>
            <v-list-item-subtitle>{{ n.message }}</v-list-item-subtitle>
            <template v-slot:append>
              <v-btn v-if="!n.is_read" size="x-small" variant="text" @click.stop="markRead(n.id)">Mark read</v-btn>
            </template>
          </v-list-item>
        </v-list>
      </v-menu>

      <!-- User Menu -->
      <v-menu v-if="authStore.isAuthenticated" offset-y>
        <template v-slot:activator="{ props }">
          <v-btn icon v-bind="props">
            <v-avatar size="32">
              <v-img v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url"
                :alt="authStore.user.username" />
              <v-icon v-else>mdi-account</v-icon>
            </v-avatar>
          </v-btn>
        </template>

        <v-list>
          <v-list-item>
            <v-list-item-title>{{ authStore.user?.username }}</v-list-item-title>
            <v-list-item-subtitle>{{ authStore.user?.email }}</v-list-item-subtitle>
          </v-list-item>
          <v-divider />
          <v-list-item @click="goToProfile">
            <v-list-item-title>
              <v-icon>mdi-account</v-icon>
              Profile</v-list-item-title>
          </v-list-item>
          <v-list-item @click="logout">
            <v-list-item-title>
              <v-icon>mdi-logout</v-icon>
              Logout</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>

    <!-- Navigation Drawer -->
    <v-navigation-drawer v-if="!isAuthPage && authStore.isAuthenticated" v-model="drawer" app temporary>
      <v-list>
        <v-list-item v-for="item in navigationItems" :key="item.title" :to="item.to" @click="drawer = false">
          <v-list-item-title>
            <v-icon>{{ item.icon }}</v-icon>
            {{ item.title }}
          </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- Main Content -->
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { notificationsAPI } from './services/api'
import { useAppStore } from './stores/app'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const drawer = ref(false)
const loading = computed(() => appStore.loading)
const notifications = ref([])
const unreadCount = ref(0)

// Navigation items
const navigationItems = [
  { title: 'Dashboard', icon: 'mdi-view-dashboard', to: '/dashboard' },
  { title: 'Transactions', icon: 'mdi-receipt', to: '/transactions' },
  { title: 'Analytics', icon: 'mdi-chart-line', to: '/analytics' },
  { title: 'Goals', icon: 'mdi-target', to: '/goals' },
  { title: 'Budgets', icon: 'mdi-wallet', to: '/budgets' },
  { title: 'Categories', icon: 'mdi-tag', to: '/categories' },
  { title: 'AI Assistant', icon: 'mdi-robot', to: '/ai-assistant' },
  { title: 'Settings', icon: 'mdi-cog', to: '/settings' },
]

// Check if current page is auth page
const isAuthPage = computed(() => {
  return route.path.startsWith('/auth')
})

// Snackbar

// Watch for route changes
watch(route, () => {
  drawer.value = false
})

// Methods
const goToProfile = () => {
  router.push('/profile')
  drawer.value = false
}

const logout = async () => {
  await authStore.logout()
  router.push('/auth/login')
}

const loadNotifications = async () => {
  try {
    const [allRes, unreadRes] = await Promise.all([
      notificationsAPI.list(false),
      notificationsAPI.list(true),
    ])
    notifications.value = allRes.data.data || []
    unreadCount.value = Array.isArray(unreadRes.data.data) ? unreadRes.data.data.length : 0
  } catch (e) {
    // silent fail in header
  }
}

const markRead = async (id) => {
  try {
    await notificationsAPI.markRead(id)
    // update local state
    notifications.value = notifications.value.map(n => n.id === id ? { ...n, is_read: true } : n)
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  } catch (e) { }
}
</script>

<style scoped>
.v-navigation-drawer {
  border-right: 1px solid rgba(0, 0, 0, 0.12);
}

.v-app-bar {
  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
}
</style>
