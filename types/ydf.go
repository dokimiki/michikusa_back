package types

type YDFFeature struct {
		ID       string `json:"Id"`
		Gid      string `json:"Gid"`
		Name     string `json:"Name"`
		Geometry struct {
			Type        string `json:"Type"`
			Coordinates string `json:"Coordinates"`
		} `json:"Geometry"`
		Category    []string      `json:"Category"`
		Description string        `json:"Description"`
		Style       []interface{} `json:"Style"`
		Property    struct {
			UID        string `json:"Uid"`
			CassetteID string `json:"CassetteId"`
			Yomi       string `json:"Yomi"`
			Country    struct {
				Code string `json:"Code"`
				Name string `json:"Name"`
			} `json:"Country"`
			Address              string `json:"Address"`
			GovernmentCode       string `json:"GovernmentCode"`
			AddressMatchingLevel string `json:"AddressMatchingLevel"`
			LandmarkCode         string `json:"LandmarkCode"`
			Tel1                 string `json:"Tel1"`
			Genre                []struct {
				Code string `json:"Code"`
				Name string `json:"Name"`
			} `json:"Genre"`
			Station []struct {
				ID       string `json:"Id"`
				SubID    string `json:"SubId"`
				Name     string `json:"Name"`
				Railway  string `json:"Railway"`
				Exit     string `json:"Exit"`
				ExitID   string `json:"ExitId"`
				Distance string `json:"Distance"`
				Time     string `json:"Time"`
				Geometry struct {
					Type        string `json:"Type"`
					Coordinates string `json:"Coordinates"`
				} `json:"Geometry"`
			} `json:"Station"`
			SmartPhoneCouponFlag string `json:"SmartPhoneCouponFlag"`
			KeepCount            string `json:"KeepCount"`
		} `json:"Property"`
	}


type YDF struct {
	ResultInfo struct {
		Count       int     `json:"Count"`
		Total       int     `json:"Total"`
		Start       int     `json:"Start"`
		Status      int     `json:"Status"`
		Description string  `json:"Description"`
		Copyright   string  `json:"Copyright"`
		Latency     float64 `json:"Latency"`
	} `json:"ResultInfo"`
	Feature []YDFFeature `json:"Feature"`
}
