import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./features/landing/landing').then(m => m.Landing),
    data: { showAuthButtons: true}
  },
  {
    path: 'auth',
    loadComponent: () => import('./features/auth/auth').then(m => m.Auth),
    data: { showAuthButtons: false }
  },
  {
    path: '**',
    redirectTo: ''
  }
];