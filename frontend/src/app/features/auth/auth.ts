import { Component, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { map } from 'rxjs';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';

@Component({
  selector: 'app-auth',
  standalone: true,
  imports: [RouterLink, Header, Footer],
  templateUrl: './auth.html',
  styleUrl: './auth.scss',
})
export class Auth {
  private route = inject(ActivatedRoute);
  private router = inject(Router);

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

  setMode(mode: 'login' | 'signup') {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { mode: mode},
      queryParamsHandling: 'merge'
    });
  }

  handleSubmit(event: Event) {
    event.preventDefault();
    
    if (this.authMode() === 'login') {
      console.log('Log in:', this.formData().email, this.formData().password);
      // todo: auth service
    } else {
      console.log('Sign up:', this.formData());
    }
  }
}
