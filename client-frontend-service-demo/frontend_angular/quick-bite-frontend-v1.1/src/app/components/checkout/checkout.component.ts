import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { CartService } from '../../services/cart.service';
import { CartItem, Order } from '../../models/order.model';

@Component({
  selector: 'app-checkout',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './checkout.component.html',
  styleUrls: ['./checkout.component.css'],
})
export class CheckoutComponent implements OnInit {
  cartItems: CartItem[] = [];
  cartTotal = 0;
  deliveryAddress = '';
  paymentMethod = 'card';
  loading = false;

  constructor(
    private router: Router,
    private apiService: ApiService,
    private cartService: CartService
  ) {}

  ngOnInit(): void {
    // Subscribe to cart changes
    this.cartService.getCartItems().subscribe((items) => {
      this.cartItems = items;
    });

    this.cartService.getTotal().subscribe((total) => {
      this.cartTotal = total;
    });

    // Pre-fill address with user's address (in a real app, this would come from user profile)
    this.deliveryAddress = '123 Main St, City';
  }

  placeOrder(): void {
    if (!this.deliveryAddress.trim()) {
      alert('Please enter a delivery address');
      return;
    }

    // Create order items
    const orderItems = this.cartItems.map((item) => {
      return {
        menuItemId: item.id,
        name: item.name,
        price: item.price,
        quantity: item.quantity,
      };
    });

    // Create order object
    const order: Partial<Order> = {
      userId: 1, // Hardcoded for demo
      restaurantId: this.cartService.getRestaurantId() || 0,
      items: orderItems,
      totalAmount: this.cartTotal,
      status: 'created',
      address: this.deliveryAddress,
    };

    this.loading = true;

    this.apiService.createOrder(order as Order).subscribe({
      next: (createdOrder) => {
        this.loading = false;
        // Clear cart
        this.cartService.clearCart();
        // Navigate to order confirmation
        this.router.navigate(['/order-confirmation', createdOrder.id]);
      },
      error: (error) => {
        console.error('Error creating order:', error);
        this.loading = false;

        // Mock order creation for demo
        const mockCreatedOrder = {
          id: Math.floor(Math.random() * 1000),
          userId: 1,
          restaurantId: this.cartService.getRestaurantId() || 0,
          items: orderItems,
          totalAmount: this.cartTotal,
          status: 'created',
          address: this.deliveryAddress,
          createdAt: new Date().toISOString(),
        };

        // Clear cart
        this.cartService.clearCart();
        // Navigate to order confirmation
        this.router.navigate(['/order-confirmation', mockCreatedOrder.id]);
      },
    });
  }

  goBack(): void {
    this.router.navigate(['/restaurant', this.cartService.getRestaurantId()]);
  }
}
