import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../models/user.model';
import { Restaurant, MenuItem } from '../models/restaurant.model';
import { Order } from '../models/order.model';
import { ConfigService } from './config.service';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private API_ENDPOINTS: any;

  constructor(private http: HttpClient, private config: ConfigService) {
    this.API_ENDPOINTS = {
      users: `${this.config.getApiBaseUrl()}/api/users`,
      restaurants: `${this.config.getRestaurantsServiceUrl()}/api/restaurants`,
      orders: `${this.config.getOrdersServiceUrl()}/api/orders`,
      payments: `${this.config.getPaymentsServiceUrl()}/api/payments`,
      deliveries: `${this.config.getDeliveriesServiceUrl()}/api/deliveries`,
      notifications: `${this.config.getNotificationsServiceUrl()}/api/notifications`,
    };
  }

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
