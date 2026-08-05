import { ChangeDetectionStrategy, Component, DestroyRef, inject, signal } from '@angular/core';

import { TokenApi } from '../../core/service/token-api';
import { UiButton } from '../../shared/ui/ui-button/ui-button';
import { UiField } from '../../shared/ui/ui-field/ui-field';
import { UiPanel } from '../../shared/ui/ui-panel/ui-panel';

/** Segundos que o cartão fica visível antes de sumir sozinho. */
const REVEAL_SECONDS = 20;

@Component({
  selector: 'app-detokenize',
  templateUrl: './detokenize.html',
  imports: [UiField, UiButton, UiPanel],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: `
    .form {
      display: flex;
      flex-direction: column;
      gap: var(--sp-4);
    }

    .exposed {
      margin-top: var(--sp-8);
      padding: var(--sp-4);
      border-radius: var(--radius-md);
      background: var(--sensitive-dim);
      border: 1px solid var(--sensitive);
      display: flex;
      flex-direction: column;
      gap: var(--sp-3);
    }

    .warn {
      display: flex;
      align-items: center;
      gap: var(--sp-2);
      font-size: var(--text-xs);
      letter-spacing: 0.08em;
      text-transform: uppercase;
      color: var(--sensitive);
    }

    .dot {
      width: 7px;
      height: 7px;
      border-radius: var(--radius-full);
      background: var(--sensitive);
      animation: blink 1.4s var(--ease-out) infinite;
    }

    @keyframes blink {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.25; }
    }

    .value {
      font-size: var(--text-xl);
      color: var(--text-primary);
      word-break: break-all;
    }
  `,
})
export class Detokenize {
  private api = inject(TokenApi);
  private destroyRef = inject(DestroyRef);

  protected readonly token = signal('');
  protected readonly pan = signal<string | null>(null);
  protected readonly error = signal<string | null>(null);
  protected readonly loading = signal(false);
  protected readonly countdown = signal(REVEAL_SECONDS);

  private timer: ReturnType<typeof setInterval> | null = null;

  constructor() {
    // Sair da tela com o cartão visível não pode deixar timer rodando.
    this.destroyRef.onDestroy(() => this.stopTimer());
  }

  protected submit(event: Event): void {
    event.preventDefault();
    if (this.loading()) return;

    this.error.set(null);
    this.pan.set(null);
    this.stopTimer();
    this.loading.set(true);

    this.api.detokenize(this.token()).subscribe({
      next: (res) => {
        this.pan.set(res.pan);
        this.loading.set(false);
        this.startTimer();
      },
      error: () => {
        // 404 aqui significa "não existe OU não é seu" — o backend não
        // distingue de propósito, e o front não pode adivinhar.
        this.error.set('Token não encontrado.');
        this.loading.set(false);
      },
    });
  }

  protected clear(): void {
    this.stopTimer();
    this.pan.set(null);
    this.token.set('');
  }

  /**
   * Auto-ocultar não é enfeite: cartão em claro esquecido numa aba aberta
   * é exposição a quem passar pela mesa. O contador deixa o prazo visível
   * em vez de sumir de surpresa.
   */
  private startTimer(): void {
    this.countdown.set(REVEAL_SECONDS);
    this.timer = setInterval(() => {
      const left = this.countdown() - 1;
      this.countdown.set(left);
      if (left <= 0) this.clear();
    }, 1000);
  }

  private stopTimer(): void {
    if (this.timer) clearInterval(this.timer);
    this.timer = null;
    this.countdown.set(REVEAL_SECONDS);
  }
}
