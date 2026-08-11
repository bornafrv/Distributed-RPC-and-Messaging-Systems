package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MemoryEvent struct {
	EventType   string `json:"event_type"`
	Service     string `json:"service"`
	MemoryMB    uint64 `json:"memory_mb"`
	ThresholdMB uint64 `json:"threshold_mb"`
	Timestamp   string `json:"timestamp"`
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event MemoryEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n[ALERT] HIGH MEMORY USAGE DETECTED\n")
	fmt.Printf("Service      : %s\n", event.Service)
	fmt.Printf("Event Type   : %s\n", event.EventType)
	fmt.Printf("Memory Usage : %d MB\n", event.MemoryMB)
	fmt.Printf("Threshold    : %d MB\n", event.ThresholdMB)
	fmt.Printf("Timestamp    : %s\n\n", event.Timestamp)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alert received"))
}

func main() {
	http.HandleFunc("/alert", alertHandler)

	fmt.Println("Subscriber is listening on port 9090...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
