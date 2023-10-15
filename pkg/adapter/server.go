package adapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/person"
	"go.uber.org/zap"
)

type Server struct {
	routes *gin.Engine
	logger *zap.SugaredLogger

	personService person.PersonService
}

func NewServer(personService person.PersonService, logger *zap.SugaredLogger) *Server {
	return &Server{
		personService: personService,
		routes:        gin.Default(),
		logger:        logger,
	}
}

func (s *Server) StartServer() {
	s.configurePersonRoutes()
	http.ListenAndServe(":9090", s.routes)
}
