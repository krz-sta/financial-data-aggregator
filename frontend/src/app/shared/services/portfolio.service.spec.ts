import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { PortfolioService } from './portfolio.service';

describe('PortfolioService', () => {
  let service: PortfolioService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        PortfolioService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(PortfolioService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should add item', () => {
    service.addItem('BTC', 1).subscribe(data => {
      expect(data).toEqual({ success: true });
    });
    const req = httpMock.expectOne('/api/protected/portfolio');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ symbol: 'BTC', amount: 1 });
    req.flush({ success: true });
  });

  it('should delete item', () => {
    service.deleteItem('id1').subscribe();

    const req = httpMock.expectOne('/api/protected/portfolio/id1');
    expect(req.request.method).toBe('DELETE');
    req.flush({ success: true });
  });
});
