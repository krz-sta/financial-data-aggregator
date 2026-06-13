import { Component, inject, OnInit, signal } from '@angular/core';
import { Header } from '../../shared/components/header/header';
import { Footer } from '../../shared/components/footer/footer';
import { UserService } from '../../shared/services/user.service';
import { UpperCasePipe } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [Header, Footer, UpperCasePipe],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home implements OnInit {
  private userService = inject(UserService);
  
  profileData = signal<any>(null);
  isLoading = signal(true);
  error = signal<string | null>(null);

  ngOnInit() {
    this.userService.getProfile().subscribe({
      next: (data) => {
        this.profileData.set(data);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.error.set('Failed to load profile.');
        this.isLoading.set(false);
      }
    });
  }
}