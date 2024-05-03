package handlers

import "github.com/labstack/echo/v4"

type ResearchElasticAndDatabase interface {
	ResearchElsWithDbHandler(c echo.Context) error
}
