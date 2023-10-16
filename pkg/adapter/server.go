package adapter

import (
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
}

func NewServer(personService person.PersonService, logger *zap.SugaredLogger) *Server {
	return &Server{
		personService: personService,
		routes:        gin.Default(),
		logger:        logger,
		rateLimiter:   rate.NewLimiter(20, 10),
	}
}

func (s *Server) StartServer() {
	s.useRateLimitMiddleware()
	s.configurePersonRoutes()
	http.ListenAndServe(":9090", s.routes)
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
