package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/time/rate"
)

type URLRequest struct {
	LongURL string `json:"long_url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
	Error    string `json:"error,omitempty"`
}

type Server struct {
	db            *sql.DB
	baseURL       string
	allowedOrigin string
	limiters      map[string]*rate.Limiter
	mu            sync.Mutex
}

func NewServer(db *sql.DB, baseURL string) *Server {
	if baseURL == "" {
		baseURL = os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
	}
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*" // Default to broad for development if not set, but we will set it in docker-compose
	}
	return &Server{
		db:            db,
		baseURL:       baseURL,
		allowedOrigin: allowedOrigin,
		limiters:      make(map[string]*rate.Limiter),
	}
}

func (s *Server) getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, exists := s.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		s.limiters[ip] = limiter
	}

	return limiter
}

func (s *Server) rateLimitMiddleware(next http.HandlerFunc, r rate.Limit, b int) http.HandlerFunc { //In addition to captcha, rate limiting is implemented to prevent abuse
	return func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			var err error
			ip, _, err = net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr // Fallback if no port
			}
		}
		limiter := s.getLimiter(ip, r, b)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, req)
	}
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func generateSlug(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		ret[i] = charset[num.Int64()]
	}
	return string(ret)
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.LongURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var slug string
	length := 6
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		slug = generateSlug(length)
		_, err := s.db.Exec("INSERT INTO urls (slug, long_url) VALUES (?, ?)", slug, req.LongURL)
		if err == nil {
			break
		}
		// If collision, try a longer slug as per readme
		length++
		if i == maxRetries-1 {
			log.Printf("Failed to generate unique slug after %d retries", maxRetries)
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}
	}

	shortURL := fmt.Sprintf("%s/%s", s.baseURL, slug)
	json.NewEncoder(w).Encode(URLResponse{ShortURL: shortURL})
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[1:]
	if slug == "" {
		// If someone just hits localhost:8080/
		fmt.Fprintf(w, "URL Shortener API")
		return
	}

	var longURL string
	err := s.db.QueryRow("SELECT long_url FROM urls WHERE slug = ?", slug).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Printf("Error querying DB: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	// Connection string
	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := os.Getenv("MARIADB_URI")
	if dsn == "" {
		dsn = "explorer:explorerpassword@tcp(127.0.0.1:3306)/url_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// Open connection
	tempDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer tempDB.Close()

	// Connection Pool Settings
	tempDB.SetMaxOpenConns(25)
	tempDB.SetMaxIdleConns(25)
	tempDB.SetConnMaxLifetime(5 * time.Minute)
	tempDB.SetConnMaxIdleTime(2 * time.Minute)

	// Verify connection
	err = tempDB.Ping()
	if err != nil {
		log.Printf("Wait for MariaDB to start... (will retry)")

		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			err = tempDB.Ping()
			if err == nil {
				fmt.Println("Successfully connected to MariaDB!")
				break
			}
		}
		if err != nil {
			log.Fatalf("Could not connect to MariaDB after retries: %v", err)
		}
	} else {
		fmt.Println("Successfully connected to MariaDB!")
	}

	// Create table if not exists, not using migrations for simplicity which could be a huge mistake later on
	_, err = tempDB.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id INT AUTO_INCREMENT PRIMARY KEY,
		slug VARCHAR(10) NOT NULL UNIQUE,
		long_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// Initialize server
	s := NewServer(tempDB, "")

	// Setup HTTP server
	// POST /shorten: 2 requests per minute (approx 0.033 req/sec)
	http.HandleFunc("/shorten", s.corsMiddleware(s.rateLimitMiddleware(s.handleShorten, rate.Every(30*time.Second), 2)))
	// GET /: 75 requests per minute (1.25 req/sec)
	http.HandleFunc("/", s.rateLimitMiddleware(s.handleRedirect, rate.Every(800*time.Millisecond), 75))

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("Server starting on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
