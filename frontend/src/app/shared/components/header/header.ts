import { Component, inject, computed } from '@angular/core';
import { Router, RouterLink, NavigationEnd } from '@angular/router';
import { LucideArrowLeft, LucideSearch, LucideTrendingUp } from '@lucide/angular';
import { filter, map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { AuthService } from '../../services/auth';
import { MarketService } from '../../services/market.service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [LucideTrendingUp, LucideSearch, RouterLink, LucideArrowLeft, FormsModule],
  templateUrl: './header.html',
  styleUrl: './header.scss',
})
export class Header {
  private router = inject(Router);
  public authService = inject(AuthService);
  public marketService = inject(MarketService);

  private currentUrl = toSignal(
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map((event: any) => event.urlAfterRedirects)
    ),
    { initialValue: this.router.url }
  );

  isAuthPage = computed(() => this.currentUrl().split('?')[0] === '/auth');
  isHomePage = computed(() => this.currentUrl().split('?')[0] === '/home');
  isLoggedIn = computed(() => this.authService.isLoggedIn());

  handleLogout() {
    this.authService.logout();
    this.marketService.searchQuery.set('');
    this.router.navigate(['/']);
  }
}