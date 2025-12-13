#!/bin/bash

# Script để test các tính năng AI sau khi đã tạo mock data
# Gọi qua Go backend API (port 8080 hoặc 3000)

set -e

# Configuration
BACKEND_URL="${2:-http://localhost:8080}"
USER_ID="${1:-15}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo ""
echo "============================================================"
echo "🧪 AI Features Test Script"
echo "============================================================"
echo ""
echo "Configuration:"
echo "   Backend URL: $BACKEND_URL"
echo "   User ID: $USER_ID"
echo ""
echo "💡 Usage:"
echo "   ./test_ai_features.sh [user_id] [backend_url]"
echo "   Examples:"
echo "     ./test_ai_features.sh 15 http://localhost:3000  # Via frontend proxy"
echo "     ./test_ai_features.sh 1 http://localhost:8080   # Direct backend"
echo ""

# Check if backend is running
echo "🔍 Checking backend health..."
if curl -s -f "${BACKEND_URL}/health" > /dev/null; then
    echo -e "${GREEN}✅ Backend is running${NC}"
else
    echo -e "${YELLOW}⚠️  Warning: Cannot connect to backend at ${BACKEND_URL}${NC}"
    echo "   Make sure the backend is running before testing"
    read -p "Continue anyway? (y/n): " response
    if [ "$response" != "y" ]; then
        exit 1
    fi
fi

# Login để lấy token
echo ""
echo "🔐 Attempting to login..."
if [ "$USER_ID" == "15" ]; then
    echo "   Using test credentials: test15@tabimoney.com / test123456"
    LOGIN_RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"test15@tabimoney.com","password":"test123456"}')
    
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    
    if [ -z "$ACCESS_TOKEN" ]; then
        echo -e "${RED}❌ Login failed${NC}"
        echo "   Response: $LOGIN_RESPONSE"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Got auth token${NC}"
else
    echo -e "${YELLOW}⚠️  User $USER_ID - Please provide credentials:${NC}"
    read -p "   Email: " email
    read -sp "   Password: " password
    echo ""
    
    LOGIN_RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$email\",\"password\":\"$password\"}")
    
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    
    if [ -z "$ACCESS_TOKEN" ]; then
        echo -e "${RED}❌ Login failed${NC}"
        echo "   Response: $LOGIN_RESPONSE"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Got auth token${NC}"
fi

# Calculate date range (6 months back)
END_DATE=$(date +%Y-%m-%d)
START_DATE=$(date -v-6m +%Y-%m-%d 2>/dev/null || date -d "6 months ago" +%Y-%m-%d)

echo ""
echo "============================================================"
echo "🧮 Testing Prediction Service"
echo "============================================================"
echo ""

PREDICTION_RESPONSE=$(curl -s -X GET "${BACKEND_URL}/api/v1/analytics/predictions?start_date=${START_DATE}&end_date=${END_DATE}" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

# Save response to file for parsing
echo "$PREDICTION_RESPONSE" | python3 -m json.tool > /tmp/prediction_response.json 2>/dev/null || echo "$PREDICTION_RESPONSE" > /tmp/prediction_response.json

if echo "$PREDICTION_RESPONSE" | grep -q "predicted_amount\|user_id"; then
    echo -e "${GREEN}✅ Prediction successful!${NC}"
    echo ""
    echo "📊 Results:"
    
    # Use python to parse JSON properly
    PREDICTED_AMOUNT=$(python3 -c "import json, sys; d=json.load(sys.stdin); print(d.get('predicted_amount', 0))" < /tmp/prediction_response.json 2>/dev/null)
    CONFIDENCE=$(python3 -c "import json, sys; d=json.load(sys.stdin); print(d.get('confidence_score', 0))" < /tmp/prediction_response.json 2>/dev/null)
    
    if [ -n "$PREDICTED_AMOUNT" ] && [ "$PREDICTED_AMOUNT" != "None" ]; then
        printf "   Predicted Amount: %.0f VND\n" "$PREDICTED_AMOUNT"
    fi
    if [ -n "$CONFIDENCE" ] && [ "$CONFIDENCE" != "None" ]; then
        CONFIDENCE_PCT=$(python3 -c "print($CONFIDENCE * 100)" 2>/dev/null)
        printf "   Confidence Score: %.1f%%\n" "$CONFIDENCE_PCT"
    fi
    
    echo ""
    echo "   Full response saved to: /tmp/prediction_response.json"
    echo "   View formatted: cat /tmp/prediction_response.json | python3 -m json.tool"
    
    PREDICTION_OK=true
else
    echo -e "${RED}❌ Prediction failed${NC}"
    echo "   Response: $PREDICTION_RESPONSE"
    PREDICTION_OK=false
fi

echo ""
echo "============================================================"
echo "🔍 Testing Anomaly Detection Service"
echo "============================================================"
echo ""

ANOMALY_RESPONSE=$(curl -s -X GET "${BACKEND_URL}/api/v1/analytics/anomalies?start_date=${START_DATE}&end_date=${END_DATE}&threshold=0.6" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

# Save response to file for parsing
echo "$ANOMALY_RESPONSE" | python3 -m json.tool > /tmp/anomaly_response.json 2>/dev/null || echo "$ANOMALY_RESPONSE" > /tmp/anomaly_response.json

if echo "$ANOMALY_RESPONSE" | grep -q "anomalies\|total_anomalies"; then
    echo -e "${GREEN}✅ Anomaly detection successful!${NC}"
    echo ""
    echo "📊 Results:"
    
    # Use python to parse JSON properly
    TOTAL_ANOMALIES=$(python3 -c "import json, sys; d=json.load(sys.stdin); print(d.get('total_anomalies', len(d.get('anomalies', []))))" < /tmp/anomaly_response.json 2>/dev/null)
    DETECTION_SCORE=$(python3 -c "import json, sys; d=json.load(sys.stdin); print(d.get('detection_score', 0))" < /tmp/anomaly_response.json 2>/dev/null)
    
    if [ -n "$TOTAL_ANOMALIES" ] && [ "$TOTAL_ANOMALIES" != "None" ]; then
        echo "   Total Anomalies: $TOTAL_ANOMALIES"
    fi
    if [ -n "$DETECTION_SCORE" ] && [ "$DETECTION_SCORE" != "None" ]; then
        SCORE_PCT=$(python3 -c "print($DETECTION_SCORE * 100)" 2>/dev/null)
        printf "   Detection Score: %.1f%%\n" "$SCORE_PCT"
    fi
    
    echo ""
    echo "   Full response saved to: /tmp/anomaly_response.json"
    echo "   View formatted: cat /tmp/anomaly_response.json | python3 -m json.tool"
    
    ANOMALY_OK=true
else
    echo -e "${RED}❌ Anomaly detection failed${NC}"
    echo "   Response: $ANOMALY_RESPONSE"
    ANOMALY_OK=false
fi

# Summary
echo ""
echo "============================================================"
echo "📋 Test Summary"
echo "============================================================"
if [ "$PREDICTION_OK" = true ]; then
    echo -e "   Prediction: ${GREEN}✅ PASS${NC}"
else
    echo -e "   Prediction: ${RED}❌ FAIL${NC}"
fi

if [ "$ANOMALY_OK" = true ]; then
    echo -e "   Anomaly Detection: ${GREEN}✅ PASS${NC}"
else
    echo -e "   Anomaly Detection: ${RED}❌ FAIL${NC}"
fi
echo "============================================================"
echo ""

# Show full responses if available
if [ -f /tmp/prediction_response.json ]; then
    echo -e "${BLUE}💡 To view full prediction response:${NC}"
    echo "   cat /tmp/prediction_response.json | python3 -m json.tool"
fi

if [ -f /tmp/anomaly_response.json ]; then
    echo -e "${BLUE}💡 To view full anomaly response:${NC}"
    echo "   cat /tmp/anomaly_response.json | python3 -m json.tool"
fi
echo ""

