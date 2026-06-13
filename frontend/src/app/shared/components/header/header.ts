import { Component, inject, computed, signal, OnInit } from '@angular/core';
import { Router, RouterLink, NavigationEnd } from '@angular/router';
import { LucideArrowLeft, LucideSearch, LucideTrendingUp } from '@lucide/angular';
import { filter, map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { AuthService } from '../../services/auth';
import { MarketService } from '../../services/market.service';
import { FormsModule } from '@angular/forms';
import { Asset } from '../../models/models';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [LucideTrendingUp, LucideSearch, RouterLink, LucideArrowLeft, FormsModule],
  templateUrl: './header.html',
  styleUrl: './header.scss',
})
export class Header implements OnInit {
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

  allAssets = signal<Asset[]>([]);

  ngOnInit() {
    this.marketService.getAssets().subscribe({
      next: (assets) => this.allAssets.set(assets),
      error: (err) => console.error('Failed to load assets for search', err)
    });
  }

  searchResults = computed(() => {
    const query = this.marketService.searchQuery().toLowerCase().trim();
    if (!query) return [];
    
    return this.allAssets()
      .filter((a: Asset) => a.type?.toLowerCase() !== 'fiat' && a.type?.toLowerCase() !== 'currency')
      .filter((a: Asset) => a.symbol.toLowerCase().includes(query) || a.name.toLowerCase().includes(query))
      .slice(0, 6);
  });

  handleResultClick(symbol: string) {
    this.marketService.searchQuery.set('');
    this.router.navigate(['/chart', symbol]);
  }

  handleLogout() {
    this.authService.logout();
    this.marketService.searchQuery.set('');
    this.router.navigate(['/']);
  }
}