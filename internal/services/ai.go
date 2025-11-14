package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"log"
	"net/http"
	"os"
	"sort"
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
    // Prefer AI-service ML anomaly detection via HTTP
    ctx := context.Background()
    payload := map[string]interface{}{
        "user_id":    req.UserID,
        "start_date": req.StartDate.Format("2006-01-02"),
        "end_date":   req.EndDate.Format("2006-01-02"),
        "threshold":  req.Threshold,
    }
    b, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    httpReq, err := http.NewRequestWithContext(ctx, "POST", s.aiServiceURL+"/anomaly/detect", bytes.NewBuffer(b))
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

    var aiResp struct {
        Anomalies      []struct {
            TransactionID   uint64  `json:"transaction_id"`
            Amount          float64 `json:"amount"`
            CategoryName    string  `json:"category_name"`
            AnomalyScore    float64 `json:"anomaly_score"`
            AnomalyType     string  `json:"anomaly_type"`
            Description     string  `json:"description"`
            TransactionDate string  `json:"transaction_date"`
        } `json:"anomalies"`
        TotalAnomalies int     `json:"total_anomalies"`
        DetectionScore float64 `json:"detection_score"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
        return nil, fmt.Errorf("failed to decode AI anomaly response: %w", err)
    }

    // Map to domain model
    anomalies := make([]models.Anomaly, 0, len(aiResp.Anomalies))
    for _, a := range aiResp.Anomalies {
        // Parse date string to time.Time
        tdate, _ := time.Parse("2006-01-02", a.TransactionDate)
        anomalies = append(anomalies, models.Anomaly{
            TransactionID:   a.TransactionID,
            Amount:          a.Amount,
            CategoryName:    a.CategoryName,
            AnomalyScore:    a.AnomalyScore,
            AnomalyType:     a.AnomalyType,
            Description:     a.Description,
            TransactionDate: tdate,
        })
    }

    out := &models.AnomalyDetectionResponse{
        UserID:         req.UserID,
        Anomalies:      anomalies,
        TotalAnomalies: len(anomalies),
        DetectionScore: aiResp.DetectionScore,
        GeneratedAt:    time.Now(),
    }

    // Persist analysis
    analysis := &models.AIAnalysis{
        UserID:          req.UserID,
        AnalysisType:    "anomaly_detection",
        Data:            s.marshalToJSON(out),
        ConfidenceScore: out.DetectionScore,
        ModelVersion:    "isolation_forest",
    }
    _ = s.db.Create(analysis).Error

    // Trigger notifications for anomalies (throttle + dedupe)
    if len(anomalies) > 0 {
        // 1) Keep only strong anomalies (score >= threshold or >= 0.8 default)
        minScore := req.Threshold
        if minScore <= 0 {
            minScore = 0.8
        }
        filtered := make([]models.Anomaly, 0, len(anomalies))
        for _, a := range anomalies {
            if a.AnomalyScore >= minScore {
                filtered = append(filtered, a)
            }
        }
        // 2) Sort by score desc and cap to top N per run
        sort.Slice(filtered, func(i, j int) bool { return filtered[i].AnomalyScore > filtered[j].AnomalyScore })
        if len(filtered) > 3 {
            filtered = filtered[:3]
        }
        // 3) Best-effort dedupe: skip if a recent notification exists for same transaction
        dispatcher := NewNotificationDispatcher()
        oneWeekAgo := time.Now().AddDate(0, 0, -7)
        for _, anomaly := range filtered {
            var recent models.Notification
            // Requires metadata to contain transaction_id; best-effort LIKE query
            if err := s.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?",
                req.UserID, "warning", oneWeekAgo, fmt.Sprintf("%%\"transaction_id\":%d%%", anomaly.TransactionID)).
                First(&recent).Error; err == nil {
                continue // similar notification exists recently
            }
            if err := dispatcher.TriggerAnomalyAlert(req.UserID, &anomaly); err != nil {
                log.Printf("Failed to trigger anomaly alert: %v", err)
            }
        }
    }

    return out, nil
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

    // Use shorter timeout for suggest to keep UX snappy
    client := &http.Client{Timeout: 8 * time.Second}
    resp, err := client.Do(httpReq)
    if err != nil {
        // Fallback: rank by recent usage + token match
        fallback := &models.CategorySuggestionResponse{
            UserID: req.UserID,
            Description: req.Description,
            Amount: req.Amount,
            Suggestions: []models.CategorySuggestion{},
            ConfidenceScore: 0.0,
            GeneratedAt: time.Now(),
        }

        // Build frequency map from recent transactions
        freq := make(map[uint64]int)
        for _, t := range recentTransactions {
            freq[t.CategoryID]++
        }

        // Tokenize description (basic)
        desc := strings.ToLower(req.Description)
        tokens := strings.FieldsFunc(desc, func(r rune) bool { return r == ' ' || r == ',' || r == '.' || r == '-' || r == '_' })
        tokenSet := make(map[string]struct{})
        for _, tk := range tokens { if tk != "" { tokenSet[tk] = struct{}{} } }

        type scored struct {
            cat models.Category
            score float64
            reason string
        }
        scoredList := []scored{}
        for _, c := range categories {
            sscore := 0.0
            reason := []string{}
            // Frequency weight
            if f, ok := freq[c.ID]; ok && f > 0 {
                sscore += float64(f) * 0.1
                reason = append(reason, "Thường xuyên sử dụng")
            }
            // Name token match
            name := strings.ToLower(c.Name)
            nameTokens := strings.FieldsFunc(name, func(r rune) bool { return r == ' ' || r == ',' || r == '.' || r == '-' || r == '_' })
            matched := 0
            for _, nt := range nameTokens {
                if _, ok := tokenSet[nt]; ok && nt != "" {
                    matched++
                }
            }
            if matched > 0 {
                sscore += float64(matched) * 0.3
                reason = append(reason, "Khớp mô tả")
            }
            if sscore > 0 {
                scoredList = append(scoredList, scored{cat: c, score: sscore, reason: strings.Join(reason, "; ")})
            }
        }

        sort.Slice(scoredList, func(i, j int) bool { return scoredList[i].score > scoredList[j].score })
        // Pick top 3
        k := 3
        if len(scoredList) < k { k = len(scoredList) }
        for i := 0; i < k; i++ {
            c := scoredList[i]
            fallback.Suggestions = append(fallback.Suggestions, models.CategorySuggestion{
                CategoryID: c.cat.ID,
                CategoryName: c.cat.Name,
                ConfidenceScore: math.Min(0.85, 0.4 + c.score*0.2),
                Reason: c.reason,
                IsUserCategory: c.cat.UserID != nil,
            })
        }
        if len(fallback.Suggestions) == 0 && len(categories) > 0 {
            // Default to most frequent category or first
            var best models.Category
            bestCount := -1
            for _, c := range categories {
                if cnt := freq[c.ID]; cnt > bestCount { best = c; bestCount = cnt }
            }
            if best.ID == 0 { best = categories[0] }
            fallback.Suggestions = append(fallback.Suggestions, models.CategorySuggestion{
                CategoryID: best.ID,
                CategoryName: best.Name,
                ConfidenceScore: 0.3,
                Reason: "Fallback: danh mục thường dùng",
                IsUserCategory: best.UserID != nil,
            })
        }
        return fallback, nil
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
    // Try cache first for this user and window
    ctx := context.Background()
    cacheKey := fmt.Sprintf("ai_analysis:%d:spending:%s:%s:%s", req.UserID, req.Granularity, req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))
    if cached, err := database.GetCache(ctx, cacheKey); err == nil && len(cached) > 0 {
        var resp models.SpendingPatternResponse
        if json.Unmarshal([]byte(cached), &resp) == nil {
            return &resp, nil
        }
    }

    // Compute fast local result
    var transactions []models.Transaction
    query := s.db.Where("user_id = ? AND transaction_date BETWEEN ? AND ?",
        req.UserID, req.StartDate, req.EndDate)
    if err := query.Preload("Category").Find(&transactions).Error; err != nil {
        return nil, fmt.Errorf("failed to get transactions: %w", err)
    }

    patterns := s.analyzeSpendingPatterns(transactions, req.Granularity)
    // Immediate rule-based insights as fallback
    rbInsights := s.generateSpendingInsights(patterns)
    rbRecs := s.generateSpendingRecommendations(patterns)
    immediate := &models.SpendingPatternResponse{
        UserID:          req.UserID,
        Patterns:        patterns,
        Insights:        rbInsights,
        Recommendations: rbRecs,
        GeneratedAt:     time.Now(),
    }

    // Return immediately, then spawn background to enrich via AI-service and cache
    go func() {
        defer func() { recover() }()
        insights, recommendations := s.fetchDynamicSpendingInsights(req.UserID, patterns)
        enriched := &models.SpendingPatternResponse{
            UserID:          req.UserID,
            Patterns:        patterns,
            Insights:        insights,
            Recommendations: recommendations,
            GeneratedAt:     time.Now(),
        }
        var bb []byte
        if b, err := json.Marshal(enriched); err == nil {
            bb = b
            // Cache for 1 hour
            _ = database.SetCache(ctx, cacheKey, bb, time.Hour)
        }
        // Persist AI analysis (best-effort)
        analysis := &models.AIAnalysis{
            UserID:          req.UserID,
            AnalysisType:    "spending_pattern",
            Data:            stringMust(bb),
            ConfidenceScore: 0.85,
            ModelVersion:    "custom",
        }
        _ = s.db.Create(analysis).Error
    }()

    return immediate, nil
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

    // Get user's financial data (both income and expense for recent 90 days)
    var transactions []models.Transaction
    windowStart := time.Now().AddDate(0, 0, -90)
    if err := s.db.Where("user_id = ? AND transaction_date >= ?", 
        req.UserID, windowStart).
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

    // AI service now handles transactions directly, so we just return the response
    // No need for backend handlers anymore

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
    needCategory := !has("category") && !has("category_id")  // Don't add category if category_id exists
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
    var categoryID uint64
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
        case "category_id":
            if v, err := strconv.ParseUint(e.Value, 10, 64); err == nil && v > 0 {
                categoryID = v
            }
        case "category":
            if categoryName == "" && categoryID == 0 {
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
    if categoryID > 0 {
        // Use category_id directly if available
        s.db.First(&category, categoryID)
    } else if categoryName != "" {
        // Fallback to category name lookup
        if err := s.db.Where("LOWER(name) = ?", strings.ToLower(categoryName)).First(&category).Error; err != nil {
            // normalize common Vietnamese variants
            lower := strings.ToLower(categoryName)
            if strings.Contains(lower, "an uong") || strings.Contains(lower, "ăn uống") {
                s.db.Where("LOWER(name) = ?", strings.ToLower("Ăn uống")).First(&category)
            }
        }
    }
    if category.ID == 0 {
        s.db.Where("name = ?", "Khác").First(&category)
    }


    today := time.Now().UTC().Format("2006-01-02")
    
    // Determine transaction type based on category
    transactionType := "expense"
    if category.ID == 8 { // Thu nhập category
        transactionType = "income"
    }
    
    createReq := &models.TransactionCreateRequest{
        CategoryID:      category.ID,
        Amount:          amount,
        Description:     req.Message,
        TransactionType: transactionType,
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
    // Upgraded heuristic:
    // - Per-category robust thresholds using Median Absolute Deviation (MAD)
    // - Fallback to global robust threshold when category has few samples
    // - Optional seasonal check: compare to previous month's category median
    var anomalies []models.Anomaly
    if len(transactions) == 0 {
        return anomalies
    }

    // Group amounts by category
    catAmounts := make(map[uint64][]float64)
    for _, t := range transactions {
        if t.TransactionType != "expense" { // focus on expenses
            continue
        }
        catAmounts[t.CategoryID] = append(catAmounts[t.CategoryID], t.Amount)
    }

    // Helper: compute median
    median := func(arr []float64) float64 {
        n := len(arr)
        if n == 0 {
            return 0
        }
        tmp := make([]float64, n)
        copy(tmp, arr)
        sort.Float64s(tmp)
        mid := n / 2
        if n%2 == 0 {
            return (tmp[mid-1] + tmp[mid]) / 2
        }
        return tmp[mid]
    }
    // Helper: compute MAD
    mad := func(arr []float64) float64 {
        n := len(arr)
        if n == 0 {
            return 0
        }
        m := median(arr)
        devs := make([]float64, 0, n)
        for _, v := range arr {
            if v >= m {
                devs = append(devs, v-m)
            } else {
                devs = append(devs, m-v)
            }
        }
        return median(devs)
    }

    // Pre-compute global robust stats as fallback
    var all []float64
    for _, t := range transactions {
        if t.TransactionType == "expense" {
            all = append(all, t.Amount)
        }
    }
    globalMed := median(all)
    globalMAD := mad(all)
    if globalMAD == 0 {
        // avoid divide-by-zero; use a small epsilon
        globalMAD = 1.0
    }

    // Seasonal baseline (previous 30 days) by category
    // Note: for simplicity, fetch from DB only when needed could be added later
    // Here we approximate by using half oldest vs newest split when enough points
    seasonalBaseline := make(map[uint64]float64)
    for catID, amounts := range catAmounts {
        if len(amounts) >= 6 { // need enough points to split
            tmp := make([]float64, len(amounts))
            copy(tmp, amounts)
            sort.Float64s(tmp)
            older := tmp[:len(tmp)/2]
            seasonalBaseline[catID] = median(older)
        } else {
            seasonalBaseline[catID] = 0
        }
    }

    // Evaluate each transaction
    for _, t := range transactions {
        if t.TransactionType != "expense" {
            continue
        }
        amounts := catAmounts[t.CategoryID]
        var m, mAD float64
        if len(amounts) >= 5 {
            m = median(amounts)
            mAD = mad(amounts)
        } else {
            m = globalMed
            mAD = globalMAD
        }
        if mAD == 0 {
            mAD = 1.0
        }

        // Robust z-score approximation: |x - median| / (1.4826 * MAD)
        z := (t.Amount - m) / (1.4826 * mAD)
        if z < 0 {
            z = -z
        }

        // Dynamic threshold: use provided threshold if >0 else default 3.5
        dyn := threshold
        if dyn <= 0 {
            dyn = 3.5
        }

        // Seasonal spike: amount is 2x seasonal median for category
        seasonalMed := seasonalBaseline[t.CategoryID]
        seasonalSpike := seasonalMed > 0 && t.Amount >= 2.0*seasonalMed

        if z >= dyn || seasonalSpike {
            score := z / dyn
            if seasonalSpike && score < 0.8 {
                score = 0.8
            }
            anomalies = append(anomalies, models.Anomaly{
                TransactionID:   t.ID,
                Amount:          t.Amount,
                CategoryName:    t.Category.Name,
                AnomalyScore:    math.Min(1.0, score),
                AnomalyType:     "amount",
                Description:     fmt.Sprintf("Chi tiêu bất thường so với mức trung vị (z=%.2f)", z),
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
    insights := []string{}
    if len(patterns) == 0 {
        return insights
    }

    // Compute totals and shares
    var total float64
    for _, p := range patterns {
        total += p.TotalAmount
    }
    // Top categories by spend
    top := make([]models.SpendingPattern, len(patterns))
    copy(top, patterns)
    sort.Slice(top, func(i, j int) bool { return top[i].TotalAmount > top[j].TotalAmount })

    // Insight 1: Top category share
    topShare := 0.0
    if total > 0 {
        topShare = (top[0].TotalAmount / total) * 100
    }
    if topShare >= 40 {
        insights = append(insights, fmt.Sprintf("Danh mục %s chiếm %.1f%% tổng chi, nên xem xét cắt giảm.", top[0].CategoryName, topShare))
    } else {
        insights = append(insights, fmt.Sprintf("Danh mục chi tiêu lớn nhất là %s (%.1f%%).", top[0].CategoryName, topShare))
    }

    // Insight 2: High average ticket categories
    hiAvg := []string{}
    for _, p := range top {
        if p.AverageAmount >= 1000000 && p.TransactionCount >= 3 { // ~1,000,000 VND
            hiAvg = append(hiAvg, p.CategoryName)
        }
        if len(hiAvg) >= 3 {
            break
        }
    }
    if len(hiAvg) > 0 {
        insights = append(insights, fmt.Sprintf("Các danh mục có giá trị giao dịch cao: %s.", strings.Join(hiAvg, ", ")))
    }

    // Insight 3: Long tail categories (small share but many transactions)
    longTail := []string{}
    for _, p := range patterns {
        if total == 0 { break }
        share := (p.TotalAmount / total) * 100
        if share < 5 && p.TransactionCount >= 5 {
            longTail = append(longTail, p.CategoryName)
        }
        if len(longTail) >= 3 {
            break
        }
    }
    if len(longTail) > 0 {
        insights = append(insights, fmt.Sprintf("Nhiều giao dịch nhỏ ở các danh mục: %s.", strings.Join(longTail, ", ")))
    }

    return insights
}

func (s *AIService) generateSpendingRecommendations(patterns []models.SpendingPattern) []string {
    recs := []string{}
    if len(patterns) == 0 {
        return []string{"Chưa đủ dữ liệu để gợi ý chi tiêu."}
    }
    // Total
    var total float64
    for _, p := range patterns { total += p.TotalAmount }
    // Recommend budget for top categories
    sorted := make([]models.SpendingPattern, len(patterns))
    copy(sorted, patterns)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i].TotalAmount > sorted[j].TotalAmount })
    topN := 3
    if len(sorted) < topN { topN = len(sorted) }
    for i := 0; i < topN; i++ {
        share := 0.0
        if total > 0 { share = (sorted[i].TotalAmount / total) * 100 }
        if share >= 20 {
            recs = append(recs, fmt.Sprintf("Thiết lập ngân sách riêng cho %s (%.1f%% tổng chi).", sorted[i].CategoryName, share))
        }
    }
    // Recommend review high average categories
    for _, p := range sorted {
        if p.AverageAmount >= 1000000 && p.TransactionCount >= 3 {
            recs = append(recs, fmt.Sprintf("Xem lại các khoản lớn ở %s, cân nhắc giảm tần suất/giá trị.", p.CategoryName))
        }
    }
    if len(recs) == 0 {
        recs = append(recs, "Theo dõi chi tiêu hàng tuần và đặt hạn mức cho 1-2 danh mục lớn nhất.")
    }
    return recs
}

func (s *AIService) calculateGoalProgress(goal models.FinancialGoal, transactions []models.Transaction) float64 {
    if goal.TargetAmount <= 0 {
        return 0
    }
    return goal.CurrentAmount / goal.TargetAmount
}

func (s *AIService) isGoalOnTrack(goal models.FinancialGoal, progress float64) bool {
    // If target date exists, compare required daily pace vs recent net savings pace
    if goal.TargetDate != nil && goal.TargetAmount > 0 {
        daysLeft := int(goal.TargetDate.Sub(time.Now()).Hours() / 24)
        if daysLeft <= 0 {
            return goal.CurrentAmount >= goal.TargetAmount
        }
        remaining := goal.TargetAmount - goal.CurrentAmount
        if remaining <= 0 {
            return true
        }
        requiredDaily := remaining / float64(daysLeft)

        // Estimate recent daily net savings from last 90 days
        var net90 float64
        var days int = 90
        // We don't have per-day aggregation here; approximate via monthly contribution: 10% target as baseline
        // Keep simple: assume user's net saving capacity equals 10% of target per month if no better data
        // This keeps logic conservative and avoids heavy queries here
        estimatedDaily := (goal.TargetAmount * 0.1) / 30.0
        if net90 > 0 { // placeholder if later we compute from transactions
            estimatedDaily = net90 / float64(days)
        }
        return estimatedDaily >= requiredDaily*0.9 // allow 10% slack
    }
    // Without target date, consider >50% as on track
    return progress >= 0.5
}

func (s *AIService) projectGoalCompletion(goal models.FinancialGoal, progress float64) *time.Time {
    if progress >= 1.0 || goal.TargetAmount <= 0 {
        t := time.Now()
        return &t
    }
    remaining := goal.TargetAmount - goal.CurrentAmount
    if remaining <= 0 {
        t := time.Now()
        return &t
    }
    // Estimate monthly saving capacity conservatively = 10% target/month
    monthlyContribution := goal.TargetAmount * 0.1
    if monthlyContribution <= 0 {
        return nil
    }
    monthsRemaining := int(math.Ceil(remaining / monthlyContribution))
    completionDate := time.Now().AddDate(0, monthsRemaining, 0)
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

// stringMust is a tiny helper to avoid panics in background goroutines
func stringMust(b []byte) string { if b == nil { return "{}" }; return string(b) }
