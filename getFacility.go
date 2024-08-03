package main

import (
	"encoding/json"
	"fmt"
	"io"
	"michikusa_back/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetFacility(i types.LocationsRequestData, yahooApiKey string) ([]types.YDFFeature, error) {
	u, _ := url.Parse("https://map.yahooapis.jp/search/local/V1/localSearch")
	q := u.Query()
	q.Set("appid", yahooApiKey)
	q.Set("lat", strconv.FormatFloat(i.Latitude, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(i.Longitude, 'f', -1, 64))
	q.Set("output", "json")
	q.Set("sort", "hybrid")
	q.Set("results", "100")
	q.Set("dist", "1.25")
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api response is not ok(%d)", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var facilities types.YDF
	if err := json.Unmarshal(body, &facilities); err != nil {
		return nil, err
	}

	// グルメ系のジャンルを除外
	// あまりいい感じになっていないので調整が必要
	genreCodes := []string{
		"0204", // デパート・百貨店・ショッピングセンター
		"0208002", // 書店
		"0209", // ファッション
		"0301001", // スポーツ施設
		"0301007", // ボウリング場
		"0303", // 遊園地・テーマパーク等
		"0305", // 映画館・美術館・博物館等
		"0418", // 銭湯・浴場
		"0424", // 寺社・寺院
	}
	var filteredFacilities []types.YDFFeature
	for _, feature := range facilities.Feature {
		if feature.Property.Genre == nil || len(feature.Property.Genre) == 0 {
			continue
		}
		CODE:
		for _, code := range genreCodes {
			for _, genre := range feature.Property.Genre {
				if strings.HasPrefix(genre.Code, code) {
					filteredFacilities = append(filteredFacilities, feature)
					break CODE
				}
			}
		}
		if len(filteredFacilities) >= 5 {
			break
		}
	}

	// なにも見つからなかった場合は5件返す
	if len(filteredFacilities) == 0 {
		return facilities.Feature[:5], nil
	}

	return filteredFacilities, nil
}
