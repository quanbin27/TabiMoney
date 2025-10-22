import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    {
      path: '/auth',
      children: [
        {
          path: 'login',
          name: 'Login',
          component: () => import('../views/auth_page/LoginView.vue'),
          meta: { requiresGuest: true },
        },
        {
          path: 'register',
          name: 'Register',
          component: () => import('../views/auth_page/RegisterView.vue'),
          meta: { requiresGuest: true },
        },
      ],
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/dashboard_page/DashboardView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions',
      name: 'Transactions',
      component: () => import('../views/transaction_page/TransactionsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/analytics',
      name: 'Analytics',
      component: () => import('../views/analytics_page/AnalyticsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/goals',
      name: 'Goals',
      component: () => import('../views/goal_page/GoalsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/budgets',
      name: 'Budgets',
      component: () => import('../views/budget_page/BudgetsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories',
      name: 'Categories',
      component: () => import('../views/category_page/CategoriesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/ai-assistant',
      name: 'AIAssistant',
      component: () => import('../views/ai_page/AIAssistantView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('../views/profile_page/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../views/setting_page/SettingsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../views/notfound_page/NotFoundView.vue'),
    },
  ],
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (!authStore.initialized) {
    await authStore.initialize()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/auth/login')
    return
  }

  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  next()
})

export default router


