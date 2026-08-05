import { ApplicationConfig, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authErrorInterceptor } from './core/interceptor/auth-error';

export const appConfig: ApplicationConfig = {
  providers: [
    provideHttpClient(withInterceptors([authErrorInterceptor])),
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes)
  ]
};
