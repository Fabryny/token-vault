import { Routes } from '@angular/router';

import { authGuard } from './core/guard/auth-guard';
import { Shell } from './layout/shell/shell';
import { Detokenize } from './features/detokenize/detokenize';
import { Login } from './features/login/login';
import { Register } from './features/register/register';
import { Tokenize } from './features/tokenize/tokenize';

export const routes: Routes = [
  { path: '', redirectTo: 'tokenize', pathMatch: 'full' },

  // Públicas — tela cheia, sem a casca autenticada.
  { path: 'login', component: Login },
  { path: 'register', component: Register },

  // Autenticadas — o Shell é rota-PAI, então a navegação não é remontada
  // a cada troca. É isso que permite a pílula deslizar em vez de piscar.
  // O guard no pai já cobre os filhos.
  {
    path: '',
    component: Shell,
    canActivate: [authGuard],
    children: [
      { path: 'tokenize', component: Tokenize },
      { path: 'detokenize', component: Detokenize },
    ],
  },

  { path: '**', redirectTo: 'tokenize' },
];
