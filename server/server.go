package server

import (
	"fmt"
	"net/http"

	middleware "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	pm "github.com/sharadregoti/devops/internal/pluginmanager"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/server/handlers"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	sm     *pm.SessionManager
	config *model.Config
}

func New(conf *model.Config) (*Server, error) {
	sm, err := pm.NewSM(conf)
	if err != nil {
		return nil, err
	}

	return &Server{
		sm:     sm,
		config: conf,
	}, nil
}

type Data struct {
	TableName string
	Table     [][]string
}

func (s *Server) routes() http.Handler {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/v1/config").HandlerFunc(handlers.HandleConfig(s.config))
	router.Methods(http.MethodGet).Path("/v1/auth/{pluginName}").HandlerFunc(handlers.HandleAuth(s.sm))
	router.Methods(http.MethodGet).Path("/v1/connect/{pluginName}/{authId}/{contextId}").HandlerFunc(handlers.HandleInfo(s.sm))
	router.Methods(http.MethodPost).Path("/v1/events/{id}").HandlerFunc(handlers.HandleEvent(s.sm))
	router.Methods(http.MethodGet).Path("/v1/ws/{id}").HandlerFunc(handlers.HandleWebsocket(s.sm))
	router.Methods(http.MethodGet).Path("/v1/ws/action/{clientId}/{id}").HandlerFunc(handlers.HandleActionWebsocket(s.sm))

	return router
}

func (s *Server) Start() error {
	fmt.Println("Starting server on port:", s.config.Server.Address)
	logger.LogInfo("Hi there")
	return http.ListenAndServe(s.config.Server.Address, utils.CreateCorsObject().Handler(middleware.LoggingHandler(logger.GetFileWriter(), s.routes())))
}
