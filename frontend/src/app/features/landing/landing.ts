import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { Header } from '../../shared/components/header/header';
import { LucideActivity, LucideArrowRight, LucideDatabase, LucideShield, LucideWifi } from '@lucide/angular';


@Component({
  selector: 'app-landing',
  imports: [RouterLink, Header, LucideActivity, LucideArrowRight, LucideWifi, LucideDatabase, LucideShield],
  templateUrl: './landing.html',
  styleUrl: './landing.scss',
})
export class Landing {
  features = [
    { icon: 'wifi', title: 'LIVE WEBSOCKET', description: 'Real-time cryptocurrency prices via Binance WebSocket streams. No polling. No delays.' },
    { icon: 'database', title: 'NBP INTEGRATION', description: 'Official Polish National Bank exchange rates. Updated daily. Accurate data.' },
    { icon: 'zap', title: 'INSTANT UPDATES', description: 'Sub-second price updates. Automatic reconnection. Always connected.' },
    { icon: 'shield', title: 'SECURE & RELIABLE', description: 'No server intermediaries. Direct API connections. Privacy-focused.' },
  ];
}
