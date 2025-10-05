# gload (Go Load Tester)

CLI em Go para testes de carga simples.

## Instalação
```bash
go mod tidy
```

## Uso (local)
```bash
go run ./cmd/gload --url=http://localhost:8080 --requests=100 --concurrency=10
```

## Relatório
- Tempo total
- Total de requests
- HTTP 200
- Distribuição de status
- Latências: média, P95, P99
