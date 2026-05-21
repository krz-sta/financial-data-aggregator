import { Component, inject, computed } from '@angular/core';
import { Router, RouterLink, NavigationEnd } from '@angular/router';
import { LucideArrowLeft, LucideSearch, LucideTrendingUp } from '@lucide/angular';
import { filter, map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { AuthService } from '../../services/auth';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [LucideTrendingUp, LucideSearch, RouterLink, LucideArrowLeft],
  templateUrl: './header.html',
  styleUrl: './header.scss',
})
export class Header {
  private router = inject(Router);
  private authService = inject(AuthService);

  private currentUrl = toSignal(
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map((event: any) => event.urlAfterRedirects)
    ),
    { initialValue: this.router.url }
  );

  // Sprawdzamy czy jesteśmy na podstronie /auth (wycinamy query params)
  isAuthPage = computed(() => this.currentUrl().split('?')[0] === '/auth');
  isHomePage = computed(() => this.currentUrl().split('?')[0] === '/home');
  
  // Pobieramy stan zalogowania bezpośrednio z sygnału w AuthService
  isLoggedIn = computed(() => this.authService.isLoggedIn());

  handleLogout() {
    this.authService.logout();
    this.router.navigate(['/']);
  }
}