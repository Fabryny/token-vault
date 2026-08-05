import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { filter, map } from 'rxjs';

import { Auth } from '../../core/service/auth';

/**
 * Casca das telas autenticadas: marca, navegação e saída.
 *
 * Fica como rota-pai de tokenize/detokenize, então a navegação NÃO é
 * remontada a cada troca — é isso que permite a pílula deslizar em vez
 * de piscar.
 */
@Component({
  selector: 'app-shell',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './shell.html',
  styles: `
    :host { display: block; }

    .shell {
      /* Centralizado de verdade: a coluna vive no meio da tela,
         com largura limitada para o texto não ficar quilométrico. */
      width: min(560px, 100%);
      margin-inline: auto;
      min-height: 100dvh;
      display: flex;
      flex-direction: column;
      gap: var(--sp-6);
      padding: var(--sp-8) var(--sp-4) var(--sp-12);
    }

    .bar {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .brand {
      display: flex;
      align-items: center;
      gap: var(--sp-3);
    }

    .mark {
      width: 12px;
      height: 12px;
      border-radius: 3px;
      background: var(--accent);
      box-shadow: 0 0 14px var(--accent-glow);
      animation: pulse 2.6s var(--ease-out) infinite;
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.45; }
    }

    .name {
      font-family: var(--font-display);
      font-weight: 700;
      font-size: var(--text-sm);
      letter-spacing: 0.24em;
      color: var(--text-secondary);
    }

    .name em {
      font-style: normal;
      color: var(--accent);
    }

    .exit {
      min-height: 44px;
      padding: 0 var(--sp-4);
      background: transparent;
      border: 1px solid transparent;
      border-radius: var(--radius-md);
      color: var(--text-muted);
      font-family: var(--font-body);
      font-size: var(--text-sm);
      cursor: pointer;
      transition: color var(--dur-fast) var(--ease-out),
                  border-color var(--dur-fast) var(--ease-out);
    }

    .exit:hover:not(:disabled) {
      color: var(--danger);
      border-color: var(--danger);
    }

    .exit:disabled { opacity: 0.5; cursor: not-allowed; }

    /* ── Segmento ───────────────────────────────────────────── */

    .seg {
      position: relative;
      display: grid;
      grid-template-columns: 1fr 1fr;
      padding: 4px;
      border-radius: var(--radius-full);
      background: rgb(255 255 255 / 3%);
      backdrop-filter: blur(12px);
      border: 1px solid var(--border);
    }

    .pill {
      position: absolute;
      top: 4px;
      left: 4px;
      width: calc(50% - 4px);
      height: calc(100% - 8px);
      border-radius: var(--radius-full);
      background: var(--accent-dim);
      border: 1px solid var(--accent);
      box-shadow: 0 0 20px var(--accent-glow);
      transition: transform var(--dur-slow) var(--ease-out);
    }

    .tab {
      position: relative; /* acima da pílula */
      z-index: var(--z-sticky);
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      border-radius: var(--radius-full);
      color: var(--text-muted);
      font-size: var(--text-sm);
      font-weight: 600;
      text-decoration: none;
      transition: color var(--dur) var(--ease-out);
    }

    .tab:hover { color: var(--text-secondary); }
    .tab.on { color: var(--accent-strong); }

    .stage { flex: 1; }
  `,
})
export class Shell {
  private readonly router = inject(Router);
  private readonly auth = inject(Auth);

  protected readonly leaving = signal(false);

  /** URL atual como signal, para a pílula reagir sem ChangeDetection manual. */
  private readonly url = toSignal(
    this.router.events.pipe(
      filter((e): e is NavigationEnd => e instanceof NavigationEnd),
      map((e) => e.urlAfterRedirects),
    ),
    { initialValue: this.router.url },
  );

  protected readonly active = computed(() =>
    this.url().includes('detokenize') ? 'detokenize' : 'tokenize',
  );

  /** 0% ou 100% da própria largura — um transform, sem tocar no layout. */
  protected readonly pillShift = computed(() =>
    this.active() === 'detokenize' ? 'translateX(100%)' : 'translateX(0)',
  );

  protected logout(): void {
    this.leaving.set(true);
    this.auth.logout().subscribe({
      next: () => this.router.navigate(['/login']),
      // Falhou? A sessão local já foi marcada como encerrada — sai mesmo assim.
      error: () => this.router.navigate(['/login']),
    });
  }
}
