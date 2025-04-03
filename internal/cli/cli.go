package cli

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/hs-zavet/comtools/logkit"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/service"
)

func Run(args []string) bool {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logkit.SetupLogger(cfg.Server.Log.Level, cfg.Server.Log.Format)
	logger.Info("Starting server...")

	var (
		application    = kingpin.New("sso-oauth", "")
		runCmd         = application.Command("run", "run command")
		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = application.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	svc, err := service.NewService(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
		return false
	}

	var wg sync.WaitGroup

	cmd, err := application.Parse(args[1:])
	if err != nil {
		logger.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		runServices(ctx, &wg, svc)
	case migrateUpCmd.FullCommand():
		err = MigrateUp(ctx, *cfg)
	case migrateDownCmd.FullCommand():
		err = MigrateDown(ctx, *cfg)
	default:
		logger.Errorf("unknown command %s", cmd)
		return false
	}
	if err != nil {
		logger.WithError(err).Error("failed to exec cmd")
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
