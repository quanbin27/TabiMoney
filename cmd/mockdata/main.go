package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

var (
	userID           = flag.Uint64("user", 15, "User ID to generate data for")
	numTxns          = flag.Int("count", 200, "Number of transactions to generate")
	months           = flag.Int("months", 6, "Number of months to spread transactions across")
	includeAnomalies = flag.Bool("anomalies", true, "Include anomaly transactions for testing")
	seed             = flag.Int64("seed", 42, "Random seed for reproducibility")
)

func main() {
	flag.Parse()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	db := database.GetDB()

	// Verify user exists
	var user models.User
	if err := db.First(&user, *userID).Error; err != nil {
		log.Fatalf("User %d not found: %v", *userID, err)
	}

	// Get system categories
	var categories []models.Category
	if err := db.Where("is_system = ? AND is_active = ?", true, true).Find(&categories).Error; err != nil {
		log.Fatalf("Failed to get categories: %v", err)
	}

	if len(categories) == 0 {
		log.Fatal("No system categories found. Please run database migrations first.")
	}

	// Find income category (usually ID 8 based on schema)
	var incomeCategory models.Category
	db.Where("name = ? OR name_en = ?", "Thu nhập", "Income").First(&incomeCategory)

	// Set random seed
	rand.Seed(*seed)

	fmt.Printf("Generating %d transactions for user %d over %d months...\n", *numTxns, *userID, *months)

	// Generate transactions
	transactions := generateTransactions(*userID, categories, incomeCategory, *numTxns, *months, *includeAnomalies)

	// Insert transactions in batches
	batchSize := 50
	for i := 0; i < len(transactions); i += batchSize {
		end := i + batchSize
		if end > len(transactions) {
			end = len(transactions)
		}
		batch := transactions[i:end]
		if err := db.Create(&batch).Error; err != nil {
			log.Fatalf("Failed to insert batch %d-%d: %v", i, end, err)
		}
		fmt.Printf("Inserted transactions %d-%d\n", i+1, end)
	}

	fmt.Printf("\n✅ Successfully generated %d transactions!\n", len(transactions))
	fmt.Printf("\nSummary:\n")
	printSummary(db, *userID)
}

func generateTransactions(userID uint64, categories []models.Category, incomeCategory models.Category, count, months int, includeAnomalies bool) []models.Transaction {
	transactions := make([]models.Transaction, 0, count)

	now := time.Now()
	startDate := now.AddDate(0, -months, 0)

	// Category distributions for expenses (more realistic spending patterns)
	expenseCategories := make([]models.Category, 0)
	for _, cat := range categories {
		if cat.Name != "Thu nhập" && cat.NameEn != "Income" {
			expenseCategories = append(expenseCategories, cat)
		}
	}

	// Typical spending ranges per category (in VND)
	categoryRanges := map[string]struct {
		min, max  float64
		frequency float64 // probability of this category
	}{
		"Ăn uống":        {min: 30000, max: 300000, frequency: 0.35},
		"Food & Dining":  {min: 30000, max: 300000, frequency: 0.35},
		"Giao thông":     {min: 20000, max: 500000, frequency: 0.20},
		"Transportation": {min: 20000, max: 500000, frequency: 0.20},
		"Mua sắm":        {min: 100000, max: 2000000, frequency: 0.15},
		"Shopping":       {min: 100000, max: 2000000, frequency: 0.15},
		"Giải trí":       {min: 50000, max: 500000, frequency: 0.10},
		"Entertainment":  {min: 50000, max: 500000, frequency: 0.10},
		"Y tế":           {min: 100000, max: 5000000, frequency: 0.05},
		"Healthcare":     {min: 100000, max: 5000000, frequency: 0.05},
		"Học tập":        {min: 50000, max: 1000000, frequency: 0.05},
		"Education":      {min: 50000, max: 1000000, frequency: 0.05},
		"Tiết kiệm":      {min: 500000, max: 5000000, frequency: 0.05},
		"Savings":        {min: 500000, max: 5000000, frequency: 0.05},
		"Khác":           {min: 20000, max: 500000, frequency: 0.05},
		"Other":          {min: 20000, max: 500000, frequency: 0.05},
	}

	// Generate income transactions (monthly salary)
	incomeCount := months
	for i := 0; i < incomeCount; i++ {
		date := startDate.AddDate(0, i, rand.Intn(3)) // Around start of month
		transactions = append(transactions, models.Transaction{
			UserID:           userID,
			CategoryID:       incomeCategory.ID,
			Amount:           15000000 + rand.Float64()*5000000, // 15-20M VND salary
			Description:      fmt.Sprintf("Lương tháng %s", date.Format("01/2006")),
			TransactionType:  "income",
			TransactionDate:  date,
			TransactionTime:  timePtr(time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)),
			Location:         "Công ty",
			Tags:             `["salary", "income"]`,
			Metadata:         `{"source": "employer", "type": "monthly"}`,
			IsRecurring:      true,
			RecurringPattern: "monthly",
		})
	}

	// Generate expense transactions
	expenseCount := count - incomeCount
	anomalyCount := 0
	if includeAnomalies {
		anomalyCount = expenseCount / 10 // 10% anomalies
	}

	for i := 0; i < expenseCount; i++ {
		isAnomaly := includeAnomalies && i < anomalyCount

		// Random date within range
		daysDiff := int(now.Sub(startDate).Hours() / 24)
		daysOffset := rand.Intn(daysDiff)
		date := startDate.AddDate(0, 0, daysOffset)

		// Select category based on frequency
		var category models.Category
		if isAnomaly {
			// Anomalies: random category, unusual amounts
			category = expenseCategories[rand.Intn(len(expenseCategories))]
		} else {
			// Normal: weighted random selection
			category = selectCategoryByFrequency(expenseCategories, categoryRanges)
		}

		// Generate amount
		var amount float64
		catRange, hasRange := categoryRanges[category.Name]
		if !hasRange {
			catRange, hasRange = categoryRanges[category.NameEn]
		}

		if isAnomaly {
			// Anomaly: very high or very low amounts, or unusual timing
			if rand.Float64() < 0.7 {
				// High amount anomaly (3-10x normal)
				if hasRange {
					amount = catRange.max * (3 + rand.Float64()*7)
				} else {
					amount = 5000000 + rand.Float64()*10000000 // 5-15M
				}
			} else {
				// Low amount anomaly (unusually small)
				amount = 5000 + rand.Float64()*10000
			}
		} else {
			// Normal amount
			if hasRange {
				amount = catRange.min + rand.Float64()*(catRange.max-catRange.min)
			} else {
				amount = 50000 + rand.Float64()*500000
			}
		}

		// Generate description
		description := generateDescription(category, isAnomaly)

		// Generate time (normal business hours vs anomalies at odd times)
		var txTime *time.Time
		if isAnomaly && rand.Float64() < 0.3 {
			// Anomaly: odd hours (2-5 AM)
			hour := 2 + rand.Intn(3)
			txTime = timePtr(time.Date(date.Year(), date.Month(), date.Day(), hour, rand.Intn(60), 0, 0, time.UTC))
		} else {
			// Normal: business hours
			hour := 8 + rand.Intn(12)
			txTime = timePtr(time.Date(date.Year(), date.Month(), date.Day(), hour, rand.Intn(60), 0, 0, time.UTC))
		}

		transactions = append(transactions, models.Transaction{
			UserID:          userID,
			CategoryID:      category.ID,
			Amount:          amount,
			Description:     description,
			TransactionType: "expense",
			TransactionDate: date,
			TransactionTime: txTime,
			Location:        generateLocation(category),
			Tags:            generateTags(category, isAnomaly),
			Metadata:        generateMetadata(category, isAnomaly),
			IsRecurring:     !isAnomaly && rand.Float64() < 0.1, // 10% recurring
			RecurringPattern: func() string {
				if rand.Float64() < 0.1 {
					patterns := []string{"daily", "weekly", "monthly"}
					return patterns[rand.Intn(len(patterns))]
				}
				return ""
			}(),
		})
	}

	return transactions
}

func selectCategoryByFrequency(categories []models.Category, ranges map[string]struct {
	min, max  float64
	frequency float64
}) models.Category {
	// Create weighted list
	weighted := make([]models.Category, 0)
	for _, cat := range categories {
		freq := 0.1 // default frequency
		if r, ok := ranges[cat.Name]; ok {
			freq = r.frequency
		} else if r, ok := ranges[cat.NameEn]; ok {
			freq = r.frequency
		}
		// Add category multiple times based on frequency
		count := int(freq * 100)
		for i := 0; i < count; i++ {
			weighted = append(weighted, cat)
		}
	}
	if len(weighted) == 0 {
		return categories[rand.Intn(len(categories))]
	}
	return weighted[rand.Intn(len(weighted))]
}

func generateDescription(category models.Category, isAnomaly bool) string {
	descriptions := map[string][]string{
		"Ăn uống":        {"Cơm trưa", "Cà phê sáng", "Nhà hàng", "Grab Food", "Siêu thị"},
		"Food & Dining":  {"Lunch", "Morning coffee", "Restaurant", "Food delivery", "Supermarket"},
		"Giao thông":     {"Xăng xe", "Grab", "Taxi", "Vé máy bay", "Bảo dưỡng xe"},
		"Transportation": {"Gas", "Ride share", "Taxi", "Flight ticket", "Car maintenance"},
		"Mua sắm":        {"Quần áo", "Điện thoại", "Laptop", "Đồ dùng nhà", "Quà tặng"},
		"Shopping":       {"Clothes", "Phone", "Laptop", "Home supplies", "Gift"},
		"Giải trí":       {"Xem phim", "Karaoke", "Du lịch", "Game", "Concert"},
		"Entertainment":  {"Movie", "Karaoke", "Travel", "Games", "Concert"},
		"Y tế":           {"Khám bệnh", "Thuốc", "Bảo hiểm", "Nha khoa", "Khám sức khỏe"},
		"Healthcare":     {"Doctor visit", "Medicine", "Insurance", "Dental", "Health check"},
		"Học tập":        {"Khóa học online", "Sách", "Tài liệu", "Workshop", "Certification"},
		"Education":      {"Online course", "Books", "Materials", "Workshop", "Certification"},
		"Tiết kiệm":      {"Tiết kiệm", "Đầu tư", "Quỹ dự phòng", "Tích lũy"},
		"Savings":        {"Savings", "Investment", "Emergency fund", "Accumulation"},
		"Khác":           {"Chi phí khác", "Phí dịch vụ", "Khác"},
		"Other":          {"Other expense", "Service fee", "Misc"},
	}

	options := descriptions[category.Name]
	if len(options) == 0 {
		options = descriptions[category.NameEn]
	}
	if len(options) == 0 {
		options = []string{"Transaction", "Payment", "Expense"}
	}

	desc := options[rand.Intn(len(options))]
	if isAnomaly {
		desc = "[ANOMALY] " + desc
	}
	return desc
}

func generateLocation(category models.Category) string {
	locations := []string{
		"Hà Nội", "TP. Hồ Chí Minh", "Đà Nẵng", "Hải Phòng", "Cần Thơ",
		"Online", "ATM", "Cửa hàng", "Siêu thị", "Nhà hàng",
	}
	return locations[rand.Intn(len(locations))]
}

func generateTags(category models.Category, isAnomaly bool) string {
	tags := []string{category.Name}
	if isAnomaly {
		tags = append(tags, "anomaly", "unusual")
	}
	if rand.Float64() < 0.3 {
		tags = append(tags, "urgent")
	}
	return fmt.Sprintf(`["%s"]`, tags[0])
}

func generateMetadata(category models.Category, isAnomaly bool) string {
	paymentMethods := []string{"cash", "card", "bank_transfer", "e_wallet"}
	method := paymentMethods[rand.Intn(len(paymentMethods))]

	meta := fmt.Sprintf(`{"payment_method": "%s"}`, method)
	if isAnomaly {
		meta = fmt.Sprintf(`{"payment_method": "%s", "flagged": true}`, method)
	}
	return meta
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func printSummary(db *gorm.DB, userID uint64) {
	var incomeStats struct {
		TotalAmount float64
		Count       int64
	}
	var expenseStats struct {
		TotalAmount float64
		Count       int64
	}

	db.Raw("SELECT COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as count FROM transactions WHERE user_id = ? AND transaction_type = ?", userID, "income").Scan(&incomeStats)
	db.Raw("SELECT COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as count FROM transactions WHERE user_id = ? AND transaction_type = ?", userID, "expense").Scan(&expenseStats)

	totalIncome := incomeStats.TotalAmount
	incomeCount := incomeStats.Count
	totalExpense := expenseStats.TotalAmount
	expenseCount := expenseStats.Count

	fmt.Printf("  Total Income:   %.0f VND (%d transactions)\n", totalIncome, incomeCount)
	fmt.Printf("  Total Expense:  %.0f VND (%d transactions)\n", totalExpense, expenseCount)
	fmt.Printf("  Net Amount:     %.0f VND\n", totalIncome-totalExpense)
	fmt.Printf("  Savings Rate:   %.1f%%\n", (totalIncome-totalExpense)/totalIncome*100)

	// Category breakdown
	var breakdown []struct {
		CategoryName string
		TotalAmount  float64
		Count        int64
	}
	db.Raw(`
		SELECT c.name as category_name, 
		       COALESCE(SUM(t.amount), 0) as total_amount,
		       COUNT(t.id) as count
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id AND t.user_id = ? AND t.transaction_type = 'expense'
		WHERE c.is_system = true
		GROUP BY c.id, c.name
		HAVING COUNT(t.id) > 0
		ORDER BY total_amount DESC
		LIMIT 5
	`, userID).Scan(&breakdown)

	fmt.Printf("\n  Top 5 Expense Categories:\n")
	for _, item := range breakdown {
		fmt.Printf("    - %s: %.0f VND (%d transactions)\n", item.CategoryName, item.TotalAmount, item.Count)
	}
}
