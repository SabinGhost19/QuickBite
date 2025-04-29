package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Adăugat pentru încărcarea fișierului .env
)

// Order represents a food order
type Order struct {
	ID           int         `json:"id"`
	UserID       int         `json:"userId"`
	RestaurantID int         `json:"restaurantId"`
	Items        []OrderItem `json:"items"`
	TotalAmount  float64     `json:"totalAmount"`
	Status       string      `json:"status"` // "created", "paid", "preparing", "out_for_delivery", "delivered", "cancelled"
	Address      string      `json:"address"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

// OrderItem represents an item in the order
type OrderItem struct {
	MenuItemID int     `json:"menuItemId"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
}

// Config holds service configuration from environment variables
type Config struct {
	Port                   string
	PaymentServiceURL      string
	DeliveryServiceURL     string
	NotificationServiceURL string
	AllowedOrigins         string
}

// Global variables
var (
	orders []Order
	nextID int = 1
	mutex  sync.Mutex
	config Config
)

// Initialize with sample data and load config
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
		Port:                   getEnv("ORDER_SERVICE_PORT", "8082"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://payment-service:8083/api/payments"),
		DeliveryServiceURL:     getEnv("DELIVERY_SERVICE_URL", "http://delivery-service:8084/api/deliveries"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8085/api/notifications"),
		AllowedOrigins:         getEnv("ALLOWED_ORIGINS", "http://localhost:3205"),
	}

	// Sample order
	now := time.Now()
	orders = append(orders, Order{
		ID:           nextID,
		UserID:       1,
		RestaurantID: 1,
		Items: []OrderItem{
			{
				MenuItemID: 1,
				Name:       "Margherita Pizza",
				Price:      12.99,
				Quantity:   2,
			},
		},
		TotalAmount: 25.98,
		Status:      "created",
		Address:     "123 Main St, City",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	nextID++
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Get all orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Get order by ID
func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	for _, order := range orders {
		if order.ID == id {
			json.NewEncoder(w).Encode(order)
			return
		}
	}
	http.Error(w, "Order not found", http.StatusNotFound)
}

// Create a new order
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate total amount from items
	totalAmount := 0.0
	for _, item := range order.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	now := time.Now()
	mutex.Lock()
	order.ID = nextID
	nextID++
	order.TotalAmount = totalAmount
	order.Status = "created"
	order.CreatedAt = now
	order.UpdatedAt = now
	orders = append(orders, order)
	mutex.Unlock()

	// Notify payment service about new order
	go notifyPaymentService(order)
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// Update order status
func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"created":          true,
		"paid":             true,
		"preparing":        true,
		"out_for_delivery": true,
		"delivered":        true,
		"cancelled":        true,
	}
	if !validStatuses[statusUpdate.Status] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = statusUpdate.Status
			orders[i].UpdatedAt = time.Now()
			
			// If status changed to paid, notify delivery service
			if statusUpdate.Status == "paid" {
				go notifyDeliveryService(orders[i])
			}
			
			// If status changed to preparing or out_for_delivery or delivered, notify notification service
			if statusUpdate.Status == "preparing" || statusUpdate.Status == "out_for_delivery" || statusUpdate.Status == "delivered" {
				go notifyNotificationService(orders[i])
			}
			
			mutex.Unlock()
			json.NewEncoder(w).Encode(orders[i])
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Order not found", http.StatusNotFound)
}

// Get orders by user ID
func getOrdersByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var userOrders []Order
	for _, order := range orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}
	json.NewEncoder(w).Encode(userOrders)
}

// Get orders by restaurant ID
func getOrdersByRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	restaurantID, err := strconv.Atoi(params["restaurantId"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var restaurantOrders []Order
	for _, order := range orders {
		if order.RestaurantID == restaurantID {
			restaurantOrders = append(restaurantOrders, order)
		}
	}
	json.NewEncoder(w).Encode(restaurantOrders)
}

// Cancel an order
func cancelOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, order := range orders {
		if order.ID == id {
			// Only allow cancellation if order is not out for delivery or delivered
			if order.Status == "out_for_delivery" || order.Status == "delivered" {
				mutex.Unlock()
				http.Error(w, "Cannot cancel order that is out for delivery or already delivered", http.StatusBadRequest)
				return
			}
			orders[i].Status = "cancelled"
			orders[i].UpdatedAt = time.Now()
			
			// Notify notification service about cancelled order
			go notifyNotificationService(orders[i])
			
			mutex.Unlock()
			json.NewEncoder(w).Encode(orders[i])
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Order not found", http.StatusNotFound)
}

// Notify payment service about new order
func notifyPaymentService(order Order) {
	paymentURL := config.PaymentServiceURL
	paymentData := map[string]interface{}{
		"orderId":     order.ID,
		"userId":      order.UserID,
		"amount":      order.TotalAmount,
		"description": fmt.Sprintf("Payment for order #%d", order.ID),
	}
	
	jsonData, err := json.Marshal(paymentData)
	if err != nil {
		log.Printf("Error marshaling payment data: %v", err)
		return
	}
	
	resp, err := http.Post(paymentURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error notifying payment service: %v", err)
		return
	}
	defer resp.Body.Close()
	
	log.Printf("Payment service notification status: %d", resp.StatusCode)
}

// Notify delivery service about paid order
func notifyDeliveryService(order Order) {
	deliveryURL := config.DeliveryServiceURL
	deliveryData := map[string]interface{}{
		"orderId":      order.ID,
		"userId":       order.UserID,
		"restaurantId": order.RestaurantID,
		"address":      order.Address,
		"status":       "pending",
	}
	
	jsonData, err := json.Marshal(deliveryData)
	if err != nil {
		log.Printf("Error marshaling delivery data: %v", err)
		return
	}
	
	resp, err := http.Post(deliveryURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error notifying delivery service: %v", err)
		return
	}
	defer resp.Body.Close()
	
	log.Printf("Delivery service notification status: %d", resp.StatusCode)
}

// Notify notification service about order status changes
func notifyNotificationService(order Order) {
	notificationURL := config.NotificationServiceURL
	notificationData := map[string]interface{}{
		"userId":    order.UserID,
		"type":      "order_update",
		"message":   fmt.Sprintf("Your order #%d status has been updated to: %s", order.ID, order.Status),
		"orderId":   order.ID,
		"status":    order.Status,
	}
	
	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		log.Printf("Error marshaling notification data: %v", err)
		return
	}
	
	resp, err := http.Post(notificationURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error notifying notification service: %v", err)
		return
	}
	defer resp.Body.Close()
	
	log.Printf("Notification service notification status: %d", resp.StatusCode)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order service is up and running"))
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

	// Order routes
	r.HandleFunc("/api/orders", getOrders).Methods("GET")
	r.HandleFunc("/api/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/api/orders", createOrder).Methods("POST")
	r.HandleFunc("/api/orders/{id}/status", updateOrderStatus).Methods("PUT")
	r.HandleFunc("/api/orders/{id}/cancel", cancelOrder).Methods("PUT")
	
	// Filtered orders
	r.HandleFunc("/api/orders/user/{userId}/orders", getOrdersByUser).Methods("GET")
	r.HandleFunc("/api/restaurants/{restaurantId}/orders", getOrdersByRestaurant).Methods("GET")

	// Apply CORS middleware
	handler := enableCORS(r)

	log.Printf("Order service configuration:")
	log.Printf("- Port: %s", config.Port)
	log.Printf("- Payment Service URL: %s", config.PaymentServiceURL)
	log.Printf("- Delivery Service URL: %s", config.DeliveryServiceURL)
	log.Printf("- Notification Service URL: %s", config.NotificationServiceURL)
	log.Printf("- Allowed Origins: %s", config.AllowedOrigins)

	log.Printf("Order service started on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}