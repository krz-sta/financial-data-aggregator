import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class MarketService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api';

  getAssets(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/assets`);
  }

  getRates(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/rates`);
  }
}