#!/bin/bash

# TabiMoney Mock Data Generator
echo "🚀 TabiMoney Mock Data Generator"
echo "=================================================="

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
  local location=$7
  
  curl -s -X POST "$BASE_URL/transactions" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"category_id\": $category_id,
      \"amount\": $amount,
      \"description\": \"$description\",
      \"transaction_type\": \"$transaction_type\",
      \"transaction_date\": \"$date\",
      \"transaction_time\": \"$time\",
      \"location\": \"$location\"
    }" > /dev/null
}

# Generate mock transactions
echo ""
echo "2. Generating mock transactions..."

# Food & Dining transactions
echo "Creating Food & Dining transactions..."
for i in {1..15}; do
  amount=$((20000 + RANDOM % 180000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((8 + RANDOM % 14)) $((RANDOM % 60)))
  descriptions=("Ăn trưa" "Ăn tối" "Cà phê" "Đồ ăn nhanh" "Nhà hàng" "Bánh mì" "Phở" "Bún chả")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 1 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Transportation transactions
echo "Creating Transportation transactions..."
for i in {1..12}; do
  amount=$((15000 + RANDOM % 85000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((7 + RANDOM % 16)) $((RANDOM % 60)))
  descriptions=("Xăng xe" "Taxi" "Grab" "Xe bus" "Đậu xe" "Uber" "Xe máy")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 2 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Shopping transactions
echo "Creating Shopping transactions..."
for i in {1..10}; do
  amount=$((50000 + RANDOM % 4950000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((9 + RANDOM % 12)) $((RANDOM % 60)))
  descriptions=("Quần áo" "Điện thoại" "Laptop" "Sách" "Đồ gia dụng" "Mỹ phẩm" "Giày dép")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 3 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Entertainment transactions
echo "Creating Entertainment transactions..."
for i in {1..8}; do
  amount=$((50000 + RANDOM % 450000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((18 + RANDOM % 6)) $((RANDOM % 60)))
  descriptions=("Xem phim" "Game" "Karaoke" "Cafe" "Concert" "Bar" "Club")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 4 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Healthcare transactions
echo "Creating Healthcare transactions..."
for i in {1..6}; do
  amount=$((100000 + RANDOM % 1900000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((8 + RANDOM % 8)) $((RANDOM % 60)))
  descriptions=("Khám bệnh" "Thuốc" "Bảo hiểm" "Nha khoa" "Khám sức khỏe" "Vaccine")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 5 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Education transactions
echo "Creating Education transactions..."
for i in {1..5}; do
  amount=$((200000 + RANDOM % 9800000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((9 + RANDOM % 8)) $((RANDOM % 60)))
  descriptions=("Học phí" "Sách giáo khoa" "Khóa học" "Thi cử" "Dụng cụ học tập" "Lớp học thêm")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 6 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Travel transactions
echo "Creating Travel transactions..."
for i in {1..4}; do
  amount=$((500000 + RANDOM % 19500000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((6 + RANDOM % 12)) $((RANDOM % 60)))
  descriptions=("Vé máy bay" "Khách sạn" "Du lịch" "Vé tàu" "Thuê xe" "Tour")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 7 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Utilities transactions
echo "Creating Utilities transactions..."
for i in {1..8}; do
  amount=$((100000 + RANDOM % 900000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((8 + RANDOM % 12)) $((RANDOM % 60)))
  descriptions=("Điện" "Nước" "Internet" "Điện thoại" "Gas" "Cáp TV")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 8 $amount "$description" "expense" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

# Income transactions
echo "Creating Income transactions..."
for i in {1..5}; do
  amount=$((5000000 + RANDOM % 15000000))
  date=$(date -d "$((RANDOM % 180)) days ago" +%Y-%m-%d)
  time=$(printf "%02d:%02d" $((9 + RANDOM % 8)) $((RANDOM % 60)))
  descriptions=("Lương" "Thưởng" "Freelance" "Đầu tư" "Bán hàng")
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  locations=("Hà Nội" "TP.HCM" "Đà Nẵng" "Hải Phòng" "")
  location=${locations[$RANDOM % ${#locations[@]}]}
  
  create_transaction 1 $amount "$description" "income" "$date" "$time" "$location"
  echo "  Created: $description - $amount VND"
done

echo ""
echo "✅ Mock data generation complete!"
echo "Total transactions created: ~73"

# Wait for data processing
echo ""
echo "3. Waiting for data processing..."
sleep 5

# Test AI service
echo ""
echo "4. Testing AI service..."
AI_RESPONSE=$(curl -s -X POST "http://localhost:8001/api/v1/prediction/expenses" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 2, "start_date": "2025-04-01T00:00:00Z", "end_date": "2025-10-01T00:00:00Z"}')

echo "AI Prediction Response:"
echo "$AI_RESPONSE" | jq '.'

echo ""
echo "🎉 Mock data generation complete!"
echo "You can now test the Analytics Dashboard at: http://localhost:3000/analytics"




