import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private http = inject(HttpClient);
  private apiUrl = '/api/protected';

  getProfile(): Observable<any> {
    return this.http.post(`${this.apiUrl}/profile`, {}); 
  }
}