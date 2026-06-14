import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';
import { UserService } from '../../shared/services/user.service';
import { PortfolioService } from '../../shared/services/portfolio.service';
import { MarketService } from '../../shared/services/market.service';
import { UpperCasePipe, CurrencyPipe, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { forkJoin } from 'rxjs';
import { CustomSelectComponent } from '../../shared/components/select/select';
import { Asset, UserProfile, RatesResponse } from '../../shared/models/models';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [Header, Footer, UpperCasePipe, CurrencyPipe, DatePipe, FormsModule, CustomSelectComponent],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home implements OnInit {
  private userService = inject(UserService);
  private portfolioService = inject(PortfolioService);
  private router = inject(Router);
  public marketService = inject(MarketService);
  
  profileData = signal<UserProfile | null>(null);
  allAssets = signal<Asset[]>([]); 
  rates = signal<RatesResponse>({});
  
  isLoading = signal(true);
  isSubmitting = signal(false);
  error = signal<string | null>(null);

  addForm = signal({ symbol: '', amount: null as number | null });

  currencies = ['USD', 'EUR', 'PLN', 'GBP'];

  cryptoAssets = computed(() => {
    return this.allAssets().filter(asset => 
      asset.type?.toLowerCase() !== 'fiat' && asset.type?.toLowerCase() !== 'currency'
    );
  });

  cryptoAssetsOptions = computed(() => {
    return this.cryptoAssets().map(asset => ({
      label: asset.symbol,
      value: asset.symbol
    }));
  });

  groupedPortfolio = computed(() => {
    const portfolio = this.profileData()?.portfolio || [];
    const map = new Map<string, { symbol: string, amount: number, ids: string[] }>();
    
    for (const item of portfolio) {
      if (map.has(item.symbol)) {
        const existing = map.get(item.symbol)!;
        existing.amount += item.amount;
        existing.ids.push(item.id);
      } else {
        map.set(item.symbol, { 
          symbol: item.symbol, 
          amount: item.amount, 
          ids: [item.id] 
        });
      }
    }
    return Array.from(map.values());
  });

  filteredPortfolio = computed(() => {
    return this.groupedPortfolio();
  });

  totalBalance = computed(() => {
    const portfolio = this.groupedPortfolio();
    const currentRates = this.rates();

    if (!portfolio.length || !currentRates) return 0;

    return portfolio.reduce((total: number, item: { symbol: string, amount: number, ids: string[] }) => {
      return total + this.marketService.convertCryptoToFiat(item.amount, item.symbol, currentRates);
    }, 0);
  });

  getItemValue(amount: number, symbol: string): number {
    return this.marketService.convertCryptoToFiat(amount, symbol, this.rates());
  }

  ngOnInit() { 
    this.fetchData(); 
  }

  fetchData() {
    this.isLoading.set(true);
    this.error.set(null);
    
    forkJoin({
      assets: this.marketService.getAssets(),
      rates: this.marketService.getRates(),
      profile: this.userService.getProfile()
    }).subscribe({
      next: ({ assets, rates, profile }) => {
        this.allAssets.set(assets);
        this.rates.set(rates);
        this.profileData.set(profile as UserProfile);
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load initial data', err);
        this.error.set('Failed to load dashboard data. Please try again.');
        this.isLoading.set(false);
      }
    });
  }

  loadProfile() {
    this.userService.getProfile().subscribe({
      next: (data) => {
        this.profileData.set(data as UserProfile);
      },
      error: () => {
        this.error.set('Failed to refresh profile data.');
      }
    });
  }

  handleAddAsset(event: Event) {
    event.preventDefault();
    const { symbol, amount } = this.addForm();

    if (!symbol || !amount || amount <= 0) {
      this.error.set('Please select an asset and enter a valid amount.');
      return;
    }

    this.isSubmitting.set(true);
    this.error.set(null);

    this.portfolioService.addItem(symbol, amount).subscribe({
      next: () => {
        this.isSubmitting.set(false);
        this.addForm.set({ symbol: '', amount: null });
        this.loadProfile();
      },
      error: () => {
        this.error.set('Failed to add item to portfolio.');
        this.isSubmitting.set(false);
      }
    });
  }

  handleDeleteGroup(ids: string[], event: Event) {
    event.stopPropagation(); 
    const deleteRequests = ids.map(id => this.portfolioService.deleteItem(id));
    
    forkJoin(deleteRequests).subscribe({
      next: () => this.loadProfile(),
      error: () => this.error.set('Failed to delete some items.')
    });
  }

  handleRowClick(symbol: string) {
    this.router.navigate(['/chart', symbol]);
  }
}