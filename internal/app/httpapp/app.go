package httpapp

import (
	"context"
	"de/internal/storage/sqlstorage"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Run(ctx context.Context, port uint16) error {
	store, err := sqlstorage.NewStore(ctx)
	if err != nil {
		return err
	}

	router := routes(store)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(), time.Second*4,
		)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("failed to gracefully shutdown http app")
		}

		if err := store.Close(ctx); err != nil {
			log.Printf("failed to gracefully shutdown db")
		}
	}()

	log.Printf("starting app bound to port %d...", port)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
