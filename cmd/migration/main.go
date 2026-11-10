package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/atmosone/oee/internal/adapters/postgres"
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()
	wdb, err := postgres.Open("")
	if err != nil {
		fmt.Println("failed to open repository", err)
		return
	}
	defer wdb.Close()
	if err := postgres.Migrate(ctx, wdb); err != nil {
		fmt.Println("repository migration failed", err)
		return
	}
	fmt.Println("repository migration completed")
}
