package cli

import (
	"context"
)

// TODO FIX NOT WORKING
func runMigration(ctx context.Context, direction string) error {
	//service, err := cifractx.GetValue[*service.Service](ctx, "TODO")
	//if err != nil {
	//	return fmt.Errorf("failed to get server from context: %v", err)
	//}
	//
	//cmd := exec.Command(
	//	"migrate",
	//	"-path", "internal/data/sql/repositories/migrations",
	//	"-database", service.Config.Database.SQL.URL,
	//	"-verbose", direction,
	//)
	//
	//cmd.Stdout = logrus.StandardLogger().Out
	//cmd.Stderr = logrus.StandardLogger().Out
	//
	//if err := cmd.Run(); err != nil {
	//	return fmt.Errorf("failed to run migration %s: %v", direction, err)
	//}

	return nil
}

func MigrateUp(ctx context.Context) error {
	return runMigration(ctx, "up")
}

func MigrateDown(ctx context.Context) error {
	return runMigration(ctx, "down")
}
