# ── Etapa 1: compilar ────────────────────────────────────────────
FROM golang:1.26.5-alpine AS build

WORKDIR /src

# go.mod/go.sum PRIMEIRO, sozinhos. Esta camada só é refeita quando as
# dependencias mudam, nao a cada .go alterado.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0  -> binario estatico, sem depender da libc do sistema;
#                   e o que permite rodar numa imagem minima.
# -trimpath      -> tira o caminho da sua maquina de dentro do binario.
# -ldflags="-s -w" -> descarta simbolos e DWARF: ~30% menor.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ── Etapa 2: rodar ───────────────────────────────────────────────
FROM alpine:3.22

# ca-certificates: sem isso, qualquer chamada HTTPS de dentro do container
#   falha com "x509: certificate signed by unknown authority".
# tzdata: sem isso o time.LoadLocation nao acha fuso nenhum.
RUN apk add --no-cache ca-certificates tzdata

# Usuario sem privilegio: container rodando como root e root no kernel do
# host se houver escape, e a API nao precisa de nada disso.
RUN adduser -D -u 10001 apiuser
USER apiuser

COPY --from=build /out/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]