import { TestBed, ComponentFixture } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { of, throwError } from 'rxjs';
import { Home } from './home';
import { MarketService } from '../../shared/services/market.service';
import { UserService } from '../../shared/services/user.service';
import { PortfolioService } from '../../shared/services/portfolio.service';
import { vi } from 'vitest';

describe('Home Component', () => {
  let component: Home;
  let fixture: ComponentFixture<Home>;
  let marketService: MarketService;
  let userService: UserService;
  let portfolioService: PortfolioService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Home],
      providers: [
        provideRouter([]),
        provideHttpClient(),
        provideHttpClientTesting(),
        MarketService,
        UserService,
        PortfolioService
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(Home);
    component = fixture.componentInstance;
    marketService = TestBed.inject(MarketService);
    userService = TestBed.inject(UserService);
    portfolioService = TestBed.inject(PortfolioService);
    router = TestBed.inject(Router);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should fetch data on init', () => {
    const mockAssets = [{ id: '1', symbol: 'BTC', name: 'Bitcoin', cmcRank: 1, type: 'crypto' }];
    const mockRates = { BTC: 50000 };
    const mockProfile = { email: 'test@test.com', portfolio: [{ id: 'p1', symbol: 'BTC', amount: 2 }] };

    vi.spyOn(marketService, 'getAssets').mockReturnValue(of(mockAssets));
    vi.spyOn(marketService, 'getRates').mockReturnValue(of(mockRates));
    vi.spyOn(userService, 'getProfile').mockReturnValue(of(mockProfile));

    component.ngOnInit();

    expect(component.allAssets()).toEqual(mockAssets);
    expect(component.rates()).toEqual(mockRates);
    expect(component.profileData()).toEqual(mockProfile);
    expect(component.isLoading()).toBe(false);
  });

  it('should handle fetch data error', () => {
    vi.spyOn(marketService, 'getAssets').mockReturnValue(throwError(() => new Error('error')));
    vi.spyOn(marketService, 'getRates').mockReturnValue(of({}));
    vi.spyOn(userService, 'getProfile').mockReturnValue(of({}));

    component.fetchData();

    expect(component.error()).toBe('Failed to load dashboard data. Please try again.');
    expect(component.isLoading()).toBe(false);
  });

  it('should load profile', () => {
    const mockProfile = { email: 'test@test.com', portfolio: [] };
    vi.spyOn(userService, 'getProfile').mockReturnValue(of(mockProfile));
    component.loadProfile();
    expect(component.profileData()).toEqual(mockProfile);
  });

  it('should handle add asset success', () => {
    component.addForm.set({ symbol: 'BTC', amount: 1 });
    vi.spyOn(portfolioService, 'addItem').mockReturnValue(of({ success: true }));
    const loadSpy = vi.spyOn(component, 'loadProfile').mockImplementation(() => {});

    component.handleAddAsset(new Event('submit'));

    expect(portfolioService.addItem).toHaveBeenCalledWith('BTC', 1);
    expect(component.addForm()).toEqual({ symbol: '', amount: null });
    expect(loadSpy).toHaveBeenCalled();
  });

  it('should handle add asset validation error', () => {
    component.addForm.set({ symbol: '', amount: 0 });
    component.handleAddAsset(new Event('submit'));
    expect(component.error()).toBe('Please select an asset and enter a valid amount.');
  });

  it('should handle add asset api error', () => {
    component.addForm.set({ symbol: 'BTC', amount: 1 });
    vi.spyOn(portfolioService, 'addItem').mockReturnValue(throwError(() => new Error('err')));
    component.handleAddAsset(new Event('submit'));
    expect(component.error()).toBe('Failed to add item to portfolio.');
  });

  it('should delete group', () => {
    vi.spyOn(portfolioService, 'deleteItem').mockReturnValue(of({ success: true }));
    const loadSpy = vi.spyOn(component, 'loadProfile').mockImplementation(() => {});
    
    component.handleDeleteGroup(['id1', 'id2'], new Event('click'));
    
    expect(portfolioService.deleteItem).toHaveBeenCalledTimes(2);
  });

  it('should navigate to chart on row click', () => {
    const navigateSpy = vi.spyOn(router, 'navigate');
    component.handleRowClick('BTC');
    expect(navigateSpy).toHaveBeenCalledWith(['/chart', 'BTC']);
  });

  it('should compute cryptoAssetsOptions correctly', () => {
    component.allAssets.set([
      { symbol: 'BTC', type: 'crypto' } as any,
      { symbol: 'USD', type: 'fiat' } as any
    ]);
    expect(component.cryptoAssetsOptions()).toEqual([{ label: 'BTC', value: 'BTC' }]);
  });

  it('should compute totalBalance correctly', () => {
    component.profileData.set({ portfolio: [{ symbol: 'BTC', amount: 2 }] } as any);
    component.rates.set({ BTC: 50000 });
    vi.spyOn(marketService, 'convertCryptoToFiat').mockReturnValue(100000);
    expect(component.totalBalance()).toBe(100000);
  });

  it('should get item value', () => {
    component.rates.set({ BTC: 50000 });
    vi.spyOn(marketService, 'convertCryptoToFiat').mockReturnValue(50000);
    expect(component.getItemValue(1, 'BTC')).toBe(50000);
  });

  it('should render the full template', () => {
    component.isLoading.set(false);
    component.error.set(null);
    component.profileData.set({ email: 'a@a.com', displayName: 'A', portfolio: [{ symbol: 'BTC', amount: 1, id: '123' }] } as any);
    component.rates.set({ BTC: 50000 });
    component.allAssets.set([{ symbol: 'BTC', type: 'crypto' } as any]);
    fixture.detectChanges();

    component.isLoading.set(true);
    fixture.detectChanges();

    component.isLoading.set(false);
    component.error.set('Test error');
    fixture.detectChanges();
  });
});
