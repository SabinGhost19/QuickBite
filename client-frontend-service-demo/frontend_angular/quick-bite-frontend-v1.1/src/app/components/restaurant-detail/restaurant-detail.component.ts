import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { CartService } from '../../services/cart.service';
import { Restaurant, MenuItem } from '../../models/restaurant.model';
import { CartItem } from '../../models/order.model';

@Component({
  selector: 'app-restaurant-detail',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './restaurant-detail.component.html',
  styleUrls: ['./restaurant-detail.component.css'],
})
export class RestaurantDetailComponent implements OnInit {
  restaurant: Restaurant | null = null;
  loading = true;
  cartItems: CartItem[] = [];
  cartTotal = 0;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiService: ApiService,
    private cartService: CartService
  ) {}

  ngOnInit(): void {
    const restaurantId = Number(this.route.snapshot.paramMap.get('id'));
    this.fetchRestaurantDetails(restaurantId);

    // Subscribe to cart changes
    this.cartService.getCartItems().subscribe((items) => {
      this.cartItems = items;
    });

    this.cartService.getTotal().subscribe((total) => {
      this.cartTotal = total;
    });
  }

  fetchRestaurantDetails(restaurantId: number): void {
    this.loading = true;
    this.apiService.getRestaurantById(restaurantId).subscribe({
      next: (data) => {
        this.restaurant = data;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error fetching restaurant details:', error);
        // Use mock data for demo
        this.restaurant = {
          id: restaurantId,
          name: 'Tasty Bites',
          address: '123 Main St',
          cuisine: 'Italian',
          rating: 4.5,
          menuItems: [
            {
              id: 1,
              name: 'Margherita Pizza',
              description:
                'Classic pizza with tomato sauce, mozzarella, and basil',
              price: 12.99,
              category: 'Main',
            },
            {
              id: 2,
              name: 'Spaghetti Carbonara',
              description: 'Pasta with egg, cheese, pancetta, and black pepper',
              price: 14.99,
              category: 'Main',
            },
            {
              id: 3,
              name: 'Tiramisu',
              description: 'Coffee-flavored Italian dessert',
              price: 7.99,
              category: 'Dessert',
            },
            {
              id: 4,
              name: 'Caesar Salad',
              description:
                'Romaine lettuce, croutons, parmesan cheese, and Caesar dressing',
              price: 8.99,
              category: 'Starter',
            },
          ],
        };
        this.loading = false;
      },
    });
  }

  addToCart(item: MenuItem): void {
    if (this.restaurant) {
      this.cartService.addToCart(
        item.id,
        item.name,
        item.price,
        this.restaurant.id
      );
    }
  }

  updateCartItemQuantity(itemId: number, newQuantity: number): void {
    this.cartService.updateCartItemQuantity(itemId, newQuantity);
  }

  goToCheckout(): void {
    this.router.navigate(['/checkout']);
  }

  goBack(): void {
    this.router.navigate(['/restaurants']);
  }

  getMenuItemsByCategory(category: string): MenuItem[] {
    if (!this.restaurant || !this.restaurant.menuItems) {
      return [];
    }
    return this.restaurant.menuItems.filter(
      (item) => item.category === category
    );
  }

  getUniqueCategories(): string[] {
    if (!this.restaurant || !this.restaurant.menuItems) {
      return [];
    }
    const categories = this.restaurant.menuItems.map((item) => item.category);
    return [...new Set(categories)];
  }
}
