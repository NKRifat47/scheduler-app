package rest

import (
	"fmt"
	"net/http"
	"os"
	"scheduler-app/config"
	"scheduler-app/rest/handlers/tasks"
	"scheduler-app/rest/handlers/user"
	"scheduler-app/rest/middleware"
	"strconv"
)

type Server struct {
	cnf *config.Config
	userHandler *user.Handler
	taskHandler *tasks.Handler
}

func NewServer(cnf *config.Config, userHandler *user.Handler, taskHandler *tasks.Handler) *Server {
	return &Server{
		cnf: cnf,
		userHandler: userHandler,
		taskHandler: taskHandler,
	}
}

func (server *Server) Start() {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Cors,
		middleware.Preflight,
		middleware.Logger,
	)

	mux := http.NewServeMux()

	// Register Routes
	server.userHandler.RegisterRoutes(mux)
	server.taskHandler.RegisterRoutes(mux, manager)

	// Static File Server
	_ = os.MkdirAll("./static", 0755)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileServer)

	wrappedMux := manager.WrapMux(mux)

	addr := ":" + strconv.Itoa(server.cnf.Port)

	fmt.Println("Server is running on", addr)

	err := http.ListenAndServe(addr, wrappedMux)
	if err != nil {
		fmt.Println("Error starting the server ", err)
		os.Exit(1)
	}
}