package main

import (
	"calendly/app/routes"
	"calendly/cmd"
	"calendly/lib/db"
	"calendly/lib/db/migration"
	"calendly/lib/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

var config *cmd.Config
var router = httprouter.New()

func init() {
	config = cmd.ParseConfig()
	logger.Init(logger.DEBUG)
	db.Connect(config.DatabaseUrl, config.DbMaxIdleConnections, config.DbMaxOpenConnections, true)
	routes.Init(router)
}

func main() {
	defer db.Close()
	migration.Run(db.Get())

	startServer(router)
}

func startServer(handler http.Handler) {
	s := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", 3000), Handler: handler}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info(context.TODO(), "starting server", nil)
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Error(context.TODO(), "unable to start server", map[string]any{"error": err})
				os.Exit(1)
			}
		}
	}()

	<-done
	logger.Info(context.TODO(), "initiated server shutdown", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Error(context.TODO(), "error while shutting down server", map[string]any{"error": err})
		os.Exit(1)
	}
	logger.Info(context.TODO(), "server successfully shutdown", nil)
}
