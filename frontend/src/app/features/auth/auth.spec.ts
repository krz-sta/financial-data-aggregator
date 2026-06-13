import { TestBed, ComponentFixture } from '@angular/core/testing';
import { provideRouter, Router, ActivatedRoute } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { of, throwError, BehaviorSubject } from 'rxjs';
import { Auth } from './auth';
import { AuthService } from '../../shared/services/auth';
import { vi } from 'vitest';

describe('Auth Component', () => {
  let component: Auth;
  let fixture: ComponentFixture<Auth>;
  let authService: AuthService;
  let router: Router;
  let queryParamsSubject: BehaviorSubject<any>;

  beforeEach(async () => {
    queryParamsSubject = new BehaviorSubject({ mode: 'login' });

    await TestBed.configureTestingModule({
      imports: [Auth],
      providers: [
        provideRouter([]),
        provideHttpClient(),
        provideHttpClientTesting(),
        AuthService,
        {
          provide: ActivatedRoute,
          useValue: { queryParams: queryParamsSubject.asObservable() }
        }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(Auth);
    component = fixture.componentInstance;
    authService = TestBed.inject(AuthService);
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize authMode to login', () => {
    expect(component.authMode()).toBe('login');
  });

  it('should initialize authMode to signup if query param says so', () => {
    queryParamsSubject.next({ mode: 'signup' });
    TestBed.resetTestingModule();
    TestBed.configureTestingModule({
      imports: [Auth],
      providers: [
        provideRouter([]),
        provideHttpClient(),
        provideHttpClientTesting(),
        AuthService,
        {
          provide: ActivatedRoute,
          useValue: { queryParams: of({ mode: 'signup' }) }
        }
      ]
    });
    const fix = TestBed.createComponent(Auth);
    expect(fix.componentInstance.authMode()).toBe('signup');
  });

  it('should set mode and navigate', () => {
    const navigateSpy = vi.spyOn(router, 'navigate').mockImplementation(async () => true);
    component.setMode('signup');
    expect(navigateSpy).toHaveBeenCalled();
    const args = navigateSpy.mock.calls[0];
    expect(args[0]).toEqual([]);
    expect((args[1] as any).queryParams).toEqual({ mode: 'signup' });
  });

  it('should handle login success', () => {
    vi.spyOn(authService, 'login').mockReturnValue(of({ token: 'abc' }));
    const navigateSpy = vi.spyOn(router, 'navigate').mockImplementation(async () => true);
    
    component.formData.set({ email: 't@t.com', password: 'pass', name: '' });
    component.handleSubmit(new Event('submit'));

    expect(authService.login).toHaveBeenCalledWith({ email: 't@t.com', password: 'pass' });
    expect(component.isLoading()).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith(['/home']);
  });

  it('should handle login error', () => {
    vi.spyOn(authService, 'login').mockReturnValue(throwError(() => ({ error: 'Bad creds' })));
    
    component.formData.set({ email: 't@t.com', password: 'pass', name: '' });
    component.handleSubmit(new Event('submit'));

    expect(component.isLoading()).toBe(false);
    expect(component.errorMessage()).toBe('Bad creds');
  });

  it('should handle signup success', () => {
    queryParamsSubject.next({ mode: 'signup' });
    fixture.detectChanges();
    expect(component.authMode()).toBe('signup');

    vi.spyOn(authService, 'register').mockReturnValue(of({ success: true }));
    const setModeSpy = vi.spyOn(component, 'setMode');

    component.formData.set({ email: 't@t.com', password: 'pass', name: 'T' });
    component.handleSubmit(new Event('submit'));

    expect(authService.register).toHaveBeenCalledWith({ email: 't@t.com', password: 'pass', name: 'T' });
    expect(component.isLoading()).toBe(false);
    expect(setModeSpy).toHaveBeenCalledWith('login');
    expect(component.formData()).toEqual({ name: '', email: '', password: '' });
  });

  it('should handle signup error', () => {
    queryParamsSubject.next({ mode: 'signup' });
    fixture.detectChanges();
    vi.spyOn(authService, 'register').mockReturnValue(throwError(() => ({ error: 'Conflict' })));

    component.formData.set({ email: 't@t.com', password: 'pass', name: 'T' });
    component.handleSubmit(new Event('submit'));

    expect(component.isLoading()).toBe(false);
    expect(component.errorMessage()).toBe('Conflict');
  });
});
