import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../models/user.model';
import { Restaurant, MenuItem } from '../models/restaurant.model';
import { Order } from '../models/order.model';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private API_BASE_URL = 'http://localhost:8080';
  private API_ENDPOINTS = {
    users: `${this.API_BASE_URL}/api/users`,
    restaurants: 'http://localhost:8081/api/restaurants',
    orders: 'http://localhost:8082/api/orders',
    payments: 'http://localhost:8083/api/payments',
    deliveries: 'http://localhost:8084/api/deliveries',
    notifications: 'http://localhost:8085/api/notifications',
  };

  constructor(private http: HttpClient) {}

  // User API calls
  getUser(userId: number): Observable<User> {
    return this.http.get<User>(`${this.API_ENDPOINTS.users}/${userId}`);
  }

  // Restaurant API calls
  getRestaurants(): Observable<Restaurant[]> {
    return this.http.get<Restaurant[]>(this.API_ENDPOINTS.restaurants);
  }

  getRestaurantById(restaurantId: number): Observable<Restaurant> {
    return this.http.get<Restaurant>(
      `${this.API_ENDPOINTS.restaurants}/${restaurantId}`
    );
  }

  // Order API calls
  getOrders(userId: number): Observable<Order[]> {
    return this.http.get<Order[]>(
      `${this.API_ENDPOINTS.orders}/user/${userId}/orders`
    );
  }

  getOrderById(orderId: number): Observable<Order> {
    return this.http.get<Order>(`${this.API_ENDPOINTS.orders}/${orderId}`);
  }

  createOrder(order: Order): Observable<Order> {
    return this.http.post<Order>(this.API_ENDPOINTS.orders, order);
  }
}
