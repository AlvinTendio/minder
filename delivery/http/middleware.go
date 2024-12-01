package http

import (
	"net/http"

	"github.com/AlvinTendio/minder/auth"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
)

const (
	cacheMaxAge = 86400
)

// CORS wraps http handler to allow cors with default options
func CORS(handler http.Handler) http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           cacheMaxAge,
	})

	return corsHandler.Handler(handler)
}

// Recover wraps http handler with panic recovery from downstream call
func Recover(handler http.Handler) http.Handler {
	return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)
}

// Auth wraps http handler to extract user info context from auth token header
func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, auth.WithUserInfoRequestContext(r))
	})
}

type Option func(http.Handler) http.Handler

// WithRecovery adds option for handling panic recovery from downstream call
func WithRecovery() Option {
	return Recover
}

// WithCompression adds option for http compression
func WithCompression() Option {
	return handlers.CompressHandler
}

// WithCORS adds option for cors handling
func WithCORS() Option {
	return CORS
}

// WithAuth adds option to extract user info context from auth token header
func WithAuth() Option {
	return Auth
}

// WithDefault adds option for compression, user info extraction and panic recovery
func WithDefault() Option {
	return func(h http.Handler) http.Handler {
		return handlers.CompressHandler(Auth(Recover(h)))
	}
}

// NewHandler returns http handler with added options
func NewHandler(handler http.Handler, options ...Option) http.Handler {
	h := handler
	for _, option := range options {
		h = option(h)
	}

	return h
}

// DefaultHandler returns http handler with default options
func DefaultHandler(handler http.Handler) http.Handler {
	return NewHandler(handler, WithDefault())
}
