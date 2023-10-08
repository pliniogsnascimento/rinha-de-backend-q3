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
	s.configureRoutes()
	http.ListenAndServe(":9090", s.routes)
}

func (s *Server) configureRoutes() {
	s.routes.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "Healthy")
	})

	s.routes.POST("/pessoas", s.CreatePerson)
	s.routes.GET("/pessoas/:id", s.GetPersonByID)
	s.routes.GET("/pessoas", s.GetAllPerson)
	s.routes.GET("/contagem-pessoas", func(ctx *gin.Context) {})
}

func (s *Server) CreatePerson(ctx *gin.Context) {
	var person person.Person

	if err := ctx.ShouldBindJSON(&person); err != nil {
		ctx.JSON(422, err)
		return
	}

	if person, err := s.personService.Insert(person); err == nil {
		ctx.JSON(201, person)
	} else {
		ctx.JSON(500, err)
	}

}

func (s *Server) GetAllPerson(ctx *gin.Context) {
	persons := s.personService.FindAll()

	if persons != nil {
		ctx.JSON(200, persons)
		return
	}

	ctx.String(404, "")
}

func (s *Server) GetPersonByID(ctx *gin.Context) {
	person := s.personService.FindByID(ctx.Param("id"))
	if person == nil {
		ctx.String(404, "")
		return
	}
	ctx.JSON(200, person)
}
