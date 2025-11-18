# CƠ CHẾ XỬ LÝ & THUẬT TOÁN HỆ THỐNG TABIMONEY

## 1. THUẬT TOÁN XỬ LÝ NLU (NATURAL LANGUAGE UNDERSTANDING)

### 1.1. Mô tả
Xử lý câu lệnh tự nhiên để trích xuất thông tin giao dịch (amount, category, description, date).

### 1.2. Luồng xử lý

```
Input: "tôi vừa ăn bún bò 50k"

1. Preprocessing
   - Chuẩn hóa text (lowercase, remove extra spaces)
   - Tokenize: ["tôi", "vừa", "ăn", "bún", "bò", "50k"]

2. Entity Extraction (sử dụng Gemini API)
   - Amount: "50k" → 50000
   - Category keywords: "ăn", "bún bò" → Food category
   - Description: "ăn bún bò"
   - Date: "vừa" → today (hoặc parse explicit date)

3. Category Matching
   - Tìm category phù hợp nhất:
     a. Tìm trong user's custom categories
     b. Tìm trong system categories
     c. So sánh với keywords/patterns
   - Tính confidence score dựa trên:
     - Keyword match
     - Historical patterns của user
     - AI confidence từ Gemini

4. Validation
   - Kiểm tra amount > 0
   - Kiểm tra category tồn tại
   - Kiểm tra date hợp lệ

5. Output
   {
     "category_id": 5,
     "amount": 50000,
     "description": "ăn bún bò",
     "transaction_date": "2024-01-15",
     "confidence": 0.95
   }
```

### 1.3. Ví dụ minh họa

**Input 1:** "mua cà phê 40k hôm qua"
- Amount: 40000
- Category: Ăn uống (từ keyword "cà phê")
- Description: "mua cà phê"
- Date: yesterday
- Confidence: 0.90

**Input 2:** "lương tháng 1 10 triệu"
- Amount: 10000000
- Category: Thu nhập
- Description: "lương tháng 1"
- Date: today (hoặc ngày nhận lương)
- Confidence: 0.85

**Input 3:** "đi taxi 50k từ nhà đến công ty"
- Amount: 50000
- Category: Giao thông
- Description: "đi taxi từ nhà đến công ty"
- Location: "nhà → công ty"
- Confidence: 0.88

---

## 2. THUẬT TOÁN PHÁT HIỆN BẤT THƯỜNG (ANOMALY DETECTION)

### 2.1. Mô tả
Phát hiện giao dịch bất thường dựa trên lịch sử chi tiêu của user.

### 2.2. Phương pháp: Statistical Analysis + Isolation Forest

#### Bước 1: Thu thập dữ liệu
```python
# Lấy lịch sử 30 ngày gần nhất
history = get_transactions(user_id, last_30_days)

# Tính toán các features:
features = {
    "amount": transaction.amount,
    "category_id": transaction.category_id,
    "day_of_week": transaction.date.weekday(),
    "hour": transaction.time.hour if transaction.time else None,
    "amount_by_category_mean": mean(history[category]),
    "amount_by_category_std": std(history[category])
}
```

#### Bước 2: Phát hiện bất thường

**Phương pháp 1: Z-Score (Statistical)**
```python
def detect_anomaly_zscore(transaction, history):
    category_history = filter_by_category(history, transaction.category_id)
    
    mean = calculate_mean(category_history.amounts)
    std = calculate_std(category_history.amounts)
    
    if std == 0:
        return False, 0.0
    
    z_score = abs((transaction.amount - mean) / std)
    
    # Nếu z-score > 2 (nằm ngoài 2 standard deviations)
    if z_score > 2:
        anomaly_score = min(z_score / 4, 1.0)  # Normalize to 0-1
        return True, anomaly_score
    
    return False, 0.0
```

**Phương pháp 2: Isolation Forest (ML)**
```python
from sklearn.ensemble import IsolationForest

def detect_anomaly_isolation_forest(transaction, history):
    # Chuẩn bị training data
    X = prepare_features(history)
    
    # Train model
    model = IsolationForest(contamination=0.1, random_state=42)
    model.fit(X)
    
    # Predict
    transaction_features = prepare_features([transaction])
    prediction = model.predict(transaction_features)
    score = model.score_samples(transaction_features)
    
    # prediction = -1: anomaly, 1: normal
    is_anomaly = (prediction[0] == -1)
    anomaly_score = 1 - normalize_score(score[0])
    
    return is_anomaly, anomaly_score
```

#### Bước 3: Kết hợp kết quả
```python
def detect_anomaly(transaction, history):
    # Nếu không đủ dữ liệu (< 10 transactions)
    if len(history) < 10:
        return False, 0.0
    
    # Phương pháp 1: Z-Score
    is_anomaly_1, score_1 = detect_anomaly_zscore(transaction, history)
    
    # Phương pháp 2: Isolation Forest
    is_anomaly_2, score_2 = detect_anomaly_isolation_forest(transaction, history)
    
    # Kết hợp (weighted average)
    final_score = (score_1 * 0.6) + (score_2 * 0.4)
    is_anomaly = (final_score > 0.7) or (is_anomaly_1 and is_anomaly_2)
    
    return is_anomaly, final_score
```

### 2.3. Ví dụ minh họa

**Scenario:** User thường chi 50k-100k cho ăn uống mỗi ngày

**Giao dịch bình thường:**
- Amount: 75000
- Category: Ăn uống
- Z-score: 0.5 (trong phạm vi bình thường)
- Anomaly: False

**Giao dịch bất thường:**
- Amount: 500000
- Category: Ăn uống
- Z-score: 4.2 (vượt quá 2 std)
- Anomaly: True, Score: 0.85
- Reason: "Amount significantly higher than average for this category"

---

## 3. THUẬT TOÁN DỰ ĐOÁN CHI TIÊU (EXPENSE PREDICTION)

### 3.1. Mô tả
Dự đoán chi tiêu tháng tới dựa trên lịch sử chi tiêu.

### 3.2. Phương pháp: Linear Regression + Time Series Analysis

#### Bước 1: Thu thập và chuẩn bị dữ liệu
```python
# Lấy dữ liệu 3-6 tháng gần nhất
history = get_transactions(user_id, months=6)

# Nhóm theo tháng và category
monthly_data = group_by_month_and_category(history)

# Features:
# - Month index (1, 2, 3, ...)
# - Total expense
# - Expense by category
# - Number of transactions
# - Average transaction amount
```

#### Bước 2: Phân tích xu hướng (Trend Analysis)
```python
def calculate_trend(monthly_data):
    """
    Tính xu hướng chi tiêu (tăng/giảm/ổn định)
    """
    amounts = [month['total'] for month in monthly_data]
    
    # Linear regression để tìm trend
    from sklearn.linear_model import LinearRegression
    
    X = [[i] for i in range(len(amounts))]
    y = amounts
    
    model = LinearRegression()
    model.fit(X, y)
    
    slope = model.coef_[0]
    
    if slope > 0:
        trend = "increasing"
        trend_percentage = (slope / np.mean(amounts)) * 100
    elif slope < 0:
        trend = "decreasing"
        trend_percentage = abs((slope / np.mean(amounts)) * 100)
    else:
        trend = "stable"
        trend_percentage = 0
    
    return trend, trend_percentage, model
```

#### Bước 3: Dự đoán
```python
def predict_expense(user_id, next_month_index):
    """
    Dự đoán chi tiêu tháng tới
    """
    history = get_transactions(user_id, months=6)
    monthly_data = group_by_month_and_category(history)
    
    # Dự đoán tổng chi tiêu
    trend, trend_pct, model = calculate_trend(monthly_data)
    
    predicted_total = model.predict([[next_month_index]])[0]
    
    # Dự đoán theo category
    predictions_by_category = {}
    for category_id in get_categories(user_id):
        category_data = filter_by_category(monthly_data, category_id)
        
        if len(category_data) >= 3:
            category_trend, _, category_model = calculate_trend(category_data)
            predicted_amount = category_model.predict([[next_month_index]])[0]
            
            # Confidence dựa trên số lượng dữ liệu và variance
            confidence = calculate_confidence(category_data)
            
            predictions_by_category[category_id] = {
                "amount": max(0, predicted_amount),  # Không âm
                "confidence": confidence,
                "trend": category_trend
            }
        else:
            # Không đủ dữ liệu, dùng trung bình
            avg = np.mean([d['amount'] for d in category_data])
            predictions_by_category[category_id] = {
                "amount": avg,
                "confidence": 0.5,
                "trend": "insufficient_data"
            }
    
    # Confidence tổng thể
    overall_confidence = np.mean([p['confidence'] for p in predictions_by_category.values()])
    
    return {
        "total_expense": predicted_total,
        "by_category": predictions_by_category,
        "confidence": overall_confidence,
        "trend": trend,
        "trend_percentage": trend_pct
    }
```

#### Bước 4: Tính confidence
```python
def calculate_confidence(data):
    """
    Tính độ tin cậy dựa trên:
    - Số lượng dữ liệu
    - Variance (độ biến thiên)
    - R-squared của model
    """
    n = len(data)
    amounts = [d['amount'] for d in data]
    
    # Variance
    variance = np.var(amounts)
    mean = np.mean(amounts)
    cv = variance / mean if mean > 0 else 1  # Coefficient of Variation
    
    # Confidence formula
    data_confidence = min(n / 6, 1.0)  # Càng nhiều dữ liệu càng tốt
    variance_confidence = 1 / (1 + cv)  # Variance thấp = confidence cao
    
    confidence = (data_confidence * 0.6) + (variance_confidence * 0.4)
    
    return min(confidence, 1.0)
```

### 3.3. Ví dụ minh họa

**Input:** User có 6 tháng dữ liệu:
- Tháng 1: 7,000,000 VND
- Tháng 2: 7,500,000 VND
- Tháng 3: 8,000,000 VND
- Tháng 4: 7,800,000 VND
- Tháng 5: 8,200,000 VND
- Tháng 6: 8,500,000 VND

**Output:**
```json
{
  "total_expense": 8900000,
  "confidence": 0.82,
  "trend": "increasing",
  "trend_percentage": 3.5,
  "by_category": {
    "5": {
      "amount": 2200000,
      "confidence": 0.85,
      "trend": "increasing"
    }
  }
}
```

---

## 4. THUẬT TOÁN ĐỀ XUẤT NGÂN SÁCH TỰ ĐỘNG

### 4.1. Mô tả
Đề xuất ngân sách cho các category dựa trên lịch sử chi tiêu.

### 4.2. Luồng xử lý

```
1. Lấy dữ liệu
   - monthly_income từ user_profile
   - Lịch sử chi tiêu 3 tháng gần nhất

2. Phân tích chi tiêu theo category
   - Tính trung bình chi tiêu mỗi category
   - Tính % so với tổng chi tiêu
   - Xác định categories chiếm % lớn

3. Đề xuất ngân sách
   For each category:
     suggested_amount = avg_spending * 1.1  // Thêm 10% buffer
     
     // Nếu category chiếm > 30% tổng chi tiêu
     if category_percentage > 30:
       suggested_amount = avg_spending * 0.9  // Giảm 10%
   
   // Đảm bảo tổng không vượt quá 90% monthly_income
   total_suggested = sum(suggested_amounts)
   if total_suggested > monthly_income * 0.9:
     scale_factor = (monthly_income * 0.9) / total_suggested
     suggested_amounts = suggested_amounts * scale_factor

4. Tạo notes và insights
   - So sánh với monthly_income
   - Đề xuất giảm chi tiêu nếu cần
   - Gợi ý categories nên ưu tiên
```

### 4.3. Code mẫu

```python
def suggest_budgets(user_id):
    profile = get_user_profile(user_id)
    monthly_income = profile.monthly_income
    
    if monthly_income == 0:
        return {"error": "Please set monthly income first"}
    
    # Lấy lịch sử 3 tháng
    history = get_transactions(user_id, months=3)
    
    # Nhóm theo category
    category_spending = {}
    for tx in history:
        if tx.transaction_type == 'expense':
            cat_id = tx.category_id
            if cat_id not in category_spending:
                category_spending[cat_id] = []
            category_spending[cat_id].append(tx.amount)
    
    # Tính trung bình
    suggestions = []
    total_avg = sum([sum(amounts) for amounts in category_spending.values()])
    
    for cat_id, amounts in category_spending.items():
        avg = sum(amounts) / len(amounts) if amounts else 0
        percentage = (sum(amounts) / total_avg * 100) if total_avg > 0 else 0
        
        # Đề xuất
        if percentage > 30:
            suggested = avg * 0.9  # Giảm 10%
        else:
            suggested = avg * 1.1  # Tăng 10% buffer
        
        suggestions.append({
            "category_id": cat_id,
            "name": get_category_name(cat_id),
            "suggested_amount": suggested,
            "current_avg": avg,
            "percentage": percentage
        })
    
    # Scale để không vượt quá 90% income
    total_suggested = sum([s['suggested_amount'] for s in suggestions])
    max_budget = monthly_income * 0.9
    
    if total_suggested > max_budget:
        scale = max_budget / total_suggested
        for s in suggestions:
            s['suggested_amount'] *= scale
    
    return {
        "suggestions": suggestions,
        "total_suggested": sum([s['suggested_amount'] for s in suggestions]),
        "monthly_income": monthly_income,
        "savings_allocated": monthly_income * 0.1
    }
```

---

## 5. THUẬT TOÁN KIỂM TRA VÀ CẢNH BÁO NGÂN SÁCH

### 5.1. Mô tả
Kiểm tra chi tiêu so với ngân sách và tạo cảnh báo khi vượt ngưỡng.

### 5.2. Luồng xử lý

```
Khi user tạo giao dịch expense:

1. Tìm budgets liên quan
   - Budgets có category_id = transaction.category_id
   - Hoặc budgets không có category_id (tổng ngân sách)
   - Budgets đang active và trong khoảng thời gian

2. Tính spent_amount
   spent = SUM(amount) WHERE 
     transaction_type = 'expense'
     AND category_id = budget.category_id
     AND transaction_date BETWEEN budget.start_date AND budget.end_date

3. Tính usage_percentage
   usage = (spent / budget.amount) * 100

4. Kiểm tra ngưỡng
   if usage >= budget.alert_threshold:
     if usage >= 100:
       // Vượt ngân sách
       create_notification(type='error', priority='urgent')
     else:
       // Gần vượt ngân sách
       create_notification(type='warning', priority='high')
   
   // Chỉ gửi 1 lần trong 24h cho mỗi budget
   if not notification_sent_in_last_24h(budget_id):
     send_notification()
```

### 5.3. Code mẫu

```python
def check_budget_alerts(user_id, transaction):
    """
    Kiểm tra và tạo cảnh báo ngân sách
    """
    if transaction.transaction_type != 'expense':
        return
    
    # Tìm budgets liên quan
    budgets = get_active_budgets(user_id, transaction.transaction_date)
    
    for budget in budgets:
        # Kiểm tra category match
        if budget.category_id and budget.category_id != transaction.category_id:
            continue
        
        # Tính spent
        spent = calculate_spent_amount(
            user_id, 
            budget.category_id,
            budget.start_date,
            budget.end_date
        )
        
        # Tính usage
        usage_pct = (spent / budget.amount) * 100
        
        # Kiểm tra ngưỡng
        if usage_pct >= budget.alert_threshold:
            # Kiểm tra đã gửi notification chưa
            last_notification = get_last_budget_notification(
                user_id, 
                budget.id,
                hours=24
            )
            
            if not last_notification:
                if usage_pct >= 100:
                    create_notification(
                        user_id=user_id,
                        title="Vượt ngân sách",
                        message=f"Bạn đã vượt {budget.name} ({usage_pct:.1f}%)",
                        type='error',
                    priority='urgent'
                    )
                else:
                    create_notification(
                        user_id=user_id,
                        title="Cảnh báo ngân sách",
                        message=f"Bạn đã sử dụng {usage_pct:.1f}% {budget.name}",
                        type='warning',
                    priority='high'
                    )
```

---

## 6. THUẬT TOÁN TÍNH TOÁN DASHBOARD ANALYTICS

### 6.1. Mô tả
Tính toán các chỉ số tài chính cho dashboard.

### 6.2. Các chỉ số

```python
def calculate_dashboard_analytics(user_id, start_date, end_date):
    """
    Tính toán analytics cho dashboard
    """
    transactions = get_transactions(user_id, start_date, end_date)
    
    # 1. Summary
    total_income = sum([t.amount for t in transactions if t.type == 'income'])
    total_expense = sum([t.amount for t in transactions if t.type == 'expense'])
    net_savings = total_income - total_expense
    savings_rate = (net_savings / total_income * 100) if total_income > 0 else 0
    
    # 2. Spending by category
    category_spending = {}
    for tx in transactions:
        if tx.type == 'expense':
            cat_id = tx.category_id
            if cat_id not in category_spending:
                category_spending[cat_id] = 0
            category_spending[cat_id] += tx.amount
    
    # Tính percentage
    spending_by_category = []
    for cat_id, amount in category_spending.items():
        percentage = (amount / total_expense * 100) if total_expense > 0 else 0
        spending_by_category.append({
            "category_id": cat_id,
            "category_name": get_category_name(cat_id),
            "amount": amount,
            "percentage": percentage
        })
    
    # 3. Trends (daily)
    daily_trends = group_by_day(transactions)
    
    # 4. Top categories
    top_categories = sorted(
        spending_by_category,
        key=lambda x: x['amount'],
        reverse=True
    )[:5]
    
    # 5. Budget status
    budgets = get_active_budgets(user_id)
    budget_status = {
        "total": len(budgets),
        "over_budget": 0,
        "at_risk": 0,
        "on_track": 0
    }
    
    for budget in budgets:
        spent = calculate_spent_amount(...)
        usage = (spent / budget.amount) * 100
        
        if usage >= 100:
            budget_status["over_budget"] += 1
        elif usage >= budget.alert_threshold:
            budget_status["at_risk"] += 1
        else:
            budget_status["on_track"] += 1
    
    return {
        "summary": {
            "total_income": total_income,
            "total_expense": total_expense,
            "net_savings": net_savings,
            "savings_rate": savings_rate
        },
        "spending_by_category": spending_by_category,
        "trends": daily_trends,
        "top_categories": top_categories,
        "budget_status": budget_status
    }
```

---

## 7. CACHE STRATEGY

### 7.1. Redis Cache cho Dashboard

```python
def get_dashboard_analytics_cached(user_id):
    """
    Lấy dashboard analytics với cache
    """
    cache_key = f"dashboard:{user_id}"
    
    # Kiểm tra cache
    cached = redis.get(cache_key)
    if cached:
        return json.loads(cached)
    
    # Tính toán
    analytics = calculate_dashboard_analytics(user_id)
    
    # Lưu cache (TTL: 5 phút)
    redis.setex(
        cache_key,
        300,  # 5 minutes
        json.dumps(analytics)
    
    return analytics

def invalidate_dashboard_cache(user_id):
    """
    Xóa cache khi có thay đổi
    """
    cache_key = f"dashboard:{user_id}"
    redis.delete(cache_key)
```

### 7.2. Cache Keys
- `dashboard:{user_id}`: Dashboard analytics (TTL: 5 phút)
- `user:{user_id}`: User profile (TTL: 1 giờ)
- `categories:{user_id}`: User categories (TTL: 1 giờ)
- `session:{token_hash}`: Session data (TTL: theo expires_at)

---

## 8. XỬ LÝ LỖI VÀ EDGE CASES

### 8.1. NLU không hiểu
- Fallback: Yêu cầu user nhập thủ công
- Learning: Lưu feedback để cải thiện

### 8.2. Không đủ dữ liệu cho prediction
- Yêu cầu user thêm dữ liệu
- Sử dụng dữ liệu mặc định nếu có

### 8.3. Anomaly false positive
- Cho phép user xác nhận/sửa
- Lưu feedback để cải thiện model

---

## TÓM TẮT

1. **NLU:** Gemini API + Category Matching
2. **Anomaly Detection:** Z-Score + Isolation Forest
3. **Prediction:** Linear Regression + Time Series
4. **Budget Suggestions:** Statistical Analysis
5. **Budget Alerts:** Real-time checking với rate limiting
6. **Analytics:** Aggregation với caching
7. **Cache:** Redis với TTL hợp lý

Tất cả các thuật toán đều được tối ưu cho performance và độ chính xác.


