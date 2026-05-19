import { Component } from '@angular/core';
import { LucideSearch, LucideTrendingUp } from '@lucide/angular';

@Component({
  selector: 'app-header',
  imports: [LucideTrendingUp, LucideSearch],
  templateUrl: './header.html',
  styleUrl: './header.scss',
})
export class Header {}
