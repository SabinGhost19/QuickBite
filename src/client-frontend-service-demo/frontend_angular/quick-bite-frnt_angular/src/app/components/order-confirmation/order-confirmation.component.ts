import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { Order } from '../../models/order.model';

@Component({
  selector: 'app-order-confirmation',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './order-confirmation.component.html',
  styleUrls: ['./order-confirmation.component.css'],
})
export class OrderConfirmationComponent implements OnInit {
  order: Order | null = null;
  loading = true;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiService: ApiService
  ) {}

  ngOnInit(): void {
    const orderId = Number(this.route.snapshot.paramMap.get('id'));
    this.fetchOrderDetails(orderId);
  }

  fetchOrderDetails(orderId: number): void {
    this.loading = true;
    this.apiService.getOrderById(orderId).subscribe({
      next: (data) => {
        this.order = data;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error fetching order details:', error);
        // Use mock data for demo
        this.order = {
          id: orderId,
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
          status: 'created',
          address: '123 Main St, City',
          createdAt: new Date().toISOString(),
        };
        this.loading = false;
      },
    });
  }

  viewOrders(): void {
    this.router.navigate(['/orders']);
  }
}
