package router

import (
	"log/slog"
	"net/http"
	"net/http/pprof"
	"watcharis/go-migrate-lotto-history-els/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRouter(e *echo.Echo, lottoHistoryHandlers handlers.ResearchElasticAndDatabase) *echo.Echo {
	slog.Info("start init router")

	InitProfiling(e)

	e.GET("/health", healthCheck)
	e.Use(middleware.Recover())

	mainGroup := e.Group("/research-els-db")
	api := mainGroup.Group("/api")
	g := api.Group("/v1")
	g.POST("/select-task", lottoHistoryHandlers.ResearchElsWithDbHandler)

	return e
}

func healthCheck(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, echo.Map{"message": "Service is Running !!"}, "	")
}

func InitProfiling(e *echo.Echo) {
	// init route profiling CPU
	e.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	// init route profiling Memory
	e.GET("/debug/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
}
