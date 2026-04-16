import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { UrlService, ShortenResponse, UrlStats } from './url.service';

@Component({
  selector: 'app-root',
  imports: [CommonModule, FormsModule],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {
  originalUrl = '';
  statsCode = '';
  loading = false;
  error = '';
  success = '';

  result: ShortenResponse | null = null;
  stats: UrlStats | null = null;
  recentUrls: ShortenResponse[] = [];

  constructor(private readonly urlService: UrlService) {}

  shortenUrl(): void {
    this.error = '';
    this.success = '';

    if (!this.originalUrl.trim()) {
      this.error = 'Please enter a URL.';
      return;
    }

    this.loading = true;
    this.urlService.shorten(this.originalUrl.trim()).subscribe({
      next: (response) => {
        this.result = response;
        this.success = 'Short URL created successfully.';
        this.statsCode = response.code;
        this.recentUrls = [response, ...this.recentUrls].slice(0, 5);
        this.loading = false;
      },
      error: (err) => {
        this.error = err?.error?.error ?? 'Failed to shorten URL.';
        this.loading = false;
      },
    });
  }

  fetchStats(): void {
    this.error = '';
    this.success = '';
    this.stats = null;

    if (!this.statsCode.trim()) {
      this.error = 'Enter a short code to fetch stats.';
      return;
    }

    this.loading = true;
    this.urlService.getStats(this.statsCode.trim()).subscribe({
      next: (response) => {
        this.stats = response;
        this.loading = false;
      },
      error: (err) => {
        this.error = err?.error?.error ?? 'Failed to fetch stats.';
        this.loading = false;
      },
    });
  }

  copyShortUrl(): void {
    if (!this.result?.short_url) {
      return;
    }
    navigator.clipboard.writeText(this.result.short_url);
    this.success = 'Copied short URL to clipboard.';
  }
}
