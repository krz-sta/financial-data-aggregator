import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient, withInterceptors, HttpClient } from '@angular/common/http';
import { authInterceptor } from './auth.interceptor';
import { AuthService } from '../services/auth';

describe('AuthInterceptor', () => {
  let authService: AuthService;
  let httpMock: HttpTestingController;
  let httpClient: HttpClient;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        AuthService,
        provideHttpClient(withInterceptors([authInterceptor])),
        provideHttpClientTesting()
      ]
    });
    authService = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
    httpClient = TestBed.inject(HttpClient);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should not add auth header if no token', () => {
    authService.getToken = () => null;
    httpClient.get('/api/test').subscribe();
    const req = httpMock.expectOne('/api/test');
    expect(req.request.headers.has('Authorization')).toBe(false);
  });

  it('should add auth header if token exists', () => {
    authService.getToken = () => 'fake-token';
    httpClient.get('/api/test').subscribe();
    const req = httpMock.expectOne('/api/test');
    expect(req.request.headers.get('Authorization')).toBe('Bearer fake-token');
  });
});
