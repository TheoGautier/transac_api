package server

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	invoiceRouter "transac_api/server/invoice"
	transactionRouter "transac_api/server/transaction"
	usersRouter "transac_api/server/users"
	"transac_api/service/invoice"
	"transac_api/service/user"
)

type Server struct {
	db     *sql.DB
	engine *gin.Engine
	port   int
	logger *log.Logger
}

func MustMakeNewServer(port int, db *sql.DB, logger *log.Logger) *Server {
	engine := gin.New()
	server := &Server{db: db, engine: engine, port: port, logger: logger}
	server.addHealthRoute()
	userService := user.MustMakeService(db)
	usersRouter.MustMakeRouter(engine, userService).AddRoutes()
	invoiceService := invoice.MustMakeService(db)
	invoiceRouter.MustMakeRouter(engine, invoiceService, logger).AddRoutes()
	transactionRouter.MustMakeRouter(engine, db, invoiceService, userService, logger).AddRoutes()
	return server
}

// Start server
//
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (s Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

func (s Server) addHealthRoute() {
	s.engine.GET("/health", health)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
