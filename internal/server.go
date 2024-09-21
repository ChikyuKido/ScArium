package internal

import (
	"ScArium/common/log"
	"ScArium/internal/backend/database"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Port    int
	Address string
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
}

func (s *Server) Start() {
	s.init()
	r := gin.Default()

	r.Run(fmt.Sprintf("%s:%d", s.Address, s.Port))
}
