import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { CartItem } from '../models/order.model';

@Injectable({
  providedIn: 'root',
})
export class CartService {
  private cartItems: CartItem[] = [];
  private restaurantId: number | null = null;
  private cartSubject = new BehaviorSubject<CartItem[]>([]);
  private totalSubject = new BehaviorSubject<number>(0);

  constructor() {}

  getCartItems(): Observable<CartItem[]> {
    return this.cartSubject.asObservable();
  }

  getTotal(): Observable<number> {
    return this.totalSubject.asObservable();
  }

  getRestaurantId(): number | null {
    return this.restaurantId;
  }

  addToCart(
    itemId: number,
    itemName: string,
    itemPrice: number,
    restaurantId: number
  ): void {
    // Check if we're adding items from a different restaurant
    if (this.restaurantId !== null && this.restaurantId !== restaurantId) {
      // Clear the cart if adding items from a different restaurant
      this.clearCart();
    }

    // Set the restaurant ID if it's not set
    if (this.restaurantId === null) {
      this.restaurantId = restaurantId;
    }

    // Check if item already in cart
    const existingItemIndex = this.cartItems.findIndex(
      (item) => item.id === itemId
    );

    if (existingItemIndex >= 0) {
      // Increment quantity if item already exists
      this.cartItems[existingItemIndex].quantity++;
    } else {
      // Add new item to cart
      this.cartItems.push({
        id: itemId,
        name: itemName,
        price: itemPrice,
        quantity: 1,
      });
    }

    this.updateCart();
  }

  updateCartItemQuantity(itemId: number, newQuantity: number): void {
    if (newQuantity <= 0) {
      // Remove item from cart if quantity is 0 or less
      this.cartItems = this.cartItems.filter((item) => item.id !== itemId);
    } else {
      // Update quantity
      const item = this.cartItems.find((item) => item.id === itemId);
      if (item) {
        item.quantity = newQuantity;
      }
    }

    this.updateCart();
  }

  clearCart(): void {
    this.cartItems = [];
    this.restaurantId = null;
    this.updateCart();
  }

  private updateCart(): void {
    this.cartSubject.next([...this.cartItems]);

    // Calculate total
    const total = this.cartItems.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );
    this.totalSubject.next(total);
  }
}
