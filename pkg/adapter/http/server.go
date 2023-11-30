package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/ports/person"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Server struct {
	routes *gin.Engine
	logger *zap.SugaredLogger

	personService person.PersonService
	rateLimiter   *rate.Limiter
	options       *ServerOptions
}

type ServerOptions struct {
	Port         int
	RateLimiting RateLimiterOptions
}

type RateLimiterOptions struct {
	Enable     bool
	Rate       int
	TokenBurst int
}

func NewServer(personService person.PersonService, logger *zap.SugaredLogger, options *ServerOptions) *Server {
	return &Server{
		personService: personService,
		routes:        gin.Default(),
		logger:        logger,
		options:       options,
	}
}

func (s *Server) StartServer() {
	if s.options.RateLimiting.Enable {
		s.logger.Infoln("Using rate limit")
		s.rateLimiter = rate.NewLimiter(rate.Limit(s.options.RateLimiting.Rate), s.options.RateLimiting.TokenBurst)
		s.useRateLimitMiddleware()
	}

	s.configurePersonRoutes()
	http.ListenAndServe(fmt.Sprintf(":%d", s.options.Port), s.routes)
}

// Simple token bucket rate limit
func (s *Server) useRateLimitMiddleware() {
	s.routes.Use(func(ctx *gin.Context) {
		s.logger.Infof("%d tokens available\n", int(s.rateLimiter.Tokens()))
		if s.rateLimiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(429)
		}
	})
}
