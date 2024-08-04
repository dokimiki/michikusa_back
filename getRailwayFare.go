package main

import (
	"encoding/json"
	"fmt"
	"io"
	"michikusa_back/types"
	"net/http"
	"net/url"
	"strings"
)

func GetRailwayFare(nearestStation types.OdptStation, odptAPIKey string) ([]types.OdptRailwayFare, error) {
	baseURL := "https://api.odpt.org/api/v4/odpt:RailwayFare"
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("odpt:operator", nearestStation.Operator)
	q.Set("odpt:fromStation", nearestStation.SameAs)
	q.Set("acl:consumerKey", odptAPIKey)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []types.OdptRailwayFare{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []types.OdptRailwayFare{}, fmt.Errorf("api response is not ok(%d)", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []types.OdptRailwayFare{}, err
	}
	var railwayFares []types.OdptRailwayFare
	if err := json.Unmarshal(body, &railwayFares); err != nil {
		return []types.OdptRailwayFare{}, err
	}

	// 出発駅と到着駅の路線が同じもののみを返す
	var filteredRailwayFares []types.OdptRailwayFare
	for _, railwayFare := range railwayFares {
		if strings.HasPrefix(strings.Split(railwayFare.ToStation, ":")[1], strings.Split(nearestStation.Railway, ":")[1]) {
			filteredRailwayFares = append(filteredRailwayFares, railwayFare)
		}
	}

	return filteredRailwayFares, nil
}
