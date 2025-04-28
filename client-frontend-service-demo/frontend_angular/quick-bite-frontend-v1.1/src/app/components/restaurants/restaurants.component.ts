import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../services/api.service';
import { Restaurant } from '../../models/restaurant.model';

@Component({
  selector: 'app-restaurants',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './restaurants.component.html',
  styleUrls: ['./restaurants.component.css'],
})
export class RestaurantsComponent implements OnInit {
  restaurants: Restaurant[] = [];
  loading = true;
  searchQuery = '';

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.fetchRestaurants();
  }

  fetchRestaurants(): void {
    this.loading = true;
    this.apiService.getRestaurants().subscribe({
      next: (data) => {
        this.restaurants = data;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error fetching restaurants:', error);
        // Use mock data for demo
        this.restaurants = [
          {
            id: 1,
            name: 'Tasty Bites',
            address: '123 Main St',
            cuisine: 'Italian',
            rating: 4.5,
          },
          {
            id: 2,
            name: 'Burger Palace',
            address: '456 Oak Ave',
            cuisine: 'American',
            rating: 4.2,
          },
          {
            id: 3,
            name: 'Sushi Heaven',
            address: '789 Maple Rd',
            cuisine: 'Japanese',
            rating: 4.7,
          },
        ];
        this.loading = false;
      },
    });
  }

  searchRestaurants(): void {
    // In a real app, we would search from the API
    // For this demo, we just log the search query
    console.log(`Searching for: ${this.searchQuery}`);
  }
}
