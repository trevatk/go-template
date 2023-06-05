package port

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"go.uber.org/zap"

	"github.com/trevatk/go-template/internal/domain"
)

// HttpServer exposed endpoints
type HttpServer struct {
	log    *zap.SugaredLogger
	bundle *domain.Bundle
}

// NewHttpServer create new http server instance
func NewHttpServer(logger *zap.Logger, bundle *domain.Bundle) *HttpServer {
	return &HttpServer{log: logger.Named("http server").Sugar(), bundle: bundle}
}

// NewRouter chi router implementation of http handler
func NewRouter(httpServer *HttpServer) *chi.Mux {

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/person", func(r chi.Router) {
			r.Post("/", httpServer.createPerson)
			r.Get("/{id}", httpServer.fetchPerson)
			r.Put("/", httpServer.updatePerson)
			r.Delete("/{id}", httpServer.deletePerson)
		})
	})

	r.Get("/health", httpServer.health)

	return r
}

func (h *HttpServer) createPerson(w http.ResponseWriter, r *http.Request) {

	request := &domain.NewPersonRequest{}
	err := render.Bind(r, request)
	if err != nil {
		h.log.Errorf("failed to bind to request %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	person, err := h.bundle.PersonService.Create(r.Context(), request.NewPerson)
	if err != nil {
		h.log.Errorf("unable to create new person %v", err)
		http.Error(w, "unable to create new person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(person); err != nil {
		h.log.Errorf("failed to encode response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *HttpServer) fetchPerson(w http.ResponseWriter, r *http.Request) {

	sid := chi.URLParam(r, "id")
	id, err := parseParamInt64(sid)
	if err != nil {
		h.log.Errorf("failed to parse param into integer %v", err)

		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "invalid url parameter", http.StatusBadRequest)
		return
	}

	person, err := h.bundle.PersonService.Read(r.Context(), id)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "person id does not exist", http.StatusNotFound)
			return
		}

		h.log.Errorf("unable to read person %v", err)
		http.Error(w, "failed to read person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(person); err != nil {
		h.log.Errorf("failed to encode response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *HttpServer) updatePerson(w http.ResponseWriter, r *http.Request) {

	request := &domain.UpdatePersonRequest{}
	err := render.Bind(r, request)
	if err != nil {
		h.log.Errorf("failed to bind request to model %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	person, err := h.bundle.PersonService.Update(r.Context(), request.UpdatePerson)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		h.log.Errorf("failed to update user %v", err)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(person); err != nil {
		h.log.Errorf("failed to encode response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *HttpServer) deletePerson(w http.ResponseWriter, r *http.Request) {

	sid := chi.URLParam(r, "id")
	id, err := parseParamInt64(sid)
	if err != nil {
		h.log.Errorf("failed to parse param %v", err)
		http.Error(w, "invalid request parameter", http.StatusBadRequest)
		return
	}

	err = h.bundle.PersonService.Delete(r.Context(), id)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "person does not exist", http.StatusNotFound)
			return
		}

		h.log.Errorf("failed to delete person %v", err)
		http.Error(w, "failed to delete person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode("SUCCESS"); err != nil {
		h.log.Errorf("unable to encode response %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

func (h *HttpServer) health(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.log.Errorf("error encoding health check response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func parseParamInt64(input string) (int64, error) {

	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse string to integer %v", err)
	}

	return value, nil
}
