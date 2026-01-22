#!/bin/bash

# Profiling Testing Script for Go E-Commerce Application
# This script helps you test different profiling scenarios

set -e

SERVICE_URL="${1:-http://localhost:6060}"
PROFILE_ENDPOINT="${SERVICE_URL}/debug/pprof"
OUTPUT_DIR="./profile-results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Go E-Commerce Profiling Testing Tool${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Menu function
show_menu() {
    echo -e "${YELLOW}Available Tests:${NC}"
    echo "1) CPU Profile (30 seconds)"
    echo "2) Memory/Heap Profile"
    echo "3) Goroutine Profile"
    echo "4) Execution Trace (5 seconds)"
    echo "5) All Profiles (Combined)"
    echo "6) Real-time Metrics Monitor (30 seconds)"
    echo "7) Memory Leak Detection (Baseline Comparison)"
    echo "8) Load Test + Profile"
    echo "9) Exit"
    echo ""
}

# Check if service is running
check_service() {
    if ! curl -s "${SERVICE_URL}" > /dev/null 2>&1; then
        echo -e "${RED}✗ Service not running at ${SERVICE_URL}${NC}"
        echo -e "${YELLOW}Make sure to set ENABLE_PPROF=true and start the service${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Service is running${NC}\n"
}

# CPU Profile Test
cpu_profile() {
    echo -e "${YELLOW}[TEST 1] CPU Profile - 30 Second Sample${NC}"
    echo "Capturing CPU profile..."
    
    OUTPUT_FILE="${OUTPUT_DIR}/cpu_${TIMESTAMP}.prof"
    
    go tool pprof -quiet \
        -output="${OUTPUT_FILE}" \
        "${SERVICE_URL}/debug/pprof/profile?seconds=30"
    
    echo -e "${GREEN}✓ CPU profile saved: ${OUTPUT_FILE}${NC}"
    echo -e "Analyze with: ${BLUE}go tool pprof ${OUTPUT_FILE}${NC}\n"
}

# Memory Profile Test
memory_profile() {
    echo -e "${YELLOW}[TEST 2] Memory/Heap Profile${NC}"
    echo "Capturing heap profile..."
    
    OUTPUT_FILE="${OUTPUT_DIR}/heap_${TIMESTAMP}.prof"
    
    curl -s "${SERVICE_URL}/debug/pprof/heap" > "${OUTPUT_FILE}"
    
    echo -e "${GREEN}✓ Heap profile saved: ${OUTPUT_FILE}${NC}"
    echo -e "Analyze with: ${BLUE}go tool pprof ${OUTPUT_FILE}${NC}\n"
    
    # Also show top allocators
    echo -e "${YELLOW}Top Memory Allocators:${NC}"
    go tool pprof -nodefraction=0.1 -top "${OUTPUT_FILE}" | head -15
    echo ""
}

# Goroutine Profile Test
goroutine_profile() {
    echo -e "${YELLOW}[TEST 3] Goroutine Profile${NC}"
    
    OUTPUT_FILE="${OUTPUT_DIR}/goroutine_${TIMESTAMP}.txt"
    
    curl -s "${SERVICE_URL}/debug/pprof/goroutine" > "${OUTPUT_FILE}"
    
    GOROUTINE_COUNT=$(grep -c "^goroutine" "${OUTPUT_FILE}" || echo "0")
    
    echo -e "${GREEN}✓ Goroutine profile saved: ${OUTPUT_FILE}${NC}"
    echo -e "${BLUE}Current goroutine count: ${GOROUTINE_COUNT}${NC}\n"
    
    echo -e "${YELLOW}First 20 goroutines:${NC}"
    head -20 "${OUTPUT_FILE}"
    echo ""
}

# Execution Trace Test
execution_trace() {
    echo -e "${YELLOW}[TEST 4] Execution Trace - 5 Seconds${NC}"
    echo "Capturing execution trace..."
    
    OUTPUT_FILE="${OUTPUT_DIR}/trace_${TIMESTAMP}.out"
    
    curl -s "${SERVICE_URL}/debug/pprof/trace?seconds=5" > "${OUTPUT_FILE}"
    
    echo -e "${GREEN}✓ Execution trace saved: ${OUTPUT_FILE}${NC}"
    echo -e "View with: ${BLUE}go tool trace ${OUTPUT_FILE}${NC}\n"
}

# All Profiles Test
all_profiles() {
    echo -e "${YELLOW}[TEST 5] Collecting All Profiles${NC}"
    
    # CPU
    echo "  Capturing CPU profile..."
    CPU_FILE="${OUTPUT_DIR}/all_cpu_${TIMESTAMP}.prof"
    go tool pprof -quiet -output="${CPU_FILE}" \
        "${SERVICE_URL}/debug/pprof/profile?seconds=10" 2>/dev/null &
    CPU_PID=$!
    
    # Heap
    echo "  Capturing heap profile..."
    HEAP_FILE="${OUTPUT_DIR}/all_heap_${TIMESTAMP}.prof"
    curl -s "${SERVICE_URL}/debug/pprof/heap" > "${HEAP_FILE}"
    
    # Goroutine
    echo "  Capturing goroutine profile..."
    GOROUTINE_FILE="${OUTPUT_DIR}/all_goroutine_${TIMESTAMP}.prof"
    curl -s "${SERVICE_URL}/debug/pprof/goroutine" > "${GOROUTINE_FILE}"
    
    # Block
    echo "  Capturing block profile..."
    BLOCK_FILE="${OUTPUT_DIR}/all_block_${TIMESTAMP}.prof"
    curl -s "${SERVICE_URL}/debug/pprof/block" > "${BLOCK_FILE}"
    
    # Wait for CPU profile
    wait $CPU_PID 2>/dev/null
    
    echo -e "${GREEN}✓ All profiles collected${NC}"
    echo "Files:"
    echo "  - CPU: ${CPU_FILE}"
    echo "  - Heap: ${HEAP_FILE}"
    echo "  - Goroutine: ${GOROUTINE_FILE}"
    echo "  - Block: ${BLOCK_FILE}"
    echo ""
}

# Real-time Metrics Monitor
monitor_metrics() {
    echo -e "${YELLOW}[TEST 6] Real-time Metrics Monitor (30 seconds)${NC}"
    echo "Monitoring resource usage..."
    
    METRICS_FILE="${OUTPUT_DIR}/metrics_${TIMESTAMP}.txt"
    
    for i in {1..6}; do
        echo -e "\n${BLUE}[Sample $i]${NC}" >> "${METRICS_FILE}"
        curl -s "${SERVICE_URL}/debug/pprof/heap" | head -15 >> "${METRICS_FILE}"
        
        echo -n "."
        sleep 5
    done
    
    echo -e "\n${GREEN}✓ Metrics saved: ${METRICS_FILE}${NC}"
    echo -e "View with: ${BLUE}cat ${METRICS_FILE}${NC}\n"
}

# Memory Leak Detection
memory_leak_detection() {
    echo -e "${YELLOW}[TEST 7] Memory Leak Detection - Baseline Comparison${NC}"
    
    BASELINE_FILE="${OUTPUT_DIR}/leak_baseline_${TIMESTAMP}.prof"
    CURRENT_FILE="${OUTPUT_DIR}/leak_current_${TIMESTAMP}.prof"
    COMPARISON_FILE="${OUTPUT_DIR}/leak_comparison_${TIMESTAMP}.txt"
    
    # Collect baseline
    echo "Collecting baseline memory profile..."
    curl -s "${SERVICE_URL}/debug/pprof/heap" > "${BASELINE_FILE}"
    
    # Wait
    echo "Waiting 30 seconds..."
    sleep 30
    
    # Collect second sample
    echo "Collecting second memory profile..."
    curl -s "${SERVICE_URL}/debug/pprof/heap" > "${CURRENT_FILE}"
    
    # Compare
    echo "Comparing profiles..."
    go tool pprof -base="${BASELINE_FILE}" "${CURRENT_FILE}" \
        -top > "${COMPARISON_FILE}" 2>&1 || true
    
    echo -e "${GREEN}✓ Comparison saved: ${COMPARISON_FILE}${NC}"
    echo -e "Results (new allocations only):${NC}"
    head -20 "${COMPARISON_FILE}"
    echo ""
}

# Load Test with Profiling
load_test() {
    echo -e "${YELLOW}[TEST 8] Load Test + Profile${NC}"
    
    # Check if hey is installed
    if ! command -v hey &> /dev/null; then
        echo -e "${YELLOW}Installing 'hey' load testing tool...${NC}"
        go install github.com/rakyll/hey@latest
    fi
    
    echo "Starting load test (100 requests, 10 concurrent)..."
    
    # Collect profile during load
    LOAD_PROFILE="${OUTPUT_DIR}/load_cpu_${TIMESTAMP}.prof"
    
    # Start profiling in background
    (sleep 1; go tool pprof -quiet -output="${LOAD_PROFILE}" \
        "${SERVICE_URL}/debug/pprof/profile?seconds=10" 2>/dev/null) &
    
    # Run load test
    hey -n 100 -c 10 -q "http://localhost:8080/" 2>&1 | tee "${OUTPUT_DIR}/load_${TIMESTAMP}.txt"
    
    wait
    
    echo -e "${GREEN}✓ Load test completed${NC}"
    echo "CPU profile during load: ${LOAD_PROFILE}"
    echo ""
}

# Main execution
check_service

while true; do
    show_menu
    read -p "Select test (1-9): " choice
    
    case $choice in
        1) cpu_profile ;;
        2) memory_profile ;;
        3) goroutine_profile ;;
        4) execution_trace ;;
        5) all_profiles ;;
        6) monitor_metrics ;;
        7) memory_leak_detection ;;
        8) load_test ;;
        9) 
            echo -e "${GREEN}Exiting...${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}Invalid choice. Please try again.${NC}\n"
            ;;
    esac
done
