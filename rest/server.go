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
	"strings"
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

type ApiRouter struct {
	mux    *http.ServeMux
	prefix string
}

func (ar *ApiRouter) prefixPattern(pattern string) string {
	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) == 2 {
		return parts[0] + " " + ar.prefix + parts[1]
	}
	return ar.prefix + pattern
}

func (ar *ApiRouter) Handle(pattern string, handler http.Handler) {
	ar.mux.Handle(ar.prefixPattern(pattern), handler)
}

func (ar *ApiRouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	ar.mux.HandleFunc(ar.prefixPattern(pattern), handler)
}

func (server *Server) Start() {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Cors,
		middleware.Preflight,
		middleware.Logger,
	)

	mux := http.NewServeMux()

	apiRouter := &ApiRouter{
		mux:    mux,
		prefix: "/api",
	}

	// Register Routes
	server.userHandler.RegisterRoutes(apiRouter)
	server.taskHandler.RegisterRoutes(apiRouter, manager)

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