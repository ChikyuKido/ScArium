package internal

import (
	"ScArium/common/log"
	"ScArium/internal/backend/database"
	"ScArium/internal/backend/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Port    int
	Address string
	Engine  *gin.Engine
}

func NewServer(port int, address string) *Server {
	return &Server{
		Port:    port,
		Address: address,
	}
}

func (s *Server) init() {
	log.I.Info("Initializing database")
	database.InitDB()
	log.I.Info("Initializing routes")
	routes.InitRoutes(s.Engine)
}

func (s *Server) Start() {
	r := gin.Default()
	r.HandleMethodNotAllowed = true
	s.Engine = r
	s.init()

	r.Run(fmt.Sprintf("%s:%d", s.Address, s.Port))
}
