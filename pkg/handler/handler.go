package handler

import (
	"medodstest/pkg/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	mux.Route("/auth", func(mux chi.Router) {
		mux.Post("/sign-up", h.signUp)
		mux.Post("/sign-in", h.signIn)
	})

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(h.userIdentity)

		mux.Route("/notes", func(mux chi.Router) {
			mux.Get("/", h.getAllNote)
			mux.Post("/", h.saveNote)

			mux.Route("/{noteId}", func(mux chi.Router) {
				mux.Get("/", h.getNoteById)
				mux.Put("/", h.updateNote)
				mux.Patch("/", h.updateNote)
				mux.Delete("/", h.deleteNote)
			})
		})
	})

	return mux
}
