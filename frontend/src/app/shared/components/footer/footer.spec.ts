import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { Footer } from './footer';

describe('Footer', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Footer],
      providers: [provideRouter([])]
    }).compileComponents();
  });

  it('should create', () => {
    const fixture = TestBed.createComponent(Footer);
    fixture.detectChanges();
    const component = fixture.componentInstance;
    expect(component).toBeTruthy();
  });
});
