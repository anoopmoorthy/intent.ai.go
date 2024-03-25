package main

type AdCampaign struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	StartDate    string     `json:"startDate"`
	EndDate      string     `json:"endDate"`
	Budget       float64    `json:"budget"`
	Status       string     `json:"status"`
	Imp          Imp        `json:"imp"`
	PlacementIDs []int      `json:"placementIds"`
	IAB          []string   `json:"iab"`
	Targeting    Targeting  `json:"targeting"`
	Creatives    []Creative `json:"creatives"`
}

type Targeting struct {
	Geo        Geo        `json:"geo"`
	DeviceType DeviceType `json:"deviceType"`
}

type DeviceType struct {
	DeviceType int    `json:"devicetype"`
	OS         string `json:"os"`
}

type Creative struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}
