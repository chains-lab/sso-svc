package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/cmd"
	"github.com/chains-lab/sso-svc/cmd/migrations"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/sirupsen/logrus"
)

func Run(args []string) bool {
	cfg, err := internal.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	log := logium.NewLogger(cfg.Server.Log.Level, cfg.Server.Log.Format)
	log.Info("Starting server...")

	var (
		service        = kingpin.New("sso-svc", "")
		runCmd         = service.Command("run", "run command")
		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c, err := service.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	var wg sync.WaitGroup

	switch c {
	case serviceCmd.FullCommand():
		cmd.StartServices(ctx, cfg, log, &wg)
	case migrateUpCmd.FullCommand():
		err = migrations.MigrateUp(cfg.Database.SQL.URL)
	case migrateDownCmd.FullCommand():
		err = migrations.MigrateDown(cfg.Database.SQL.URL)
	default:
		log.Errorf("unknown command %s", c)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Print("All services stopped")
	}

	return true
}
