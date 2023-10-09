package adapter

import (
	"fmt"
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
	s.routes.GET("/pessoas", s.GetPersonBySearchTerm)
	s.routes.GET("/contagem-pessoas", s.GetPeopleCount)
}

func (s *Server) CreatePerson(ctx *gin.Context) {
	var person person.Person

	if err := ctx.ShouldBindJSON(&person); err != nil {
		ctx.AbortWithError(422, err)
		return
	}

	if person, err := s.personService.Insert(person); err == nil {
		ctx.Header("Location", fmt.Sprintf("/pessoas/%s", person.ID))
		ctx.JSON(201, person)
	} else {
		ctx.JSON(500, err)
	}

}

func (s *Server) GetPersonBySearchTerm(ctx *gin.Context) {
	term, exists := ctx.GetQuery("t")
	if !exists {
		ctx.Status(400)
		return
	}

	personList, err := s.personService.FindByTerm(term)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, personList)
}

func (s *Server) GetAllPerson(ctx *gin.Context) {
	persons, err := s.personService.FindAll()

	if err != nil {
		ctx.JSON(500, err)
		return
	}

	if persons != nil {
		ctx.JSON(200, persons)
		return
	}

	ctx.Status(404)
}

func (s *Server) GetPersonByID(ctx *gin.Context) {
	person := s.personService.FindByID(ctx.Param("id"))
	if person == nil {
		ctx.Status(404)
		return
	}
	ctx.JSON(200, person)
}

func (s *Server) GetPeopleCount(ctx *gin.Context) {
	count, err := s.personService.Count()
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, gin.H{"count": count})
}
