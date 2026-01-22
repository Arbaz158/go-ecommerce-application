# Go E-Commerce Application - Profiling Guide

## Overview

Your application now has **production-level profiling** integrated using Go's `pprof` (CPU, Memory, Goroutine, Block, Mutex profiling).

---

## 1. Current Profiling Integration

### **Method Used: HTTP-based pprof**
- Built-in Go profiling framework
- Exposes profiles via HTTP endpoints
- Zero-cost when disabled
- Supported profile types:
  - ✅ CPU Profile
  - ✅ Memory/Heap Profile
  - ✅ Goroutine Profile
  - ✅ Block Profile (synchronization blocking)
  - ✅ Mutex Profile (lock contention)
  - ✅ Trace (execution trace)

---

## 2. How to Enable Profiling

### Configuration:

Set environment variables before running your services:

```bash
# Enable profiling
export ENABLE_PPROF=true

# For auth-service
export HTTP_ADDR=:8080

# For user-service  
export HTTP_ADDR=:8081
```

### In `.env` file:
```
ENABLE_PPROF=true
HTTP_ADDR=:8080
GIN_MODE=debug  # or release
```

---

## 3. Profile Types & How to Test Them

### **A. CPU Profile** (30-second sample)
Shows where your application spends CPU time.

```bash
# Method 1: Live sampling (30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Method 2: In interactive mode, type:
# (pprof) top10          # Show top 10 CPU consuming functions
# (pprof) list <func>    # Show source code of specific function
# (pprof) pdf > cpu.pdf  # Generate visualization
```

### **B. Memory/Heap Profile**
Shows memory allocation patterns and potential leaks.

```bash
# Live heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# In pprof mode:
# (pprof) top10          # Top memory allocators
# (pprof) alloc_space    # Total allocations
# (pprof) alloc_objects  # Number of allocated objects
# (pprof) inuse_space    # Currently in-use memory
```

### **C. Goroutine Profile**
Shows all running goroutines and potential goroutine leaks.

```bash
# Check goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Quick view in browser
curl http://localhost:6060/debug/pprof/goroutine

# Text format
go tool pprof -http=:8888 http://localhost:6060/debug/pprof/goroutine
```

### **D. Mutex Profile**
Shows lock contention issues (which mutexes are slowing down your app).

```bash
# Requires runtime.SetMutexProfileFraction() to be set
go tool pprof http://localhost:6060/debug/pprof/mutex

# In interactive mode:
# (pprof) top          # Top mutex contentions
```

### **E. Block Profile**
Shows where goroutines block on channels, mutexes, or I/O.

```bash
go tool pprof http://localhost:6060/debug/pprof/block

# Check channel/mutex blocking points
```

### **F. Execution Trace** (5-second trace)
Shows detailed timeline of execution.

```bash
# Download trace
wget -O trace.out 'http://localhost:6060/debug/pprof/trace?seconds=5'

# Analyze trace
go tool trace trace.out

# Opens browser with interactive timeline showing:
# - Goroutine execution
# - GC events
# - Network blocking
# - Synchronization events
```

---

## 4. Real-Time Resource Monitoring

### **A. Real-Time Metrics in Logs**
Your application logs resource usage every 10 seconds:

```
[METRICS] Alloc=45MB TotalAlloc=320MB Sys=78MB NumGC=42 Goroutines=15
```

**Breakdown:**
- `Alloc`: Current heap allocation
- `TotalAlloc`: Total allocated (including freed)
- `Sys`: Total system memory
- `NumGC`: Number of GC runs
- `Goroutines`: Current goroutine count

### **B. All Available Endpoints**

When profiling is enabled, check the startup logs:

```
[PROFILING] Available Endpoints:
=====================================
  Heap (Memory): http://localhost:6060/debug/pprof/heap
  CPU (30s sample): go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
  Goroutines: http://localhost:6060/debug/pprof/goroutine
  Mutex Contention: http://localhost:6060/debug/pprof/mutex
  Block Profile: http://localhost:6060/debug/pprof/block
  All Profiles: http://localhost:6060/debug/pprof/
  Trace (5s): wget -O trace.out 'http://localhost:6060/debug/pprof/trace?seconds=5'
=====================================
```

### **C. Browser Web UI**
Access the pprof web interface:

```
http://localhost:6060/debug/pprof/
```

Shows:
- 📊 Graph visualizations
- 📈 Performance statistics
- 🔍 Interactive analysis

---

## 5. Testing Workflow

### **Step 1: Start Services with Profiling Enabled**
```bash
export ENABLE_PPROF=true
cd services/auth-service && go run cmd/main.go
```

### **Step 2: Generate Load** (in another terminal)
```bash
# Using ab (Apache Bench)
ab -n 1000 -c 50 http://localhost:8080/api/health

# Or using hey
go install github.com/rakyll/hey@latest
hey -n 5000 -c 100 http://localhost:8080/api/health
```

### **Step 3: Profile During Load**

```bash
# CPU Profile (capture while load is running)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory Profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutines (see if they're growing)
curl http://localhost:6060/debug/pprof/goroutine | grep -c goroutine
```

### **Step 4: View Results**

In pprof interactive shell:
```
(pprof) top10            # Top 10 functions
(pprof) list <function>  # Source code view
(pprof) web              # Generate SVG graph
(pprof) png              # Save as PNG
```

---

## 6. Practical Example: Finding Memory Leaks

### **Scenario: Suspected Memory Leak**

```bash
# Terminal 1: Start service
export ENABLE_PPROF=true
go run cmd/main.go

# Terminal 2: Generate continuous load
while true; do
  curl -s http://localhost:8080/api/endpoint > /dev/null
  sleep 0.1
done

# Terminal 3: Monitor memory over time
# Take first baseline
curl -s http://localhost:6060/debug/pprof/heap > heap_1.prof

# Wait 5 minutes
sleep 300

# Take second sample
curl -s http://localhost:6060/debug/pprof/heap > heap_2.prof

# Compare profiles
go tool pprof -base heap_1.prof heap_2.prof
# Shows only new allocations - indicates leaks if growing
```

---

## 7. Production-Level Enhancements Checklist

| Feature | Status | Notes |
|---------|--------|-------|
| HTTP pprof endpoints | ✅ Enabled | All profiles available |
| Real-time metrics logging | ✅ Enabled | Every 10 seconds |
| CPU profiling | ✅ Available | 30-second samples |
| Memory profiling | ✅ Available | Heap analysis |
| Goroutine tracking | ✅ Available | Leak detection |
| Execution trace | ✅ Available | 5-second traces |
| Mutex contention | ✅ Available | Lock analysis |
| Block profile | ✅ Available | Channel/I/O blocking |

---

## 8. Interpreting Results

### **High CPU Usage**
```
(pprof) top
Shows which functions consume most CPU
→ Profile during load for accurate results
```

### **Memory Growing**
```
Check: inuse_space (currently allocated)
vs: alloc_space (total allocations)
→ If gap increases = memory leak
```

### **Goroutine Leak**
```
curl http://localhost:6060/debug/pprof/goroutine
→ Count should stabilize, not continuously grow
```

### **Lock Contention**
```
(pprof) top -cum    # Cumulative time in mutexes
→ High values = serialization bottleneck
```

---

## 9. Advanced: File-Based Profiling (Optional)

For longer-running profiles saved to disk:

```go
profiling.Start(profiling.Config{
    Enabled:            true,
    Addr:               ":6060",
    EnableFileProfiles: true,
    ProfileDir:         "./profiles",  // Creates CPU, Memory, Goroutine profiles
})
```

This creates:
- `cpu_<timestamp>.prof` - 30-second CPU profile
- `mem_<timestamp>.prof` - Memory snapshot
- `goroutine_<timestamp>.prof` - Goroutine dump

Analyze: `go tool pprof ./profiles/cpu_1234567890.prof`

---

## 10. Quick Reference Commands

```bash
# View all available profiles
curl http://localhost:6060/debug/pprof/

# CPU: Find hot spots
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory: Find allocators
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutines: Count and stack traces
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# Real-time watch memory
watch 'curl -s http://localhost:6060/debug/pprof/heap | head -20'

# Generate comparison (leak detection)
go tool pprof -base baseline.prof current.prof

# Export as interactive web UI
go tool pprof -http=:8888 http://localhost:6060/debug/pprof/heap
```

---

## Summary

Your application now has:
- ✅ **Production-grade HTTP profiling** via pprof
- ✅ **Real-time resource metrics** logged every 10 seconds
- ✅ **All profile types**: CPU, Memory, Goroutine, Block, Mutex, Trace
- ✅ **Multiple analysis methods**: CLI, web UI, file export
- ✅ **Load testing capability** to capture realistic profiles

This setup is suitable for production monitoring and debugging!