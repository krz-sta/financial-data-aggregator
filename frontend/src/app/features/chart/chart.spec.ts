import { TestBed, ComponentFixture } from '@angular/core/testing';
import { provideRouter, ActivatedRoute } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { of, throwError } from 'rxjs';
import { Chart } from './chart';
import { MarketService } from '../../shared/services/market.service';
import { vi } from 'vitest';

describe('Chart Component', () => {
  let component: Chart;
  let fixture: ComponentFixture<Chart>;
  let marketService: MarketService;

  beforeEach(async () => {
    (window as any).ResizeObserver = class {
      observe() {}
      unobserve() {}
      disconnect() {}
    };
    await TestBed.configureTestingModule({
      imports: [Chart],
      providers: [
        provideRouter([]),
        provideHttpClient(),
        provideHttpClientTesting(),
        MarketService,
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { paramMap: { get: () => 'BTC' } } }
        }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(Chart);
    component = fixture.componentInstance;
    marketService = TestBed.inject(MarketService);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should load chart data', () => {
    const mockHistory = [{ timestamp: 1000, price: 50000 }, { timestamp: 2000, price: 60000 }];
    const mockRates = { BTC: 60000 };

    vi.spyOn(marketService, 'getHistory').mockReturnValue(of(mockHistory));
    vi.spyOn(marketService, 'getRates').mockReturnValue(of(mockRates));

    component.ngOnInit();
    
    expect(component.symbol()).toBe('BTC');
    expect(component.currentPrice()).toBe(60000);
    expect(component.isLoading()).toBe(false);
    expect(component.highPrice()).toBe(60000);
    expect(component.lowPrice()).toBe(50000);
    expect(component.priceChange()).toBe(20);
  });

  it('should handle chart data error', () => {
    vi.spyOn(marketService, 'getHistory').mockReturnValue(throwError(() => new Error('err')));
    vi.spyOn(marketService, 'getRates').mockReturnValue(of({}));

    component.ngOnInit();

    expect(component.error()).toBe('Failed to load chart data. Please try again later.');
    expect(component.isLoading()).toBe(false);
  });

  it('should handle empty history', () => {
    vi.spyOn(marketService, 'getHistory').mockReturnValue(of([]));
    vi.spyOn(marketService, 'getRates').mockReturnValue(of({}));

    component.ngOnInit();

    expect(component.error()).toBe('No historical data available for this asset.');
    expect(component.isLoading()).toBe(false);
  });

  it('should set period and filter data', () => {
    const mockHistory = [
      { timestamp: Date.now() - 10 * 24 * 60 * 60 * 1000, price: 50000 },
      { timestamp: Date.now(), price: 60000 }
    ];
    vi.spyOn(marketService, 'getHistory').mockReturnValue(of(mockHistory));
    vi.spyOn(marketService, 'getRates').mockReturnValue(of({}));
    
    component.ngOnInit();

    component.setPeriod('7d');
    expect(component.selectedPeriod()).toBe('7d');
    expect(component.filteredChartData().length).toBe(1);

    component.setPeriod('30d');
    expect(component.filteredChartData().length).toBe(2);
  });

  it('should unbind listeners on destroy', () => {
    component.ngAfterViewInit();
    const removeSpy = vi.spyOn(component.chartCanvas.nativeElement, 'removeEventListener');
    component.ngOnDestroy();
    expect(removeSpy).toHaveBeenCalledTimes(2);
  });

  it('should render the template properly', () => {
    component.isLoading.set(false);
    component.error.set(null);
    fixture.detectChanges();

    component.isLoading.set(true);
    fixture.detectChanges();

    component.isLoading.set(false);
    component.error.set('Test error');
    fixture.detectChanges();
  });

  it('should format price label', () => {
    expect((component as any).formatPriceLabel(1005, '$')).toBe('$1,005');
    expect((component as any).formatPriceLabel(50, '$')).toBe('$50.00');
    expect((component as any).formatPriceLabel(0.5, '$')).toBe('$0.5000');
    expect((component as any).formatPriceLabel(0.005, '$')).toBe('$0.005000');
  });

  it('should handle canvas mouse events', () => {
    const mockHistory = [{ timestamp: 1000, price: 50000 }, { timestamp: 2000, price: 60000 }];
    (component as any).allChartData.set(mockHistory);
    
    const mockEvent = { clientX: 100, clientY: 100 } as MouseEvent;
    
    component.chartCanvas.nativeElement.getBoundingClientRect = () => ({
      left: 0, top: 0, width: 800, height: 600, right: 800, bottom: 600, x: 0, y: 0, toJSON: () => {}
    });

    component.onCanvasMouseMove(mockEvent);
    expect(component.showTooltip()).toBe(true);

    component.onCanvasMouseLeave();
    expect(component.showTooltip()).toBe(false);
  });

  it('should handle internal methods for coverage', () => {
    const mockHistory = [
      { timestamp: 1000, price: 50000 },
      { timestamp: 2000, price: 60000 }
    ];
    (component as any).allChartData.set(mockHistory);
    (component as any).currentPrice.set(55000);
    (component as any).computeStats();
    expect(component.priceChange()).toBeDefined();

    const ctx = {
      fillRect: () => {},
      scale: () => {},
      beginPath: () => {},
      moveTo: () => {},
      lineTo: () => {},
      stroke: () => {},
      fillText: () => {},
      createLinearGradient: () => ({ addColorStop: () => {} }),
      closePath: () => {},
      fill: () => {},
      arc: () => {},
      setLineDash: () => {},
    };
    component.chartCanvas.nativeElement.getContext = () => ctx as any;
    (component as any).drawChart(1);
    (component as any).drawChart(-1);
  });
});
