# Instruções para uso da lib de logger

Ao escrever código Go neste projeto, siga rigorosamente as regras abaixo para uso do logger.

## Importação

```go
import "github.com/LeoRBlume/go-libs/logger"
```

## Regras obrigatórias

### 1. `context.Context` é sempre obrigatório

Nunca chame uma função de log sem um contexto real. Jamais use `context.Background()` ou `context.TODO()` em código de produção — o contexto deve vir do caller (handler HTTP, consumer de fila, etc.).

```go
// ERRADO
logger.Info(context.Background(), "src", "msg")

// CORRETO
func (s *UserService) Create(ctx context.Context, ...) {
    logger.Info(ctx, "UserService.Create", "usuário criado")
}
```

### 2. `src` identifica o call site

O parâmetro `src` deve seguir o padrão `"NomeDoTipo.NomeDoMetodo"`. Nunca deixe vazio ou genérico.

```go
// ERRADO
logger.Info(ctx, "", "msg")
logger.Info(ctx, "aqui", "msg")

// CORRETO
logger.Info(ctx, "OrderService.Checkout", "pedido finalizado")
logger.Info(ctx, "PaymentGateway.Charge", "cobrança enviada")
```

### 3. Use o nível correto

| Nível | Quando usar |
|-------|-------------|
| `Debug` | Detalhes internos úteis apenas em desenvolvimento |
| `Info` | Eventos de negócio relevantes (criação, conclusão) |
| `Warn` | Situações inesperadas que não interrompem o fluxo |
| `Error` | Falhas que precisam de atenção, com `err` preenchido |
| `Fatal` | Falhas irrecuperáveis na inicialização — chama `os.Exit(1)` |

```go
logger.Debug(ctx, "src", "entrando no handler")
logger.Info(ctx, "src", "pedido criado")
logger.Warn(ctx, "src", "retry tentativa 2")
logger.Error(ctx, "src", "falha ao salvar no banco", err)
logger.Fatal(ctx, "src", "não foi possível conectar ao banco", err)
```

### 4. `Fatal` encerra o processo

`Fatal` e `Fatalf` chamam `os.Exit(1)` após logar. Use somente durante a inicialização da aplicação (conexão com banco, leitura de config, etc.). Nunca use dentro de handlers HTTP ou goroutines de processamento.

### 5. Prefira a variante sem formatação quando não há interpolação

```go
// ERRADO — formatação desnecessária
logger.Infof(ctx, "src", "usuário criado")

// CORRETO
logger.Info(ctx, "src", "usuário criado")

// Variante formatada apenas quando há valores dinâmicos
logger.Infof(ctx, "src", "usuário %s criado com id %d", name, id)
```

### 6. Inicialização global — uma única vez no `main`

```go
func main() {
    logger.Setup(logger.Config{
        ServiceName: "nome-do-servico",
        Level:       logger.LevelInfo,
    })
    // ...
}
```

Nunca chame `Setup` fora do `main`. Em workers ou libs, receba o logger como dependência via `*logger.Logger`.

### 7. Middleware HTTP — sempre na raiz do roteador

```go
mux := http.NewServeMux()
// registre as rotas...
http.ListenAndServe(":8080", logger.TraceMiddleware(mux))
```

O middleware injeta `trace_id` no contexto automaticamente. Não injete `trace_id` manualmente em handlers que passam pelo middleware.

### 8. Injete `user_id` após autenticação

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := extractUserID(r)
        ctx := logger.WithUserID(r.Context(), userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 9. Testes — sempre use `NewNop()`

Em testes unitários, injete `logger.NewNop()` no lugar do logger real. Nunca deixe logs reais aparecerem na saída de testes.

```go
func TestUserService_Create(t *testing.T) {
    svc := NewUserService(logger.NewNop(), repo)
    // ...
}
```

### 10. Prefira funções globais — injete instância apenas quando necessário

O padrão recomendado é configurar o logger uma vez no `main` e usar as funções globais em qualquer lugar da aplicação. Não é necessário passar `*logger.Logger` para cada struct.

```go
// main.go — configura uma vez
logger.Setup(logger.Config{
    ServiceName: "meu-servico",
    Level:       logger.LevelInfo,
})

// qualquer serviço — usa diretamente, sem campo na struct
func (s *UserService) Create(ctx context.Context) {
    logger.Info(ctx, "UserService.Create", "criado")
}
```

Use injeção de dependência (`logger.New`) apenas quando precisar de loggers com configurações diferentes por componente (ex: nível diferente por módulo). Nesse caso, injete `*logger.Logger` como campo da struct e use `logger.NewNop()` nos testes.

## O que NÃO fazer

- Não use `fmt.Println`, `fmt.Printf` ou `log.Println` — use sempre o logger
- Não logue dados sensíveis (senhas, tokens, CPF, cartão)
- Não passe `nil` como contexto
- Não crie um `*logger.Logger` manualmente com `&logger.Logger{}` — use `logger.NewNop()` para testes ou `logger.Setup` + funções globais para produção
