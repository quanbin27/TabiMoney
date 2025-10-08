package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type AIService struct {
	config      *config.Config
	db          *gorm.DB
	httpClient  *http.Client
	aiServiceURL string
    intentHandlers map[string]func(*models.ChatRequest, *models.ChatResponse)
}

func NewAIService(cfg *config.Config) *AIService {
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		aiURL = "http://localhost:8001"
	}
    svc := &AIService{
		config:       cfg,
		db:           database.GetDB(),
		httpClient:   &http.Client{Timeout: 60 * time.Second},
		aiServiceURL: aiURL + "/api/v1",
    }
    svc.intentHandlers = map[string]func(*models.ChatRequest, *models.ChatResponse){
        "query_balance": svc.handleQueryBalance,
        "add_transaction": svc.handleAddTransaction,
        "create_goal": svc.handleCreateGoal,
        "list_goals": svc.handleListGoals,
        "update_goal": svc.handleUpdateGoal,
        "create_budget": svc.handleCreateBudget,
        "list_budgets": svc.handleListBudgets,
        "update_budget": svc.handleUpdateBudget,
    }
    return svc
}

// NLU (Natural Language Understanding) Service
func (s *AIService) ProcessNLU(req *models.NLURequest) (*models.NLUResponse, error) {
	// Check cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("nlu:%d:%s", req.UserID, req.Text)
	if cached, err := database.GetCache(ctx, cacheKey); err == nil {
		var response models.NLUResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Call AI Service
	requestBody, err := json.Marshal(map[string]interface{}{
		"text":     req.Text,
		"user_id":  req.UserID,
		"context":  req.Context,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.aiServiceURL+"/nlu/process", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
	}

    var aiResponse struct {
		UserID          int                    `json:"user_id"`
		Intent          string                 `json:"intent"`
		Entities        []models.Entity        `json:"entities"`
		Confidence      float64                `json:"confidence"`
		SuggestedAction string                 `json:"suggested_action"`
		Response        string                 `json:"response"`
		GeneratedAt     string                 `json:"generated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

    response := &models.NLUResponse{
		UserID:          uint64(aiResponse.UserID),
		Intent:          aiResponse.Intent,
		Entities:        aiResponse.Entities,
		Confidence:      aiResponse.Confidence,
		SuggestedAction: aiResponse.SuggestedAction,
		Response:        aiResponse.Response,
		GeneratedAt:     time.Now(),
	}

	// Cache response
	if responseJSON, err := json.Marshal(response); err == nil {
		database.SetCache(ctx, cacheKey, responseJSON, 1*time.Hour)
	}

    return response, nil
}

// Expense Prediction Service
func (s *AIService) PredictExpenses(req *models.ExpensePredictionRequest) (*models.ExpensePredictionResponse, error) {
	// Check cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("expense_prediction:%d:%s:%s", req.UserID, req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))
	if cached, err := database.GetCache(ctx, cacheKey); err == nil {
		var response models.ExpensePredictionResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

    // Call AI Service (Prediction)
    requestBody, err := json.Marshal(map[string]interface{}{
        "user_id":     req.UserID,
        "start_date":  req.StartDate.Format("2006-01-02T15:04:05Z07:00"),
        "end_date":    req.EndDate.Format("2006-01-02T15:04:05Z07:00"),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", s.aiServiceURL+"/prediction/expenses", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    httpReq.Header.Set("Content-Type", "application/json")

    log.Printf("Calling AI Service: %s", s.aiServiceURL+"/prediction/expenses")
    log.Printf("Request body: %s", string(requestBody))

    resp, err := s.httpClient.Do(httpReq)
    if err != nil {
        log.Printf("AI Service call failed: %v", err)
        return nil, fmt.Errorf("failed to call AI service: %w", err)
    }
    defer resp.Body.Close()

    log.Printf("AI Service response status: %d", resp.StatusCode)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
    }

    var response models.ExpensePredictionResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        log.Printf("Failed to decode AI service response: %v", err)
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    log.Printf("AI Service response decoded successfully: %+v", response)

    // Cache response
    if b, err := json.Marshal(response); err == nil {
        database.SetCache(ctx, cacheKey, b, 6*time.Hour)
        // Save to database
        analysis := &models.AIAnalysis{
            UserID:          req.UserID,
            AnalysisType:    "expense_prediction",
            Data:            string(b),
            ConfidenceScore: response.ConfidenceScore,
            ModelVersion:    "ai-service",
        }
        s.db.Create(analysis)
    }

    return &response, nil
}

// Anomaly Detection Service
func (s *AIService) DetectAnomalies(req *models.AnomalyDetectionRequest) (*models.AnomalyDetectionResponse, error) {
    // Get transaction data (preload Category to avoid nil dereference)
    var transactions []models.Transaction
    query := s.db.Where("user_id = ? AND transaction_date BETWEEN ? AND ?", 
        req.UserID, req.StartDate, req.EndDate).Preload("Category")
    
    if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Analyze patterns
	anomalies := s.analyzeTransactionPatterns(transactions, req.Threshold)
	
	response := &models.AnomalyDetectionResponse{
		UserID:         req.UserID,
		Anomalies:      anomalies,
		TotalAnomalies: len(anomalies),
		DetectionScore: s.calculateDetectionScore(anomalies),
		GeneratedAt:    time.Now(),
	}

	// Save to database
	analysis := &models.AIAnalysis{
		UserID:          req.UserID,
		AnalysisType:    "anomaly_detection",
		Data:            s.marshalToJSON(response),
		ConfidenceScore: response.DetectionScore,
		ModelVersion:    "custom",
	}
	s.db.Create(analysis)

	return response, nil
}

// Category Suggestion Service
func (s *AIService) SuggestCategory(req *models.CategorySuggestionRequest) (*models.CategorySuggestionResponse, error) {
	// Get user's categories
	var categories []models.Category
	if err := s.db.Where("user_id = ? OR is_system = ?", req.UserID, true).
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Get user's transaction history for context
	var recentTransactions []models.Transaction
	if err := s.db.Where("user_id = ? AND transaction_date >= ?", 
		req.UserID, time.Now().AddDate(0, -3, 0)).
		Order("transaction_date DESC").
		Limit(50).
		Find(&recentTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}

    // Call AI Service (Categorization)
    ctx := context.Background()
    requestBody, err := json.Marshal(map[string]interface{}{
        "user_id":     req.UserID,
        "description": req.Description,
        "amount":      req.Amount,
        "location":    req.Location,
        "tags":        req.Tags,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", s.aiServiceURL+"/categorization/suggest", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := s.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to call AI service: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
    }

    var response models.CategorySuggestionResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &response, nil
}

// Spending Pattern Analysis
func (s *AIService) AnalyzeSpendingPattern(req *models.SpendingPatternRequest) (*models.SpendingPatternResponse, error) {
	// Get transaction data
	var transactions []models.Transaction
	query := s.db.Where("user_id = ? AND transaction_date BETWEEN ? AND ?", 
		req.UserID, req.StartDate, req.EndDate)
	
	if err := query.Preload("Category").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

    // Analyze patterns (aggregate locally)
    patterns := s.analyzeSpendingPatterns(transactions, req.Granularity)
    // Ask AI service for dynamic insights
    insights, recommendations := s.fetchDynamicSpendingInsights(req.UserID, patterns)

	response := &models.SpendingPatternResponse{
		UserID:          req.UserID,
		Patterns:        patterns,
		Insights:        insights,
		Recommendations: recommendations,
		GeneratedAt:     time.Now(),
	}

	// Save to database
	analysis := &models.AIAnalysis{
		UserID:          req.UserID,
		AnalysisType:    "spending_pattern",
		Data:            s.marshalToJSON(response),
		ConfidenceScore: 0.85, // Based on pattern analysis
		ModelVersion:    "custom",
	}
	s.db.Create(analysis)

	return response, nil
}

// fetchDynamicSpendingInsights calls AI service to get dynamic insights; falls back to rule-based
func (s *AIService) fetchDynamicSpendingInsights(userID uint64, patterns []models.SpendingPattern) ([]string, []string) {
    // Build minimal payload for AI-service /analysis/spending
    type cs struct {
        CategoryID       uint64  `json:"category_id"`
        CategoryName     string  `json:"category_name"`
        TotalAmount      float64 `json:"total_amount"`
        TransactionCount int     `json:"transaction_count"`
    }
    payload := map[string]interface{}{
        "user_id": userID,
        "patterns": []cs{},
    }
    arr := make([]cs, 0, len(patterns))
    for _, p := range patterns {
        arr = append(arr, cs{CategoryID: p.CategoryID, CategoryName: p.CategoryName, TotalAmount: p.TotalAmount, TransactionCount: p.TransactionCount})
    }
    payload["patterns"] = arr

    b, _ := json.Marshal(payload)
    req, err := http.NewRequest("POST", s.aiServiceURL+"/analysis/spending", bytes.NewBuffer(b))
    if err != nil {
        return s.generateSpendingInsights(patterns), s.generateSpendingRecommendations(patterns)
    }
    req.Header.Set("Content-Type", "application/json")
    resp, err := s.httpClient.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        return s.generateSpendingInsights(patterns), s.generateSpendingRecommendations(patterns)
    }
    defer resp.Body.Close()
    var out struct {
        Insights        []string `json:"insights"`
        Recommendations []string `json:"recommendations"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return s.generateSpendingInsights(patterns), s.generateSpendingRecommendations(patterns)
    }
    if out.Insights == nil {
        out.Insights = []string{}
    }
    if out.Recommendations == nil {
        out.Recommendations = []string{}
    }
    return out.Insights, out.Recommendations
}

// Goal Analysis
func (s *AIService) AnalyzeGoal(req *models.GoalAnalysisRequest) (*models.GoalAnalysisResponse, error) {
	// Get goal data
	var goal models.FinancialGoal
	if err := s.db.First(&goal, req.GoalID).Error; err != nil {
		return nil, fmt.Errorf("goal not found: %w", err)
	}

	// Get user's financial data
	var transactions []models.Transaction
	if err := s.db.Where("user_id = ? AND transaction_type = ?", 
		req.UserID, "expense").
		Where("transaction_date >= ?", goal.CreatedAt).
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Analyze goal progress
	progress := s.calculateGoalProgress(goal, transactions)
	onTrack := s.isGoalOnTrack(goal, progress)
	projectedDate := s.projectGoalCompletion(goal, progress)
	recommendations := s.generateGoalRecommendations(goal, progress)
	riskFactors := s.identifyGoalRiskFactors(goal, progress)

	response := &models.GoalAnalysisResponse{
		UserID:          req.UserID,
		GoalID:          req.GoalID,
		Progress:        progress,
		OnTrack:         onTrack,
		ProjectedDate:   projectedDate,
		Recommendations: recommendations,
		RiskFactors:     riskFactors,
		GeneratedAt:     time.Now(),
	}

	// Save to database
	analysis := &models.AIAnalysis{
		UserID:          req.UserID,
		AnalysisType:    "goal_analysis",
		Data:            s.marshalToJSON(response),
		ConfidenceScore: 0.9,
		ModelVersion:    "custom",
	}
	s.db.Create(analysis)

	return response, nil
}

// Chat Service
func (s *AIService) ProcessChat(req *models.ChatRequest) (*models.ChatResponse, error) {
	// Get user context
	var user models.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get recent transactions for context
	var recentTransactions []models.Transaction
	if err := s.db.Where("user_id = ?", req.UserID).
		Order("transaction_date DESC").
		Limit(10).
		Preload("Category").
		Find(&recentTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}

    // Call AI Service (Chat)
    ctx := context.Background()
    requestBody, err := json.Marshal(map[string]interface{}{
        "user_id":  req.UserID,
        "message":  req.Message,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", s.aiServiceURL+"/chat/process", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := s.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to call AI service: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
    }

    var response models.ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    // Apply short-term context memory to backfill missing entities
    s.applyShortTermContext(req.UserID, &response)

    // Route to handler if available
    if handler, ok := s.intentHandlers[response.Intent]; ok {
        handler(req, &response)
    }

	// Save chat message
    entitiesJSON, _ := json.Marshal(response.Entities)
    chatMessage := &models.ChatMessage{
		UserID:      req.UserID,
		Message:     req.Message,
		Response:    response.Response,
		Intent:      response.Intent,
        Entities:    string(entitiesJSON),
		IsProcessed: true,
	}
	s.db.Create(chatMessage)

    return &response, nil
}

// applyShortTermContext backfills missing common entities (amount, category, title, period)
func (s *AIService) applyShortTermContext(userID uint64, resp *models.ChatResponse) {
    has := func(t string) bool {
        for _, e := range resp.Entities {
            if e.Type == t && e.Value != "" {
                return true
            }
        }
        return false
    }

    needAmount := !has("amount")
    needCategory := !has("category")
    needTitle := !has("title")
    needPeriod := !has("period")
    if !(needAmount || needCategory || needTitle || needPeriod) {
        return
    }

    var chats []models.ChatMessage
    if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(5).Find(&chats).Error; err != nil || len(chats) == 0 {
        return
    }

    for _, cm := range chats {
        if cm.Entities == "" {
            continue
        }
        var ents []models.Entity
        if err := json.Unmarshal([]byte(cm.Entities), &ents); err != nil {
            continue
        }
        for _, e := range ents {
            if needAmount && e.Type == "amount" && e.Value != "" {
                resp.Entities = append(resp.Entities, e)
                needAmount = false
            }
            if needCategory && e.Type == "category" && e.Value != "" {
                resp.Entities = append(resp.Entities, e)
                needCategory = false
            }
            if needTitle && e.Type == "title" && e.Value != "" {
                resp.Entities = append(resp.Entities, e)
                needTitle = false
            }
            if needPeriod && e.Type == "period" && e.Value != "" {
                resp.Entities = append(resp.Entities, e)
                needPeriod = false
            }
        }
        if !(needAmount || needCategory || needTitle || needPeriod) {
            break
        }
    }
}

// Handlers
func (s *AIService) handleQueryBalance(req *models.ChatRequest, resp *models.ChatResponse) {
    now := time.Now().UTC()
    startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

    type totalsRow struct {
        TotalIncome  float64 `gorm:"column:total_income"`
        TotalExpense float64 `gorm:"column:total_expense"`
    }
    var totals totalsRow
    if err := s.db.Table("transactions").
        Select("SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) as total_income, SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) as total_expense").
        Where("user_id = ? AND transaction_date BETWEEN ? AND ?", req.UserID, startDate, endDate).
        Scan(&totals).Error; err == nil {
        netAmount := totals.TotalIncome - totals.TotalExpense
        resp.Response = fmt.Sprintf("Tháng %02d/%d: Tổng chi tiêu %.2f VND, tổng thu %.2f VND, chênh lệch %.2f VND.",
            now.Month(), now.Year(), totals.TotalExpense, totals.TotalIncome, netAmount)
    }
    return
}

func (s *AIService) handleAddTransaction(req *models.ChatRequest, resp *models.ChatResponse) {
    var amount float64
    var amountMaxConfidence float64
    var categoryName string
    for _, e := range resp.Entities {
        switch e.Type {
        case "amount":
            if v, err := strconv.ParseFloat(e.Value, 64); err == nil && v > 0 {
                if v > amount {
                    amount = v
                }
                if e.Confidence > amountMaxConfidence {
                    amountMaxConfidence = e.Confidence
                }
            }
        case "category":
            if categoryName == "" {
                categoryName = e.Value
            }
        }
    }

    if amount <= 0 {
        return
    }

    // Confidence threshold for confirmation
    if amountMaxConfidence > 0 && amountMaxConfidence < 0.6 {
        resp.Response = fmt.Sprintf("Tôi hiểu bạn muốn thêm giao dịch khoảng %.0f VND. Bạn xác nhận số tiền này chứ?", amount)
        return
    }

    var category models.Category
    if err := s.db.Where("LOWER(name) = ?", strings.ToLower(categoryName)).First(&category).Error; err != nil {
        // normalize common Vietnamese variants
        lower := strings.ToLower(categoryName)
        if strings.Contains(lower, "an uong") || strings.Contains(lower, "ăn uống") {
            s.db.Where("LOWER(name) = ?", strings.ToLower("Ăn uống")).First(&category)
        }
    }
    if category.ID == 0 {
        s.db.Where("name = ?", "Khác").First(&category)
    }

    today := time.Now().UTC().Format("2006-01-02")
    createReq := &models.TransactionCreateRequest{
        CategoryID:      category.ID,
        Amount:          amount,
        Description:     req.Message,
        TransactionType: "expense",
        TransactionDate: today,
    }

    txnSvc := NewTransactionService()
    if _, err := txnSvc.CreateTransaction(req.UserID, createReq); err != nil {
        log.Printf("AI chat add_transaction: failed to create transaction: %v", err)
        return
    }
    if category.ID != 0 {
        resp.Response = fmt.Sprintf("Đã thêm giao dịch %.0f VND cho danh mục %s.", amount, category.Name)
    } else {
        resp.Response = fmt.Sprintf("Đã thêm giao dịch %.0f VND.", amount)
    }
    return
}

// Helper methods for prompt building and response parsing

func (s *AIService) buildNLUPrompt(text, context string) string {
	return fmt.Sprintf(`You are a financial assistant for TabiMoney. Analyze the following user input and extract:

1. Intent (add_transaction, query_balance, ask_question, etc.)
2. Entities (amounts, categories, dates, descriptions)
3. Confidence score (0-1)
4. Suggested action
5. Response

User input: %s
Context: %s

Return JSON format:
{
  "intent": "add_transaction",
  "entities": [
    {"type": "amount", "value": "50000", "confidence": 0.95},
    {"type": "category", "value": "food", "confidence": 0.9}
  ],
  "confidence": 0.92,
  "suggested_action": "create_transaction",
  "response": "I'll help you add this transaction..."
}`, text, context)
}

func (s *AIService) parseNLUResponse(content string, userID uint64) (*models.NLUResponse, error) {
	var response models.NLUResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, err
	}
	response.UserID = userID
	response.GeneratedAt = time.Now()
	return &response, nil
}

func (s *AIService) buildExpensePredictionPrompt(historicalData []models.Transaction, req *models.ExpensePredictionRequest) string {
	// Build comprehensive prompt with historical data
	return fmt.Sprintf(`Analyze the following financial data and predict expenses for the period %s to %s:

Historical Data: %s

Provide predictions in JSON format with:
- Predicted total amount
- Category breakdown
- Trends analysis
- Recommendations

Focus on patterns, seasonality, and user behavior.`, 
		req.StartDate.Format("2006-01-02"), 
		req.EndDate.Format("2006-01-02"),
		s.marshalToJSON(historicalData))
}

func (s *AIService) parseExpensePredictionResponse(content string, userID uint64) (*models.ExpensePredictionResponse, error) {
	var response models.ExpensePredictionResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, err
	}
	response.UserID = userID
	response.GeneratedAt = time.Now()
	return &response, nil
}

func (s *AIService) buildCategorySuggestionPrompt(req *models.CategorySuggestionRequest, categories []models.Category, transactions []models.Transaction) string {
	return fmt.Sprintf(`Suggest the most appropriate category for this transaction:

Description: %s
Amount: %.2f
Location: %s
Tags: %v

Available Categories: %s
Recent Transactions: %s

Return JSON with suggestions ranked by confidence.`, 
		req.Description, req.Amount, req.Location, req.Tags,
		s.marshalToJSON(categories),
		s.marshalToJSON(transactions))
}

func (s *AIService) parseCategorySuggestionResponse(content string, userID uint64) (*models.CategorySuggestionResponse, error) {
	var response models.CategorySuggestionResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, err
	}
	response.UserID = userID
	response.GeneratedAt = time.Now()
	return &response, nil
}

func (s *AIService) buildChatPrompt(message string, user models.User, transactions []models.Transaction) string {
	return fmt.Sprintf(`You are TabiMoney, a helpful financial assistant. Respond to the user's message about their personal finances.

User: %s
Recent Transactions: %s

Provide helpful, accurate financial advice. Be conversational and supportive.`, 
		message, s.marshalToJSON(transactions))
}

func (s *AIService) parseChatResponse(content string, userID uint64) (*models.ChatResponse, error) {
	var response models.ChatResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		// If JSON parsing fails, create a simple response
		response = models.ChatResponse{
			UserID:      userID,
			Message:     content,
			Response:    content,
			Intent:      "general",
			Entities:    []models.Entity{},
			Suggestions: []string{},
			GeneratedAt: time.Now(),
		}
	}
	return &response, nil
}

// Additional helper methods for data analysis

func (s *AIService) getHistoricalExpenseData(userID uint64, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := s.db.Where("user_id = ? AND transaction_type = ? AND transaction_date BETWEEN ? AND ?", 
		userID, "expense", startDate, endDate).
		Preload("Category").
		Find(&transactions).Error
	return transactions, err
}

func (s *AIService) analyzeTransactionPatterns(transactions []models.Transaction, threshold float64) []models.Anomaly {
	// Implement anomaly detection logic
	// This is a simplified version - in production, use more sophisticated algorithms
	var anomalies []models.Anomaly
	
	// Calculate average amount
	var totalAmount float64
	for _, t := range transactions {
		totalAmount += t.Amount
	}
	avgAmount := totalAmount / float64(len(transactions))

	// Find anomalies
	for _, t := range transactions {
		if t.Amount > avgAmount*2 { // Simple threshold-based detection
			anomalies = append(anomalies, models.Anomaly{
				TransactionID:   t.ID,
				Amount:         t.Amount,
				CategoryName:   t.Category.Name,
				AnomalyScore:   0.8,
				AnomalyType:    "amount",
				Description:    "Unusually high transaction amount",
				TransactionDate: t.TransactionDate,
			})
		}
	}

	return anomalies
}

func (s *AIService) calculateDetectionScore(anomalies []models.Anomaly) float64 {
	if len(anomalies) == 0 {
		return 0.0
	}
	
	var totalScore float64
	for _, a := range anomalies {
		totalScore += a.AnomalyScore
	}
	
	return totalScore / float64(len(anomalies))
}

func (s *AIService) analyzeSpendingPatterns(transactions []models.Transaction, granularity string) []models.SpendingPattern {
	// Implement spending pattern analysis
	// This is a simplified version
	var patterns []models.SpendingPattern
	
	// Group by category
	categoryMap := make(map[uint64]*models.SpendingPattern)
	
	for _, t := range transactions {
		if pattern, exists := categoryMap[t.CategoryID]; exists {
			pattern.TotalAmount += t.Amount
			pattern.TransactionCount++
		} else {
			pattern = &models.SpendingPattern{
				CategoryID:       t.CategoryID,
				CategoryName:     t.Category.Name,
				TotalAmount:      t.Amount,
				TransactionCount: 1,
			}
			categoryMap[t.CategoryID] = pattern
		}
	}
	
	// Convert to slice
	for _, pattern := range categoryMap {
		pattern.AverageAmount = pattern.TotalAmount / float64(pattern.TransactionCount)
		patterns = append(patterns, *pattern)
	}
	
	return patterns
}

func (s *AIService) generateSpendingInsights(patterns []models.SpendingPattern) []string {
    // Generate insights based on patterns (Vietnamese)
    insights := []string{
        "Mẫu chi tiêu của bạn cho thấy hành vi ổn định giữa các danh mục.",
        "Hãy xem lại các danh mục chi tiêu hàng đầu để tối ưu hoá.",
    }
    return insights
}

func (s *AIService) generateSpendingRecommendations(patterns []models.SpendingPattern) []string {
    // Generate recommendations based on patterns (Vietnamese)
    recommendations := []string{
        "Thiết lập cảnh báo ngân sách cho các danh mục chi tiêu hàng đầu.",
        "Cân nhắc đặt hạn mức chi tiêu hàng tháng cho các danh mục tuỳ ý.",
    }
    return recommendations
}

func (s *AIService) calculateGoalProgress(goal models.FinancialGoal, transactions []models.Transaction) float64 {
	// Calculate progress based on transactions
	// This is simplified - in reality, you'd need to consider income, savings, etc.
	return goal.CurrentAmount / goal.TargetAmount
}

func (s *AIService) isGoalOnTrack(goal models.FinancialGoal, progress float64) bool {
	// Determine if goal is on track based on progress and time
	// This is simplified logic
	return progress > 0.5 // Example threshold
}

func (s *AIService) projectGoalCompletion(goal models.FinancialGoal, progress float64) *time.Time {
	// Project when goal will be completed
	// This is simplified logic
	if progress >= 1.0 {
		return &time.Time{}
	}
	
	// Simple linear projection
	remaining := goal.TargetAmount - goal.CurrentAmount
	monthlyContribution := goal.TargetAmount * 0.1 // Example: 10% per month
	monthsRemaining := remaining / monthlyContribution
	
	completionDate := time.Now().AddDate(0, int(monthsRemaining), 0)
	return &completionDate
}

func (s *AIService) generateGoalRecommendations(goal models.FinancialGoal, progress float64) []string {
    recommendations := []string{
        "Hãy tăng mức tiết kiệm hàng tháng để đạt mục tiêu nhanh hơn.",
        "Xem lại các khoản chi để tìm cơ hội tiết kiệm thêm.",
    }
    return recommendations
}

func (s *AIService) identifyGoalRiskFactors(goal models.FinancialGoal, progress float64) []string {
    riskFactors := []string{
        "Tiến độ thấp có thể cho thấy tỷ lệ tiết kiệm chưa đủ.",
        "Biến động thị trường có thể ảnh hưởng đến mục tiêu đầu tư.",
    }
    return riskFactors
}

func (s *AIService) marshalToJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return "{}"
	}
	return string(jsonData)
}

func formatCurrency(amount float64) string {
    // Format number with thousand separators
    s := fmt.Sprintf("%.0f", amount)
    n := len(s)
    if n <= 3 {
        return s
    }
    var b strings.Builder
    pre := n % 3
    if pre == 0 {
        pre = 3
    }
    b.WriteString(s[:pre])
    for i := pre; i < n; i += 3 {
        b.WriteString(",")
        b.WriteString(s[i : i+3])
    }
    return b.String()
}

func (s *AIService) handleCreateGoal(req *models.ChatRequest, resp *models.ChatResponse) {
    var amount float64
    var amountMaxConfidence float64
    var title, goalType string
    var titleConfidence float64

    // Extract entities
    for _, e := range resp.Entities {
        switch e.Type {
        case "amount":
            if v, err := strconv.ParseFloat(e.Value, 64); err == nil && v > 0 {
                if v > amount {
                    amount = v
                }
                if e.Confidence > amountMaxConfidence {
                    amountMaxConfidence = e.Confidence
                }
            }
        case "goal_type":
            goalType = e.Value
        case "title":
            title = e.Value
            if e.Confidence > titleConfidence {
                titleConfidence = e.Confidence
            }
        }
    }

    if amount <= 0 {
        resp.Response = "Vui lòng cung cấp số tiền mục tiêu để tạo mục tiêu tài chính."
        return
    }

    // Confidence threshold for confirmation
    if amountMaxConfidence > 0 && amountMaxConfidence < 0.6 {
        resp.Response = fmt.Sprintf("Bạn muốn tạo mục tiêu với số tiền %s VND? Vui lòng xác nhận hoặc cung cấp lại số tiền.", formatCurrency(amount))
        return
    }

    if title == "" {
        title = "Mục tiêu tài chính"
    }
    // Ask confirmation if title detected with low confidence
    if titleConfidence > 0 && titleConfidence < 0.6 {
        resp.Response = fmt.Sprintf("Bạn muốn đặt tên mục tiêu là '%s' chứ? Vui lòng xác nhận hoặc cung cấp tên khác.", title)
        return
    }
    if goalType == "" {
        goalType = "savings"
    }

    // Map Vietnamese goal types
    goalTypeMap := map[string]string{
        "tiết kiệm": "savings",
        "mua sắm": "purchase",
        "đầu tư": "investment",
        "trả nợ": "debt_payment",
    }
    if mapped, exists := goalTypeMap[goalType]; exists {
        goalType = mapped
    }

    createReq := &models.FinancialGoalCreateRequest{
        Title:       title,
        Description: req.Message,
        TargetAmount: amount,
        GoalType:    goalType,
        Priority:    "medium",
    }

    goalService := NewGoalService()
    if goal, err := goalService.CreateGoal(req.UserID, createReq); err != nil {
        log.Printf("AI chat create_goal: failed to create goal: %v", err)
        resp.Response = "Không thể tạo mục tiêu tài chính. Vui lòng thử lại."
    } else {
        resp.Response = fmt.Sprintf("Đã tạo mục tiêu '%s' với số tiền %s VND.", goal.Title, formatCurrency(goal.TargetAmount))
    }
}

func (s *AIService) handleListGoals(req *models.ChatRequest, resp *models.ChatResponse) {
    goalService := NewGoalService()
    goals, err := goalService.GetGoals(req.UserID)
    if err != nil {
        log.Printf("AI chat list_goals: failed to get goals: %v", err)
        resp.Response = "Không thể lấy danh sách mục tiêu. Vui lòng thử lại."
        return
    }

    if len(goals) == 0 {
        resp.Response = "Bạn chưa có mục tiêu tài chính nào. Hãy tạo mục tiêu đầu tiên!"
        return
    }

    resp.Response = fmt.Sprintf("Bạn có %d mục tiêu tài chính:\n", len(goals))
    for i, goal := range goals {
        if i >= 5 { // Limit to 5 goals for readability
            resp.Response += "..."
            break
        }
        progress := (goal.CurrentAmount / goal.TargetAmount) * 100
        resp.Response += fmt.Sprintf("- %s: %s/%s VND (%.1f%%)\n", 
            goal.Title, formatCurrency(goal.CurrentAmount), formatCurrency(goal.TargetAmount), progress)
    }
}

func (s *AIService) handleUpdateGoal(req *models.ChatRequest, resp *models.ChatResponse) {
    var amount float64
    var goalID uint64

    // Extract entities
    for _, e := range resp.Entities {
        switch e.Type {
        case "amount":
            if v, err := strconv.ParseFloat(e.Value, 64); err == nil && v > 0 {
                amount = v
            }
        case "goal_id":
            if v, err := strconv.ParseUint(e.Value, 10, 64); err == nil {
                goalID = v
            }
        }
    }

    if amount <= 0 {
        resp.Response = "Vui lòng cung cấp số tiền để cập nhật mục tiêu."
        return
    }

    // If no goal ID specified, get the first goal
    if goalID == 0 {
        goalService := NewGoalService()
        goals, err := goalService.GetGoals(req.UserID)
        if err != nil || len(goals) == 0 {
            resp.Response = "Không tìm thấy mục tiêu nào để cập nhật."
            return
        }
        goalID = goals[0].ID
    }

    goalService := NewGoalService()
    if goal, err := goalService.AddContribution(req.UserID, goalID, amount, req.Message); err != nil {
        log.Printf("AI chat update_goal: failed to update goal: %v", err)
        resp.Response = "Không thể cập nhật mục tiêu. Vui lòng thử lại."
    } else {
        progress := (goal.CurrentAmount / goal.TargetAmount) * 100
        resp.Response = fmt.Sprintf("Đã thêm %s VND vào mục tiêu '%s'. Tiến độ: %.1f%%", 
            formatCurrency(amount), goal.Title, progress)
    }
}

func (s *AIService) handleCreateBudget(req *models.ChatRequest, resp *models.ChatResponse) {
    var amount float64
    var amountMaxConfidence float64
    var name, period string
    var categoryConfidence, periodConfidence float64
    var categoryID *uint64
    var alertThreshold *float64
    var categoryName string

    // Extract entities
    for _, e := range resp.Entities {
        switch e.Type {
        case "amount":
            if v, err := strconv.ParseFloat(e.Value, 64); err == nil && v > 0 {
                if v > amount {
                    amount = v
                }
                if e.Confidence > amountMaxConfidence {
                    amountMaxConfidence = e.Confidence
                }
            }
        case "category":
            // Try resolve category by name from DB instead of static map
            if categoryName == "" {
                categoryName = e.Value
            }
            if e.Confidence > categoryConfidence {
                categoryConfidence = e.Confidence
            }
        case "period":
            period = e.Value
            if e.Confidence > periodConfidence {
                periodConfidence = e.Confidence
            }
        case "title":
            name = e.Value
        }
    }

    // Parse alert threshold phrases like "cảnh báo 70%"
    lower := strings.ToLower(req.Message)
    if strings.Contains(lower, "cảnh báo") && strings.Contains(lower, "%") {
        // extract number before %
        for i := 0; i < len(lower); i++ {
            if lower[i] >= '0' && lower[i] <= '9' {
                j := i
                for j < len(lower) && ((lower[j] >= '0' && lower[j] <= '9') || lower[j] == '.') {
                    j++
                }
                if j < len(lower) && lower[j] == '%' {
                    if v, err := strconv.ParseFloat(lower[i:j], 64); err == nil && v > 0 {
                        alertThreshold = &v
                    }
                    break
                }
            }
        }
    }

    if amount <= 0 {
        resp.Response = "Vui lòng cung cấp số tiền ngân sách để tạo ngân sách."
        return
    }

    // Confidence threshold for confirmation
    if amountMaxConfidence > 0 && amountMaxConfidence < 0.6 {
        resp.Response = fmt.Sprintf("Bạn muốn tạo ngân sách %s VND? Vui lòng xác nhận hoặc cung cấp lại số tiền.", formatCurrency(amount))
        return
    }

    if name == "" {
        name = "Ngân sách hàng tháng"
    }
    if period == "" {
        period = "monthly"
    }
    // Confirmation prompts for low-confidence entities
    if categoryID == nil && categoryConfidence > 0 && categoryConfidence < 0.6 {
        resp.Response = "Bạn muốn tạo ngân sách cho danh mục nào? Vui lòng xác nhận danh mục."
        return
    }
    if periodConfidence > 0 && periodConfidence < 0.6 {
        resp.Response = "Kỳ ngân sách bạn mong muốn là hàng tuần, hàng tháng hay hàng năm?"
        return
    }

    // Map Vietnamese periods
    periodMap := map[string]string{
        "hàng tuần": "weekly",
        "hàng tháng": "monthly", 
        "hàng năm": "yearly",
    }
    if mapped, exists := periodMap[period]; exists {
        period = mapped
    }

    // Resolve categoryID by name if provided
    if categoryID == nil && categoryName != "" {
        var cat models.Category
        // Try exact and normalized matching
        if err := s.db.Where("LOWER(name) = ?", strings.ToLower(categoryName)).First(&cat).Error; err != nil {
            // fallback common name variants
            if strings.Contains(strings.ToLower(categoryName), "an uong") || strings.Contains(strings.ToLower(categoryName), "ăn uống") {
                s.db.Where("LOWER(name) = ?", strings.ToLower("Ăn uống")).First(&cat)
            }
        }
        if cat.ID != 0 {
            id := cat.ID
            categoryID = &id
        }
    }

    // Set default dates for monthly budget
    now := time.Now()
    startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    endDate := startDate.AddDate(0, 1, -1)

    createReq := &models.BudgetCreateRequest{
        CategoryID:     categoryID,
        Name:           name,
        Amount:         amount,
        Period:         period,
        StartDate:      startDate,
        EndDate:        endDate,
        AlertThreshold: func() float64 { if alertThreshold != nil { return *alertThreshold } ; return 80.0 }(),
    }

    budgetService := NewBudgetService()
    if budget, err := budgetService.CreateBudget(req.UserID, createReq); err != nil {
        log.Printf("AI chat create_budget: failed to create budget: %v", err)
        resp.Response = "Không thể tạo ngân sách. Vui lòng thử lại."
    } else {
        resp.Response = fmt.Sprintf("Đã tạo ngân sách '%s' với số tiền %s VND cho kỳ %s.", 
            budget.Name, formatCurrency(budget.Amount), period)
    }
}

func (s *AIService) handleListBudgets(req *models.ChatRequest, resp *models.ChatResponse) {
    budgetService := NewBudgetService()
    budgets, err := budgetService.GetBudgets(req.UserID)
    if err != nil {
        log.Printf("AI chat list_budgets: failed to get budgets: %v", err)
        resp.Response = "Không thể lấy danh sách ngân sách. Vui lòng thử lại."
        return
    }

    if len(budgets) == 0 {
        resp.Response = "Bạn chưa có ngân sách nào. Hãy tạo ngân sách đầu tiên!"
        return
    }

    resp.Response = fmt.Sprintf("Bạn có %d ngân sách:\n", len(budgets))
    for i, budget := range budgets {
        if i >= 5 { // Limit to 5 budgets for readability
            resp.Response += "..."
            break
        }
        categoryName := "Tất cả"
        if budget.Category != nil {
            categoryName = budget.Category.Name
        }
        resp.Response += fmt.Sprintf("- %s (%s): %s/%s VND (%.1f%%)\n", 
            budget.Name, categoryName, formatCurrency(budget.SpentAmount), 
            formatCurrency(budget.Amount), budget.UsagePercentage)
    }
}

func (s *AIService) handleUpdateBudget(req *models.ChatRequest, resp *models.ChatResponse) {
    var amount float64
    var budgetID uint64

    // Extract entities
    for _, e := range resp.Entities {
        switch e.Type {
        case "amount":
            if v, err := strconv.ParseFloat(e.Value, 64); err == nil && v > 0 {
                amount = v
            }
        case "budget_id":
            if v, err := strconv.ParseUint(e.Value, 10, 64); err == nil {
                budgetID = v
            }
        }
    }

    if amount <= 0 {
        resp.Response = "Vui lòng cung cấp số tiền để cập nhật ngân sách."
        return
    }

    // If no budget ID specified, get the first budget
    if budgetID == 0 {
        budgetService := NewBudgetService()
        budgets, err := budgetService.GetBudgets(req.UserID)
        if err != nil || len(budgets) == 0 {
            resp.Response = "Không tìm thấy ngân sách nào để cập nhật."
            return
        }
        budgetID = budgets[0].ID
    }

    // For now, just show current budget status
    budgetService := NewBudgetService()
    budgets, err := budgetService.GetBudgets(req.UserID)
    if err != nil {
        log.Printf("AI chat update_budget: failed to get budgets: %v", err)
        resp.Response = "Không thể cập nhật ngân sách. Vui lòng thử lại."
        return
    }

    var budget *models.Budget
    for i := range budgets {
        if budgets[i].ID == budgetID {
            budget = &budgets[i]
            break
        }
    }

    if budget == nil {
        resp.Response = "Không tìm thấy ngân sách để cập nhật."
        return
    }

    categoryName := "Tất cả"
    if budget.Category != nil {
        categoryName = budget.Category.Name
    }

    resp.Response = fmt.Sprintf("Ngân sách '%s' (%s): Đã chi %s/%s VND (%.1f%%). Còn lại %s VND.", 
        budget.Name, categoryName, formatCurrency(budget.SpentAmount), 
        formatCurrency(budget.Amount), budget.UsagePercentage, formatCurrency(budget.RemainingAmount))
}
