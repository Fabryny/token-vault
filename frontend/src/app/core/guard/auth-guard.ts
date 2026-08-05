import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { map } from 'rxjs';

import { Auth } from '../service/auth';

export const authGuard: CanActivateFn = () => {
  const auth = inject(Auth);
  const router = inject(Router);

  return auth.checkSession().pipe(
    map((ok) => (ok ? true : router.parseUrl('/login'))),
  );
};