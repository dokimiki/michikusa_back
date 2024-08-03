package types

// referece: 3.3.3. odpt:Railway
type OdptRailway struct {
	Context string `json:"@context"`
	ID      string `json:"@id"`
	Type    string `json:"@type"`
	Date   string `json:"dc:date"`
	SameAs string `json:"owl:sameAs"`
	Title  string `json:"dc:title"`
	RailwayTitle map[string]string `json:"odpt:railwayTitle,omitempty"`
	Kana string `json:"odpt:kana,omitempty"`
	LineCode string `json:"odpt:lineCode,omitempty"`
	Color string `json:"odpt:color,omitempty"`
	// Region
	AscendingRailDirection string `json:"odpt:ascendingRailDirection,omitempty"`
	DescendingRailDirection string `json:"odpt:descendingRailDirection,omitempty"`
	StationOrder []map[string]string `json:"odpt:stationOrder,omitempty"`
}

// referece: 3.3.4. odpt:RailwayFare
type OdptRailwayFare struct {
	Context string `json:"@context"`
	ID      string `json:"@id"`
	Type    string `json:"@type"`
	Date   string `json:"dc:date"`
	Issued string `json:"dct:issued,omitempty"`
	Valid string `json:"dct:valid,omitempty"`
	SameAs string `json:"owl:sameAs"`
	Operator string `json:"odpt:operator"`
	FromStation string `json:"odpt:fromStation"`
	ToStation string `json:"odpt:toStation"`
	TicketFare int `json:"odpt:ticketFare"`
	IcCardFare int `json:"odpt:icCardFare,omitempty"`
	ChildTicketFare int `json:"odpt:childTicketFare,omitempty"`
	ChildIcCardFare int `json:"odpt:childIcCardFare,omitempty"`
	ViaStation []string `json:"odpt:viaStation,omitempty"`
	ViaRailway []string `json:"odpt:viaRailway,omitempty"`
	TicketType string `json:"odpt:ticketType,omitempty"`
	PaymentMethod []string `json:"odpt:paymentMethod,omitempty"`
}

// referece: 3.3.5. odpt:Station
type OdptStation struct {
	Context string `json:"@context"`
	ID      string `json:"@id"`
	Type    string `json:"@type"`
	Date   string `json:"dc:date"`
	SameAs string `json:"owl:sameAs"`
	Title  string `json:"dc:title,omitempty"`
	StationTitle map[string]string `json:"odpt:stationTitle,omitempty"`
	Operator string `json:"odpt:operator"`
	Railway string `json:"odpt:railway"`
	StationCode string `json:"odpt:stationCode,omitempty"`
	Long float64 `json:"geo:long,omitempty"`
	Lat float64 `json:"geo:lat,omitempty"`
	// Region
	Exit []string `json:"odpt:exit,omitempty"`
	ConnectingRailway []string `json:"odpt:connectingRailway,omitempty"`
	ConnectionStation []string `json:"odpt:connectionStation,omitempty"`
	StationTimetable []string `json:"odpt:stationTimetable,omitempty"`
	PassengerSurvey []string `json:"odpt:passengerSurvey.omitempty"`
}
