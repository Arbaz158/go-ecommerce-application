# Profiling Quick Reference

## Enable Profiling

```bash
export ENABLE_PPROF=true
go run cmd/main.go
```

## Profile Types & Commands

### 1. CPU Profile (Find Hot Spots)
```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
(pprof) top10              # Top 10 CPU consumers
(pprof) list main.handler  # Source code view
```

### 2. Memory Profile (Find Allocators)
```bash
go tool pprof http://localhost:6060/debug/pprof/heap
(pprof) top                # Top allocators
(pprof) alloc_space        # Total allocations
(pprof) inuse_space        # Current memory
```

### 3. Goroutine Profile (Leak Detection)
```bash
curl http://localhost:6060/debug/pprof/goroutine
# or
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### 4. Execution Trace (Timeline)
```bash
wget -O trace.out 'http://localhost:6060/debug/pprof/trace?seconds=5'
go tool trace trace.out    # Opens browser
```

### 5. Block Profile (Blocking Points)
```bash
go tool pprof http://localhost:6060/debug/pprof/block
```

### 6. Mutex Profile (Lock Contention)
```bash
go tool pprof http://localhost:6060/debug/pprof/mutex
```

## Web Interface
```
http://localhost:6060/debug/pprof/
```

## Real-Time Monitoring
```bash
# Watch logs for metrics (auto-logged every 10s)
tail -f application.log | grep METRICS

# Manual check
curl http://localhost:6060/debug/pprof/heap | head -20
```

## Memory Leak Detection
```bash
# Get baseline
curl http://localhost:6060/debug/pprof/heap > heap1.prof

# Wait and generate load
sleep 300

# Get second sample
curl http://localhost:6060/debug/pprof/heap > heap2.prof

# Compare (shows new allocations)
go tool pprof -base heap1.prof heap2.prof
```

## Testing Script
```bash
./test-profiling.sh
# Interactive menu with 8 different profile tests
```

## Log Metrics Format
```
[METRICS] Alloc=45MB TotalAlloc=320MB Sys=78MB NumGC=42 Goroutines=15

Alloc:      Current heap allocation (MB)
TotalAlloc: Total allocated memory (MB)
Sys:        Total system memory (MB)
NumGC:      Garbage collection count
Goroutines: Current running goroutines
```

## Available Endpoints
- **Heap**: `http://localhost:6060/debug/pprof/heap`
- **CPU**: `go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30`
- **Goroutines**: `http://localhost:6060/debug/pprof/goroutine`
- **Block**: `http://localhost:6060/debug/pprof/block`
- **Mutex**: `http://localhost:6060/debug/pprof/mutex`
- **Trace**: `http://localhost:6060/debug/pprof/trace?seconds=5`
- **All**: `http://localhost:6060/debug/pprof/`

## pprof Commands
```
(pprof) top         # Top functions/allocators
(pprof) list <fn>   # Source code view
(pprof) web         # Generate SVG graph
(pprof) png         # Save PNG graph
(pprof) pdf         # Save PDF graph
(pprof) peek <fn>   # Quick source peek
(pprof) quit        # Exit
```

## Visual Analysis
```bash
# Interactive web UI on port 8888
go tool pprof -http=:8888 http://localhost:6060/debug/pprof/heap

# Download and analyze locally
curl http://localhost:6060/debug/pprof/heap > local.prof
go tool pprof -http=:8888 local.prof
```

## Example Workflow

```bash
# 1. Start service
export ENABLE_PPROF=true
go run cmd/main.go

# 2. Generate load (another terminal)
hey -n 5000 -c 100 http://localhost:8080/api/endpoint

# 3. Profile during load
go tool pprof -http=:8888 http://localhost:6060/debug/pprof/profile?seconds=30

# 4. Analyze results
# - Graph view shows call hierarchy
# - Top shows most expensive functions
# - Source view shows actual code
```

## Expected Values for Health Check
```
Alloc:      10-100 MB (depends on load)
Goroutines: 5-50 (should not grow indefinitely)
NumGC:      Gradually increasing (normal)
Memory gap: alloc_space - inuse_space should be stable
```
