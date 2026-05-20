import { Component, inject } from '@angular/core';
import { ActivatedRoute, NavigationEnd, Router, RouterLink } from '@angular/router';
import { LucideArrowLeft, LucideSearch, LucideTrendingUp } from '@lucide/angular';
import { filter, map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-header',
  imports: [LucideTrendingUp, LucideSearch, RouterLink, LucideArrowLeft],
  templateUrl: './header.html',
  styleUrl: './header.scss',
})
export class Header {
  private router = inject(Router);
  private activatedRoute = inject(ActivatedRoute);

  showAuthButtons = toSignal(
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map(() => {
        let route = this.activatedRoute;
        while (route.firstChild) {
          route = route.firstChild;
        }
        return route.snapshot.data['showAuthButtons'] !== false;
      })
    ),
    { initialValue: true}
  );
}
