#!/bin/bash

# Script để debug lỗi Analytics page
# Kiểm tra tất cả API calls mà frontend thực hiện

set -e

BACKEND_URL="${1:-http://localhost:8080}"
USER_ID="${2:-15}"

echo "============================================================"
echo "🔍 Debug Analytics Page APIs"
echo "============================================================"
echo "Backend: $BACKEND_URL"
echo "User ID: $USER_ID"
echo ""

# Login
echo "🔐 Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"test15@tabimoney.com\",\"password\":\"test123456\"}")

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | python3 -c "import json,sys; print(json.load(sys.stdin)['access_token'])" 2>/dev/null)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ Login failed"
    exit 1
fi

echo "✅ Logged in"
echo ""

# Calculate date range (last 6 months)
END_DATE=$(date +%Y-%m-%d)
START_DATE=$(date -v-6m +%Y-%m-%d 2>/dev/null || date -d "6 months ago" +%Y-%m-%d)

echo "Date range: $START_DATE to $END_DATE"
echo ""

# Test all APIs that frontend calls
APIS=(
    "dashboard:GET:/api/v1/analytics/dashboard"
    "category-spending:GET:/api/v1/analytics/category-spending?start_date=${START_DATE}&end_date=${END_DATE}"
    "spending-patterns:GET:/api/v1/analytics/spending-patterns?start_date=${START_DATE}&end_date=${END_DATE}"
    "anomalies:GET:/api/v1/analytics/anomalies?start_date=${START_DATE}&end_date=${END_DATE}&threshold=0.6"
    "predictions:GET:/api/v1/analytics/predictions?start_date=${START_DATE}&end_date=${END_DATE}"
)

for api_info in "${APIS[@]}"; do
    IFS=':' read -r name method path <<< "$api_info"
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ "$method" == "GET" ]; then
        RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "${BACKEND_URL}${path}" \
            -H "Authorization: Bearer ${ACCESS_TOKEN}" \
            -H "Content-Type: application/json")
    else
        RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "${BACKEND_URL}${path}" \
            -H "Authorization: Bearer ${ACCESS_TOKEN}" \
            -H "Content-Type: application/json")
    fi
    
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
    BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')
    
    echo "HTTP Status: $HTTP_STATUS"
    
    if [ "$HTTP_STATUS" == "200" ]; then
        # Validate JSON
        if echo "$BODY" | python3 -m json.tool > /dev/null 2>&1; then
            echo "✅ Valid JSON"
            
            # Check for common issues
            if echo "$BODY" | grep -q "null"; then
                echo "⚠️  Contains null values:"
                echo "$BODY" | python3 -c "import json,sys; d=json.load(sys.stdin); [print(f'  - {k}: {v}') for k,v in d.items() if v is None]" 2>/dev/null || true
            fi
            
            # Show structure
            echo "📊 Response structure:"
            echo "$BODY" | python3 -c "import json,sys; d=json.load(sys.stdin); print('  Keys:', list(d.keys())[:10])" 2>/dev/null || echo "  (Could not parse)"
        else
            echo "❌ Invalid JSON!"
            echo "Response: $BODY"
        fi
    else
        echo "❌ Failed!"
        echo "Response: $BODY"
    fi
    
    echo ""
done

echo "============================================================"
echo "✅ Debug complete"
echo "============================================================"




