package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/rest"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, srv *service.Service) {
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

	run(func() { rest.Run(ctx, srv) })
}
