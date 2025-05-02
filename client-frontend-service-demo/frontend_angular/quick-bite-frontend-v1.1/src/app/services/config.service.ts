import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {
  private config = {
    apiBaseUrl: environment.apiBaseUrl,
    restaurantsServiceUrl: environment.restaurantsServiceUrl,
    ordersServiceUrl: environment.ordersServiceUrl,
    paymentsServiceUrl: environment.paymentsServiceUrl,
    deliveriesServiceUrl: environment.deliveriesServiceUrl,
    notificationsServiceUrl: environment.notificationsServiceUrl,
    frontendPort: environment.frontendPort,
  };

  constructor() {
    // In production, override with window variables if they exist
    if (typeof window !== 'undefined') {
      const windowEnv = (window as any).__env || {};

      Object.keys(this.config).forEach((key) => {
        if (windowEnv[key]) {
          (this.config as any)[key] = windowEnv[key];
        }
      });
    }
  }

  get(key: string): any {
    return (this.config as any)[key];
  }

  getApiBaseUrl(): string {
    return this.config.apiBaseUrl;
  }

  getRestaurantsServiceUrl(): string {
    return this.config.restaurantsServiceUrl;
  }

  getOrdersServiceUrl(): string {
    return this.config.ordersServiceUrl;
  }

  getPaymentsServiceUrl(): string {
    return this.config.paymentsServiceUrl;
  }

  getDeliveriesServiceUrl(): string {
    return this.config.deliveriesServiceUrl;
  }

  getNotificationsServiceUrl(): string {
    return this.config.notificationsServiceUrl;
  }

  getFrontendPort(): string | number {
    return this.config.frontendPort;
  }
}
