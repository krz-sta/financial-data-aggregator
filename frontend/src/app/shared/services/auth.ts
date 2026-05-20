import { HttpClient } from '@angular/common/http';
import { inject, Injectable, signal } from '@angular/core';
import { Observable, tap } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api';

  isLoggedIn = signal<boolean>(!!localStorage.getItem('auth_token'));

  login(credentials: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/auth/login`, credentials).pipe(
      tap((res: any) => {
        if (res.token) {
          localStorage.setItem('auth_token', res.token);
          this.isLoggedIn.set(true);
        }
      })
    );
  }

  register(formData: any): Observable<any> {
    const payload = {
      email: formData.email,
      displayName: formData.name, 
      password: formData.password
    };
    return this.http.post(`${this.apiUrl}/auth/register`, payload);
  }

  logout() {
    localStorage.removeItem('auth_token');
  }

  getToken() {
    return localStorage.getItem('auth_token');
    this.isLoggedIn.set(false);
  }
}
