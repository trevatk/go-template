package port

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trevatk/go-template/internal/domain"
	"go.uber.org/zap"
)

// HttpServer
type HttpServer struct {
	log    *zap.SugaredLogger
	bundle *domain.Bundle
}

// NewHttpServer
func NewHttpServer(logger *zap.Logger, bundle *domain.Bundle) *HttpServer {
	return &HttpServer{log: logger.Named("http server").Sugar(), bundle: bundle}
}

// NewRouter
func NewRouter(httpServer *HttpServer) *chi.Mux {

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/person", func(r chi.Router) {
			r.Post("", httpServer.createPerson)
			r.Get("/{id}", httpServer.fetchPerson)
			r.Put("", httpServer.updatePerson)
			r.Delete("/{id}", httpServer.deletePerson)
		})
	})

	r.Get("/health", httpServer.health)

	return r
}

func (h *HttpServer) createPerson(w http.ResponseWriter, r *http.Request) {}

func (h *HttpServer) fetchPerson(w http.ResponseWriter, r *http.Request) {}

func (h *HttpServer) updatePerson(w http.ResponseWriter, r *http.Request) {}

func (h *HttpServer) deletePerson(w http.ResponseWriter, r *http.Request) {}

func (h *HttpServer) health(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.log.Errorf("error encoding health check response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
