package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// NewRouter creates and configures the Gin router with all API endpoints.
// It sets up public endpoints (no auth required), client-authenticated endpoints,
// and admin endpoints that require API key authentication.
func NewRouter(cnt *Controller) *gin.Engine {
	r := gin.Default()

	// add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"*"},
		ExposeHeaders:   []string{"*"},
	}))

	base := r.Group(fmt.Sprintf("/%s", cnt.config.Version))

	// public endpoints that do not require authentication
	public := base.Group("/public")

	// admin endpoints that require authentication
	admin := base.Group("/admin")
	admin.Use(AuthMiddleware(cnt.db))

	public.GET("/version", func(c *gin.Context) {
		response := cnt.Version(c)
		response.Send(c)
	})

	public.GET("/health", func(c *gin.Context) {
		response := cnt.HealthCheck(c)
		response.Send(c)
	})

	public.GET("/simulations", func(c *gin.Context) {
		response := cnt.ListSimulations(c)
		response.Send(c)
	})

	public.GET("/simulations/:id/meta", func(c *gin.Context) {
		response := cnt.GetSimulationMeta(c)
		response.Send(c)
	})

	public.GET("/simulations/:id/binary/:cpu_architecture", func(c *gin.Context) {
		cnt.GetSimulationBinary(c)
	})

	public.GET("/client", func(c *gin.Context) {
		response := cnt.GetClient(c)
		response.Send(c)
	})

	public.POST("/client/init", func(c *gin.Context) {
		response := cnt.InitClient(c)
		response.Send(c)
	})

	// add client id middleware to all further public endpoints. client IDs
	// are effectively temporary user IDs that are stored in the browser
	// local storage and are used to identify the "user".
	public.Use(ClientIdMiddleware(cnt.db))

	public.POST("/simulations/run", func(c *gin.Context) {
		response := cnt.RunSimulation(c)
		response.Send(c)
	})

	public.GET("/simulations/run", func(c *gin.Context) {
		response := cnt.GetSimulationRun(c)
		response.Send(c)
	})

	public.GET("/simulations/output", func(c *gin.Context) {
		cnt.GetSimulationOutput(c)
	})

	// admin endpoints that require authentication
	admin.POST("/simulations", func(c *gin.Context) {
		response := cnt.CreateSimulation(c)
		response.Send(c)
	})

	admin.PUT("/simulations/:id/meta", func(c *gin.Context) {
		response := cnt.UpdateSimulationMeta(c)
		response.Send(c)
	})

	admin.PUT("/simulations/:id/binary/:cpu_architecture", func(c *gin.Context) {
		response := cnt.UpdateSimulationBinary(c)
		response.Send(c)
	})

	admin.DELETE("/simulations/:id", func(c *gin.Context) {
		response := cnt.DeleteSimulation(c)
		response.Send(c)
	})

	return r
}

// main is the entry point for the simulation catalogue API server.
// It loads configuration, establishes database connections, and starts the HTTP server.
func main() {
	config := LoadConfig()
	// set log level
	log.SetLevel(ParseLogLevel(config.LogLevel))

	// connect to database
	db, err := NewPostgresDB(config)
	if err != nil {
		log.WithError(err).Error("failed to connect to database")
		panic(err)
	}
	defer db.pool.Close()

	events := make(chan SimulationRunEvent, config.MaxConcurrentRuns)

	for i := 0; i < config.MaxConcurrentRuns; i++ {
		go func() {
			runner := NewSimulationRunner(db, config)
			if err := runner.Start(events); err != nil {
				log.WithError(err).Error("failed to start runner")
			}
		}()
	}

	eventBroker := NewLocalEventBroker(events)
	// create controller
	controller := NewController(config, db, eventBroker)
	// create router
	router := NewRouter(controller)

	log.WithFields(log.Fields{
		"version": config.Version,
		"port":    config.Port,
	}).Info("starting app")

	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.WithError(err).Error("failed to start server")
		panic(err)
	}
}
