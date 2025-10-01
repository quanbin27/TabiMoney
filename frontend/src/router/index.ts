import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/auth',
      children: [
        {
          path: 'login',
          name: 'Login',
          component: () => import('../views/auth/LoginView.vue'),
          meta: { requiresGuest: true },
        },
        {
          path: 'register',
          name: 'Register',
          component: () => import('../views/auth/RegisterView.vue'),
          meta: { requiresGuest: true },
        },
      ],
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions',
      name: 'Transactions',
      component: () => import('../views/TransactionsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions/add',
      name: 'AddTransaction',
      component: () => import('../views/AddTransactionView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions/:id/edit',
      name: 'EditTransaction',
      component: () => import('../views/EditTransactionView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/analytics',
      name: 'Analytics',
      component: () => import('../views/AnalyticsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/goals',
      name: 'Goals',
      component: () => import('../views/GoalsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/goals/add',
      name: 'AddGoal',
      component: () => import('../views/AddGoalView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/goals/:id/edit',
      name: 'EditGoal',
      component: () => import('../views/EditGoalView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/budgets',
      name: 'Budgets',
      component: () => import('../views/BudgetsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/budgets/add',
      name: 'AddBudget',
      component: () => import('../views/AddBudgetView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/budgets/:id/edit',
      name: 'EditBudget',
      component: () => import('../views/EditBudgetView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories',
      name: 'Categories',
      component: () => import('../views/CategoriesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories/add',
      name: 'AddCategory',
      component: () => import('../views/AddCategoryView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories/:id/edit',
      name: 'EditCategory',
      component: () => import('../views/EditCategoryView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/ai-assistant',
      name: 'AIAssistant',
      component: () => import('../views/AIAssistantView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../views/SettingsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../views/NotFoundView.vue'),
    },
  ],
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Initialize auth store if not already done
  if (!authStore.initialized) {
    await authStore.initialize()
  }

  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/auth/login')
    return
  }

  // Check if route requires guest (not authenticated)
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  next()
})

export default router
