#!/bin/bash

# Mock Data Generator Script
# Wrapper script để dễ dàng tạo mock data

set -e

# Default values
USER_ID=15
COUNT=200
MONTHS=6
ANOMALIES=true
SEED=42

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--user)
            USER_ID="$2"
            shift 2
            ;;
        -c|--count)
            COUNT="$2"
            shift 2
            ;;
        -m|--months)
            MONTHS="$2"
            shift 2
            ;;
        -a|--anomalies)
            ANOMALIES="$2"
            shift 2
            ;;
        -s|--seed)
            SEED="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -u, --user ID       User ID (default: 1)"
            echo "  -c, --count NUM     Number of transactions (default: 200)"
            echo "  -m, --months NUM    Number of months (default: 6)"
            echo "  -a, --anomalies     Include anomalies (default: true)"
            echo "  -s, --seed NUM      Random seed (default: 42)"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Examples:"
            echo "  $0 -u 1 -c 500 -m 12"
            echo "  $0 --user 1 --count 200 --months 6 --anomalies false"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"

echo "🚀 Generating mock data..."
echo "   User ID: $USER_ID"
echo "   Count: $COUNT transactions"
echo "   Months: $MONTHS"
echo "   Anomalies: $ANOMALIES"
echo "   Seed: $SEED"
echo ""

cd "$PROJECT_ROOT"

# Run the generator
go run cmd/mockdata/main.go \
    -user="$USER_ID" \
    -count="$COUNT" \
    -months="$MONTHS" \
    -anomalies="$ANOMALIES" \
    -seed="$SEED"

echo ""
echo "✅ Done!"

