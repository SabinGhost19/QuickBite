import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { Order } from '../../models/order.model';

@Component({
  selector: 'app-orders',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './orders.component.html',
  styleUrls: ['./orders.component.css'],
})
export class OrdersComponent implements OnInit {
  orders: Order[] = [];
  loading = true;

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.fetchOrders();
  }

  fetchOrders(): void {
    this.loading = true;
    this.apiService.getOrders(1).subscribe({
      next: (data) => {
        this.orders = data;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error fetching orders:', error);
        // Use mock data for demo
        this.orders = [
          {
            id: 1,
            userId: 1,
            restaurantId: 1,
            items: [
              {
                menuItemId: 1,
                name: 'Margherita Pizza',
                price: 12.99,
                quantity: 2,
              },
            ],
            totalAmount: 25.98,
            status: 'delivered',
            address: '123 Main St, City',
            createdAt: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
          },
          {
            id: 2,
            userId: 1,
            restaurantId: 2,
            items: [
              {
                menuItemId: 3,
                name: 'Chicken Burger',
                price: 9.99,
                quantity: 1,
              },
              {
                menuItemId: 4,
                name: 'French Fries',
                price: 4.99,
                quantity: 1,
              },
            ],
            totalAmount: 14.98,
            status: 'processing',
            address: '123 Main St, City',
            createdAt: new Date().toISOString(),
          },
        ];
        this.loading = false;
      },
    });
  }

  getStatusClass(status: string): string {
    switch (status.toLowerCase()) {
      case 'created':
        return 'status-created';
      case 'processing':
        return 'status-processing';
      case 'delivered':
        return 'status-delivered';
      case 'cancelled':
        return 'status-cancelled';
      default:
        return '';
    }
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
}
