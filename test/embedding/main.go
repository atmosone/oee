package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/atmosone/oee/web"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGINT,
	)
	defer cancel()
	target, err := url.Parse("http://localhost:8880")
	if err != nil {
		return
	}
	router := http.NewServeMux()
	router.Handle(web.API(target))
	router.Handle(web.UI())
	server := &http.Server{
		Addr:    ":9990",
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
			return
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(context.TODO()); err != nil {
		fmt.Println(err)
	}
}
