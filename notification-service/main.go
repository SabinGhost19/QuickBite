// notification-service/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Notification represents a notification entity
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Type      string    `json:"type"` // "order_update", "delivery_update", "payment_update", etc.
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	OrderID   int       `json:"orderId,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// Config holds service configuration from environment variables
type Config struct {
	Port           string
	AllowedOrigins string
}

var (
	notifications []Notification
	nextID        int = 1
	mutex         sync.Mutex
	config        Config
)

// Initialize configuration
func init() {
	// Încarcă variabilele de mediu din fișierul .env
	// Nu va genera eroare dacă fișierul nu există
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: No .env file found or error loading it. Using environment variables and defaults.")
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Load configuration
	config = Config{
		// Default values
		Port:           getEnv("NOTIFICATION_SERVICE_PORT", "8085"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3205"),
	}
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Get all notifications
func getNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// Get notification by ID
func getNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	for _, notification := range notifications {
		if notification.ID == id {
			json.NewEncoder(w).Encode(notification)
			return
		}
	}
	http.Error(w, "Notification not found", http.StatusNotFound)
}

// Create a new notification
func createNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	notification.ID = nextID
	nextID++
	notification.Read = false
	notification.CreatedAt = time.Now()
	notifications = append(notifications, notification)
	mutex.Unlock()

	// Here in a real application we would:
	// 1. Save to a database
	// 2. Send push notification, email, SMS, etc.
	// For demo purposes, we just log it
	log.Printf("New notification for user %d: %s", notification.UserID, notification.Message)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

// Mark notification as read
func markAsRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, notification := range notifications {
		if notification.ID == id {
			notifications[i].Read = true
			mutex.Unlock()
			json.NewEncoder(w).Encode(notifications[i])
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Notification not found", http.StatusNotFound)
}

// Get notifications by user ID
func getNotificationsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var userNotifications []Notification
	for _, notification := range notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}
	json.NewEncoder(w).Encode(userNotifications)
}

// Get unread notifications by user ID
func getUnreadNotificationsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var unreadNotifications []Notification
	for _, notification := range notifications {
		if notification.UserID == userID && !notification.Read {
			unreadNotifications = append(unreadNotifications, notification)
		}
	}
	json.NewEncoder(w).Encode(unreadNotifications)
}

// Get notifications by order ID
func getNotificationsByOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var orderNotifications []Notification
	for _, notification := range notifications {
		if notification.OrderID == orderID {
			orderNotifications = append(orderNotifications, notification)
		}
	}
	json.NewEncoder(w).Encode(orderNotifications)
}

// Mark all user notifications as read
func markAllAsRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, notification := range notifications {
		if notification.UserID == userID && !notification.Read {
			notifications[i].Read = true
		}
	}
	mutex.Unlock()

	// Return the updated user notifications
	var userNotifications []Notification
	for _, notification := range notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}
	json.NewEncoder(w).Encode(userNotifications)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification service is up and running"))
}

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", config.AllowedOrigins)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Notification routes
	r.HandleFunc("/api/notifications", getNotifications).Methods("GET")
	r.HandleFunc("/api/notifications/{id}", getNotification).Methods("GET")
	r.HandleFunc("/api/notifications", createNotification).Methods("POST")
	r.HandleFunc("/api/notifications/{id}/read", markAsRead).Methods("PUT")
	
	// Filtered notifications
	r.HandleFunc("/api/users/{userId}/notifications", getNotificationsByUser).Methods("GET")
	r.HandleFunc("/api/users/{userId}/notifications/unread", getUnreadNotificationsByUser).Methods("GET")
	r.HandleFunc("/api/orders/{orderId}/notifications", getNotificationsByOrder).Methods("GET")
	r.HandleFunc("/api/users/{userId}/notifications/read-all", markAllAsRead).Methods("PUT")

	// Apply CORS middleware
	handler := enableCORS(r)

	log.Printf("Notification service configuration:")
	log.Printf("- Port: %s", config.Port)
	log.Printf("- Allowed Origins: %s", config.AllowedOrigins)

	log.Printf("Notification service started on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}