import { ChangeDetectionStrategy, Component, input } from '@angular/core';

export type UiButtonVariant = 'primary' | 'ghost' | 'danger';

/**
 * Botão base do app.
 *
 * Regras embutidas para não serem esquecidas em cada uso:
 * - altura mínima de 44px (alvo de toque)
 * - desabilita sozinho enquanto `loading`, com aria-busy
 * - feedback de pressão por `scale`, que NÃO desloca o layout ao redor
 *
 * Uso: <ui-button variant="ghost" [loading]="salvando()">Salvar</ui-button>
 */
@Component({
  selector: 'ui-button',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <button
      class="btn"
      [class]="variant()"
      [type]="type()"
      [disabled]="disabled() || loading()"
      [attr.aria-busy]="loading() ? 'true' : null"
    >
      @if (loading()) {
        <span class="spinner" aria-hidden="true"></span>
      }
      <span class="label"><ng-content /></span>
    </button>
  `,
  styles: `
    :host { display: inline-block; }
    :host(.full) { display: block; }

    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: var(--sp-2);
      width: 100%;
      min-height: 44px;
      padding: var(--sp-3) var(--sp-6);
      border: 1px solid transparent;
      border-radius: var(--radius-md);
      font-family: var(--font-body);
      font-size: var(--text-md);
      font-weight: 600;
      cursor: pointer;
      transition:
        background-color var(--dur-fast) var(--ease-out),
        border-color var(--dur-fast) var(--ease-out),
        transform var(--dur-fast) var(--ease-out),
        box-shadow var(--dur-fast) var(--ease-out);
    }

    .btn:active:not(:disabled) { transform: scale(0.98); }

    /* 0.5 de opacidade + cursor: o desabilitado tem que PARECER
       desabilitado, não só ignorar o clique. */
    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .primary {
      background: var(--accent);
      color: #04121a;
    }
    .primary:hover:not(:disabled) {
      background: var(--accent-strong);
      box-shadow: 0 0 24px var(--accent-glow);
    }

    .ghost {
      background: transparent;
      border-color: var(--border-strong);
      color: var(--text-primary);
    }
    .ghost:hover:not(:disabled) {
      background: var(--accent-dim);
      border-color: var(--accent);
    }

    .danger {
      background: var(--danger-dim);
      border-color: var(--danger);
      color: var(--danger);
    }
    .danger:hover:not(:disabled) {
      background: var(--danger);
      color: var(--surface-base);
    }

    .spinner {
      width: 16px;
      height: 16px;
      border: 2px solid currentColor;
      border-top-color: transparent;
      border-radius: var(--radius-full);
      animation: spin 700ms linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(1turn); }
    }
  `,
})
export class UiButton {
  readonly variant = input<UiButtonVariant>('primary');
  readonly type = input<'button' | 'submit'>('button');
  readonly disabled = input(false);
  readonly loading = input(false);
}
