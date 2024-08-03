package main

import (
	"encoding/json"
	"fmt"
	"io"
	"michikusa_back/types"
	"net/http"
	"net/url"
)

func GetStationInfo(odptRailway string, odptAPIKey string) ([]types.OdptStation, error) {
	u, _ := url.Parse("https://api.odpt.org/api/v4/odpt:Railway")
	q := u.Query()
	q.Set("owl:sameAs", odptRailway)
	q.Set("acl:consumerKey", odptAPIKey)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []types.OdptStation{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []types.OdptStation{}, fmt.Errorf("api response is not ok(%d)", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []types.OdptStation{}, err
	}

	fmt.Println(string(body))

	var stationInfo []types.OdptStation
	if err := json.Unmarshal(body, &stationInfo); err != nil {
		return []types.OdptStation{}, err
	}

	return stationInfo, nil
}
