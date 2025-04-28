export interface Order {
  id: number;
  userId: number;
  restaurantId: number;
  items: OrderItem[];
  totalAmount: number;
  status: string;
  address: string;
  createdAt: string;
  updatedAt?: string;
}

export interface OrderItem {
  menuItemId: number;
  name: string;
  price: number;
  quantity: number;
}

export interface CartItem {
  id: number;
  name: string;
  price: number;
  quantity: number;
}
