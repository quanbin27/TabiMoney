// User types
export interface User {
  id: number
  email: string
  username: string
  first_name: string
  last_name: string
  phone: string
  avatar_url: string
  is_verified: boolean
  last_login_at: string | null
  created_at: string
  updated_at: string
  profile?: UserProfile
}

export interface UserProfile {
  id: number
  user_id: number
  monthly_income: number
  currency: string
  timezone: string
  language: string
  notification_settings: string
  ai_settings: string
  created_at: string
  updated_at: string
}

// Auth types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
  first_name?: string
  last_name?: string
  phone?: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  expires_at: string
}

// Transaction types
export interface Transaction {
  id: number
  user_id: number
  category_id: number
  amount: number
  description: string
  transaction_type: 'income' | 'expense' | 'transfer'
  transaction_date: string
  transaction_time?: string
  location: string
  tags: string[]
  metadata: Record<string, any>
  is_recurring: boolean
  recurring_pattern?: string
  parent_transaction_id?: number
  ai_confidence?: number
  ai_suggested_category_id?: number
  created_at: string
  updated_at: string
  category?: Category
  ai_suggested_category?: Category
}

export interface TransactionCreateRequest {
  category_id: number
  amount: number
  description: string
  transaction_type: 'income' | 'expense' | 'transfer'
  transaction_date: string
  transaction_time?: string
  location?: string
  tags?: string[]
  metadata?: Record<string, any>
  is_recurring?: boolean
  recurring_pattern?: string
}

export interface TransactionUpdateRequest {
  category_id: number
  amount: number
  description: string
  transaction_type: 'income' | 'expense' | 'transfer'
  transaction_date: string
  transaction_time?: string
  location?: string
  tags?: string[]
  metadata?: Record<string, any>
}

export interface TransactionQueryRequest {
  page: number
  limit: number
  category_id?: number
  transaction_type?: 'income' | 'expense' | 'transfer'
  start_date?: string
  end_date?: string
  min_amount?: number
  max_amount?: number
  search?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

// Category types
export interface Category {
  id: number
  user_id?: number
  name: string
  name_en: string
  description: string
  icon: string
  color: string
  parent_id?: number
  is_system: boolean
  is_active: boolean
  sort_order: number
  created_at: string
  updated_at: string
}

export interface CategoryCreateRequest {
  name: string
  name_en?: string
  description?: string
  icon?: string
  color?: string
  parent_id?: number
}

export interface CategoryUpdateRequest {
  name: string
  name_en?: string
  description?: string
  icon?: string
  color?: string
  parent_id?: number
  is_active?: boolean
}

// Financial Goal types
export interface FinancialGoal {
  id: number
  user_id: number
  title: string
  description: string
  target_amount: number
  current_amount: number
  target_date?: string
  goal_type: 'savings' | 'debt_payment' | 'investment' | 'purchase' | 'other'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  is_achieved: boolean
  achieved_at?: string
  created_at: string
  updated_at: string
  progress?: number
}

export interface FinancialGoalCreateRequest {
  title: string
  description?: string
  target_amount: number
  target_date?: string
  goal_type: 'savings' | 'debt_payment' | 'investment' | 'purchase' | 'other'
  priority?: 'low' | 'medium' | 'high' | 'urgent'
}

export interface FinancialGoalUpdateRequest {
  title: string
  description?: string
  target_amount: number
  current_amount?: number
  target_date?: string
  goal_type: 'savings' | 'debt_payment' | 'investment' | 'purchase' | 'other'
  priority?: 'low' | 'medium' | 'high' | 'urgent'
}

// Budget types
export interface Budget {
  id: number
  user_id: number
  category_id?: number
  name: string
  amount: number
  period: 'weekly' | 'monthly' | 'yearly'
  start_date: string
  end_date: string
  is_active: boolean
  alert_threshold: number
  created_at: string
  updated_at: string
  category?: Category
  spent_amount?: number
  remaining_amount?: number
  usage_percentage?: number
}

export interface BudgetCreateRequest {
  category_id?: number
  name: string
  amount: number
  period: 'weekly' | 'monthly' | 'yearly'
  start_date: string
  end_date: string
  alert_threshold?: number
}

export interface BudgetUpdateRequest {
  category_id?: number
  name: string
  amount: number
  period: 'weekly' | 'monthly' | 'yearly'
  start_date: string
  end_date: string
  is_active?: boolean
  alert_threshold?: number
}

// Analytics types
export interface DashboardAnalytics {
  user_id: number
  period: string
  total_income: number
  total_expense: number
  net_amount: number
  transaction_count: number
  category_breakdown: CategoryAnalytics[]
  monthly_trends: MonthlyTrend[]
  top_categories: CategoryAnalytics[]
  financial_health: FinancialHealth
  generated_at: string
}

export interface CategoryAnalytics {
  category_id: number
  category_name: string
  amount: number
  transaction_count: number
  percentage: number
  trend: string
  average_amount: number
}

export interface MonthlyTrend {
  month: string
  income: number
  expense: number
  net_amount: number
  transaction_count: number
}

export interface FinancialHealth {
  score: number
  level: 'excellent' | 'good' | 'fair' | 'poor'
  income_ratio: number
  savings_rate: number
  debt_ratio: number
  recommendations: string[]
}

// AI types
export interface NLURequest {
  text: string
  user_id: number
  context?: string
}

export interface NLUResponse {
  user_id: number
  intent: string
  entities: Entity[]
  confidence: number
  suggested_action: string
  response: string
  generated_at: string
}

export interface Entity {
  type: string
  value: string
  confidence: number
  start_pos: number
  end_pos: number
}

export interface ExpensePredictionRequest {
  user_id: number
  start_date: string
  end_date: string
  category_id?: number
}

export interface ExpensePredictionResponse {
  user_id: number
  predicted_amount: number
  confidence_score: number
  category_breakdown: CategoryPrediction[]
  trends: ExpenseTrend[]
  recommendations: string[]
  generated_at: string
}

export interface CategoryPrediction {
  category_id: number
  category_name: string
  predicted_amount: number
  confidence_score: number
  trend: string
}

export interface ExpenseTrend {
  period: string
  amount: number
  change_percentage: number
  trend: string
}

export interface AnomalyDetectionRequest {
  user_id: number
  start_date: string
  end_date: string
  threshold: number
}

export interface AnomalyDetectionResponse {
  user_id: number
  anomalies: Anomaly[]
  total_anomalies: number
  detection_score: number
  generated_at: string
}

export interface Anomaly {
  transaction_id: number
  amount: number
  category_name: string
  anomaly_score: number
  anomaly_type: string
  description: string
  transaction_date: string
}

export interface CategorySuggestionRequest {
  user_id: number
  description: string
  amount: number
  location?: string
  tags?: string[]
}

export interface CategorySuggestionResponse {
  user_id: number
  description: string
  amount: number
  suggestions: CategorySuggestion[]
  confidence_score: number
  generated_at: string
}

export interface CategorySuggestion {
  category_id: number
  category_name: string
  confidence_score: number
  reason: string
  is_user_category: boolean
}

export interface SpendingPatternRequest {
  user_id: number
  start_date: string
  end_date: string
  granularity?: 'daily' | 'weekly' | 'monthly'
}

export interface SpendingPatternResponse {
  user_id: number
  patterns: SpendingPattern[]
  insights: string[]
  recommendations: string[]
  generated_at: string
}

export interface SpendingPattern {
  category_id: number
  category_name: string
  total_amount: number
  transaction_count: number
  average_amount: number
  frequency: string
  trend: string
  peak_days: string[]
  peak_times: string[]
}

export interface GoalAnalysisRequest {
  user_id: number
  goal_id: number
}

export interface GoalAnalysisResponse {
  user_id: number
  goal_id: number
  progress: number
  on_track: boolean
  projected_date?: string
  recommendations: string[]
  risk_factors: string[]
  generated_at: string
}

export interface ChatRequest {
  message: string
  user_id: number
}

export interface ChatResponse {
  user_id: number
  message: string
  response: string
  intent: string
  entities: Entity[]
  suggestions: string[]
  generated_at: string
}

// API Response types
export interface ApiResponse<T = any> {
  data: T
  message?: string
  status: number
}

export interface PaginatedResponse<T = any> {
  data: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}

// Error types
export interface ApiError {
  error: string
  message: string
  status: number
}
