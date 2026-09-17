package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"astroleap/web"
)

type ScoreEntry struct {
	PlayerName string    `json:"player_name"`
	Score      int       `json:"score"`
	Timestamp  time.Time `json:"timestamp"`
}

type HighScoreStore struct {
	mu     sync.RWMutex
	scores []ScoreEntry
}

func (s *HighScoreStore) Add(entry ScoreEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry.Timestamp = time.Now()
	s.scores = append(s.scores, entry)
}

func (s *HighScoreStore) GetTop(limit int) []ScoreEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copied := make([]ScoreEntry, len(s.scores))
	copy(copied, s.scores)
	if len(copied) > limit {
		return copied[:limit]
	}
	return copied
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := &HighScoreStore{}

	http.Handle("/", http.FileServer(http.FS(web.FS)))

	http.HandleFunc("/api/v1/scores", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(store.GetTop(10))
		case http.MethodPost:
			var entry ScoreEntry
			if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
				http.Error(w, `{"error":"invalid payload"}`, http.StatusBadRequest)
				return
			}
			store.Add(entry)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Printf("AstroLeap game server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
