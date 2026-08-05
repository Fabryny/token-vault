import { Injectable, computed, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, map, of, tap } from 'rxjs';

import { Credentials, MeResponse } from '../../models/api';

@Injectable({ providedIn: 'root' })
export class Auth {
  private http = inject(HttpClient);

  private readonly _authenticated = signal<boolean | null>(null);

  readonly isAuthenticated = computed(() => this._authenticated() === true);

  register(credentials: Credentials): Observable<void> {
    return this.http.post<void>('/api/auth/register', credentials);
  }

  login(credentials: Credentials): Observable<void> {
    // Resposta 204 sem corpo: o que importa veio no header Set-Cookie.
    return this.http
      .post<void>('/api/auth/login', credentials)
      .pipe(tap(() => this._authenticated.set(true)));
  }

  logout(): Observable<void> {
    // Tem que ser requisição: JavaScript não apaga cookie HttpOnly.
    return this.http
      .post<void>('/api/auth/logout', {})
      .pipe(tap(() => this._authenticated.set(false)));
  }

  checkSession(): Observable<boolean> {
    const known = this._authenticated();
    if (known !== null) {
      return of(known); 
    }

    return this.http.get<MeResponse>('/api/auth/me').pipe(
      map(() => {
        this._authenticated.set(true);
        return true;
      }),
      catchError(() => {
        this._authenticated.set(false);
        return of(false);
      }),
    );
  }

  // Chamado pelo interceptor quando o backend devolve 401.
  markLoggedOut(): void {
    this._authenticated.set(false);
  }
}