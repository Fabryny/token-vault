import {
  ChangeDetectionStrategy,
  Component,
  DestroyRef,
  ElementRef,
  effect,
  inject,
  signal,
  viewChild,
} from '@angular/core';
import { animate, createScope, type Scope } from 'animejs';

import { TokenApi } from '../../core/service/token-api';
import { TokenizeResponse } from '../../models/api';
import { UiButton } from '../../shared/ui/ui-button/ui-button';
import { UiField } from '../../shared/ui/ui-field/ui-field';
import { UiPanel } from '../../shared/ui/ui-panel/ui-panel';

/** Alfabeto base32 do crypto/rand.Text() do backend — o embaralhado usa os
 *  mesmos glifos do token real, senão o efeito denuncia que é fachada. */
const GLYPHS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';

@Component({
  selector: 'app-tokenize',
  templateUrl: './tokenize.html',
  imports: [UiField, UiButton, UiPanel],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: `
    .form {
      display: flex;
      flex-direction: column;
      gap: var(--sp-4);
    }

    .vault {
      margin-top: var(--sp-8);
      padding-top: var(--sp-6);
      border-top: 1px solid var(--border);
      display: flex;
      flex-direction: column;
      gap: var(--sp-4);
    }

    .row {
      display: flex;
      align-items: baseline;
      justify-content: space-between;
      gap: var(--sp-4);
    }

    .tag {
      font-size: var(--text-xs);
      letter-spacing: 0.18em;
      text-transform: uppercase;
      color: var(--text-muted);
    }

    .masked { color: var(--text-secondary); }

    .token { flex-direction: column; align-items: stretch; gap: var(--sp-2); }

    .value {
      padding: var(--sp-3) var(--sp-4);
      border-radius: var(--radius-md);
      background: var(--accent-dim);
      border: 1px solid var(--accent);
      color: var(--accent-strong);
      font-size: var(--text-sm);
      word-break: break-all;
      /* min-height evita o painel "pular" antes do texto ser escrito. */
      min-height: 44px;
    }
  `,
})
export class Tokenize {
  private api = inject(TokenApi);
  private host = inject(ElementRef<HTMLElement>);
  private destroyRef = inject(DestroyRef);

  private readonly tokenText = viewChild<ElementRef<HTMLElement>>('tokenText');
  private scope: Scope | null = null;

  protected readonly pan = signal('');
  protected readonly result = signal<TokenizeResponse | null>(null);
  protected readonly error = signal<string | null>(null);
  protected readonly loading = signal(false);
  protected readonly copied = signal(false);

  constructor() {
    // O nó do token só existe depois do @if renderizar. O effect roda
    // quando o viewChild aparece, e aí a animação tem alvo.
    effect(() => {
      const el = this.tokenText()?.nativeElement;
      const res = this.result();
      if (el && res) this.reveal(el, res.token);
    });

    this.destroyRef.onDestroy(() => this.scope?.revert());
  }

  protected submit(event: Event): void {
    event.preventDefault();
    if (this.loading()) return;

    this.error.set(null);
    this.result.set(null);
    this.copied.set(false);
    this.loading.set(true);

    this.api.tokenize(this.pan()).subscribe({
      next: (res) => {
        this.result.set(res);
        this.pan.set(''); // o PAN sai da memória do componente
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Não foi possível tokenizar. Confira o número.');
        this.loading.set(false);
      },
    });
  }

  protected async copy(token: string): Promise<void> {
    await navigator.clipboard.writeText(token);
    this.copied.set(true);
  }

  /**
   * Materialização do token: os caracteres chegam embaralhados e vão se
   * fixando da esquerda para a direita.
   *
   * A animação roda sobre um OBJETO JavaScript ({ p: 0 → 1 }), não sobre
   * o DOM — é o onUpdate que redesenha o texto. Sem esse recurso da v4,
   * o efeito exigiria setInterval na mão.
   *
   * Ela representa substituição, não decoração: é o conceito do projeto
   * acontecendo na tela.
   */
  private reveal(el: HTMLElement, token: string): void {
    this.scope?.revert();

    const reduced = matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (reduced) {
      el.textContent = token; // sem movimento, resultado imediato
      return;
    }

    this.scope = createScope({ root: this.host.nativeElement }).add(() => {
      const state = { p: 0 };

      animate(state, {
        p: 1,
        duration: 900,
        ease: 'outExpo',
        onUpdate: () => {
          const fixed = Math.floor(state.p * token.length);
          let out = token.slice(0, fixed);
          for (let i = fixed; i < token.length; i++) {
            out += GLYPHS[Math.floor(Math.random() * GLYPHS.length)];
          }
          el.textContent = out;
        },
        onComplete: () => {
          el.textContent = token; // garante o valor exato no fim
        },
      });
    });
  }
}
