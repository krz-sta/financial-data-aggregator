import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { MarketService } from './market.service';

describe('MarketService', () => {
  let service: MarketService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        MarketService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(MarketService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get assets', () => {
    const mockAssets = [{ id: 'bitcoin', symbol: 'BTC', name: 'Bitcoin', cmcRank: 1, currentPrice: 50000, priceChangePercentage24h: 5, type: 'crypto' }];
    service.getAssets().subscribe(data => {
      expect(data).toEqual(mockAssets);
    });
    const req = httpMock.expectOne('http://localhost:8080/api/assets');
    expect(req.request.method).toBe('GET');
    req.flush(mockAssets);
  });

  it('should get history', () => {
    const mockData = [{ timestamp: 123, price: 50000 }];
    service.getHistory('BTC').subscribe(data => {
      expect(data).toEqual(mockData);
    });
    const req = httpMock.expectOne('http://localhost:8080/api/rates/history/BTC');
    expect(req.request.method).toBe('GET');
    req.flush(mockData);
  });

  it('should get rates', () => {
    const mockRates = { BTC: 50000, ETH: 3000 };
    service.getRates().subscribe(data => {
      expect(data).toEqual(mockRates);
    });
    const req = httpMock.expectOne('http://localhost:8080/api/rates');
    expect(req.request.method).toBe('GET');
    req.flush(mockRates);
  });

  it('should convert usd to fiat', () => {
    service.selectedCurrency.set('EUR');
    service.fiatRates.set({ USD: 1, EUR: 0.9, PLN: 4.0, GBP: 0.8 });
    expect(service.convertUsdToFiat(100)).toBe(90);
  });

  it('should convert crypto to fiat', () => {
    service.selectedCurrency.set('EUR');
    service.fiatRates.set({ USD: 1, EUR: 0.9, PLN: 4.0, GBP: 0.8 });
    const rates = { BTC: 50000 };
    expect(service.convertCryptoToFiat(2, 'BTC', rates)).toBe(90000);
  });
});
