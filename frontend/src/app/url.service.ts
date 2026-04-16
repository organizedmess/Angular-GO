import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface ShortenResponse {
  code: string;
  short_url: string;
  original_url: string;
}

export interface UrlStats {
  code: string;
  short_url: string;
  original_url: string;
  click_count: number;
  created_at: string;
}

@Injectable({ providedIn: 'root' })
export class UrlService {
  private readonly api = 'http://localhost:8080/api';

  constructor(private readonly http: HttpClient) {}

  shorten(originalUrl: string) {
    return this.http.post<ShortenResponse>(`${this.api}/shorten`, {
      original_url: originalUrl,
    });
  }

  getStats(code: string) {
    return this.http.get<UrlStats>(`${this.api}/stats/${code}`);
  }
}
