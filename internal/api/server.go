package api

import (
	"com.github/asdsec/planny/configs"
	"com.github/asdsec/planny/internal/security"
	"com.github/asdsec/planny/internal/store"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// Server represents the server
type Server struct {
	conf   configs.Config
	router *echo.Echo
	token  security.TokenGenerator
	store  db.Store
}

// NewServer creates a new server
func NewServer(conf configs.Config, store db.Store) (*Server, error) {
	tokenGen, err := security.NewJWTGenerator(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token generator: %w", err)
	}

	serv := &Server{
		conf:  conf,
		token: tokenGen,
		store: store,
	}
	serv.setupRouter()
	return serv, nil
}

func (serv *Server) setupRouter() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	v1 := e.Group("/api/v1")
	v1.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	v1.POST("/login", serv.login)
	v1.POST("/register", serv.register)
	v1.POST("/renew_access", serv.renewAccess)

	authorized := v1.Group("")
	authorized.Use(serv.authMiddleware)
	authorized.POST("/plans", serv.createPlan)
	authorized.GET("/plans", serv.retrievePlans)
	authorized.PATCH("/plans/:id", serv.updatePlan)
	authorized.DELETE("/plans/:id", serv.deletePlan)

	serv.router = e
}

// Start starts the server
func (serv *Server) Start() error {
	if serv.conf.Environment == "development" {
		serv.router.Debug = true
	}
	return serv.router.Start(serv.conf.ServerAddress)
}

func (serv *Server) err(ctx echo.Context, code int, msg string) error {
	return ctx.JSON(code, echo.Map{"error": msg})
}
