// delivery-service/main.go
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Delivery represents a delivery entity
type Delivery struct {
	ID           int       `json:"id"`
	OrderID      int       `json:"orderId"`
	UserID       int       `json:"userId"`
	RestaurantID int       `json:"restaurantId"`
	CourierID    int       `json:"courierId"`
	Status       string    `json:"status"` // "pending", "assigned", "picked_up", "delivered", "cancelled"
	Address      string    `json:"address"`
	EstimatedTime int      `json:"estimatedTime"` // in minutes
	ActualTime    int      `json:"actualTime"`    // in minutes
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Courier represents a courier entity
type Courier struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Available bool  `json:"available"`
	Location  string `json:"location"`
}

var (
	deliveries []Delivery
	couriers   []Courier
	nextDelID  int = 1
	nextCourID int = 1
	mutex      sync.Mutex
)

// Initialize with sample data
func init() {
	// Sample courier
	couriers = append(couriers, Courier{
		ID:        nextCourID,
		Name:      "John Doe",
		Phone:     "555-1234",
		Available: true,
		Location:  "Downtown",
	})
	nextCourID++
}

// Get all deliveries
func getDeliveries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

// Get delivery by ID
func getDelivery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}

	for _, delivery := range deliveries {
		if delivery.ID == id {
			json.NewEncoder(w).Encode(delivery)
			return
		}
	}
	http.Error(w, "Delivery not found", http.StatusNotFound)
}

// Create a new delivery
func createDelivery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var delivery Delivery
	err := json.NewDecoder(r.Body).Decode(&delivery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	mutex.Lock()
	delivery.ID = nextDelID
	nextDelID++
	
	// Default to pending status if not provided
	if delivery.Status == "" {
		delivery.Status = "pending"
	}
	
	// Default estimated delivery time (random between 30-60 minutes)
	if delivery.EstimatedTime == 0 {
		delivery.EstimatedTime = 30 + (int(now.UnixNano() % 30))
	}
	
	delivery.CreatedAt = now
	delivery.UpdatedAt = now
	deliveries = append(deliveries, delivery)
	
	// Try to assign an available courier
	var assignedCourier *Courier
	for i := range couriers {
		if couriers[i].Available {
			assignedCourier = &couriers[i]
			couriers[i].Available = false
			break
		}
	}
	
	// If a courier was found, update the delivery
	if assignedCourier != nil {
		for i := range deliveries {
			if deliveries[i].ID == delivery.ID {
				deliveries[i].CourierID = assignedCourier.ID
				deliveries[i].Status = "assigned"
				
				// Notify order service about the assignment
				go updateOrderStatus(delivery.OrderID, "preparing")
			}
		}
	}
	
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(delivery)
}

// Update delivery status
func updateDeliveryStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
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
		"pending":   true,
		"assigned":  true,
		"picked_up": true,
		"delivered": true,
		"cancelled": true,
	}
	if !validStatuses[statusUpdate.Status] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	var delivery *Delivery
	for i := range deliveries {
		if deliveries[i].ID == id {
			deliveries[i].Status = statusUpdate.Status
			deliveries[i].UpdatedAt = time.Now()
			
			// If delivery is completed or cancelled, make courier available again
			if statusUpdate.Status == "delivered" || statusUpdate.Status == "cancelled" {
				for j := range couriers {
					if couriers[j].ID == deliveries[i].CourierID {
						couriers[j].Available = true
						break
					}
				}
			}
			
			// If delivery is picked up or delivered, update order status
			if statusUpdate.Status == "picked_up" {
				go updateOrderStatus(deliveries[i].OrderID, "out_for_delivery")
			} else if statusUpdate.Status == "delivered" {
				// Calculate actual delivery time (random between estimated-10 and estimated+10)
				deliveries[i].ActualTime = deliveries[i].EstimatedTime + (int(time.Now().UnixNano() % 20) - 10)
				if deliveries[i].ActualTime < 10 {
					deliveries[i].ActualTime = 10
				}
				
				go updateOrderStatus(deliveries[i].OrderID, "delivered")
			}
			
			delivery = &deliveries[i]
			break
		}
	}
	mutex.Unlock()

	if delivery == nil {
		http.Error(w, "Delivery not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(delivery)
}

// Update order status
func updateOrderStatus(orderID int, status string) {
	orderURL := "http://order-service:8082/api/orders/" + strconv.Itoa(orderID) + "/status"
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

// Get deliveries by order ID
func getDeliveriesByOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var orderDeliveries []Delivery
	for _, delivery := range deliveries {
		if delivery.OrderID == orderID {
			orderDeliveries = append(orderDeliveries, delivery)
		}
	}
	json.NewEncoder(w).Encode(orderDeliveries)
}

// Get deliveries by courier ID
func getDeliveriesByCourier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	courierID, err := strconv.Atoi(params["courierId"])
	if err != nil {
		http.Error(w, "Invalid courier ID", http.StatusBadRequest)
		return
	}

	var courierDeliveries []Delivery
	for _, delivery := range deliveries {
		if delivery.CourierID == courierID {
			courierDeliveries = append(courierDeliveries, delivery)
		}
	}
	json.NewEncoder(w).Encode(courierDeliveries)
}

// Get all couriers
func getCouriers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(couriers)
}

// Get courier by ID
func getCourier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid courier ID", http.StatusBadRequest)
		return
	}

	for _, courier := range couriers {
		if courier.ID == id {
			json.NewEncoder(w).Encode(courier)
			return
		}
	}
	http.Error(w, "Courier not found", http.StatusNotFound)
}

// Create a new courier
func createCourier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var courier Courier
	err := json.NewDecoder(r.Body).Decode(&courier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	courier.ID = nextCourID
	nextCourID++
	// Default to available if not specified
	if !courier.Available {
		courier.Available = true
	}
	couriers = append(couriers, courier)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(courier)
}

// Update courier availability
func updateCourierAvailability(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid courier ID", http.StatusBadRequest)
		return
	}

	var availabilityUpdate struct {
		Available bool `json:"available"`
	}
	err = json.NewDecoder(r.Body).Decode(&availabilityUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	var courier *Courier
	for i := range couriers {
		if couriers[i].ID == id {
			couriers[i].Available = availabilityUpdate.Available
			courier = &couriers[i]
			break
		}
	}
	mutex.Unlock()

	if courier == nil {
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(courier)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delivery service is up and running"))
}

func main() {
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Delivery routes
	r.HandleFunc("/api/deliveries", getDeliveries).Methods("GET")
	r.HandleFunc("/api/deliveries/{id}", getDelivery).Methods("GET")
	r.HandleFunc("/api/deliveries", createDelivery).Methods("POST")
	r.HandleFunc("/api/deliveries/{id}/status", updateDeliveryStatus).Methods("PUT")
	
	// Filtered deliveries
	r.HandleFunc("/api/orders/{orderId}/deliveries", getDeliveriesByOrder).Methods("GET")
	r.HandleFunc("/api/couriers/{courierId}/deliveries", getDeliveriesByCourier).Methods("GET")
	
	// Courier routes
	r.HandleFunc("/api/couriers", getCouriers).Methods("GET")
	r.HandleFunc("/api/couriers/{id}", getCourier).Methods("GET")
	r.HandleFunc("/api/couriers", createCourier).Methods("POST")
	r.HandleFunc("/api/couriers/{id}/availability", updateCourierAvailability).Methods("PUT")

	log.Println("Delivery service started on :8084")
	log.Fatal(http.ListenAndServe(":8084", r))
}