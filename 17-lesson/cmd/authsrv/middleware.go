package main

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// logMiddleware логирует все запросы в формате Apache Common Log Format
func (s *Service) logMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(s.logger, next)
}
