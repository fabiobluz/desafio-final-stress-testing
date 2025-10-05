# gload - Go Load Tester

Um sistema CLI em Go para realizar testes de carga em serviços web de forma eficiente e com relatórios detalhados.

## 🚀 Características

- **Testes de carga HTTP**: Suporte completo para métodos HTTP (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS)
- **Controle de concorrência**: Configuração precisa do número de workers simultâneos
- **Headers customizados**: Suporte para múltiplos headers HTTP
- **Body de requisição**: Envio de dados no corpo da requisição (texto ou arquivo)
- **Múltiplos formatos de relatório**: Texto legível ou JSON estruturado
- **Métricas avançadas**: Latência média, P95, P99 e distribuição de códigos de status
- **Docker ready**: Imagem otimizada para execução em containers
- **Performance otimizada**: Pool de conexões HTTP configurado para alta performance

## 📦 Instalação

### Opção 1: Build local
```bash
git clone <repository-url>
cd desafio-final-stress-testing
go mod tidy
make build
```

### Opção 2: Docker (Recomendado)
```bash
docker build -t gload:latest .
```

## 🛠️ Uso

### Parâmetros obrigatórios
- `--url`: URL do serviço a ser testado (http/https)
- `--requests`: Número total de requests a serem executados
- `--concurrency`: Número de workers simultâneos (concorrência)

### Parâmetros opcionais
- `--timeout`: Timeout por request (padrão: 10s)
- `--method`: Método HTTP (padrão: GET)
- `--format`: Formato do relatório - `text` ou `json` (padrão: text)
- `--body`: Corpo da requisição (texto ou @arquivo)
- `--header`: Headers HTTP no formato 'Chave:Valor' (pode repetir)

### Exemplos de uso

#### Teste básico
```bash
# Local
./bin/gload --url=http://google.com --requests=1000 --concurrency=10

# Docker
docker run gload:latest --url=http://google.com --requests=1000 --concurrency=10
```

#### Teste com POST e JSON
```bash
./bin/gload \
  --url=https://api.exemplo.com/users \
  --method=POST \
  --requests=500 \
  --concurrency=5 \
  --header="Content-Type:application/json" \
  --header="Authorization:Bearer token123" \
  --body='{"name":"João","email":"joao@exemplo.com"}' \
  --format=json
```

#### Teste com arquivo no body
```bash
./bin/gload \
  --url=https://api.exemplo.com/upload \
  --method=POST \
  --requests=100 \
  --concurrency=3 \
  --header="Content-Type:application/octet-stream" \
  --body=@arquivo.pdf
```

#### Teste com timeout customizado
```bash
./bin/gload \
  --url=http://servico-lento.com \
  --requests=200 \
  --concurrency=5 \
  --timeout=30s
```

## 📊 Relatório de Resultados

O sistema gera relatórios detalhados com as seguintes métricas:

### Formato Texto (padrão)
```
Total time: 2.345s
Total requests: 1000
HTTP 200: 987
Status distribution: map[200:987 404:10 500:3]
Avg latency: 45ms
P95: 120ms
P99: 250ms
```

### Formato JSON
```json
{
  "total_requests": 1000,
  "success_200": 987,
  "status_distribution": {
    "200": 987,
    "404": 10,
    "500": 3
  },
  "total_time": "2.345s",
  "avg_latency": "45ms",
  "p95_latency": "120ms",
  "p99_latency": "250ms"
}
```

### Métricas explicadas
- **Total time**: Tempo total de execução do teste
- **Total requests**: Quantidade total de requests executados
- **HTTP 200**: Quantidade de requests com status 200 (sucesso)
- **Status distribution**: Distribuição de todos os códigos de status HTTP retornados
- **Avg latency**: Latência média dos requests
- **P95**: 95º percentil de latência (95% dos requests foram mais rápidos)
- **P99**: 99º percentil de latência (99% dos requests foram mais rápidos)

## 🐳 Docker

### Build da imagem
```bash
make docker-build
# ou
docker build -t gload:latest .
```

### Execução via Docker
```bash
# Teste básico
docker run gload:latest --url=http://google.com --requests=1000 --concurrency=10

# Com bind mount para arquivos
docker run -v $(pwd):/data gload:latest \
  --url=https://api.exemplo.com/upload \
  --method=POST \
  --requests=100 \
  --concurrency=5 \
  --body=@/data/arquivo.pdf
```

### Imagem otimizada
A imagem Docker utiliza `distroless/base-debian12` para:
- Tamanho mínimo (~20MB)
- Segurança aprimorada (sem shell, usuário não-root)
- Performance otimizada

## 🏗️ Arquitetura

O sistema é organizado em módulos especializados:

- **cmd/gload**: Ponto de entrada da aplicação
- **internal/cli**: Parsing e validação de argumentos CLI
- **internal/runner**: Orquestração dos testes de carga
- **internal/worker**: Workers para execução de requests HTTP
- **internal/httpc**: Cliente HTTP otimizado com pool de conexões
- **internal/report**: Geração de relatórios e métricas
- **internal/version**: Controle de versão

### Fluxo de execução
1. **CLI**: Parsing e validação dos parâmetros
2. **Runner**: Criação de workers e distribuição de jobs
3. **Workers**: Execução paralela de requests HTTP
4. **Report**: Coleta e análise dos resultados
5. **Output**: Renderização do relatório final

## ⚡ Performance

O sistema foi otimizado para alta performance:

- **Pool de conexões**: Reutilização de conexões TCP
- **Workers concorrentes**: Execução paralela de requests
- **Streaming de respostas**: Descartar bodies para economizar memória
- **Timeouts configuráveis**: Controle preciso de tempo limite
- **Métricas eficientes**: Cálculo otimizado de percentis

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Executar teste de exemplo
make run
```

## 📝 Makefile

Comandos disponíveis:

```bash
make build      # Compila o binário
make test       # Executa os testes
make run        # Executa exemplo
make docker-build # Constrói imagem Docker
```

## 🔧 Desenvolvimento

### Pré-requisitos
- Go 1.22+
- Docker (opcional)

### Estrutura do projeto
```
.
├── cmd/gload/          # Aplicação principal
├── internal/           # Módulos internos
│   ├── cli/           # CLI e parsing de flags
│   ├── runner/        # Orquestração dos testes
│   ├── worker/        # Workers HTTP
│   ├── httpc/         # Cliente HTTP otimizado
│   ├── report/        # Relatórios e métricas
│   └── version/       # Controle de versão
├── docs/              # Documentação adicional
├── Dockerfile         # Imagem Docker
└── Makefile          # Comandos de build
```

## 📄 Licença

Este projeto é parte de um desafio técnico e está disponível para fins educacionais.

## 🤝 Contribuição

Para contribuir com o projeto:

1. Fork o repositório
2. Crie uma branch para sua feature
3. Implemente as mudanças
4. Adicione testes
5. Submeta um pull request
