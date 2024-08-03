package main

import (
	"encoding/json"
	"fmt"
	"io"
	"michikusa_back/types"
	"net/http"
	"net/url"
	"strconv"
)

func GetFacility(i types.LocationsRequestData, yahooApiKey string) (types.YDF, error) {
	u, _ := url.Parse("https://map.yahooapis.jp/search/local/V1/localSearch")
	q := u.Query()
	q.Set("appid", yahooApiKey)
	q.Set("lat", strconv.FormatFloat(i.Latitude, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(i.Longitude, 'f', -1, 64))
	q.Set("output", "json")
	q.Set("sort", "hybrid")
	q.Set("dist", "1")
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.YDF{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return types.YDF{}, fmt.Errorf("api response is not ok(%d)", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.YDF{}, err
	}

	var facilities types.YDF
	if err := json.Unmarshal(body, &facilities); err != nil {
		return types.YDF{}, err
	}

	return facilities, nil
}
