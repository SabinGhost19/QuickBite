// payment-service/main.go
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
	"github.com/joho/godotenv"
)

// Payment represents a payment transaction
type Payment struct {
	ID          int       `json:"id"`
	OrderID     int       `json:"orderId"`
	UserID      int       `json:"userId"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // "pending", "completed", "failed", "refunded"
	Method      string    `json:"method"` // "card", "cash", etc.
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var (
	payments []Payment
	nextID   int = 1
	mutex    sync.Mutex
)

// Get all payments
func getPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// Get payment by ID
func getPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	for _, payment := range payments {
		if payment.ID == id {
			json.NewEncoder(w).Encode(payment)
			return
		}
	}
	http.Error(w, "Payment not found", http.StatusNotFound)
}

// Create a new payment
func createPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var payment Payment
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	mutex.Lock()
	payment.ID = nextID
	nextID++
	// Default to pending status if not provided
	if payment.Status == "" {
		payment.Status = "pending"
	}
	// Default to card method if not provided
	if payment.Method == "" {
		payment.Method = "card"
	}
	payment.CreatedAt = now
	payment.UpdatedAt = now
	payments = append(payments, payment)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// Process payment (simulate payment processing)
func processPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	var processRequest struct {
		Method string `json:"method"`
	}
	err = json.NewDecoder(r.Body).Decode(&processRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate payment method
	validMethods := map[string]bool{
		"card": true,
		"cash": true,
	}
	if !validMethods[processRequest.Method] {
		http.Error(w, "Invalid payment method", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	var payment *Payment
	for i := range payments {
		if payments[i].ID == id {
			payment = &payments[i]
			break
		}
	}

	if payment == nil {
		mutex.Unlock()
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	// Process payment (in a real system, this would integrate with payment gateways)
	// Here we're simulating payment processing - 90% success rate
	success := time.Now().UnixNano()%10 != 0 // 90% success rate
	
	payment.Method = processRequest.Method
	payment.Status = "completed"
	if !success {
		payment.Status = "failed"
	}
	payment.UpdatedAt = time.Now()
	mutex.Unlock()

	// If payment is successful, update order status
	if success {
		go updateOrderStatus(payment.OrderID, "paid")
	}

	json.NewEncoder(w).Encode(payment)
}

// Update order status after payment processing
func updateOrderStatus(orderID int, status string) {
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	orderURL := fmt.Sprintf("%s/api/orders/%d/status", orderServiceURL, orderID)
	
	statusUpdate := map[string]string{
		"status": status,
	}
	
	jsonData, err := json.Marshal(statusUpdate)
	if err != nil {
		log.Printf("Error marshaling order status update: %v", err)
		return
	}
	
	req, err := http.NewRequest("PUT", orderURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating order status update request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		return
	}
	defer resp.Body.Close()
	
	log.Printf("Order status update response: %d", resp.StatusCode)
}

// Get payments by order ID
func getPaymentsByOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var orderPayments []Payment
	for _, payment := range payments {
		if payment.OrderID == orderID {
			orderPayments = append(orderPayments, payment)
		}
	}
	json.NewEncoder(w).Encode(orderPayments)
}

// Refund a payment
func refundPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, payment := range payments {
		if payment.ID == id {
			// Only allow refunds for completed payments
			if payment.Status != "completed" {
				mutex.Unlock()
				http.Error(w, "Only completed payments can be refunded", http.StatusBadRequest)
				return
			}
			payments[i].Status = "refunded"
			payments[i].UpdatedAt = time.Now()
			
			// Update order status to cancelled when payment is refunded
			go updateOrderStatus(payment.OrderID, "cancelled")
			
			mutex.Unlock()
			json.NewEncoder(w).Encode(payments[i])
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Payment not found", http.StatusNotFound)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment service is up and running"))
}

func loadEnv() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	// Log environment variables (for debugging)
	log.Println("Environment loaded successfully")
	log.Printf("Server running on %s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Printf("Order service URL: %s", os.Getenv("ORDER_SERVICE_URL"))
}

func main() {
	// Load environment variables
	loadEnv()
	
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Payment routes
	r.HandleFunc("/api/payments", getPayments).Methods("GET")
	r.HandleFunc("/api/payments/{id}", getPayment).Methods("GET")
	r.HandleFunc("/api/payments", createPayment).Methods("POST")
	r.HandleFunc("/api/payments/{id}/process", processPayment).Methods("PUT")
	r.HandleFunc("/api/payments/{id}/refund", refundPayment).Methods("PUT")
	
	// Filtered payments
	r.HandleFunc("/api/orders/{orderId}/payments", getPaymentsByOrder).Methods("GET")

	// Get server address from environment variables
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	
	log.Printf("Payment service started on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}