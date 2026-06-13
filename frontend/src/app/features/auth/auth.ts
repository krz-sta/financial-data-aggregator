import { Component, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { map } from 'rxjs';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../shared/services/auth'
import { email } from '@angular/forms/signals';

@Component({
  selector: 'app-auth',
  standalone: true,
  imports: [RouterLink, Header, Footer, FormsModule],
  templateUrl: './auth.html',
  styleUrl: './auth.scss',
})
export class Auth {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private authService = inject(AuthService)

  authMode = toSignal(
    this.route.queryParams.pipe(
      map(params => params['mode'] === 'signup' ? 'signup' : 'login')
    ),
    { initialValue: 'login' }
  );

  formData = signal({
    name: '',
    email: '',
    password: ''
  });

  isLoading = signal(false);
  errorMessage = signal<string | null>(null);
  private errorTimer: any;

  setMode(mode: 'login' | 'signup') {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { mode: mode},
      queryParamsHandling: 'merge'
    });
  }

  handleSubmit(event: Event) {
    event.preventDefault();
    this.isLoading.set(true);
    this.errorMessage.set(null);
    
    if (this.authMode() === 'login') {
      
      this.authService.login({
        email: this.formData().email,
        password: this.formData().password
      }).subscribe({
        next: (res) => {
          this.isLoading.set(false);
          this.router.navigate(['/home']);
        },
        error: (err) => {
          this.isLoading.set(false);
          this.showError(err?.error || 'ERROR.');
        }
      });

    } else {

      this.authService.register(this.formData()).subscribe({
        next: (res) => {
          this.isLoading.set(false);
          this.setMode('login');
          this.formData.set({ name: '', email: '', password: '' });
        },
        error: (err) => {
          this.isLoading.set(false);
          this.showError(err?.error || 'ERROR.');
        }
      });
    }
  }

  private showError(msg: string) {
    this.errorMessage.set(msg);
    if (this.errorTimer) clearTimeout(this.errorTimer);
    this.errorTimer = setTimeout(() => {
      this.errorMessage.set(null);
    }, 2000);
  }
}
