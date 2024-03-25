package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var banners = []Banner{{320, 240}, {320, 480}, {300, 600}, {300, 250}}
var ua = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.0) Gecko/20100101 Firefox/60.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
}
var campaignNames = []string{"Summer Sale 2024", "Back-to-School Deals", "Winter Clearance", "Holiday Specials",
	"New Year's Discounts", "Spring Collection Launch", "Tech Gadgets Showcase", "Fitness Challenge Promo",
	"Travel Destination Offers", "Food Festival Discounts"}
var os = []string{"Windows", "Mac OS X"}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func between(min, max int) int {
	return random.Intn(max-min+1) + min
}
func generateAdCampaigns() []AdCampaign {
	var n int = between(1, 50)
	var campaigns []AdCampaign

	for i := 0; i < n; i++ {
		campaigns = append(campaigns, newAdCampaign())
	}

	return campaigns
}
func newAdCampaign() AdCampaign {
	return AdCampaign{
		ID:        fmt.Sprintf("%d", between(1, 10)),
		Name:      campaignNames[between(0, len(campaignNames)-1)],
		StartDate: "",
		EndDate:   "",
		Budget:    0.50 + (0.80-0.50)*random.Float64(),
		Status:    "",
		Imp: Imp{
			ID:     "1",
			Banner: banners[between(1, len(banners)-1)],
			Ext: Ext{
				Intent: Intent{
					PlacementID: strconv.Itoa(between(1, 2)),
				},
			},
		},
		PlacementIDs: []int{between(1, 2)},
		IAB:          iab(),
		Targeting: Targeting{
			Geo: Geo{
				Lat:     between(1, 99),
				Lon:     between(1, 99),
				Country: "US",
				Region:  "CA",
				Metro:   "SF",
				City:    "San Francisco",
				Zip:     "94107",
			},
			DeviceType: DeviceType{
				DeviceType: between(1, 2),
				OS:         os[between(0, 1)],
			},
		},
	}
}
func iab() []string {

	ofCodes := [][]int{
		{},
		{1, 7},
		{1, 20},
		{1, 12},
		{0, 11},
		{1, 15},
		{0, 9},
		{0, 15},
		{0, 18},
		{0, 31},
	}

	iab := random.Intn(9) + 1
	codeRange := ofCodes[iab]
	code := codeRange[0] + random.Intn(codeRange[1]-codeRange[0]+1)

	return []string{fmt.Sprintf("IAB%d-%d", iab, code)}
}
func adExchange() string {
	var winning AdCampaign
	bidRequestJson := popFromCache(miniRedis, &readWriteMutex)
	var adCampaigns []AdCampaign = generateAdCampaigns()
	var bidRequest BidRequest
	jsonToBidRequest(&bidRequestJson, &bidRequest)
	filteredCampaigns := filterAdCampaigns(bidRequest, adCampaigns)
	var index int = findCampaignWithLargestBid(filteredCampaigns)
	if index > 0 {
		winning = adCampaigns[index]
	}
	var exchange = Exchange{
		BidRequest:      bidRequest,
		AdCampaigns:     adCampaigns,
		WinningCampaign: winning,
	}
	exchangeJson, _ := json.Marshal(exchange)
	json := string(exchangeJson)
	return json
}
func filterAdCampaigns(bid BidRequest, campaigns []AdCampaign) []AdCampaign {
	var filtered []AdCampaign

	for _, campaign := range campaigns {
		if !containsAny(campaign.IAB, bid.Bcat) {
			if bannersMatch(bid.Imp[0].Banner, campaign.Imp.Banner) {
				filtered = append(filtered, campaign)
			}
		}
	}

	return filtered
}
func findCampaignWithLargestBid(campaigns []AdCampaign) int {
	if len(campaigns) == 0 {
		return -1
	}
	var selected int = -1
	largestBidCampaign := &campaigns[0]
	for index, campaign := range campaigns {
		if campaign.Budget > largestBidCampaign.Budget {
			selected = index
		}
	}

	return selected
}
func containsAny(slice []string, targets []string) bool {
	for _, target := range targets {
		for _, item := range slice {
			if target == item {
				return true
			}
		}
	}
	return false
}

func bannersMatch(banner1 Banner, banner2 Banner) bool {
	return banner1.W == banner2.W && banner1.H == banner2.H
}
