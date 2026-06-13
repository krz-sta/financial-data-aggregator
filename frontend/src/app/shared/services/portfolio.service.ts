import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class PortfolioService {
  private http = inject(HttpClient);
  private apiUrl = '/api/protected/portfolio';

  addItem(symbol: string, amount: number): Observable<any> {
    return this.http.post(this.apiUrl, { symbol, amount });
  }

  deleteItem(id: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}