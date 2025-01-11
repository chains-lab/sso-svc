package cli

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/sirupsen/logrus"
)

func runMigration(ctx context.Context, direction string) error {
	service, err := cifractx.GetValue[*config.Server](ctx, config.SERVER)
	if err != nil {
		return fmt.Errorf("failed to get server from context: %v", err)
	}

	cmd := exec.Command(
		"migrate",
		"-path", "internal/data/sql/repositories/migration",
		"-database", service.Config.Database.URL,
		"-verbose", direction,
	)

	cmd.Stdout = logrus.StandardLogger().Out
	cmd.Stderr = logrus.StandardLogger().Out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run migration %s: %v", direction, err)
	}

	return nil
}

func MigrateUp(ctx context.Context) error {
	return runMigration(ctx, "up")
}

func MigrateDown(ctx context.Context) error {
	return runMigration(ctx, "down")
}
