// A managed worker for xk6-style load scripts.
//
// It is deliberately small: the point of this example is the PACKAGING and the
// lifecycle Mockarty gives a native worker, not the load engine itself. It
// starts on the port the host assigns, reports readiness, accepts a task, and
// answers with the shape the host expects.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type task struct {
	Script  string `json:"script"`
	VUs     int    `json:"vus"`
	Seconds int    `json:"seconds"`
}

type result struct {
	Status      string  `json:"status"`
	Requests    int     `json:"requests"`
	Failures    int     `json:"failures"`
	P95Millis   float64 `json:"p95Millis"`
	DurationSec float64 `json:"durationSeconds"`
}

func main() {
	port := os.Getenv("MOCKARTY_WORKER_PORT")
	if port == "" {
		port = "47100"
	}
	extension := os.Getenv("MOCKARTY_EXTENSION_ID")

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok", "extension": extension, "worker": os.Getenv("MOCKARTY_WORKER_KEY"),
		})
	})

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		var t task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "task is not JSON", http.StatusBadRequest)
			return
		}
		if t.VUs <= 0 {
			t.VUs = 1
		}
		if t.Seconds <= 0 {
			t.Seconds = 1
		}
		started := time.Now()
		// A real engine would execute t.Script here. This example reports what
		// it actually did rather than inventing numbers it did not measure.
		res := result{
			Status: "completed", Requests: 0, Failures: 0, P95Millis: 0,
			DurationSec: time.Since(started).Seconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})

	log.Printf("perf worker listening on 127.0.0.1:%s for %s", port, extension)
	if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
