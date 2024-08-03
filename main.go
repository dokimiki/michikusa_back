package main

import (
	"michikusa_back/types"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 環境変数の読み込み
	odptAPIKey := os.Getenv("ODPT_API_KEY")
	if odptAPIKey == "" {
		panic("Error loading ODPT_API_KEY")
	}
	yahooAPIKey := os.Getenv("YAHOO_API_KEY")
	if yahooAPIKey == "" {
		panic("Error loading YAHOO_API_KEY")
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://michikusa-front.pages.dev/", "http://localhost:5173"},
		AllowMethods: []string{http.MethodGet},
	}))

	// 「現在地」の情報を受け取って、最寄駅や行き先駅、施設の情報を返すAPI
	e.GET("/", func(c echo.Context) error {
		var req types.InitialRequestData
		if err := c.Bind(&req); err != nil {
			// todo: fix エラーハンドリング
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Bad Request",
			})
		}

		res := types.InitialResponseData{
			NearestStation: types.Station{
				Name:      "新宿",
				Latitude:  35.69237,
				Longitude: 139.70121,
			},
			DestinationStation: types.Station{
				Name:      "東京",
				Latitude:  35.6818,
				Longitude: 139.7647,
			},
			RailwayName:  "丸ノ内線",
			RailwayColor: "#F62E36",
			Facilities: []types.Facility{
				{
					Name:      "東京国際フォーラム",
					Distance:  500,
					Genre:     "コンベンションセンター",
					Latitude:  35.6784,
					Longitude: 139.7636,
				},
				{
					Name:      "皇居外苑",
					Distance:  1000,
					Genre:     "公園",
					Latitude:  35.6825,
					Longitude: 139.7521,
				},
				{
					Name:      "東京ステーションギャラリー",
					Distance:  0,
					Genre:     "美術館",
					Latitude:  35.6812,
					Longitude: 139.7671,
				},
			},
		}
		return c.JSON(http.StatusOK, res)
	})

	// 「行き先駅」の情報を受け取って、付近の施設の情報を返すAPI
	e.GET("/locations-list", func(c echo.Context) error {
		var req types.LocationsRequestData
		if err := c.Bind(&req); err != nil {
			// todo: fix エラーハンドリング
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Bad Request",
			})
		}

		res := types.LocationsResponseData{
			Facilities: []types.Facility{
				{
					Name:      "東京国際フォーラム",
					Distance:  500,
					Genre:     "コンベンションセンター",
					Latitude:  35.6784,
					Longitude: 139.7636,
				},
				{
					Name:      "皇居外苑",
					Distance:  1000,
					Genre:     "公園",
					Latitude:  35.6825,
					Longitude: 139.7521,
				},
				{
					Name:      "東京ステーションギャラリー",
					Distance:  0,
					Genre:     "美術館",
					Latitude:  35.6812,
					Longitude: 139.7671,
				},
			},
		}
		return c.JSON(http.StatusOK, res)
	})

	e.GET("/get-facility", GetFacility)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
