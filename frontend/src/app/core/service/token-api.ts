import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { DetokenizeResponse, TokenizeResponse } from '../../models/api';

@Injectable({ providedIn: 'root' })
export class TokenApi {
  private http = inject(HttpClient);

  tokenize(pan: string): Observable<TokenizeResponse> {
    return this.http.post<TokenizeResponse>('/api/tokenize', { pan });
  }

  detokenize(token: string): Observable<DetokenizeResponse> {
    return this.http.post<DetokenizeResponse>('/api/detokenize', { token });
  }
}