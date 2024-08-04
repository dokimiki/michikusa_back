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
	"time"
)

const (
	facilityCount = 10
	maxCount      = 500
	maxDistance   = 2
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
	q.Set("dist", strconv.FormatFloat(maxDistance, 'f', -1, 64))

	var filteredFacilities []types.YDFFeature
	start := 0

	for {
		fmt.Println("start:", start)
		q.Set("start", strconv.Itoa(start))
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
			"0118",    // スイーツ
			"0204",    // デパート・百貨店・ショッピングセンター
			"0301001", // スポーツ施設
			"0301007", // ボウリング場
			"0303001",
			"0303002",
			"0303003",
			"0303004",
			"0303005",
			"0303006",
			"0303007",
			"0303008",
			"0303009",
			"0303010",
			"0303011",
			"0303012",
			"0305001", // 映画館・美術館・博物館等
			"0305002",
			"0305003",
			"0305007",
			"0418", // 銭湯・浴場
			"0424", // 寺社・寺院
		}
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
			if len(filteredFacilities) >= facilityCount {
				break
			}
		}

		if len(filteredFacilities) >= facilityCount || start >= maxCount {
			break
		}
		start += 100
		time.Sleep(100 * time.Millisecond)
	}

	return filteredFacilities, nil
}
