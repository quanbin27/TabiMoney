# CƠ CHẾ XỬ LÝ & THUẬT TOÁN HỆ THỐNG TABIMONEY

## 1. THUẬT TOÁN XỬ LÝ NLU (NATURAL LANGUAGE UNDERSTANDING)

### 1.1. Mô tả
Xử lý câu lệnh tự nhiên để trích xuất thông tin giao dịch (amount, category, description, date) sử dụng Google Gemini API với fallback rule-based.

### 1.2. Luồng xử lý

```
Input: "tôi vừa ăn bún bò 50k"

1. Preprocessing & Initialization
   - Kiểm tra Gemini API availability
   - Nếu không có → fallback to rule-based NLU
   - Nếu có → sử dụng Gemini với structured prompt

2. Entity Extraction (Primary: Gemini API)
   
   a. Build Prompt với Context:
      - Lấy danh sách categories của user (top 30 categories)
      - Format: "id|name (name_en)"
      - Đưa vào prompt để AI chọn đúng category_id
   
   b. Gemini Processing:
      - Model: gemini-1.5-flash
      - Temperature: 0.3 (low để đảm bảo consistency)
      - Response format: JSON strict schema
      - Prompt yêu cầu:
        * Parse amount về VND đầy đủ (16 triệu → 16000000)
        * Trả về category_id (không phải category name)
        * Xác định intent: add_transaction, query_balance, analyze_data, etc.
        * Tự quyết định needs_confirmation
   
   c. Parse Response:
      - Extract JSON từ response
      - Validate schema
      - Normalize amount (đảm bảo là số VND đầy đủ)
      - Resolve category_id nếu chỉ có category name

3. Category Resolution
   - Nếu AI trả về category name → query DB để resolve category_id
   - Match theo: exact name (VI/EN), partial match, fallback to "Khác"
   - Remove category entity, add category_id entity

4. Fallback: Rule-based NLU (khi Gemini không available)
   
   a. Amount Extraction:
      - Regex patterns: (\d+)k, (\d+)tr, (\d+)tỷ, (\d+)đ
      - Parse và convert: 50k → 50000, 16tr → 16000000
   
   b. Category Matching:
      - Keyword dictionary cho các categories phổ biến
      - Match keywords trong text
      - Map to category_id
   
   c. Date Parsing:
      - "hôm nay" → today
      - "hôm qua" → yesterday
      - DD/MM/YYYY format parsing

5. Intent Classification
   - add_transaction: có amount + category
   - query_balance: keywords "bao nhiêu", "tổng", "số dư"
   - budget_management: keywords "ngân sách", "budget"
   - goal_tracking: keywords "mục tiêu", "goal"
   - analyze_data: keywords "phân tích", "xu hướng"
   - expense_forecasting: keywords "dự đoán", "forecast"
   - general: default

6. Action Execution (tự động nếu confidence cao)
   - add_transaction: Tạo transaction trực tiếp vào DB
   - query_balance: Query và trả về số dư
   - analyze_data: Phân tích dữ liệu 30 ngày gần nhất
   - budget_management: Kiểm tra tình hình ngân sách
   - goal_tracking: Hiển thị tiến độ mục tiêu
   - smart_recommendations: Đưa ra gợi ý thông minh
   - expense_forecasting: Dự đoán chi tiêu tương lai

7. Output
   {
     "intent": "add_transaction",
     "entities": [
       {"type": "amount", "value": "50000", "confidence": 0.95},
       {"type": "category_id", "value": "5", "confidence": 0.90},
       {"type": "date", "value": "2024-01-15", "confidence": 0.85}
     ],
     "confidence": 0.90,
     "needs_confirmation": false,
     "response": "Đã thêm giao dịch ăn bún bò 50,000 VND",
     "action": "create_transaction"
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
Phát hiện giao dịch bất thường dựa trên lịch sử chi tiêu của user sử dụng Isolation Forest algorithm.

### 2.2. Phương pháp: Isolation Forest (ML-based)

#### Bước 1: Thu thập dữ liệu
```python
# Lấy lịch sử giao dịch trong khoảng thời gian
query = (
    "SELECT t.id, t.amount, t.transaction_type, t.transaction_date, "
    "t.category_id, c.name as category_name "
    "FROM transactions t LEFT JOIN categories c ON c.id = t.category_id "
    "WHERE t.user_id = %s AND t.transaction_type = 'expense' "
    "AND t.transaction_date BETWEEN %s AND %s "
    "ORDER BY t.transaction_date ASC"
)
```

#### Bước 2: Chuẩn bị Features
```python
# Xây dựng feature matrix với các đặc trưng:
# - amount (log-transformed để xử lý skewness)
# - day_of_week (0-6)
# - month (1-12)
# - category_id (integer)

X = []
for transaction in transactions:
    dt = datetime.strptime(transaction["transaction_date"], "%Y-%m-%d")
    amount = float(transaction["amount"] or 0.0)
    amount_log = np.log1p(max(amount, 0.0))  # Log transform để giảm skewness
    day_of_week = dt.weekday()
    month = dt.month
    category_id = int(transaction["category_id"] or 0)
    
    X.append([amount_log, day_of_week, month, category_id])
```

#### Bước 3: Phát hiện bất thường với Isolation Forest
```python
from sklearn.ensemble import IsolationForest

def detect_anomalies(transactions, threshold=0.6):
    """
    Phát hiện anomalies sử dụng Isolation Forest
    
    Args:
        transactions: List các giao dịch
        threshold: Contamination rate (0-1), điều chỉnh độ nhạy
    
    Returns:
        List các anomalies với scores
    """
    # Yêu cầu tối thiểu 10 giao dịch để train model
    if len(transactions) < 10:
        return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0}
    
    # Chuẩn bị feature matrix
    X = prepare_features(transactions)
    
    # Điều chỉnh contamination rate (0.01 - 0.4)
    contamination = min(max(threshold, 0.01), 0.4)
    
    # Train Isolation Forest model
    model = IsolationForest(
        n_estimators=200,      # Số cây quyết định
        contamination=contamination,  # Tỷ lệ outliers mong đợi
        random_state=42        # Đảm bảo reproducibility
    )
    model.fit(X)
    
    # Predict và tính scores
    predictions = model.predict(X)  # -1: anomaly, 1: normal
    scores = model.decision_function(X)  # Higher = more normal, Lower = more anomalous
    
    # Xác định anomalies
    anomalies = []
    for i, pred in enumerate(predictions):
        if pred == -1:  # Là anomaly
            anomaly_score = float(-scores[i])  # Invert để score cao = anomalous hơn
            anomalies.append({
                "transaction_id": transactions[i]["id"],
                "amount": float(transactions[i]["amount"]),
                "category_name": transactions[i].get("category_name", ""),
                "anomaly_score": anomaly_score,
                "anomaly_type": "amount_pattern",
                "description": "Giao dịch có mẫu khác thường theo mô hình IsolationForest",
                "transaction_date": str(transactions[i]["transaction_date"])
            })
    
    # Tính detection score tổng thể
    if anomalies:
        detection_score = float(np.clip(
            np.mean([-scores[i] for i, p in enumerate(predictions) if p == -1]),
            0.0, 1.0
        ))
    else:
        detection_score = 0.0
    
    return {
        "anomalies": anomalies,
        "total_anomalies": len(anomalies),
        "detection_score": detection_score
    }
```

#### Bước 4: Giải thích Isolation Forest
- **Isolation Forest** là thuật toán unsupervised learning dựa trên nguyên lý: outliers dễ bị "cô lập" hơn normal points
- Mỗi tree trong forest cố gắng isolate một điểm bằng cách chọn random feature và split value
- Anomalies cần ít splits hơn để isolate → có path ngắn hơn trong tree
- Kết hợp nhiều trees (200 estimators) để tăng độ chính xác

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
Dự đoán chi tiêu tháng tới dựa trên lịch sử chi tiêu sử dụng Random Forest Regressor kết hợp với Exponential Moving Average (EMA).

### 3.2. Phương pháp: Random Forest Regressor + EMA (Ensemble Approach)

#### Bước 1: Thu thập và chuẩn bị dữ liệu
```python
# Lấy dữ liệu lịch sử giao dịch
query = (
    "SELECT t.amount, t.transaction_type AS type, t.transaction_date AS date, "
    "COALESCE(c.name, 'Unknown') AS category "
    "FROM transactions t "
    "LEFT JOIN categories c ON c.id = t.category_id "
    "WHERE t.user_id = %s AND t.transaction_date BETWEEN %s AND %s "
    "ORDER BY t.transaction_date ASC"
)

# Yêu cầu tối thiểu 5 transactions để bắt đầu dự đoán
if len(historical_data) < 5:
    return error_response("Cần thêm dữ liệu giao dịch để đưa ra dự đoán")
```

#### Bước 2: Xây dựng Monthly Time Series
```python
def _build_monthly_series(historical_data):
    """
    Tổng hợp giao dịch theo tháng và tính toán các features
    """
    df = pd.DataFrame(historical_data)
    df['date'] = pd.to_datetime(df['date'])
    df = df[df['type'] == 'expense'].copy()  # Chỉ lấy chi tiêu
    
    # Nhóm theo tháng
    df['year_month'] = df['date'].dt.to_period('M').astype(str)
    monthly = df.groupby('year_month')['amount'].sum().reset_index()
    monthly = monthly.sort_values('year_month')
    monthly['total_expense'] = monthly['amount'].astype(float)
    
    # Tính toán temporal features
    monthly['month'] = pd.to_datetime(monthly['year_month']).dt.month
    monthly['year'] = pd.to_datetime(monthly['year_month']).dt.year
    
    # Rolling statistics (moving averages)
    monthly['roll_mean_3'] = monthly['total_expense'].rolling(window=3, min_periods=1).mean()
    monthly['roll_mean_6'] = monthly['total_expense'].rolling(window=6, min_periods=1).mean()
    monthly['roll_std_6'] = monthly['total_expense'].rolling(window=6, min_periods=1).std().fillna(0.0)
    monthly['count_seen'] = np.arange(1, len(monthly) + 1)
    
    return monthly[['year_month', 'year', 'month', 'total_expense', 
                    'roll_mean_3', 'roll_mean_6', 'roll_std_6', 'count_seen']]
```

#### Bước 3: Chuẩn bị Training Data cho Random Forest
```python
def _prepare_monthly_training_data(monthly_df):
    """
    Xây dựng features và target để train Random Forest
    
    Features:
    - month: Tháng (1-12) - để capture seasonality
    - rm3_prev: Rolling mean 3 tháng trước đó
    - rm6_prev: Rolling mean 6 tháng trước đó
    - rs6_prev: Rolling std 6 tháng trước đó
    - count_prev: Số tháng đã quan sát
    
    Target: total_expense của tháng hiện tại
    """
    df = monthly_df.copy().reset_index(drop=True)
    
    # Shift rolling stats để tránh data leakage
    df['rm3_prev'] = df['roll_mean_3'].shift(1)
    df['rm6_prev'] = df['roll_mean_6'].shift(1)
    df['rs6_prev'] = df['roll_std_6'].shift(1)
    df['count_prev'] = (df['count_seen'] - 1).clip(lower=0)
    
    # Bỏ row đầu tiên do shift
    df = df.iloc[1:].reset_index(drop=True)
    
    X = df[['month', 'rm3_prev', 'rm6_prev', 'rs6_prev', 'count_prev']].fillna(0.0).astype(float).values
    y = df['total_expense'].astype(float).values
    
    return X, y
```

#### Bước 4: Train Random Forest Model (Per-User Caching)
```python
from sklearn.ensemble import RandomForestRegressor

def _get_or_train_user_model(user_id, X, y, monthly_df):
    """
    Cache model cho mỗi user để tránh retrain không cần thiết
    Sử dụng fingerprint của time series để detect changes
    """
    # Tính fingerprint của time series
    fingerprint = _fingerprint_monthly_series(monthly_df)
    
    # Kiểm tra cache
    cached_model = self._user_models.get(user_id)
    cached_fp = self._user_series_fp.get(user_id)
    
    # Nếu model chưa có hoặc data đã thay đổi → train mới
    if cached_model is None or cached_fp != fingerprint:
        model = RandomForestRegressor(
            n_estimators=200,      # Số cây quyết định
            random_state=42,        # Reproducibility
            max_depth=12           # Giới hạn độ sâu để tránh overfitting
        )
        model.fit(X, y)
        
        # Cache model và fingerprint
        self._user_models[user_id] = model
        self._user_series_fp[user_id] = fingerprint
    
    return self._user_models[user_id]
```

#### Bước 5: Dự đoán với EMA (Exponential Moving Average)
```python
def _predict_with_ema(historical_data):
    """
    Dự đoán sử dụng Exponential Moving Average trên daily expenses
    EMA phản ánh xu hướng gần đây tốt hơn simple moving average
    """
    df = pd.DataFrame(historical_data)
    df['date'] = pd.to_datetime(df['date'])
    df = df[df['type'] == 'expense']
    
    # Tổng hợp theo ngày
    daily = df.groupby(df['date'].dt.date)['amount'].sum().astype(float)
    
    if len(daily) < 5:
        return None
    
    # Tính EMA với span động (5-20 ngày tùy data)
    span = max(5, min(20, len(daily) // 2))
    ema = daily.ewm(span=span, adjust=False).mean()
    
    # Dự đoán tháng tới từ daily EMA
    last_ema = float(ema.iloc[-1])
    projected_monthly = last_ema * 30.0  # Scale lên monthly
    
    return max(0.0, projected_monthly)
```

#### Bước 6: Kết hợp Predictions (Ensemble)
```python
def predict_expenses(user_id, start_date, end_date):
    """
    Dự đoán chi tiêu tháng tới bằng cách kết hợp:
    - Random Forest prediction (60% weight)
    - EMA prediction (40% weight)
    """
    # Lấy dữ liệu lịch sử
    historical_data = await _get_historical_data(user_id, start_date, end_date)
    
    # Xây dựng monthly series
    monthly_df = _build_monthly_series(historical_data)
    
    if len(monthly_df) < 3:
        return error_response("Dữ liệu không đủ để dự đoán theo tháng")
    
    # Train và predict với Random Forest
    X, y = _prepare_monthly_training_data(monthly_df)
    model = _get_or_train_user_model(user_id, X, y, monthly_df)
    
    # Chuẩn bị features cho tháng tới
    next_features = _prepare_next_month_features(monthly_df)
    ml_prediction = float(model.predict([next_features])[0])
    
    # Dự đoán với EMA
    ema_prediction = _predict_with_ema(historical_data)
    
    # Kết hợp predictions (weighted ensemble)
    if ema_prediction is not None:
        final_prediction = max(0.0, 0.6 * ml_prediction + 0.4 * ema_prediction)
        
        # Confidence cao hơn nếu 2 predictions đồng thuận
        agreement = 1.0 - min(1.0, abs(ml_prediction - ema_prediction) / 
                             max(1.0, abs(ml_prediction) + abs(ema_prediction)))
        confidence = min(0.99, base_confidence * (0.9 + 0.1 * agreement))
    else:
        final_prediction = max(0.0, ml_prediction)
        confidence = base_confidence
    
    return {
        "predicted_amount": final_prediction,
        "confidence_score": confidence,
        "category_breakdown": _generate_category_breakdown(historical_data),
        "trends": _generate_trends(historical_data),
        "recommendations": _generate_recommendations(historical_data, final_prediction)
    }
```

#### Bước 7: Tính Confidence Score
```python
def calculate_confidence(monthly_df, ml_pred, ema_pred):
    """
    Tính độ tin cậy dựa trên:
    - Số lượng tháng dữ liệu (càng nhiều càng tốt)
    - Sự đồng thuận giữa ML và EMA predictions
    """
    base_confidence = min(0.95, len(monthly_df) / 36.0)  # Scale đến ~3 năm
    
    if ema_pred is not None:
        # Agreement factor: predictions càng gần nhau → confidence càng cao
        denom = max(1.0, abs(ml_pred) + abs(ema_pred))
        agreement = 1.0 - min(1.0, abs(ml_pred - ema_pred) / denom)
        confidence = min(0.99, base_confidence * (0.9 + 0.1 * agreement))
    else:
        confidence = base_confidence
    
    return confidence
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

### Các Thuật toán Thực tế Được Triển khai:

1. **NLU (Natural Language Understanding):**
   - **Primary:** Google Gemini API (gemini-1.5-flash) với structured JSON response
   - **Fallback:** Rule-based NLU với regex patterns và keyword matching
   - **Features:** Intent classification, entity extraction, category resolution, auto-action execution

2. **Anomaly Detection:**
   - **Method:** Isolation Forest (scikit-learn)
   - **Features:** log(amount), day_of_week, month, category_id
   - **Parameters:** n_estimators=200, contamination=0.01-0.4 (adjustable)
   - **Output:** Anomaly score và transaction details

3. **Expense Prediction:**
   - **Primary Method:** Random Forest Regressor (n_estimators=200, max_depth=12)
   - **Secondary Method:** Exponential Moving Average (EMA) với dynamic span
   - **Ensemble:** Weighted combination (60% RF + 40% EMA)
   - **Features:** month, rolling_mean_3, rolling_mean_6, rolling_std_6, count_seen
   - **Caching:** Per-user model caching với fingerprint-based invalidation

4. **Budget Suggestions:**
   - **Method:** Statistical Analysis (median of last 3 months) với 10% safety margin
   - **Fallback:** 50/30/20 rule (Needs/Wants/Savings) nếu không có lịch sử
   - **Scaling:** Đảm bảo tổng không vượt quá 90% monthly income

5. **Budget Alerts:**
   - **Method:** Real-time percentage-based checking
   - **Triggers:** Usage >= alert_threshold (default 80%)
   - **Rate Limiting:** 1 notification per 24 hours per budget
   - **Types:** Warning (threshold reached), Error (budget exceeded)

6. **Dashboard Analytics:**
   - **Method:** SQL aggregation với Redis caching
   - **Metrics:** Total income/expense, net savings, savings rate, category breakdown
   - **Financial Health Score:** Based on savings rate (0-100 scale)
   - **Cache TTL:** 1 hour for monthly summaries, 5 minutes for dashboard

7. **Cache Strategy:**
   - **Redis Keys:**
     - `dashboard:{user_id}:{period}` - Dashboard analytics (TTL: 1 hour)
     - `user:{user_id}` - User profile (TTL: 1 hour)
     - `categories:{user_id}` - User categories (TTL: 1 hour)
     - `session:{token_hash}` - Session data (TTL: theo expires_at)
   - **Invalidation:** Automatic on transaction create/update/delete

### Công nghệ và Thư viện:

- **ML Libraries:** scikit-learn (RandomForestRegressor, IsolationForest), pandas, numpy
- **AI API:** Google Gemini API (gemini-1.5-flash)
- **Database:** MySQL 8.0 (primary), Redis 7.0 (cache)
- **Backend:** Golang (Echo framework) + Python (FastAPI for AI service)
- **Time Series:** pandas với rolling statistics và EMA

### Tối ưu hóa Performance:

1. **Model Caching:** Per-user Random Forest models với fingerprint-based invalidation
2. **Query Optimization:** Indexed database queries, efficient aggregations
3. **Cache Strategy:** Multi-level caching (Redis) với appropriate TTLs
4. **Batch Processing:** Monthly aggregation thay vì real-time calculation
5. **Fallback Mechanisms:** Rule-based fallbacks khi AI services unavailable

Tất cả các thuật toán đều được tối ưu cho performance, accuracy và reliability.


