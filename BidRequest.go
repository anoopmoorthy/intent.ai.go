package main

type BidRequest struct {
	ID     string   `json:"id"`
	Imp    []Imp    `json:"imp"`
	Device Device   `json:"device"`
	User   User     `json:"user"`
	At     int      `json:"at"`
	Bcat   []string `json:"bcat"`
}

type Imp struct {
	ID     string `json:"id"`
	Banner Banner `json:"banner"`
	Ext    Ext    `json:"ext"`
}

type Ext struct {
	Intent Intent `json:"intent"`
}

type Intent struct {
	PlacementID string `json:"placementId"`
}

type Banner struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Device struct {
	UA         string `json:"ua"`
	Geo        Geo    `json:"geo"`
	IP         string `json:"ip"`
	DeviceType int    `json:"devicetype"`
	OS         string `json:"os"`
}

type Geo struct {
	Lat     int    `json:"lat"`
	Lon     int    `json:"lon"`
	Country string `json:"country"`
	Region  string `json:"region"`
	Metro   string `json:"metro"`
	City    string `json:"city"`
	Zip     string `json:"zip"`
}

type User struct {
	ID string `json:"id"`
}

type Exchange struct {
	BidRequest      BidRequest   `json:"bidRequest"`
	AdCampaigns     []AdCampaign `json:"adCampaigns"`
	WinningCampaign AdCampaign   `json:"winningCampaign"`
}
