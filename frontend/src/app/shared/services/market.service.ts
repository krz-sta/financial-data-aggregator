import { HttpClient } from '@angular/common/http';
import { inject, Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class MarketService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api';

  selectedCurrency = signal<string>('USD');
  searchQuery = signal<string>('');
  
  fiatRates = signal<any>({ USD: 1, EUR: 0.92, PLN: 4.00, GBP: 0.79 });

  constructor() {
    this.fetchFiatRates();
  }

  getAssets(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/assets`);
  }

  getRates(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/rates`);
  }

  private fetchFiatRates() {
    this.http.get('https://api.exchangerate-api.com/v4/latest/USD').subscribe({
      next: (res: any) => {
        if (res && res.rates) {
          this.fiatRates.set(res.rates);
        }
      },
      error: () => console.warn('Failed to load fiat rates, using fallback values.')
    });
  }
}