package main

import (
	"fmt"
	"michikusa_back/types"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

var API_KEY string

func init() {
	API_KEY = os.Getenv("YAHOO_API_KEY")
}

func GetFacility(c echo.Context) error {
	var i types.LocationsRequestData
	if err := c.Bind(&i); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad Request",
		})
	}

	i.Latitude = 35.6907446
	i.Longitude = 139.6881503

	u, _ := url.Parse("https://map.yahooapis.jp/search/local/V1/localSearch")
	q := u.Query()
	q.Set("appid", API_KEY)
	q.Set("lat", strconv.FormatFloat(i.Latitude, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(i.Longitude, 'f', -1, 64))
	u.RawQuery = q.Encode()
	fmt.Println(API_KEY)
	fmt.Println(u.String())
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "failed to send request",
		})
	}
	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "api response is not ok(" + strconv.Itoa(resp.StatusCode) + ")",
		})
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}
