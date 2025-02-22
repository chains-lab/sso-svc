package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/api"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/rabbit"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, svc *service.Service) {
	var (
	// signals indicate the finished initialization of each worker
	)

	//TODO INITIALIZATION OF INFRASTRUCTURE AND DOMAIN SERVICES

	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	run(func() { api.Run(ctx, svc) })

	run(func() { rabbit.Listener(ctx, svc) })
}
