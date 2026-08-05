import { ChangeDetectionStrategy, Component, computed, input, model } from '@angular/core';

/**
 * Campo com float label — sem PrimeNG, só CSS.
 *
 * O truque: placeholder=" " (um espaço). Isso faz :placeholder-shown ser
 * verdadeiro só enquanto o campo está vazio, e é o que permite a label
 * subir sozinha, sem JavaScript nenhum.
 *
 * Regras de UX embutidas:
 * - label sempre presente (sobe, nunca some — placeholder não é label)
 * - erro ABAIXO do campo, com role="alert" e aria-describedby
 * - altura mínima de 56px (alvo de toque com folga)
 */
@Component({
  selector: 'ui-field',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="field" [class.invalid]="!!error()">
      <input
        class="input"
        placeholder=" "
        [id]="id"
        [class.mono]="mono()"
        [type]="type()"
        [value]="value()"
        [required]="required()"
        [attr.inputmode]="inputmode()"
        [attr.autocomplete]="autocomplete()"
        [attr.maxlength]="maxlength()"
        [attr.aria-invalid]="error() ? 'true' : null"
        [attr.aria-describedby]="describedBy()"
        (input)="onInput($event)"
      />

      <label class="label" [for]="id">
        {{ label() }}
        @if (required()) {
          <span aria-hidden="true">*</span>
          <span class="sr-only">(obrigatório)</span>
        }
      </label>

      <span class="beam" aria-hidden="true"></span>
    </div>

    @if (error()) {
      <p class="msg error" [id]="id + '-err'" role="alert">{{ error() }}</p>
    } @else if (hint()) {
      <p class="msg hint" [id]="id + '-hint'">{{ hint() }}</p>
    }
  `,
  styles: `
    :host { display: block; }

    .field {
      position: relative;
      border-radius: var(--radius-md);
      background: rgb(255 255 255 / 3%);
      backdrop-filter: blur(12px);
      border: 1px solid var(--border);
      transition:
        border-color var(--dur) var(--ease-out),
        background-color var(--dur) var(--ease-out);
      overflow: hidden;
    }

    .field:hover { border-color: var(--border-strong); }

    .field:has(.input:focus) {
      border-color: var(--accent);
      background: rgb(34 211 238 / 6%);
    }

    .field.invalid { border-color: var(--danger); }

    .input {
      width: 100%;
      min-height: 56px;
      /* Espaço extra no topo é onde a label pousa quando sobe. */
      padding: 22px var(--sp-4) 8px;
      background: transparent;
      border: 0;
      outline: none;
      color: var(--text-primary);
      font-family: var(--font-body);
      font-size: var(--text-md);
    }

    .input.mono {
      font-family: var(--font-mono);
      font-variant-numeric: tabular-nums;
      letter-spacing: 0.1em;
    }

    .label {
      position: absolute;
      left: var(--sp-4);
      top: 50%;
      transform: translateY(-50%);
      color: var(--text-muted);
      font-size: var(--text-md);
      pointer-events: none;
      transition:
        top var(--dur) var(--ease-out),
        transform var(--dur) var(--ease-out),
        font-size var(--dur) var(--ease-out),
        color var(--dur) var(--ease-out);
    }

    /* Sobe quando o campo tem foco OU já tem conteúdo. */
    .input:focus + .label,
    .input:not(:placeholder-shown) + .label {
      top: 9px;
      transform: none;
      font-size: var(--text-xs);
      letter-spacing: 0.08em;
      text-transform: uppercase;
      color: var(--accent);
    }

    .field.invalid .input:focus + .label,
    .field.invalid .input:not(:placeholder-shown) + .label {
      color: var(--danger);
    }

    /* Feixe que corre da esquerda para a direita ao focar. */
    .beam {
      position: absolute;
      left: 0;
      bottom: 0;
      height: 1px;
      width: 100%;
      background: linear-gradient(90deg, transparent, var(--accent), transparent);
      transform: scaleX(0);
      transition: transform var(--dur-slow) var(--ease-out);
    }

    .field:has(.input:focus) .beam { transform: scaleX(1); }
    .field.invalid .beam {
      background: linear-gradient(90deg, transparent, var(--danger), transparent);
    }

    .msg {
      margin-top: var(--sp-2);
      padding-left: var(--sp-1);
      font-size: var(--text-sm);
      line-height: 1.4;
    }

    /* A cor não é o único sinal: há o texto e o aria-invalid. */
    .error { color: var(--danger); }
    .hint { color: var(--text-muted); }
  `,
})
export class UiField {
  private static seq = 0;
  protected readonly id = `uif-${UiField.seq++}`;

  readonly label = input.required<string>();
  readonly value = model<string>('');

  readonly type = input<'text' | 'email' | 'password'>('text');
  readonly required = input(false);
  readonly mono = input(false);

  readonly inputmode = input<string | null>(null);
  readonly autocomplete = input<string | null>(null);
  readonly maxlength = input<number | null>(null);

  /** Descarta tudo que não for dígito — inclusive espaço, letra e sinal. */
  readonly digitsOnly = input(false);

  readonly error = input<string | null>(null);
  readonly hint = input<string | null>(null);

  protected readonly describedBy = computed(() => {
    if (this.error()) return `${this.id}-err`;
    if (this.hint()) return `${this.id}-hint`;
    return null;
  });

  protected onInput(event: Event): void {
    const el = event.target as HTMLInputElement;
    const raw = el.value;
    const clean = this.digitsOnly() ? raw.replace(/\D/g, '') : raw;

    if (clean !== raw) {
      // Só atualizar o signal NÃO limpa a tela: se o valor limpo for igual
      // ao anterior, o Angular não vê mudança e não reescreve o input —
      // o caractere inválido fica visível. Por isso mexemos no elemento.
      //
      // E ao reescrever, o cursor pularia para o fim. Descontamos quantos
      // caracteres caíram para devolvê-lo ao lugar, senão editar o meio
      // do número vira tortura.
      const pos = (el.selectionStart ?? raw.length) - (raw.length - clean.length);
      el.value = clean;
      el.setSelectionRange(pos, pos);
    }

    this.value.set(clean);
  }
}
