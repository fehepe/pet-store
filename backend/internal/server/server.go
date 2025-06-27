package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fehepe/pet-store/backend/internal/app"
	"github.com/fehepe/pet-store/backend/internal/auth"
	"github.com/fehepe/pet-store/backend/internal/graph"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP server
type Server struct {
	*http.Server
	deps *app.Dependencies
}

// New creates a new HTTP server with all routes configured
func New(deps *app.Dependencies) *Server {
	router := setupRouter(deps)

	return &Server{
		Server: &http.Server{
			Addr:         ":" + deps.Config.Port,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		deps: deps,
	}
}

// setupRouter configures all routes and middleware
func setupRouter(deps *app.Dependencies) chi.Router {
	router := chi.NewRouter()

	// Global middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(CORS())

	// Health check endpoint
	router.Get("/health", healthCheckHandler())

	// Static file server for uploads (keep for serving uploaded files)
	fileServer := http.FileServer(http.Dir(deps.Config.UploadDir))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	// GraphQL endpoints with conditional authentication
	router.Route("/graphql", func(r chi.Router) {
		r.Use(auth.ConditionalAuthMiddleware) // Enable authentication for mutations
		srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: deps.Resolver}))
		srv.AddTransport(transport.POST{})
		srv.AddTransport(transport.GET{})
		r.Handle("/", srv)
	})

	// GraphQL playground (only in development)
	if deps.Config.Env == "development" {
		router.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	}

	return router
}

// CORS returns a middleware that handles CORS headers with configurable origins
func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get allowed origins from environment variable, default to localhost for development
			allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
			if allowedOrigins == "" {
				allowedOrigins = "http://localhost:3000,http://localhost:3001"
			}

			origin := r.Header.Get("Origin")
			if origin != "" {
				origins := strings.Split(allowedOrigins, ",")
				for _, allowedOrigin := range origins {
					if strings.TrimSpace(allowedOrigin) == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Requested-With")
			w.Header().Set("Access-Control-Expose-Headers", "Link")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// healthCheckHandler returns a simple health check endpoint
func healthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "pet-store-backend",
			"version":   "1.0.0",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
