package httpapp

import (
	"de/internal/storage/sqlstorage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(store *sqlstorage.Store) chi.Router {
	mux := chi.NewMux()

	mux.Get("/", handleIndexPage(store))
	mux.Route("/ui", func(r chi.Router) {
		r.Handle("/", http.RedirectHandler("/", http.StatusFound))
		r.Get("/isolation", handleIsolationPage(store))
		r.Get("/indices", handleIndexingPage(store))
	})
	mux.Get("/isolation", handleIsolation(store))
	mux.Post("/refresh", handleRefreshDB(store))
	mux.Post("/transfer", handleTransfer(store))

	return mux
}
