package server

import (
	"fmt"
	"net/http"
	"strconv"

	middleware "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	core "github.com/sharadregoti/devops"
	"github.com/sharadregoti/devops/server/handlers"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	sm *core.SessionManager
}

func New() (*Server, error) {
	sm, err := core.NewSM()
	if err != nil {
		return nil, err
	}

	return &Server{
		sm: sm,
	}, nil
}

type Data struct {
	TableName string
	Table     [][]string
}

func (s *Server) routes() http.Handler {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/v1/config").HandlerFunc(handlers.HandleConfig())
	router.Methods(http.MethodGet).Path("/v1/auth/{pluginName}").HandlerFunc(handlers.HandleAuth(s.sm))
	router.Methods(http.MethodGet).Path("/v1/info/{authId}/{contextId}").HandlerFunc(handlers.HandleInfo(s.sm))
	router.Methods(http.MethodPost).Path("/v1/events/{id}").HandlerFunc(handlers.HandleEvent(s.sm))
	router.Methods(http.MethodGet).Path("/v1/ws/{id}").HandlerFunc(handlers.HandleWebsocket(s.sm))
	router.Methods(http.MethodGet).Path("/v1/ws/action/{clientId}/{id}").HandlerFunc(handlers.HandleActionWebsocket(s.sm))

	return router
}

func (s *Server) Start(port int) error {
	fmt.Println("Starting server on port:", port)
	return http.ListenAndServe(":"+strconv.Itoa(port), utils.CreateCorsObject().Handler(middleware.LoggingHandler(logger.GetFileWriter(), s.routes())))
}
