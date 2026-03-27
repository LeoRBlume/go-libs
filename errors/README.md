# errors

Traduz domain errors para HTTP status codes, desacoplando a camada de domínio da camada HTTP.

## Como funciona

O domínio define seus próprios sentinel errors. O `Translator` mapeia cada um para um status code usando `errors.Is` internamente, então erros wrapped são resolvidos corretamente. O handler HTTP chama `Translate(err)` e obtém o status code sem conhecer o domínio.

## Importar

```bash
go get github.com/LeoRBlume/go-libs/errors@v1.0.0
```

## Uso

```go
import (
    "errors"
    "net/http"

    apperrors "github.com/LeoRBlume/go-libs/errors"
)

// 1. Defina os sentinel errors no domínio
var (
    ErrNotFound  = errors.New("not found")
    ErrForbidden = errors.New("forbidden")
)

// 2. Configure o translator (uma vez, na inicialização)
var translator = apperrors.NewTranslator().
    Register(ErrNotFound, http.StatusNotFound).
    Register(ErrForbidden, http.StatusForbidden)

// 3. Use no handler HTTP
func GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := service.FindUser(r.Context(), id)
    if err != nil {
        w.WriteHeader(translator.Translate(err)) // 404, 403, ou 500
        return
    }
    // ...
}
```

Erros não mapeados e `nil` retornam `500 Internal Server Error`.

Erros wrapped com `fmt.Errorf("...: %w", ErrNotFound)` são resolvidos corretamente.
