import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { AuthService } from './auth';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        AuthService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
    localStorage.clear();
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should login and set token', () => {
    service.login({ email: 'test@test.com', password: 'password' }).subscribe();
    const req = httpMock.expectOne('http://localhost:8080/api/auth/login');
    expect(req.request.method).toBe('POST');
    req.flush({ token: 'fake-token' });
    expect(localStorage.getItem('auth_token')).toBe('fake-token');
    expect(service.isLoggedIn()).toBe(true);
  });

  it('should register', () => {
    service.register({ email: 'test@test.com', name: 'Test', password: 'password' }).subscribe();
    const req = httpMock.expectOne('http://localhost:8080/api/auth/register');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ email: 'test@test.com', displayName: 'Test', password: 'password' });
    req.flush({ success: true });
  });

  it('should logout', () => {
    localStorage.setItem('auth_token', 'test-token');
    service.logout();
    expect(localStorage.getItem('auth_token')).toBeNull();
    expect(service.isLoggedIn()).toBe(false);
  });

  it('should get token', () => {
    localStorage.setItem('auth_token', 'my-token');
    expect(service.getToken()).toBe('my-token');
  });
});
