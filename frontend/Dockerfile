# ── Etapa 1: compilar ────────────────────────────────────────────
FROM node:24-alpine AS build

WORKDIR /app

# npm ci (nao npm install): instala EXATAMENTE o package-lock.json.
# O install pode resolver versoes diferentes e quebrar a build de
# forma nao reproduzivel.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# ── Etapa 2: servir ──────────────────────────────────────────────
FROM nginx:1.29-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf

# 🔑 O CAMINHO TEM /browser NO FIM.
# O builder @angular/build:application gera dist/<projeto>/browser/.
# Copiar dist/frontend daria uma pasta browser/ aninhada e todo request
# viraria 404 — e o log do nginx nao explica o motivo.
COPY --from=build /app/dist/frontend/browser /usr/share/nginx/html

EXPOSE 80