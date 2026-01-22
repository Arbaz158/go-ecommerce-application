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
	Enabled            bool
	Addr               string
	EnableFileProfiles bool
	ProfileDir         string // Directory to save profile files
}

func Start(cfg Config) {
	if !cfg.Enabled {
		return
	}

	go startHTTPProfiling(cfg.Addr)

	if cfg.EnableFileProfiles {
		go startFileProfiles(cfg.ProfileDir)
	}

	go collectRuntimeMetrics()
}

func startHTTPProfiling(addr string) {
	log.Printf("[PROFILING] HTTP pprof server starting on http://%s/debug/pprof/\n", addr)

	printPprofEndpoints(addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("[PROFILING] HTTP pprof error: %v\n", err)
	}
}

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

func startFileProfiles(profileDir string) {
	if err := os.MkdirAll(profileDir, os.ModePerm); err != nil {
		log.Printf("[PROFILING] Error creating profile directory: %v\n", err)
		return
	}

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

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func EnableGCMetrics() {
	debug.SetGCPercent(100) // Trigger GC when heap grows 100% larger
	log.Println("[PROFILING] GC metrics enabled with 100% heap growth threshold")
}
