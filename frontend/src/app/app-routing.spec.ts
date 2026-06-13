import { TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { Location } from '@angular/common';
import { routes } from './app.routes';
import { App } from './app';
import { AuthService } from './shared/services/auth';

describe('App Routing Integration', () => {
  let router: Router;
  let location: Location;
  let authService: AuthService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [App],
      providers: [
        provideRouter(routes),
        provideHttpClient(),
        provideHttpClientTesting(),
        AuthService
      ]
    });

    router = TestBed.inject(Router);
    location = TestBed.inject(Location);
    authService = TestBed.inject(AuthService);
    localStorage.clear();
  });

  it('should navigate to /landing by default', async () => {
    await router.navigate(['']);
    expect(location.path()).toBe('');
  }, 10000);

  it('should allow navigation to /auth', async () => {
    await router.navigate(['/auth']);
    expect(location.path()).toBe('/auth');
  });

  it('should redirect to /auth if accessing /home without being logged in', async () => {
    authService.isLoggedIn.set(false);
    await router.navigate(['/home']);
    expect(location.path()).toBe('/auth');
  });

  it('should allow navigation to /home if logged in', async () => {
    authService.isLoggedIn.set(true);
    await router.navigate(['/home']);
    expect(location.path()).toBe('/home');
  });
});
