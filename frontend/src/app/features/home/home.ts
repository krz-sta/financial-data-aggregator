import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';
import { UserService } from '../../shared/services/user.service';
import { PortfolioService } from '../../shared/services/portfolio.service';
import { MarketService } from '../../shared/services/market.service';
import { UpperCasePipe, CurrencyPipe, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { forkJoin } from 'rxjs';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [Header, Footer, UpperCasePipe, CurrencyPipe, DatePipe, FormsModule],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home implements OnInit {
  private userService = inject(UserService);
  private portfolioService = inject(PortfolioService);
  public marketService = inject(MarketService);
  
  profileData = signal<any>(null);
  allAssets = signal<any[]>([]); 
  rates = signal<any>({});
  
  isLoading = signal(true);
  isSubmitting = signal(false);
  error = signal<string | null>(null);

  addForm = signal({ symbol: '', amount: null as number | null });

  cryptoAssets = computed(() => {
    return this.allAssets().filter(asset => 
      asset.type?.toLowerCase() !== 'fiat' && asset.type?.toLowerCase() !== 'currency'
    );
  });

  groupedPortfolio = computed(() => {
    const portfolio = this.profileData()?.portfolio || [];
    const map = new Map<string, any>();
    
    for (const item of portfolio) {
      if (map.has(item.symbol)) {
        const existing = map.get(item.symbol);
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
    const portfolio = this.groupedPortfolio();
    const query = this.marketService.searchQuery().toLowerCase().trim();
    
    if (!query) return portfolio;
    return portfolio.filter((item: any) => item.symbol.toLowerCase().includes(query));
  });

  totalBalance = computed(() => {
    const portfolio = this.groupedPortfolio();
    const currentRates = this.rates();
    const fiat = this.marketService.selectedCurrency();
    const fiatRates = this.marketService.fiatRates();

    if (!portfolio.length || !currentRates) return 0;

    let totalUsd = portfolio.reduce((total: number, item: any) => {
      const cryptoRateInUsd = currentRates[item.symbol] || 0;
      return total + (item.amount * cryptoRateInUsd);
    }, 0);

    return totalUsd * (fiatRates[fiat] || 1);
  });

  getItemValue(amount: number, symbol: string): number {
    const currentRates = this.rates();
    const fiat = this.marketService.selectedCurrency();
    const fiatRates = this.marketService.fiatRates();
    
    const cryptoUsdValue = amount * (currentRates[symbol] || 0);
    return cryptoUsdValue * (fiatRates[fiat] || 1);
  }

  ngOnInit() { 
    this.fetchData(); 
  }

  fetchData() {
    this.isLoading.set(true);
    
    this.marketService.getAssets().subscribe({
      next: (assets) => this.allAssets.set(assets),
      error: (err) => console.error('Failed to load assets', err)
    });
    
    this.marketService.getRates().subscribe({
      next: (ratesData) => this.rates.set(ratesData),
      error: (err) => console.error('Failed to load rates', err)
    });
    
    this.loadProfile();
  }

  loadProfile() {
    this.userService.getProfile().subscribe({
      next: (data) => {
        this.profileData.set(data);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.error.set('Failed to load profile data.');
        this.isLoading.set(false);
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
    alert(`todo`);
  }
}