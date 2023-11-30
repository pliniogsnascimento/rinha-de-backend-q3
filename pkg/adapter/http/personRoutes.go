package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/ports/person"
)

func (s *Server) configurePersonRoutes() {
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

	if personList, err := s.personService.FindByTerm(term); err != nil {
		ctx.JSON(500, err)
	} else {
		ctx.JSON(200, personList)
	}
}

func (s *Server) GetAllPerson(ctx *gin.Context) {
	if persons, err := s.personService.FindAll(); err != nil {
		ctx.JSON(500, err)
	} else if persons != nil {
		ctx.JSON(200, persons)
	} else {
		ctx.Status(404)
	}
}

func (s *Server) GetPersonByID(ctx *gin.Context) {
	if person := s.personService.FindByID(ctx.Param("id")); person != nil {
		ctx.JSON(200, person)
	} else {
		ctx.Status(404)
	}
}

func (s *Server) GetPeopleCount(ctx *gin.Context) {
	if count, err := s.personService.Count(); err != nil {
		ctx.Status(500)
		return
	} else {
		ctx.JSON(200, gin.H{"count": count})
	}
}
