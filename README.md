# Graceful

Graceful is a tiny but powerful package to ensures your application shuts down cleanly, allowing registered closers to complete their tasks before exiting.

## Force Shutdowns

Force shutdown are possible in case a graceful shutdown fails. If you send a kill signal (`Ctrl+C`) once, the application begins its shutdown process, allowing all registered closers to complete. If send a kill signal (`Ctrl+C`) again, the application will force quit immediately.

```go
package main

import (
    "context"
    "github.com/cloudproud/graceful"
)

func main() {
    ctx := graceful.NewContext(context.Background())

    go ctx.Closer(func() {
        // gracefully shutdown connection
    })

    var err error
    if err != nil {
        ctx.Shutdown() // initiates a graceful shutdown of the application due to a critical error
    }

    ctx.AwaitKillSignal()
}
```
