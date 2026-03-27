# logger

Logging estruturado em JSON construído sobre `log/slog` (stdlib Go 1.21+).

Todo log exige um `context.Context` — rastreabilidade é obrigatória. `trace_id` e `user_id` fluem automaticamente pelo contexto.

## Importar

```bash
go get github.com/LeoRBlume/go-libs/logger@v1.0.0
```

## Inicialização

```go
import "github.com/LeoRBlume/go-libs/logger"

func main() {
    logger.Setup(logger.Config{
        ServiceName: "meu-servico",
        Level:       logger.LevelInfo,
    })
}
```

Níveis disponíveis: `LevelDebug`, `LevelInfo`, `LevelWarn`, `LevelError`.

## Logging

Todas as funções recebem `ctx`, `src` (identificador do call site) e `msg`.

```go
logger.Info(ctx, "UserService.Create", "usuário criado")
logger.Warn(ctx, "UserService.Create", "tentativa duplicada")
logger.Error(ctx, "UserService.Create", "falha ao salvar", err)
logger.Fatal(ctx, "UserService.Create", "erro irrecuperável", err) // chama os.Exit(1)

// Variantes com formatação
logger.Infof(ctx, "UserService.Create", "usuário %s criado com id %d", name, id)
logger.Errorf(ctx, "UserService.Create", "falha para usuário %s", err, name)
```

Exemplo de saída JSON:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "ERROR",
  "service": "meu-servico",
  "src": "UserService.Create",
  "message": "falha ao salvar",
  "trace_id": "abc-123",
  "user_id": "usr-456",
  "error": "connection refused"
}
```

`trace_id` e `user_id` são omitidos quando ausentes no contexto.

## Injetando trace_id e user_id no contexto

```go
ctx = logger.WithTraceID(ctx, "abc-123")
ctx = logger.WithUserID(ctx, "usr-456")

logger.Info(ctx, "UserService.Create", "usuário criado") // inclui ambos no JSON
```

## Middleware HTTP (trace automático)

Lê o header `X-Correlation-ID` da requisição; gera um UUID v4 se ausente. Injeta o valor no contexto e devolve no header de resposta.

```go
import "net/http"

mux := http.NewServeMux()
http.ListenAndServe(":8080", logger.TraceMiddleware(mux))
```

## Logger de instância (injeção de dependência)

Use `logger.New()` para criar uma instância configurada e injetá-la nos serviços.

```go
// No main
log := logger.New(logger.Config{ServiceName: "meu-servico", Level: logger.LevelInfo})
svc := NewUserService(log)

// No serviço
type UserService struct {
    log *logger.Logger
}

func (s *UserService) Create(ctx context.Context) {
    s.log.Info(ctx, "UserService.Create", "usuário criado")
}
```

## Testes

Use `NewNop()` para silenciar o logger em testes unitários sem nenhuma configuração extra:

```go
func TestMeuServico(t *testing.T) {
    svc := NewMeuServico(logger.NewNop())
    // ...
}
```
