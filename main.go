package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://michikusa-front.pages.dev/","http://localhost:5173"},
		AllowMethods: []string{http.MethodGet},
	}))

	// 「現在地」の情報を受け取って、最寄駅や行き先駅、施設の情報を返すAPI
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello, World!",
		})
	})
	
	// 「行き先駅」の情報を受け取って、付近の施設の情報を返すAPI
	e.GET("/locations-list", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Locations List",
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
