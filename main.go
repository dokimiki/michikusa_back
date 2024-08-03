package main

import (
	"log"
	"math/rand"
	"michikusa_back/types"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
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
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Bad Request",
			})
		}

		log.Println("GET / with", req)
		nearestStation, err := GetNearestStation(req.Longitude, req.Latitude, odptAPIKey)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error (getNearestStation)",
			})
		}

		stationList, err := GetStationList(nearestStation, odptAPIKey)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error (getStationList)",
			})
		}

		railwayInfo, err := GetRailwayInfo(nearestStation.Railway, odptAPIKey)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error (getRailwayInfo)",
			})
		}
		// 乱数で行き先駅を選択
		destinationStation := stationList[rand.Intn(len(stationList))]

		// neaest->destinationまでの駅数
		var neaestStationIndex int
		var destinationStationIndex int
		for i, station := range railwayInfo.StationOrder {
			if station.Station == nearestStation.SameAs {
				neaestStationIndex = i
			}
			if station.Station == destinationStation.SameAs {
				destinationStationIndex = i
			}
		}
		var numberOfStations int
		if neaestStationIndex < destinationStationIndex {
			numberOfStations = destinationStationIndex - neaestStationIndex
		} else {
			numberOfStations = neaestStationIndex - destinationStationIndex
		}

		facilityList, err := GetFacility(types.LocationsRequestData{
			Latitude:  destinationStation.Lat,
			Longitude: destinationStation.Long,
		}, yahooAPIKey)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error (getFacility)",
			})
		}

		var facilities []types.Facility
		for _, facility := range facilityList {
			lat, _ := strconv.ParseFloat(strings.Split(facility.Geometry.Coordinates, ",")[1], 64)
			long, _ := strconv.ParseFloat(strings.Split(facility.Geometry.Coordinates, ",")[0], 64)
			var genre string
			if facility.Property.Genre == nil || len(facility.Property.Genre) == 0 {
				genre = "その他"
			} else {
				genre = facility.Property.Genre[0].Name
			}
			facilities = append(facilities, types.Facility{
				Name:      facility.Name,
				Distance:  int(geo.Distance(orb.Point{destinationStation.Lat, destinationStation.Long}, orb.Point{lat, long})),
				Genre:     genre,
				Latitude:  lat,
				Longitude: long,
				MapURL: "https://map.yahoo.co.jp/v2/place/" + facility.Gid,
			})
		}

		var res types.InitialResponseData
		res.NearestStation = types.Station{
			Name:      nearestStation.Title,
			Latitude:  nearestStation.Lat,
			Longitude: nearestStation.Long,
		}
		res.DestinationStation = types.Station{
			Name:      destinationStation.Title,
			Latitude:  destinationStation.Lat,
			Longitude: destinationStation.Long,
		}
		res.RailwayName = railwayInfo.Title
		res.RailwayColor = railwayInfo.Color
		res.Facilities = facilities
		res.NumerOfStations = numberOfStations
		return c.JSON(http.StatusOK, res)
	})

	// 「行き先駅」の情報を受け取って、付近の施設の情報を返すAPI
	e.GET("/location-list", func(c echo.Context) error {
		var req types.LocationsRequestData
		if err := c.Bind(&req); err != nil {
			// todo: fix エラーハンドリング
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Bad Request",
			})
		}

		log.Println("GET /location-list with", req)
		facilityList, err := GetFacility(req, yahooAPIKey)
		if err != nil {
			log.Println(err)
			if err.Error() == "no facilities found" {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": "No facilities found",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error (getFacility)",
			})
		}

		var facilities []types.Facility
		for _, facility := range facilityList {
			coords := strings.Split(facility.Geometry.Coordinates, ",")
			lat, _ := strconv.ParseFloat(coords[1], 64)
			long, _ := strconv.ParseFloat(coords[0], 64)
			facilities = append(facilities, types.Facility{
				Name:      facility.Name,
				Distance:  int(geo.Distance(orb.Point{req.Latitude, req.Longitude}, orb.Point{lat, long})),
				Genre:     facility.Property.Genre[0].Name,
				Latitude:  lat,
				Longitude: long,
				MapURL: "https://map.yahoo.co.jp/v2/place/" + facility.Gid,
			})
		}
		return c.JSON(http.StatusOK, types.LocationsResponseData{
			Facilities: facilities,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
