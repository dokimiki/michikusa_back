package main

import (
	"encoding/json"
	"errors"
	"io"
	"michikusa_back/types"
	"net/http"
	"net/url"
	"strconv"
)

// 最寄駅を通る路線の駅一覧を取得する関数
// 最寄り駅はフィルタされる
func getStationList(neaestStation types.OdptStation, odptAPIKey string) ([]types.OdptStation, error) {
	baseURL := "https://api.odpt.org/api/v4/odpt:Station"
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("odpt:railway", neaestStation.Railway)
	q.Set("acl:consumerKey", odptAPIKey)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get station list. status code: " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()

	body , err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stations []types.OdptStation
	if err := json.Unmarshal(body, &stations); err != nil {
		return nil, err
	}
	// 最寄り駅をフィルタする
	for i, station := range stations {
		if station.ID == neaestStation.ID {
			stations = append(stations[:i], stations[i+1:]...)
			break
		}
	}
	return stations, nil
}
