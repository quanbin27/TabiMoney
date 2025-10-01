package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
}

func NewAIService(cfg *config.Config) *AIService {
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		aiURL = "http://localhost:8001"
	}
	return &AIService{
		config:       cfg,
		db:           database.GetDB(),
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		aiServiceURL: aiURL + "/api/v1",
	}
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

	// Analyze patterns
	patterns := s.analyzeSpendingPatterns(transactions, req.Granularity)
	insights := s.generateSpendingInsights(patterns)
	recommendations := s.generateSpendingRecommendations(patterns)

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

	// Save chat message
	chatMessage := &models.ChatMessage{
		UserID:      req.UserID,
		Message:     req.Message,
		Response:    response.Response,
		Intent:      response.Intent,
		IsProcessed: true,
	}
	s.db.Create(chatMessage)

    return &response, nil
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
		"Consider increasing your monthly savings to reach your goal faster.",
		"Review your spending to identify areas where you can save more.",
	}
	return recommendations
}

func (s *AIService) identifyGoalRiskFactors(goal models.FinancialGoal, progress float64) []string {
	riskFactors := []string{
		"Low progress may indicate insufficient savings rate.",
		"Market volatility could affect investment goals.",
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
