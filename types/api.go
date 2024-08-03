package types

type Facility struct {
	Name string `json:"name"`
	Distance int `json:"distance"` // バックエンド側で計算した距離、単位はメートル
	Genre string `json:"genre"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Station struct {
	Name string `json:"name"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
// todo rename
type InitialRequestData struct {
	Latitude  float64 `json:"latitude" query:"latitude"`
	Longitude float64 `json:"longitude" query:"longitude"`
	Price	 int `json:"price,omitempty" query:"price"`
}

type InitialResponseData struct {
	NearestStation Station `json:"nearest_station"`
	DestinationStation Station `json:"destination_station"`
	RailwayName string `json:"railway_name"`
	RailwayColor string `json:"railway_color"`
	NumerOfStations int `json:"number_of_stations"`
	Facilities []Facility `json:"facilities"`
}

type LocationsRequestData struct {
	Latitude  float64 `json:"latitude" query:"latitude"`
	Longitude float64 `json:"longitude" query:"longitude"`
}

type LocationsResponseData struct {
	Facilities []Facility `json:"facilities"`
}
