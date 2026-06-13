import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';
import { UserService } from '../../shared/services/user.service';
import { PortfolioService } from '../../shared/services/portfolio.service';
import { MarketService } from '../../shared/services/market.service';
import { UpperCasePipe, CurrencyPipe } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [Header, Footer, UpperCasePipe, CurrencyPipe, FormsModule],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home implements OnInit {
  private userService = inject(UserService);
  private portfolioService = inject(PortfolioService);
  private marketService = inject(MarketService);
  
  profileData = signal<any>(null);
  availableAssets = signal<any[]>([]);
  rates = signal<any>({});
  
  isLoading = signal(true);
  isSubmitting = signal(false);
  error = signal<string | null>(null);

  addForm = signal({
    symbol: '',
    amount: null as number | null
  });

  totalBalance = computed(() => {
    const portfolio = this.profileData()?.portfolio;
    const currentRates = this.rates();

    if (!portfolio || !currentRates) return 0;

    return portfolio.reduce((total: number, item: any) => {
      const rate = currentRates[item.symbol] || 0;
      return total + (item.amount * rate);
    }, 0);
  });

  ngOnInit() {
    this.fetchData();
  }

  fetchData() {
    this.isLoading.set(true);
    
    // Pobranie listy aktywów
    this.marketService.getAssets().subscribe({
      next: (assets) => this.availableAssets.set(assets),
      error: (err) => console.error('Failed to load assets', err)
    });

    // Pobranie kursów rynkowych
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

  handleDelete(id: string) {
    this.portfolioService.deleteItem(id).subscribe({
      next: () => this.loadProfile(),
      error: () => this.error.set('Failed to delete item.')
    });
  }
}