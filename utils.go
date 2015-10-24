package reporter

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"time"
)

// dateForFilename is a simple helper function to return a Time from a filename
func dateForFilename(path string) (time.Time, error) {
	return time.Parse("2006-01-02-reporter-export.json", filepath.Base(path))
}

// googleTimezoneResponse is a struct to contain the response from Google with the timezone for the given latitude and longitude
type googleTimezoneResponse struct {
	DstOffset    int    `json:"dstOffset"`
	RawOffset    int    `json:"rawOffset"`
	Status       string `json:"status"`
	TimeZoneID   string `json:"timeZoneId"`
	TimeZoneName string `json:"timeZoneName"`
}

// getTimezoneForLocation returns the timezone identifier (i.e. America/Los_Angeles) for the given latitude/longitude
func getTimezoneForLocation(timestamp int64, lat, long float64) (string, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/timezone/json?location=%f,%f&timestamp=%d", lat, long, timestamp)

	var gResp googleTimezoneResponse

	request, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer request.Body.Close()

	err = json.NewDecoder(request.Body).Decode(&gResp)
	if err != nil {
		return "", err
	}

	return gResp.TimeZoneID, nil
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}

func roundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return round(f*shift) / shift
}
