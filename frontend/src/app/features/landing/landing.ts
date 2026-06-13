import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { Header } from '../../shared/components/header/header';
import { LucideActivity, LucideArrowRight, LucideDatabase, LucideShield, LucideWifi, LucideZap } from '@lucide/angular';
import { Footer } from '../../shared/components/footer/footer';


@Component({
  selector: 'app-landing',
  imports: [RouterLink, Header, Footer, LucideActivity, LucideArrowRight, LucideWifi, LucideDatabase, LucideZap, LucideShield],
  templateUrl: './landing.html',
  styleUrl: './landing.scss',
})
export class Landing {
  private router = inject(Router);
  features = [
    { icon: 'zap', title: 'COINGECKO INTEGRATION', description: 'Accurate cryptocurrency prices and historical data via CoinGecko API.' },
    { icon: 'database', title: 'FRANKFURTER API', description: 'Reliable fiat currency exchange rates tracking against PLN via Frankfurter API.' },
    { icon: 'wifi', title: 'PERIODIC UPDATES', description: 'Background workers fetch and update market data periodically.' },
    { icon: 'shield', title: 'CACHED FOR SPEED', description: 'Redis-backed caching ensures fast data delivery and limits API requests.' },
  ];
}
