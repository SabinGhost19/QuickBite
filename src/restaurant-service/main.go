// restaurant-service/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Restaurant represents a restaurant entity
type Restaurant struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Cuisine  string   `json:"cuisine"`
	Rating   float64  `json:"rating"`
	MenuItems []MenuItem `json:"menuItems"`
}

// MenuItem represents a menu item
type MenuItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

var (
	restaurants []Restaurant
	nextRestID  int = 1
	nextItemID  int = 1
	mutex       sync.Mutex
)

// Initialize with sample data
func init() {
	// Sample restaurant with menu items
	restaurants = append(restaurants, Restaurant{
		ID:      nextRestID,
		Name:    "Tasty Bites",
		Address: "123 Main St",
		Cuisine: "Italian",
		Rating:  4.5,
		MenuItems: []MenuItem{
			{
				ID:          nextItemID,
				Name:        "Margherita Pizza",
				Description: "Classic pizza with tomato sauce, mozzarella, and basil",
				Price:       12.99,
				Category:    "Main",
			},
		},
	})
	nextRestID++
	nextItemID++
}

// Get all restaurants
func getRestaurants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

// Get restaurant by ID
func getRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	for _, restaurant := range restaurants {
		if restaurant.ID == id {
			json.NewEncoder(w).Encode(restaurant)
			return
		}
	}
	http.Error(w, "Restaurant not found", http.StatusNotFound)
}

// Create a new restaurant
func createRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var restaurant Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	restaurant.ID = nextRestID
	nextRestID++
	restaurants = append(restaurants, restaurant)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(restaurant)
}

// Update a restaurant
func updateRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var updatedRestaurant Restaurant
	err = json.NewDecoder(r.Body).Decode(&updatedRestaurant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, restaurant := range restaurants {
		if restaurant.ID == id {
			updatedRestaurant.ID = id
			restaurants[i] = updatedRestaurant
			mutex.Unlock()
			json.NewEncoder(w).Encode(updatedRestaurant)
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Restaurant not found", http.StatusNotFound)
}

// Delete a restaurant
func deleteRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, restaurant := range restaurants {
		if restaurant.ID == id {
			restaurants = append(restaurants[:i], restaurants[i+1:]...)
			mutex.Unlock()
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Restaurant not found", http.StatusNotFound)
}

// Get all menu items for a restaurant
func getMenuItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	for _, restaurant := range restaurants {
		if restaurant.ID == id {
			json.NewEncoder(w).Encode(restaurant.MenuItems)
			return
		}
	}
	http.Error(w, "Restaurant not found", http.StatusNotFound)
}

// Add a menu item to a restaurant
func addMenuItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var menuItem MenuItem
	err = json.NewDecoder(r.Body).Decode(&menuItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	for i, restaurant := range restaurants {
		if restaurant.ID == id {
			menuItem.ID = nextItemID
			nextItemID++
			restaurants[i].MenuItems = append(restaurants[i].MenuItems, menuItem)
			mutex.Unlock()
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(menuItem)
			return
		}
	}
	mutex.Unlock()
	http.Error(w, "Restaurant not found", http.StatusNotFound)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Restaurant service is up and running"))
}

func main() {
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Restaurant routes
	r.HandleFunc("/api/restaurants", getRestaurants).Methods("GET")
	r.HandleFunc("/api/restaurants/{id}", getRestaurant).Methods("GET")
	r.HandleFunc("/api/restaurants", createRestaurant).Methods("POST")
	r.HandleFunc("/api/restaurants/{id}", updateRestaurant).Methods("PUT")
	r.HandleFunc("/api/restaurants/{id}", deleteRestaurant).Methods("DELETE")

	// Menu item routes
	r.HandleFunc("/api/restaurants/{id}/menu", getMenuItems).Methods("GET")
	r.HandleFunc("/api/restaurants/{id}/menu", addMenuItem).Methods("POST")

	log.Println("Restaurant service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}