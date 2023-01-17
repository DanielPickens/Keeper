package http

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/danielpickens/keeper/pkg/api"
	"github.com/sirupsen/logrus"
)

// Handler actually handles the http requests.
// It use a router to map uri to HandlerFunc
type Handler struct {
	api        api.Api
	configPath string

	engine *gin.Engine
}

// NewHandler creates a Handler using defined routes.
// It takes a client parameter as an argument in order to pass to the handler and be accessible to the HandlerFunc
// Typically in a CRUD API, the client manages it's own connections to a storage system.
func NewHandler(api api.Api, configPath string, corsEnable bool) *Handler {
	v := &Handler{
		api:        api,
		configPath: configPath,
	}

	v.engine = gin.New()
	v.engine.Use(jsonLogMiddleware(), gin.Recovery())

	if corsEnable == true {
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AddAllowHeaders("authorization")
		v.engine.Use(cors.New(config))
		logrus.Info("CORS are enabled")
	}

	v.engine.GET("/ready", v.HealthCheck)
	v.engine.GET("/alive", v.HealthCheck)
	v.engine.POST("/inventories", v.Create)
	v.engine.GET("/inventories/:namespace", v.Get)
	v.engine.GET("/inventories/:namespace/status", v.GetStatus)
	v.engine.POST("/inventories/:namespace/reset", v.Reset)
	v.engine.GET("/inventories/:namespace/services", v.ListServices)
	v.engine.GET("/inventories", v.List)
	//v.engine.GET("/inventories/status", v.GetStatuses)
	v.engine.GET("/defaults", v.GetDefaults)
	v.engine.PUT("/inventories/:namespace", v.Update)
	v.engine.DELETE("/inventories/:namespace", v.Delete)
	v.engine.DELETE("/resources/:namespace/jobs/:resource", v.DeleteResource)
	v.engine.GET("/version", v.Version)

	return v
}

// Engine returns the defined router for the Handler
func (v *Handler) Engine() *gin.Engine { return v.engine }

// Server represents a http server that handles api requests to server
type Server struct {
	handler *Handler
}

// NewServer returns a http server with a given handler 
func NewServer(v *Handler) *Server {
	return &Server{
		handler: v,
	}
}

// Serve launch the webserver
func (s *Server) Serve(port int) {
	s.handler.Engine().Run(fmt.Sprintf(":%d", port))
}
