package adapter

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/person"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Server struct {
	routes *gin.Engine
	logger *zap.SugaredLogger

	personService person.PersonService
	rateLimiter   *rate.Limiter
	options       map[string]string
}

func NewServer(personService person.PersonService, logger *zap.SugaredLogger, options map[string]string) *Server {
	return &Server{
		personService: personService,
		routes:        gin.Default(),
		logger:        logger,
		rateLimiter:   rate.NewLimiter(20, 10),
		options:       options,
	}
}

func (s *Server) StartServer() {
	s.useRateLimitMiddleware()
	s.configurePersonRoutes()

	if _, ok := s.options["port"]; !ok {
		s.options["port"] = "8080"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", s.options["port"]), s.routes)
}

// Simple token bucket rate limit
func (s *Server) useRateLimitMiddleware() {
	s.routes.Use(func(ctx *gin.Context) {
		if s.rateLimiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(429)
		}
	})
}
