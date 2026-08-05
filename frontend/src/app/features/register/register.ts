import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Router, RouterLink } from '@angular/router';
import { switchMap } from 'rxjs';

import { Auth } from '../../core/service/auth';
import { UiButton } from '../../shared/ui/ui-button/ui-button';
import { UiField } from '../../shared/ui/ui-field/ui-field';
import { UiPanel } from '../../shared/ui/ui-panel/ui-panel';

/** Mesmo limite do bcrypt no backend: acima de 72 BYTES ele rejeita. */
const MAX_PASSWORD = 72;
const MIN_PASSWORD = 8;

@Component({
  selector: 'app-register',
  templateUrl: './register.html',
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
export class Register {
  private auth = inject(Auth);
  private router = inject(Router);

  protected readonly email = signal('');
  protected readonly password = signal('');
  protected readonly confirm = signal('');
  protected readonly error = signal<string | null>(null);
  protected readonly loading = signal(false);

  /** Só valida depois da primeira tentativa: não acusa erro enquanto digita. */
  private readonly tried = signal(false);

  protected readonly emailError = computed(() =>
    this.tried() && !this.email().includes('@') ? 'Informe um e-mail válido.' : null,
  );

  protected readonly passwordError = computed(() => {
    if (!this.tried()) return null;
    const len = new TextEncoder().encode(this.password()).length; // BYTES, não caracteres
    if (len < MIN_PASSWORD) return `Mínimo de ${MIN_PASSWORD} caracteres.`;
    if (len > MAX_PASSWORD) return `Máximo de ${MAX_PASSWORD} bytes (acentos contam mais).`;
    return null;
  });

  protected readonly confirmError = computed(() =>
    this.tried() && this.confirm() !== this.password() ? 'As senhas não conferem.' : null,
  );

  protected submit(event: Event): void {
    event.preventDefault();
    if (this.loading()) return;

    this.tried.set(true);
    this.error.set(null);

    if (this.emailError() || this.passwordError() || this.confirmError()) return;

    this.loading.set(true);
    const credentials = { email: this.email(), password: this.password() };

    // Cadastra e já entra: pedir para logar em seguida seria repetir
    // dados que o usuário acabou de digitar.
    this.auth
      .register(credentials)
      .pipe(switchMap(() => this.auth.login(credentials)))
      .subscribe({
        next: () => this.router.navigate(['/tokenize']),
        error: (err: HttpErrorResponse) => {
          this.error.set(
            err.status === 409
              ? 'Este e-mail já está cadastrado.'
              : 'Não foi possível criar a conta. Tente de novo.',
          );
          this.loading.set(false);
        },
      });
  }
}
