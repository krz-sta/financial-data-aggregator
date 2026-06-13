import { HttpClient } from '@angular/common/http';
import { inject, Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { Asset, RatesResponse, HistoryPoint, FiatRates } from '../models/models';

@Injectable({
  providedIn: 'root'
})
export class MarketService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api';

  selectedCurrency = signal<string>('USD');
  searchQuery = signal<string>('');
  
  fiatRates = signal<FiatRates>({ USD: 1, EUR: 0.92, PLN: 4.00, GBP: 0.79 });

  getAssets(): Observable<Asset[]> {
    return this.http.get<Asset[]>(`${this.apiUrl}/assets`);
  }

  getRates(): Observable<RatesResponse> {
    return this.http.get<RatesResponse>(`${this.apiUrl}/rates`);
  }

  getHistory(symbol: string): Observable<HistoryPoint[]> {
    return this.http.get<HistoryPoint[]>(`${this.apiUrl}/rates/history/${symbol}`);
  }

  convertCryptoToFiat(amount: number, cryptoSymbol: string, currentRates: RatesResponse): number {
    const cryptoRateInUsd = currentRates[cryptoSymbol] || 0;
    const valueInUsd = amount * cryptoRateInUsd;
    return this.convertUsdToFiat(valueInUsd);
  }

  convertUsdToFiat(usdValue: number): number {
    const fiat = this.selectedCurrency();
    const rates = this.fiatRates() as any;
    return usdValue * (rates[fiat] || 1);
  }
}