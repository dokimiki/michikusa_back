package main

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"michikusa_back/types"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// 現在地の緯度経度から最寄りの駅を取得する関数
// 200m, 400m, 800m, 1600m, 3200m の範囲で検索を行い、見つかればその駅を返す
// 複数見つかった場合はランダムで1つ選ぶ
func GetNearestStation(longitude float64, latitude float64, odptAPIKey string) (types.OdptStation, error) {
	baseURL := "https://api.odpt.org/api/v4/places/odpt:Station"
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("lon", strconv.FormatFloat(longitude, 'f', -1, 64))
	q.Set("lat", strconv.FormatFloat(latitude, 'f', -1, 64))
	q.Set("acl:consumerKey", odptAPIKey)
	radiuses := []int{200, 400, 800, 1600, 3200}

	for _, radius := range radiuses {
		q.Set("radius", strconv.Itoa(radius))
		u.RawQuery = q.Encode()
		req, _ := http.NewRequest("GET", u.String(), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return types.OdptStation{}, err
		}
		if resp.StatusCode != 200 {
			return types.OdptStation{}, errors.New("failed to get nearest station. status code: " + strconv.Itoa(resp.StatusCode))
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return types.OdptStation{}, err
		}

		var stations []types.OdptStation
		if err := json.Unmarshal(body, &stations); err != nil {
			return types.OdptStation{}, err
		}
		if len(stations) > 0 {
			return stations[rand.Intn(len(stations))], nil
		}
		// 外部APIに対する負荷を軽減するため100ms待機する
		time.Sleep(100 * time.Millisecond)
	}
	return types.OdptStation{}, errors.New("no stations exist within searchable area")
}
