package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// getenv возвращает значение переменной окружения или значение по умолчанию.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// newProxy создаёт обратный прокси на указанный бэкенд.
func newProxy(rawURL string) *httputil.ReverseProxy {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("некорректный URL бэкенда %q: %v", rawURL, err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	port := getenv("PORT", "8000")

	monolithURL := getenv("MONOLITH_URL", "http://monolith:8080")
	moviesURL := getenv("MOVIES_SERVICE_URL", "http://movies-service:8081")
	eventsURL := getenv("EVENTS_SERVICE_URL", "http://events-service:8082")

	monolith := newProxy(monolithURL)
	movies := newProxy(moviesURL)
	events := newProxy(eventsURL)

	// Фиче-флаг постепенной миграции.
	gradual := strings.EqualFold(getenv("GRADUAL_MIGRATION", "false"), "true")
	percent, err := strconv.Atoi(getenv("MOVIES_MIGRATION_PERCENT", "0"))
	if err != nil || percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	log.Printf("proxy: порт=%s монолит=%s movies=%s events=%s", port, monolithURL, moviesURL, eventsURL)
	log.Printf("proxy: постепенная миграция=%v, процент трафика на movies=%d%%", gradual, percent)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","gradual_migration":%v,"movies_migration_percent":%d}`, gradual, percent)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var target *httputil.ReverseProxy
		var name string

		switch {
		// События уходят в events-service
		case strings.HasPrefix(path, "/api/events"):
			target, name = events, "events-service"

		case path == "/api/movies/health":
			target, name = movies, "movies-service"

		// Strangler Fig фича-тоглер
		case strings.HasPrefix(path, "/api/movies"):
			if gradual && rand.Intn(100) < percent {
				target, name = movies, "movies-service"
			} else {
				target, name = monolith, "monolith"
			}

		// Всё остальное пока живёт в монолите.
		default:
			target, name = monolith, "monolith"
		}

		log.Printf("%s %s -> %s", r.Method, path, name)
		w.Header().Set("X-Upstream", name)
		target.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("proxy: сервер остановлен: %v", err)
	}
}
