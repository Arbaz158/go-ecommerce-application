package profiling

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

type Config struct {
	Enabled bool
	Addr    string
	// Optional: enable file-based profiling for CPU, memory, goroutine
	EnableFileProfiles bool
	ProfileDir         string // Directory to save profile files
}

func Start(cfg Config) {
	if !cfg.Enabled {
		return
	}

	// Start HTTP pprof server
	go startHTTPProfiling(cfg.Addr)

	// Optional: Start file-based profiling for metrics collection
	if cfg.EnableFileProfiles {
		go startFileProfiles(cfg.ProfileDir)
	}

	// Start runtime metrics collection
	go collectRuntimeMetrics()
}

// startHTTPProfiling starts the HTTP pprof server with all profiles exposed
func startHTTPProfiling(addr string) {
	log.Printf("[PROFILING] HTTP pprof server starting on http://%s/debug/pprof/\n", addr)

	// Print available endpoints
	printPprofEndpoints(addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("[PROFILING] HTTP pprof error: %v\n", err)
	}
}

// printPprofEndpoints shows all available profiling endpoints
func printPprofEndpoints(addr string) {
	endpoints := map[string]string{
		"Heap (Memory)":    fmt.Sprintf("http://%s/debug/pprof/heap", addr),
		"CPU (30s sample)": fmt.Sprintf("go tool pprof http://%s/debug/pprof/profile?seconds=30", addr),
		"Goroutines":       fmt.Sprintf("http://%s/debug/pprof/goroutine", addr),
		"Mutex Contention": fmt.Sprintf("http://%s/debug/pprof/mutex", addr),
		"Block Profile":    fmt.Sprintf("http://%s/debug/pprof/block", addr),
		"All Profiles":     fmt.Sprintf("http://%s/debug/pprof/", addr),
		"Trace (5s)":       fmt.Sprintf("wget -O trace.out http://%s/debug/pprof/trace?seconds=5", addr),
	}

	log.Println("\n[PROFILING] Available Endpoints:")
	log.Println("=====================================")
	for name, endpoint := range endpoints {
		log.Printf("  %s: %s\n", name, endpoint)
	}
	log.Println("\n=====================================")
}

// startFileProfiles creates CPU, memory, and goroutine profile files
func startFileProfiles(profileDir string) {
	// Create profile directory if it doesn't exist
	if err := os.MkdirAll(profileDir, os.ModePerm); err != nil {
		log.Printf("[PROFILING] Error creating profile directory: %v\n", err)
		return
	}

	// CPU Profile (30 seconds)
	go func() {
		cpuFile := fmt.Sprintf("%s/cpu_%d.prof", profileDir, time.Now().Unix())
		f, err := os.Create(cpuFile)
		if err != nil {
			log.Printf("[PROFILING] Error creating CPU profile: %v\n", err)
			return
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Printf("[PROFILING] Error starting CPU profile: %v\n", err)
			f.Close()
			return
		}
		defer pprof.StopCPUProfile()

		log.Printf("[PROFILING] CPU profiling started for 30 seconds: %s\n", cpuFile)
		time.Sleep(30 * time.Second)
		log.Printf("[PROFILING] CPU profile saved: %s\n", cpuFile)
	}()

	// Memory Profile (after 30 seconds)
	go func() {
		time.Sleep(40 * time.Second)
		memFile := fmt.Sprintf("%s/mem_%d.prof", profileDir, time.Now().Unix())
		f, err := os.Create(memFile)
		if err != nil {
			log.Printf("[PROFILING] Error creating memory profile: %v\n", err)
			return
		}
		defer f.Close()

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Printf("[PROFILING] Error writing heap profile: %v\n", err)
			return
		}
		log.Printf("[PROFILING] Memory profile saved: %s\n", memFile)
	}()

	// Goroutine Profile
	go func() {
		time.Sleep(50 * time.Second)
		grFile := fmt.Sprintf("%s/goroutine_%d.prof", profileDir, time.Now().Unix())
		f, err := os.Create(grFile)
		if err != nil {
			log.Printf("[PROFILING] Error creating goroutine profile: %v\n", err)
			return
		}
		defer f.Close()

		if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
			log.Printf("[PROFILING] Error writing goroutine profile: %v\n", err)
			return
		}
		log.Printf("[PROFILING] Goroutine profile saved: %s\n", grFile)
	}()
}

// collectRuntimeMetrics collects and logs runtime metrics every 10 seconds
func collectRuntimeMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		log.Printf("[METRICS] Alloc=%vMB TotalAlloc=%vMB Sys=%vMB NumGC=%v Goroutines=%v\n",
			bToMb(m.Alloc),
			bToMb(m.TotalAlloc),
			bToMb(m.Sys),
			m.NumGC,
			runtime.NumGoroutine(),
		)
	}
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// EnableGCMetrics enables detailed GC metrics logging
func EnableGCMetrics() {
	debug.SetGCPercent(100) // Trigger GC when heap grows 100% larger
	log.Println("[PROFILING] GC metrics enabled with 100% heap growth threshold")
}
