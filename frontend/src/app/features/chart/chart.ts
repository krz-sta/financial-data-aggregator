import { Component, inject, OnInit, OnDestroy, signal, ElementRef, ViewChild, effect, computed, NgZone, AfterViewInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';
import { MarketService } from '../../shared/services/market.service';
import { CurrencyPipe, DecimalPipe } from '@angular/common';
import { CustomSelectComponent } from '../../shared/components/select/select';
import { forkJoin } from 'rxjs';
import { HistoryPoint, RatesResponse } from '../../shared/models/models';

@Component({
  selector: 'app-chart',
  standalone: true,
  imports: [Header, Footer, CurrencyPipe, DecimalPipe, RouterLink, CustomSelectComponent],
  templateUrl: './chart.html',
  styleUrl: './chart.scss',
})
export class Chart implements OnInit, OnDestroy, AfterViewInit {
  private route = inject(ActivatedRoute);
  public marketService = inject(MarketService);
  private ngZone = inject(NgZone);

  symbol = signal<string>('');
  isLoading = signal(true);
  error = signal<string | null>(null);
  currentPrice = signal<number | null>(null);
  priceChange = signal<number | null>(null);
  highPrice = signal<number | null>(null);
  lowPrice = signal<number | null>(null);

  selectedPeriod = signal<'7d' | '30d'>('30d');
  currencies = ['USD', 'EUR', 'PLN', 'GBP'];


  hoverPrice = signal<number | null>(null);
  hoverDate = signal<string | null>(null);
  hoverX = signal<number>(0);
  hoverY = signal<number>(0);
  showTooltip = signal(false);

  @ViewChild('chartCanvas', { static: true }) chartCanvas!: ElementRef<HTMLCanvasElement>;
  @ViewChild('chartWrapper', { static: true }) chartWrapper!: ElementRef<HTMLDivElement>;

  private allChartData = signal<HistoryPoint[]>([]);
  

  filteredChartData = computed(() => {
    const data = this.allChartData();
    if (this.selectedPeriod() === '7d') {
      const cutoff = (Date.now() / 1000) - (7 * 24 * 60 * 60);
      return data.filter(p => (p.timestamp / 1000) >= cutoff);
    }
    return data;
  });

  private resizeObserver: ResizeObserver | null = null;
  private dpr = 1;
  private dataLoaded = false;
  private mouseMoveListener: ((e: MouseEvent) => void) | null = null;
  private mouseLeaveListener: (() => void) | null = null;
  private paramSub: any;


  private readonly PADDING_TOP = 24;
  private readonly PADDING_BOTTOM = 40;
  private readonly PADDING_LEFT = 16;
  private readonly PADDING_RIGHT = 80;

  constructor() {

    effect(() => {
      this.marketService.selectedCurrency();
      this.marketService.fiatRates();
      this.selectedPeriod();
      if (this.dataLoaded) {
        this.computeStats();
        this.drawChart();
      }
    });
  }

  ngOnInit() {
    this.dpr = window.devicePixelRatio || 1;
    this.paramSub = this.route.paramMap.subscribe(params => {
      const sym = params.get('symbol') || '';
      this.symbol.set(sym.toUpperCase());
      this.loadChartData();
    });
  }

  ngAfterViewInit() {
    const canvas = this.chartCanvas.nativeElement;
    
    this.ngZone.runOutsideAngular(() => {
      this.mouseMoveListener = (e: MouseEvent) => this.onCanvasMouseMove(e);
      this.mouseLeaveListener = () => this.onCanvasMouseLeave();
      
      canvas.addEventListener('mousemove', this.mouseMoveListener);
      canvas.addEventListener('mouseleave', this.mouseLeaveListener);
    });
  }

  ngOnDestroy() {
    if (this.paramSub) {
      this.paramSub.unsubscribe();
    }
    if (this.resizeObserver) {
      this.resizeObserver.disconnect();
    }
    const canvas = this.chartCanvas?.nativeElement;
    if (canvas && this.mouseMoveListener && this.mouseLeaveListener) {
      canvas.removeEventListener('mousemove', this.mouseMoveListener);
      canvas.removeEventListener('mouseleave', this.mouseLeaveListener);
    }
  }

  setPeriod(period: '7d' | '30d') {
    this.selectedPeriod.set(period);
  }


  private getFiatMultiplier(): number {
    const fiat = this.marketService.selectedCurrency();
    const fiatRates = this.marketService.fiatRates();
    return fiatRates[fiat] || 1;
  }


  private getCurrencySymbol(): string {
    const symbols: Record<string, string> = {
      USD: '$', EUR: '€', PLN: 'zł', GBP: '£'
    };
    return symbols[this.marketService.selectedCurrency()] || '$';
  }

  onCanvasMouseMove(event: MouseEvent) {
    const data = this.filteredChartData();
    if (!data.length) return;

    const canvas = this.chartCanvas.nativeElement;
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const width = rect.width;
    const height = rect.height;

    const plotLeft = this.PADDING_LEFT;
    const plotRight = width - this.PADDING_RIGHT;
    const plotTop = this.PADDING_TOP;
    const plotBottom = height - this.PADDING_BOTTOM;

    if (x < plotLeft || x > plotRight || y < plotTop || y > plotBottom) {
      this.showTooltip.set(false);
      this.drawChart();
      return;
    }

    const ratio = (x - plotLeft) / (plotRight - plotLeft);
    const idx = Math.round(ratio * (data.length - 1));
    const clamped = Math.max(0, Math.min(data.length - 1, idx));
    const point = data[clamped];

    const date = new Date(point.timestamp);
    

    this.ngZone.run(() => {
      this.hoverPrice.set(point.price);
      this.hoverDate.set(date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }));
      
      const mult = this.getFiatMultiplier();
      const prices = data.map(p => p.price * mult);
      const minPrice = Math.min(...prices);
      const maxPrice = Math.max(...prices);
      const pricePadding = (maxPrice - minPrice) * 0.05;
      const adjMin = minPrice - pricePadding;
      const adjRange = (maxPrice + pricePadding) - adjMin;

      const pointX = plotLeft + (clamped / (data.length - 1)) * (plotRight - plotLeft);
      const pointY = plotBottom - ((point.price * mult - adjMin) / adjRange) * (plotBottom - plotTop);

      this.hoverX.set(pointX);
      this.hoverY.set(pointY);
      this.showTooltip.set(true);
    });

    this.drawChart(clamped);
  }

  onCanvasMouseLeave() {
    this.ngZone.run(() => {
      this.showTooltip.set(false);
    });
    this.drawChart();
  }

  private loadChartData() {
    this.isLoading.set(true);
    this.error.set(null);

    forkJoin({
      history: this.marketService.getHistory(this.symbol()),
      rates: this.marketService.getRates()
    }).subscribe({
      next: ({ history, rates }) => {
        if (!history || history.length === 0) {
          this.error.set('No historical data available for this asset.');
          this.isLoading.set(false);
          return;
        }

        const livePrice = rates[this.symbol()];
        if (livePrice != null) {
          this.currentPrice.set(livePrice);
        }

        this.allChartData.set(history.sort((a, b) => a.timestamp - b.timestamp));

        this.dataLoaded = true;
        this.computeStats();
        this.isLoading.set(false);

        setTimeout(() => {
          this.setupResize();
          this.drawChart();
        }, 0);
      },
      error: () => {
        this.error.set('Failed to load chart data. Please try again later.');
        this.isLoading.set(false);
      }
    });
  }

  private computeStats() {
    const data = this.filteredChartData();
    if (data.length === 0) return;

    const earliest = data[0].price;
    const latest = this.currentPrice() ?? data[data.length - 1].price;

    if (this.currentPrice() == null) {
      this.currentPrice.set(data[data.length - 1].price);
    }

    this.priceChange.set(((latest - earliest) / earliest) * 100);
    this.highPrice.set(Math.max(...data.map(p => p.price)));
    this.lowPrice.set(Math.min(...data.map(p => p.price)));
  }

  private setupResize() {
    const wrapper = this.chartWrapper?.nativeElement;
    if (!wrapper) return;

    this.resizeObserver = new ResizeObserver(() => {
      this.drawChart();
    });
    this.resizeObserver.observe(wrapper);
  }

  private drawChart(hoveredIdx: number = -1) {
    const canvas = this.chartCanvas?.nativeElement;
    const wrapper = this.chartWrapper?.nativeElement;
    const data = this.filteredChartData();
    if (!canvas || !wrapper || data.length === 0) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const mult = this.getFiatMultiplier();
    const currSymbol = this.getCurrencySymbol();


    const displayWidth = wrapper.clientWidth;
    const displayHeight = 500;
    canvas.style.width = displayWidth + 'px';
    canvas.style.height = displayHeight + 'px';
    canvas.width = displayWidth * this.dpr;
    canvas.height = displayHeight * this.dpr;
    ctx.scale(this.dpr, this.dpr);

    const w = displayWidth;
    const h = displayHeight;
    const plotLeft = this.PADDING_LEFT;
    const plotRight = w - this.PADDING_RIGHT;
    const plotTop = this.PADDING_TOP;
    const plotBottom = h - this.PADDING_BOTTOM;
    const plotWidth = plotRight - plotLeft;
    const plotHeight = plotBottom - plotTop;


    ctx.fillStyle = '#000000';
    ctx.fillRect(0, 0, w, h);


    const prices = data.map(p => p.price * mult);
    const minPrice = Math.min(...prices);
    const maxPrice = Math.max(...prices);
    const priceRange = maxPrice - minPrice || 1;
    const pricePadding = priceRange * 0.05;
    const adjMin = minPrice - pricePadding;
    const adjMax = maxPrice + pricePadding;
    const adjRange = adjMax - adjMin;

    const toX = (i: number) => plotLeft + (i / (data.length - 1)) * plotWidth;
    const toY = (convertedPrice: number) => plotBottom - ((convertedPrice - adjMin) / adjRange) * plotHeight;


    const gridLines = 5;
    ctx.strokeStyle = '#18181b';
    ctx.lineWidth = 1;
    ctx.fillStyle = '#52525b';
    ctx.font = '11px monospace';
    ctx.textAlign = 'right';
    ctx.textBaseline = 'middle';

    for (let i = 0; i <= gridLines; i++) {
      const price = adjMin + (i / gridLines) * adjRange;
      const y = toY(price);

      ctx.beginPath();
      ctx.moveTo(plotLeft, y);
      ctx.lineTo(plotRight, y);
      ctx.stroke();

      ctx.fillText(this.formatPriceLabel(price, currSymbol), w - 8, y);
    }


    const labelCount = Math.min(6, data.length);
    ctx.textAlign = 'center';
    ctx.textBaseline = 'top';
    ctx.fillStyle = '#52525b';

    for (let i = 0; i < labelCount; i++) {
      const dataIdx = Math.round((i / (labelCount - 1)) * (data.length - 1));
      const x = toX(dataIdx);
      const date = new Date(data[dataIdx].timestamp);
      const label = date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

      ctx.strokeStyle = '#18181b';
      ctx.beginPath();
      ctx.moveTo(x, plotTop);
      ctx.lineTo(x, plotBottom);
      ctx.stroke();

      ctx.fillStyle = '#52525b';
      ctx.fillText(label, x, plotBottom + 10);
    }


    const periodChange = prices[prices.length - 1] - prices[0];
    const isPositive = periodChange >= 0;
    const lineColor = isPositive ? '#22c55e' : '#ef4444';
    const gradientTop = isPositive ? 'rgba(34, 197, 94, 0.18)' : 'rgba(239, 68, 68, 0.18)';
    const gradientBottom = isPositive ? 'rgba(34, 197, 94, 0.0)' : 'rgba(239, 68, 68, 0.0)';


    const gradient = ctx.createLinearGradient(0, plotTop, 0, plotBottom);
    gradient.addColorStop(0, gradientTop);
    gradient.addColorStop(1, gradientBottom);

    ctx.beginPath();
    ctx.moveTo(toX(0), toY(prices[0]));
    for (let i = 1; i < data.length; i++) {
      ctx.lineTo(toX(i), toY(prices[i]));
    }
    ctx.lineTo(toX(data.length - 1), plotBottom);
    ctx.lineTo(toX(0), plotBottom);
    ctx.closePath();
    ctx.fillStyle = gradient;
    ctx.fill();


    ctx.beginPath();
    ctx.moveTo(toX(0), toY(prices[0]));
    for (let i = 1; i < data.length; i++) {
      ctx.lineTo(toX(i), toY(prices[i]));
    }
    ctx.strokeStyle = lineColor;
    ctx.lineWidth = 2;
    ctx.lineJoin = 'round';
    ctx.lineCap = 'round';
    ctx.stroke();


    if (hoveredIdx >= 0 && hoveredIdx < data.length) {
      const hx = toX(hoveredIdx);
      const hy = toY(prices[hoveredIdx]);

      ctx.strokeStyle = '#3f3f46';
      ctx.lineWidth = 1;
      ctx.setLineDash([4, 4]);

      ctx.beginPath();
      ctx.moveTo(hx, plotTop);
      ctx.lineTo(hx, plotBottom);
      ctx.stroke();

      ctx.beginPath();
      ctx.moveTo(plotLeft, hy);
      ctx.lineTo(plotRight, hy);
      ctx.stroke();
      ctx.setLineDash([]);


      ctx.beginPath();
      ctx.arc(hx, hy, 5, 0, Math.PI * 2);
      ctx.fillStyle = lineColor;
      ctx.fill();
      ctx.strokeStyle = '#000000';
      ctx.lineWidth = 2;
      ctx.stroke();


      const labelH = 20;
      ctx.fillStyle = '#27272a';
      ctx.fillRect(plotRight + 2, hy - labelH / 2, this.PADDING_RIGHT - 10, labelH);
      ctx.fillStyle = '#e4e4e7';
      ctx.font = '11px monospace';
      ctx.textAlign = 'right';
      ctx.textBaseline = 'middle';
      ctx.fillText(this.formatPriceLabel(prices[hoveredIdx], currSymbol), w - 8, hy);
    }
  }

  private formatPriceLabel(price: number, symbol: string): string {
    if (price >= 1000) return symbol + price.toFixed(0).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
    if (price >= 1) return symbol + price.toFixed(2);
    if (price >= 0.01) return symbol + price.toFixed(4);
    return symbol + price.toFixed(6);
  }
}
