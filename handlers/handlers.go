package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pocketbase/pocketbase"
)

// App struct will hold PocketBase instance and other services
type AppContext struct {
	DB *pocketbase.PocketBase
}

func (a *AppContext) SubmitGastromatic(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{
		"ok": true,
	})

}

func (a *AppContext) SubmitGastronovi(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{
		"ok": true,
	})

}
