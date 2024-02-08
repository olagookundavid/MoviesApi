package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olagookundavid/movie-api/internal/jsonlog"
)

func (app *application) serve() error {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)
	go func() {
		// Intercept the signals, as before.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.logger.PrintInfo("shutting down server", map[string]string{"signal": s.String()})
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.logger.PrintInfo("completing background tasks", map[string]string{"addr": srv.Addr})
		// Call Wait() to block until our WaitGroup counter is zero --- essentially
		// blocking until the background goroutines have finished. Then we return nil on // the shutdownError channel, to indicate that the shutdown completed without // any issues.
		app.wg.Wait()
		shutdownError <- nil
	}()

	logger.PrintInfo("starting server", map[string]string{"addr": srv.Addr, "env": app.config.env})
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.PrintInfo("stopped server", map[string]string{"addr": srv.Addr})
	return nil

}
