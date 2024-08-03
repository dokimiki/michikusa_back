package main

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"michikusa_back/types"
	"net/http"
	"strconv"
	"time"
)

// 現在地の緯度経度から最寄りの駅を取得する関数
// 200m, 400m, 800m, 1600m, 3200m の範囲で検索を行い、見つかればその駅を返す
// 複数見つかった場合はランダムで1つ選ぶ
// 見つからなかった場合は空の types.Station を返すため、呼び出し元でエラーハンドリングを行う必要がある
func getNearestStation(longitude float64, latitude float64, odptAPIKey string) (types.OdptStation, error) {
	lon := strconv.FormatFloat(longitude, 'f', -1, 64)
	lat := strconv.FormatFloat(latitude, 'f', -1, 64)
	radiuses := []int{200, 400, 800, 1600, 3200}
	client := &http.Client{}

	for _, radius := range radiuses {
		url := "https://api.odpt.org/api/v4/places/odpt:Station?lon=" + lon + "&lat=" + lat + "&radius=" + strconv.Itoa(radius) + "&acl:consumerKey=" + odptAPIKey

		resp, err := client.Get(url)
		if err != nil {
			return types.OdptStation{}, err
		}
		defer resp.Body.Close()

		body , err := io.ReadAll(resp.Body)
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
	return types.OdptStation{}, errors.New("no stations exit within searchable area")
}
