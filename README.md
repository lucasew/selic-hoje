# Selic Hoje

Taxa Selic direto ao ponto, desenterrada do site do [Banco Central do Brasil](https://www.bcb.gov.br/).

Site: [selic-hoje.vercel.app](https://selic-hoje.vercel.app)

## API

`GET /api/selichoje` — texto puro com a taxa (% a.a.). Resposta cacheada por 1 hora (`Cache-Control: max-age=3600`).

## Desenvolvimento

Ferramentas via [mise](https://mise.jdx.dev/) (`mise.toml`):

```bash
mise install
go test ./...
```

Ou entre no shell Nix (`shell.nix`) para `go` e `gopls`. Deploy é o projeto Vercel ligado a este repositório.
