import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';

import { Auth } from '../../core/service/auth';
import { UiButton } from '../../shared/ui/ui-button/ui-button';
import { UiField } from '../../shared/ui/ui-field/ui-field';
import { UiPanel } from '../../shared/ui/ui-panel/ui-panel';

@Component({
  selector: 'app-login',
  templateUrl: './login.html',
  imports: [UiField, UiButton, UiPanel, RouterLink],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: `
    .form {
      display: flex;
      flex-direction: column;
      gap: var(--sp-4);
    }

    .foot {
      margin-top: var(--sp-6);
      text-align: center;
      font-size: var(--text-sm);
      color: var(--text-muted);
    }

    .foot a {
      color: var(--accent);
      text-decoration: none;
      font-weight: 500;
    }

    .foot a:hover { text-decoration: underline; }
  `,
})
export class Login {
  private auth = inject(Auth);
  private router = inject(Router);

  protected readonly email = signal('');
  protected readonly password = signal('');
  protected readonly error = signal<string | null>(null);
  protected readonly loading = signal(false);

  protected submit(event: Event): void {
    event.preventDefault();
    if (this.loading()) return;

    this.error.set(null);
    this.loading.set(true);

    this.auth.login({ email: this.email(), password: this.password() }).subscribe({
      next: () => this.router.navigate(['/tokenize']),
      error: () => {
        // Mensagem única: o backend não diz se o e-mail existe, e o
        // front não pode inventar essa distinção.
        this.error.set('E-mail ou senha inválidos.');
        this.loading.set(false);
      },
    });
  }
}
