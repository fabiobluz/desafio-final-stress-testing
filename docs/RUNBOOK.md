# RUNBOOK

## Compilar
```bash
go build -o bin/gload ./cmd/gload
```

## Executar
```bash
./bin/gload --url=http://localhost:8080 --requests=100 --concurrency=10
```

## Docker
```bash
docker build -t acme/gload:latest .
docker run --rm acme/gload:latest --url=http://google.com --requests=1000 --concurrency=10
```
