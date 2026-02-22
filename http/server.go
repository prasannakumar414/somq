package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	port   int
	logger zap.Logger
}

func NewServer(port int, logger zap.Logger) *Server {
	return &Server{
		port:   port,
		logger: logger,
	}
}

func (s *Server) Serve() error {
	router := chi.NewRouter()
	router.Get("/health", ToHttpHandler(healthHandler))
	s.logger.Info("Starting server", zap.Int("port", s.port))
	return http.ListenAndServe(":8090", router)
}

func ToHttpHandler(handler func(w http.ResponseWriter, r *http.Request) (error error, status int, body any)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, status, body := handler(w, r)
		if err != nil {
			status = http.StatusInternalServerError
			body = err.Error()
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(body)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) (error error, status int, body any) {
	return nil, http.StatusOK, "OK"
}
