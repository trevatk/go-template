package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/trevatk/go-template/internal/db"
	"github.com/trevatk/go-template/internal/domain"
	"github.com/trevatk/go-template/internal/logging"
	"github.com/trevatk/go-template/internal/port"
)

func main() {

	fxApp := fx.New(
		fx.Provide(logging.New),
		fx.Provide(db.NewSQLite),
		fx.Provide(domain.NewPersonService),
		fx.Provide(domain.NewBundle),
		fx.Provide(port.NewHttpServer),
		fx.Provide(port.NewRouter),
		fx.Invoke(registerHooks),
	)

	start, cancel := context.WithTimeout(context.TODO(), time.Second*15)
	defer cancel()

	if err := fxApp.Start(start); err != nil {
		log.Fatalf("error starting service %v", err)
	}

	<-fxApp.Done()

	stop, cancel := context.WithTimeout(context.TODO(), time.Second*15)
	defer cancel()

	if err := fxApp.Stop(stop); err != nil {
		log.Fatalf("error stopping service %v", err)
	}
}

func registerHooks(lc fx.Lifecycle, log *zap.Logger, handler http.Handler, sqlite *sql.DB) error {

	logger := log.Named("lifecycle").Sugar()

	port := os.Getenv("HTTP_SERVER_PORT")
	if port == "" {
		return errors.New("$HTTP_SERVER_PORT is unset")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {

				logger.Info("execute database migration")

				err := db.Migrate(sqlite)
				if err != nil {
					return fmt.Errorf("failed to execute database migration %v", err)
				}

				logger.Infof("start http server http://localhost:%s" + port)

				go func() {
					if err := srv.ListenAndServe(); err != nil {
						logger.Fatalf("failed to start http server %v", err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {

				var err error

				logger.Info("close database connection")

				err = sqlite.Close()
				if err != nil {
					logger.Errorf("failed to close database connection %v", err)
				}

				logger.Info("shutdown http server")

				err = srv.Close()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Errorf("failed to shutdown http server %v", err)
				}

				// redudant logging
				return err
			},
		},
	)

	return nil
}
