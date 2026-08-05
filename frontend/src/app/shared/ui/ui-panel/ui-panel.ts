import { ChangeDetectionStrategy, Component, input } from '@angular/core';

/**
 * Painel de vidro — a superfície onde tudo acontece.
 *
 * A borda superior luminosa é o detalhe que dá o ar "futurista" sem
 * atrapalhar leitura: é um gradiente de 1px, não uma sombra colorida
 * espalhada por trás do texto.
 */
@Component({
  selector: 'ui-panel',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="panel">
      @if (eyebrow() || heading()) {
        <header class="head">
          @if (eyebrow()) {
            <p class="eyebrow">{{ eyebrow() }}</p>
          }
          @if (heading()) {
            <h1 class="heading">{{ heading() }}</h1>
          }
          @if (sub()) {
            <p class="sub">{{ sub() }}</p>
          }
        </header>
      }
      <ng-content />
    </section>
  `,
  styles: `
    :host { display: block; }

    .panel {
      position: relative;
      padding: var(--sp-8);
      border-radius: var(--radius-lg);
      background: rgb(15 21 36 / 72%);
      backdrop-filter: blur(20px);
      border: 1px solid var(--border);
      box-shadow: 0 24px 64px rgb(0 0 0 / 45%);
      overflow: hidden;
    }

    /* Filete luminoso no topo. 1px, não compete com o conteúdo. */
    .panel::before {
      content: '';
      position: absolute;
      inset: 0 0 auto;
      height: 1px;
      background: linear-gradient(
        90deg,
        transparent,
        var(--accent) 35%,
        var(--accent-strong) 50%,
        var(--accent) 65%,
        transparent
      );
      opacity: 0.7;
    }

    .head {
      margin-bottom: var(--sp-8);
    }

    .eyebrow {
      margin-bottom: var(--sp-2);
      font-family: var(--font-mono);
      font-size: var(--text-xs);
      letter-spacing: 0.22em;
      text-transform: uppercase;
      color: var(--accent);
    }

    .heading {
      font-size: var(--text-2xl);
    }

    .sub {
      margin-top: var(--sp-2);
      color: var(--text-secondary);
      font-size: var(--text-sm);
      line-height: var(--leading-body);
    }

    @media (max-width: 520px) {
      .panel { padding: var(--sp-6); }
    }
  `,
})
export class UiPanel {
  readonly eyebrow = input<string | null>(null);
  readonly heading = input<string | null>(null);
  readonly sub = input<string | null>(null);
}
