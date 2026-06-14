import { TestBed, ComponentFixture } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { of, throwError } from 'rxjs';
import { Header } from './header';
import { AuthService } from '../../services/auth';
import { MarketService } from '../../services/market.service';
import { vi } from 'vitest';

describe('Header Component', () => {
  let component: Header;
  let fixture: ComponentFixture<Header>;
  let authService: AuthService;
  let marketService: MarketService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Header],
      providers: [
        provideRouter([]),
        provideHttpClient(),
        provideHttpClientTesting(),
        AuthService,
        MarketService
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(Header);
    component = fixture.componentInstance;
    authService = TestBed.inject(AuthService);
    marketService = TestBed.inject(MarketService);
    router = TestBed.inject(Router);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should fetch assets on init', () => {
    const mockAssets = [{ id: '1', symbol: 'BTC', name: 'Bitcoin', cmcRank: 1, type: 'crypto' }];
    vi.spyOn(marketService, 'getAssets').mockReturnValue(of(mockAssets));
    component.ngOnInit();
    expect(component.allAssets()).toEqual(mockAssets);
  });

  it('should handle fetch assets error on init', () => {
    vi.spyOn(marketService, 'getAssets').mockReturnValue(throwError(() => new Error('error')));
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(vi.fn());
    component.ngOnInit();
    expect(consoleSpy).toHaveBeenCalled();
  });

  it('should compute search results', () => {
    component.allAssets.set([
      { symbol: 'BTC', name: 'Bitcoin', type: 'crypto' } as any,
      { symbol: 'ETH', name: 'Ethereum', type: 'crypto' } as any,
      { symbol: 'USD', name: 'Dollar', type: 'fiat' } as any
    ]);
    marketService.searchQuery.set('bit');
    expect(component.searchResults().length).toBe(1);
    expect(component.searchResults()[0].symbol).toBe('BTC');
  });

  it('should handle result click', () => {
    const navigateSpy = vi.spyOn(router, 'navigate');
    marketService.searchQuery.set('bit');
    component.handleResultClick('BTC');
    expect(marketService.searchQuery()).toBe('');
    expect(navigateSpy).toHaveBeenCalledWith(['/chart', 'BTC']);
  });

  it('should handle logout', () => {
    const logoutSpy = vi.spyOn(authService, 'logout');
    const navigateSpy = vi.spyOn(router, 'navigate');
    component.handleLogout();
    expect(logoutSpy).toHaveBeenCalled();
    expect(marketService.searchQuery()).toBe('');
    expect(navigateSpy).toHaveBeenCalledWith(['/']);
  });
});
