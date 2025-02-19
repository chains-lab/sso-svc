package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/sso-oauth/internal/service/transport"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/handlers"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, app *handlers.App) {
	var (
	// signals indicate the finished initialization of each worker
	)

	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	run(func() { transport.Run(ctx, app) })
}
