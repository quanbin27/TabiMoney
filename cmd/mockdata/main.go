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

	// Find income category
	var incomeCategory models.Category
	db.Where("name = ? OR name_en = ?", "Thu nhập", "Income").First(&incomeCategory)

	// Set random seed
	rand.Seed(*seed)

	fmt.Printf("Generating realistic Vietnamese user data: %d transactions for user %d over %d months...\n", *numTxns, *userID, *months)

	// Generate transactions
	transactions := generateRealisticTransactions(*userID, categories, incomeCategory, *numTxns, *months, *includeAnomalies)

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

func generateRealisticTransactions(userID uint64, categories []models.Category, incomeCategory models.Category, count, months int, includeAnomalies bool) []models.Transaction {
	transactions := make([]models.Transaction, 0, count)

	now := time.Now()
	startDate := now.AddDate(0, -months, 0)

	// Find categories by name
	categoryMap := make(map[string]models.Category)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat
		if cat.NameEn != "" {
			categoryMap[cat.NameEn] = cat
		}
	}

	// Monthly salary (realistic for Vietnamese office worker)
	monthlySalary := 18000000.0 // 18M VND/month

	// Generate income transactions (monthly salary on 5th of each month)
	for i := 0; i < months; i++ {
		date := startDate.AddDate(0, i, 0)
		// Salary on 5th of month
		salaryDate := time.Date(date.Year(), date.Month(), 5, 9, 0, 0, 0, time.UTC)
		
		// Sometimes late salary (5th-10th)
		if rand.Float64() < 0.2 {
			salaryDate = salaryDate.AddDate(0, 0, rand.Intn(6))
		}

		transactions = append(transactions, models.Transaction{
			UserID:           userID,
			CategoryID:       incomeCategory.ID,
			Amount:           monthlySalary,
			Description:      fmt.Sprintf("Lương tháng %s", salaryDate.Format("01/2006")),
			TransactionType:  "income",
			TransactionDate:  salaryDate,
			TransactionTime:  timePtr(time.Date(salaryDate.Year(), salaryDate.Month(), salaryDate.Day(), 9, 0, 0, 0, time.UTC)),
			Location:         "Ngân hàng",
			Tags:             `["salary", "income", "lương"]`,
			Metadata:         `{"source": "employer", "type": "monthly", "payment_method": "bank_transfer"}`,
			IsRecurring:      true,
			RecurringPattern: "monthly",
		})

		// Bonus sometimes (end of quarter)
		if date.Month()%3 == 0 && rand.Float64() < 0.3 {
			bonusDate := time.Date(date.Year(), date.Month(), 28, 15, 0, 0, 0, time.UTC)
			transactions = append(transactions, models.Transaction{
				UserID:          userID,
				CategoryID:      incomeCategory.ID,
				Amount:          monthlySalary * 0.5, // 50% bonus
				Description:     fmt.Sprintf("Thưởng quý %d", (int(date.Month())-1)/3+1),
				TransactionType: "income",
				TransactionDate: bonusDate,
				TransactionTime: timePtr(bonusDate),
				Location:        "Ngân hàng",
				Tags:            `["bonus", "thưởng"]`,
				Metadata:        `{"source": "employer", "type": "bonus"}`,
			})
		}
	}

	// Generate realistic expense patterns
	expenseCount := count - months - (months/3) // Reserve space for bonuses
	anomalyCount := 0
	if includeAnomalies {
		anomalyCount = expenseCount / 15 // ~6-7% anomalies (more realistic)
	}

	// Daily patterns
	dayCount := int(now.Sub(startDate).Hours() / 24)
	normalExpenseCount := expenseCount - anomalyCount

	// Generate recurring bills (monthly)
	for i := 0; i < months; i++ {
		monthStart := startDate.AddDate(0, i, 0)
		
		// Electricity bill (15th of month, 500k-1.5M)
		if cat, ok := categoryMap["Khác"]; ok {
			transactions = append(transactions, models.Transaction{
				UserID:          userID,
				CategoryID:      cat.ID,
				Amount:          500000 + rand.Float64()*1000000,
				Description:     "Tiền điện",
				TransactionType: "expense",
				TransactionDate: time.Date(monthStart.Year(), monthStart.Month(), 15, 18, 0, 0, 0, time.UTC),
				TransactionTime: timePtr(time.Date(monthStart.Year(), monthStart.Month(), 15, 18, 0, 0, 0, time.UTC)),
				Location:        "EVN",
				Tags:            `["tiện ích", "điện"]`,
				Metadata:        `{"payment_method": "bank_transfer", "recurring": true}`,
				IsRecurring:     true,
				RecurringPattern: "monthly",
			})
		}

		// Water bill (20th of month, 100k-300k)
		if cat, ok := categoryMap["Khác"]; ok {
			transactions = append(transactions, models.Transaction{
				UserID:          userID,
				CategoryID:      cat.ID,
				Amount:          100000 + rand.Float64()*200000,
				Description:     "Tiền nước",
				TransactionType: "expense",
				TransactionDate: time.Date(monthStart.Year(), monthStart.Month(), 20, 19, 0, 0, 0, time.UTC),
				TransactionTime: timePtr(time.Date(monthStart.Year(), monthStart.Month(), 20, 19, 0, 0, 0, time.UTC)),
				Location:        "Cấp nước",
				Tags:            `["tiện ích", "nước"]`,
				Metadata:        `{"payment_method": "bank_transfer", "recurring": true}`,
				IsRecurring:     true,
				RecurringPattern: "monthly",
			})
		}

		// Internet bill (1st of month, 200k-500k)
		if cat, ok := categoryMap["Khác"]; ok {
			transactions = append(transactions, models.Transaction{
				UserID:          userID,
				CategoryID:      cat.ID,
				Amount:          200000 + rand.Float64()*300000,
				Description:     "Tiền internet",
				TransactionType: "expense",
				TransactionDate: time.Date(monthStart.Year(), monthStart.Month(), 1, 20, 0, 0, 0, time.UTC),
				TransactionTime: timePtr(time.Date(monthStart.Year(), monthStart.Month(), 1, 20, 0, 0, 0, time.UTC)),
				Location:        "VNPT/FPT",
				Tags:            `["tiện ích", "internet"]`,
				Metadata:        `{"payment_method": "e_wallet", "recurring": true}`,
				IsRecurring:     true,
				RecurringPattern: "monthly",
			})
		}
	}

	// Generate daily expenses (food, transport, etc.)
	dailyExpenseCount := normalExpenseCount - (months * 3) // Subtract bills
	expensesPerDay := float64(dailyExpenseCount) / float64(dayCount)

	for day := 0; day < dayCount; day++ {
		currentDate := startDate.AddDate(0, 0, day)
		weekday := currentDate.Weekday()

		// Skip some days (not every day has expenses)
		if rand.Float64() > expensesPerDay {
			continue
		}

		// Morning coffee (Mon-Fri, 7-8 AM, 25k-50k)
		if weekday >= time.Monday && weekday <= time.Friday && rand.Float64() < 0.7 {
			if cat, ok := categoryMap["Ăn uống"]; ok {
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          25000 + rand.Float64()*25000,
					Description:     getRandomFood([]string{"Cà phê sáng", "Trà đá", "Cà phê đen", "Nước cam"}),
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 7+rand.Intn(2), rand.Intn(30), 0, 0, time.UTC)),
					Location:        getRandomLocation([]string{"Highlands Coffee", "Starbucks", "Cà phê vỉa hè", "Circle K"}),
					Tags:            `["cà phê", "sáng"]`,
					Metadata:        `{"payment_method": "e_wallet"}`,
				})
			}
		}

		// Lunch (Mon-Fri, 11:30-13:00, 50k-150k)
		if weekday >= time.Monday && weekday <= time.Friday && rand.Float64() < 0.9 {
			if cat, ok := categoryMap["Ăn uống"]; ok {
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          50000 + rand.Float64()*100000,
					Description:     getRandomFood([]string{"Cơm trưa", "Bún chả", "Phở", "Bánh mì", "Cơm tấm", "Bún bò"}),
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 11+rand.Intn(2), 30+rand.Intn(30), 0, 0, time.UTC)),
					Location:        getRandomLocation([]string{"Quán cơm", "Nhà hàng", "Canteen công ty", "Food court"}),
					Tags:            `["ăn trưa", "cơm"]`,
					Metadata:        `{"payment_method": "cash"}`,
				})
			}
		}

		// Dinner (daily, 18:00-20:00, varies)
		if rand.Float64() < 0.8 {
			if cat, ok := categoryMap["Ăn uống"]; ok {
				var amount float64
				var desc string
				if weekday >= time.Saturday && rand.Float64() < 0.5 {
					// Weekend: Grab Food or restaurant (more expensive)
					amount = 100000 + rand.Float64()*200000
					desc = getRandomFood([]string{"Grab Food", "Now Food", "Baemin", "Nhà hàng", "Buffet"})
				} else {
					// Weekday: home cooking or cheap food
					amount = 30000 + rand.Float64()*70000
					desc = getRandomFood([]string{"Cơm tối", "Chợ", "Siêu thị", "Bánh mì", "Bún"})
				}
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          amount,
					Description:     desc,
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 18+rand.Intn(2), rand.Intn(60), 0, 0, time.UTC)),
					Location:        getRandomLocation([]string{"Nhà", "Grab", "Siêu thị", "Chợ", "Nhà hàng"}),
					Tags:            `["ăn tối"]`,
					Metadata:        `{"payment_method": "e_wallet"}`,
				})
			}
		}

		// Transport (Mon-Fri, morning and evening)
		if weekday >= time.Monday && weekday <= time.Friday {
			// Morning commute (7-8 AM)
			if rand.Float64() < 0.8 {
				if cat, ok := categoryMap["Giao thông"]; ok {
					transactions = append(transactions, models.Transaction{
						UserID:          userID,
						CategoryID:      cat.ID,
						Amount:          15000 + rand.Float64()*35000, // Grab/Gojek 15k-50k
						Description:     getRandomTransport([]string{"Grab đi làm", "Gojek", "Xe buýt", "Xe máy"}),
						TransactionType: "expense",
						TransactionDate: currentDate,
						TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 7+rand.Intn(2), rand.Intn(30), 0, 0, time.UTC)),
						Location:        "Đường phố",
						Tags:            `["đi làm"]`,
						Metadata:        `{"payment_method": "e_wallet"}`,
					})
				}
			}

			// Evening commute (17-18 PM)
			if rand.Float64() < 0.7 {
				if cat, ok := categoryMap["Giao thông"]; ok {
					transactions = append(transactions, models.Transaction{
						UserID:          userID,
						CategoryID:      cat.ID,
						Amount:          15000 + rand.Float64()*35000,
						Description:     getRandomTransport([]string{"Grab về nhà", "Gojek", "Xe buýt"}),
						TransactionType: "expense",
						TransactionDate: currentDate,
						TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 17+rand.Intn(2), rand.Intn(30), 0, 0, time.UTC)),
						Location:        "Đường phố",
						Tags:            `["về nhà"]`,
						Metadata:        `{"payment_method": "e_wallet"}`,
					})
				}
			}
		}

		// Gas (weekly, 200k-500k)
		if weekday == time.Sunday && rand.Float64() < 0.3 {
			if cat, ok := categoryMap["Giao thông"]; ok {
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          200000 + rand.Float64()*300000,
					Description:     "Đổ xăng",
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 10+rand.Intn(4), rand.Intn(60), 0, 0, time.UTC)),
					Location:        "Cây xăng",
					Tags:            `["xăng", "xe máy"]`,
					Metadata:        `{"payment_method": "cash"}`,
				})
			}
		}

		// Shopping (occasional, weekends more likely)
		if (weekday == time.Saturday || weekday == time.Sunday) && rand.Float64() < 0.2 {
			if cat, ok := categoryMap["Mua sắm"]; ok {
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          200000 + rand.Float64()*800000,
					Description:     getRandomShopping([]string{"Quần áo", "Đồ dùng nhà", "Mỹ phẩm", "Điện thoại phụ kiện", "Sách"}),
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 14+rand.Intn(4), rand.Intn(60), 0, 0, time.UTC)),
					Location:        getRandomLocation([]string{"Vincom", "Big C", "AEON Mall", "Lotte", "Online"}),
					Tags:            `["mua sắm"]`,
					Metadata:        `{"payment_method": "card"}`,
				})
			}
		}

		// Entertainment (weekends, 100k-500k)
		if (weekday == time.Saturday || weekday == time.Sunday) && rand.Float64() < 0.3 {
			if cat, ok := categoryMap["Giải trí"]; ok {
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          100000 + rand.Float64()*400000,
					Description:     getRandomEntertainment([]string{"Xem phim", "Cà phê", "Karaoke", "Game", "Công viên"}),
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 15+rand.Intn(6), rand.Intn(60), 0, 0, time.UTC)),
					Location:        getRandomLocation([]string{"CGV", "Lotte Cinema", "Karaoke", "Game center"}),
					Tags:            `["giải trí"]`,
					Metadata:        `{"payment_method": "card"}`,
				})
			}
		}

		// Healthcare (occasional, 200k-2M)
		if rand.Float64() < 0.05 { // ~5% chance per day
			if cat, ok := categoryMap["Y tế"]; ok {
				transactions = append(transactions, models.Transaction{
					UserID:          userID,
					CategoryID:      cat.ID,
					Amount:          200000 + rand.Float64()*1800000,
					Description:     getRandomHealthcare([]string{"Khám bệnh", "Thuốc", "Khám sức khỏe", "Nha khoa"}),
					TransactionType: "expense",
					TransactionDate: currentDate,
					TransactionTime: timePtr(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 9+rand.Intn(4), rand.Intn(60), 0, 0, time.UTC)),
					Location:        getRandomLocation([]string{"Bệnh viện", "Phòng khám", "Nhà thuốc"}),
					Tags:            `["y tế"]`,
					Metadata:        `{"payment_method": "cash"}`,
				})
			}
		}
	}

	// Generate anomalies (if enabled)
	if includeAnomalies {
		for i := 0; i < anomalyCount; i++ {
			daysOffset := rand.Intn(dayCount)
			date := startDate.AddDate(0, 0, daysOffset)

			// Anomaly types: very high amount, unusual category, odd time
			anomalyType := rand.Intn(3)
			var cat models.Category
			var amount float64
			var desc string

			switch anomalyType {
			case 0: // Very high amount
				cat = categoryMap["Mua sắm"]
				amount = 5000000 + rand.Float64()*10000000 // 5-15M
				desc = "Mua laptop/điện thoại"
			case 1: // Unusual category
				cat = categoryMap["Y tế"]
				amount = 3000000 + rand.Float64()*7000000 // 3-10M
				desc = "Khám bệnh/Phẫu thuật"
			case 2: // Odd time + high amount
				cat = categoryMap["Giải trí"]
				amount = 2000000 + rand.Float64()*5000000 // 2-7M
				desc = "Du lịch/Celebration"
			}

			transactions = append(transactions, models.Transaction{
				UserID:          userID,
				CategoryID:      cat.ID,
				Amount:          amount,
				Description:     desc,
				TransactionType: "expense",
				TransactionDate: date,
				TransactionTime: timePtr(time.Date(date.Year(), date.Month(), date.Day(), rand.Intn(24), rand.Intn(60), 0, 0, time.UTC)),
				Location:        "Various",
				Tags:            `["anomaly", "unusual"]`,
				Metadata:        `{"payment_method": "card", "flagged": true}`,
			})
		}
	}

	return transactions
}

func getRandomFood(options []string) string {
	return options[rand.Intn(len(options))]
}

func getRandomTransport(options []string) string {
	return options[rand.Intn(len(options))]
}

func getRandomShopping(options []string) string {
	return options[rand.Intn(len(options))]
}

func getRandomEntertainment(options []string) string {
	return options[rand.Intn(len(options))]
}

func getRandomHealthcare(options []string) string {
	return options[rand.Intn(len(options))]
}

func getRandomLocation(options []string) string {
	return options[rand.Intn(len(options))]
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
	if totalIncome > 0 {
		fmt.Printf("  Savings Rate:   %.1f%%\n", (totalIncome-totalExpense)/totalIncome*100)
	} else {
		fmt.Printf("  Savings Rate:   N/A (No income)\n")
	}

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
