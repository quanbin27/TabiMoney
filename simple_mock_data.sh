#!/bin/bash

echo "🚀 Simple Mock Data Generator"
echo "=============================="

# Configuration
BASE_URL="http://localhost:3000/api/v1"
EMAIL="test+1@example.com"
PASSWORD="secret123"

# Get authentication token
echo "1. Authenticating..."
TOKEN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" | \
  jq -r '.access_token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Authentication failed!"
  exit 1
fi

echo "✅ Authentication successful!"

# Function to create a transaction
create_transaction() {
  local category_id=$1
  local amount=$2
  local description=$3
  local transaction_type=$4
  local date=$5
  local time=$6
  
  curl -s -X POST "$BASE_URL/transactions" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"category_id\": $category_id,
      \"amount\": $amount,
      \"description\": \"$description\",
      \"transaction_type\": \"$transaction_type\",
      \"transaction_date\": \"$date\",
      \"transaction_time\": \"$time\"
    }" > /dev/null
}

echo ""
echo "2. Creating mock transactions..."

# Create transactions for the last 3 months with fixed dates
create_transaction 1 50000 "Ăn trưa" "expense" "2025-09-15" "12:30"
create_transaction 1 75000 "Ăn tối" "expense" "2025-09-14" "19:00"
create_transaction 1 25000 "Cà phê" "expense" "2025-09-13" "15:30"
create_transaction 1 120000 "Nhà hàng" "expense" "2025-09-12" "18:45"
create_transaction 1 35000 "Bánh mì" "expense" "2025-09-11" "08:15"

create_transaction 2 80000 "Xăng xe" "expense" "2025-09-10" "07:30"
create_transaction 2 45000 "Taxi" "expense" "2025-09-09" "14:20"
create_transaction 2 60000 "Grab" "expense" "2025-09-08" "16:45"
create_transaction 2 30000 "Xe bus" "expense" "2025-09-07" "09:15"
create_transaction 2 25000 "Đậu xe" "expense" "2025-09-06" "10:30"

create_transaction 3 500000 "Quần áo" "expense" "2025-09-05" "14:00"
create_transaction 3 15000000 "Laptop" "expense" "2025-09-04" "11:30"
create_transaction 3 200000 "Sách" "expense" "2025-09-03" "16:00"
create_transaction 3 300000 "Đồ gia dụng" "expense" "2025-09-02" "13:45"
create_transaction 3 800000 "Điện thoại" "expense" "2025-09-01" "15:20"

create_transaction 4 150000 "Xem phim" "expense" "2025-08-31" "20:00"
create_transaction 4 200000 "Game" "expense" "2025-08-30" "19:30"
create_transaction 4 300000 "Karaoke" "expense" "2025-08-29" "21:00"
create_transaction 4 80000 "Cafe" "expense" "2025-08-28" "16:30"
create_transaction 4 500000 "Concert" "expense" "2025-08-27" "19:45"

create_transaction 5 500000 "Khám bệnh" "expense" "2025-08-26" "09:00"
create_transaction 5 150000 "Thuốc" "expense" "2025-08-25" "10:30"
create_transaction 5 200000 "Bảo hiểm" "expense" "2025-08-24" "14:00"
create_transaction 5 800000 "Nha khoa" "expense" "2025-08-23" "11:15"
create_transaction 5 300000 "Khám sức khỏe" "expense" "2025-08-22" "08:30"

create_transaction 6 2000000 "Học phí" "expense" "2025-08-21" "09:00"
create_transaction 6 500000 "Sách giáo khoa" "expense" "2025-08-20" "15:00"
create_transaction 6 1000000 "Khóa học" "expense" "2025-08-19" "10:30"
create_transaction 6 300000 "Thi cử" "expense" "2025-08-18" "08:00"
create_transaction 6 200000 "Dụng cụ học tập" "expense" "2025-08-17" "14:30"

create_transaction 7 5000000 "Vé máy bay" "expense" "2025-08-16" "06:00"
create_transaction 7 2000000 "Khách sạn" "expense" "2025-08-15" "15:00"
create_transaction 7 3000000 "Du lịch" "expense" "2025-08-14" "12:00"
create_transaction 7 1500000 "Vé tàu" "expense" "2025-08-13" "08:30"
create_transaction 7 800000 "Thuê xe" "expense" "2025-08-12" "10:00"

create_transaction 8 400000 "Điện" "expense" "2025-08-11" "09:00"
create_transaction 8 200000 "Nước" "expense" "2025-08-10" "10:00"
create_transaction 8 300000 "Internet" "expense" "2025-08-09" "11:00"
create_transaction 8 150000 "Điện thoại" "expense" "2025-08-08" "12:00"
create_transaction 8 250000 "Gas" "expense" "2025-08-07" "13:00"

# Income transactions
create_transaction 1 15000000 "Lương" "income" "2025-08-01" "09:00"
create_transaction 1 5000000 "Thưởng" "income" "2025-07-15" "10:00"
create_transaction 1 8000000 "Freelance" "income" "2025-07-01" "11:00"
create_transaction 1 12000000 "Đầu tư" "income" "2025-06-15" "12:00"
create_transaction 1 6000000 "Bán hàng" "income" "2025-06-01" "13:00"

echo "✅ Created 40 mock transactions!"

# Wait for data processing
echo ""
echo "3. Waiting for data processing..."
sleep 5

# Check total transactions
echo ""
echo "4. Checking total transactions..."
TOTAL=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/transactions?limit=100" | jq '.data | length')
echo "Total transactions: $TOTAL"

# Test AI service
echo ""
echo "5. Testing AI service..."
AI_RESPONSE=$(curl -s -X POST "http://localhost:8001/api/v1/prediction/expenses" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 2, "start_date": "2025-06-01T00:00:00Z", "end_date": "2025-10-01T00:00:00Z"}')

echo "AI Prediction Response:"
echo "$AI_RESPONSE" | jq '.'

echo ""
echo "🎉 Mock data generation complete!"
echo "You can now test the Analytics Dashboard at: http://localhost:3000/analytics"




