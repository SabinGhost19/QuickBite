export interface Restaurant {
  id: number;
  name: string;
  address: string;
  cuisine: string;
  rating: number;
  menuItems?: MenuItem[];
}

export interface MenuItem {
  id: number;
  name: string;
  description: string;
  price: number;
  category: string;
}
