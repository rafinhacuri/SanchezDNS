package util

import (
	"strings"
	"time"

	"resty.dev/v3"
)

type ipAPIRes struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	City       string `json:"city"`
	RegionName string `json:"regionName"`
	Country    string `json:"country"`
}

func GeoLocation(ip string) string {
	ipAPIClient := resty.New().
		SetBaseURL("http://ip-api.com").
		SetTimeout(5*time.Second).
		SetRetryCount(1).
		SetRetryWaitTime(200*time.Millisecond).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "cbpf-sso/1.0")

	ip = strings.TrimSpace(ip)
	if ip == "" {
		return ""
	}

	var res ipAPIRes

	resp, err := ipAPIClient.R().
		SetQueryParams(map[string]string{
			"fields": "status,message,city,regionName,country",
			"lang":   "en",
		}).
		SetResult(&res).
		Get("/json/" + ip)

	if err != nil || resp == nil {
		return ""
	}

	if resp.StatusCode() != 200 {
		return ""
	}

	if res.Status != "success" {
		return ""
	}

	city := strings.TrimSpace(res.City)
	region := strings.TrimSpace(res.RegionName)
	country := strings.TrimSpace(res.Country)

	if city != "" && country != "" {
		return city + " - " + country
	}

	if region != "" && country != "" {
		return region + " - " + country
	}

	return ""
}
