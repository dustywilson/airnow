package airnow

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"strings"
	"time"
)

// AirNow provides access to the AirNow AQI API
type AirNow string

// Point is a pair of Latitude and Longitude
type Point struct {
	Latitude  float64
	Longitude float64
}

// Category is ...
type Category struct {
	Num   int `json:"Number"`
	Name  string
	Color color.RGBA
}

// CategoryColor represents colors of AQI categories
var CategoryColor = []color.RGBA{
	color.RGBA{0, 228, 0, 0},   // Green, Good
	color.RGBA{255, 255, 0, 0}, // Yellow, Moderate
	color.RGBA{255, 126, 0, 0}, // Orange, Unhealthy for Sensitive Groups
	color.RGBA{255, 0, 0, 0},   // Red, Unhealthy
	color.RGBA{153, 0, 76, 0},  // Purple, Very Unhealthy
	color.RGBA{76, 0, 38, 0},   // Maroon, Hazardous
	color.RGBA{0, 0, 0, 0},     // Black, unknown
}

// Observation is a single set of observation details
type Observation struct {
	Time     time.Time `json:"DateObserved"`
	Area     string
	State    string
	LatLng   Point
	AQI      int
	Category Category
}

type rawObservation struct {
	DateObserved  string
	HourObserved  int
	LocalTimeZone string
	ReportingArea string
	StateCode     string
	Latitude      float64
	Longitude     float64
	ParameterName string
	AQI           int
	Category      Category
}

// New takes an API key and returns an AirNow
func New(key string) AirNow {
	return AirNow(key)
}

// Errors
var (
	ErrBadServerResponse = errors.New("received unexpected response from AirNow source")
)

// NowByZIP returns a current Observation based on the zipcode and radiusMiles
func (a AirNow) NowByZIP(zipcode string, radiusMiles int) (*Observation, error) {
	r, err := http.Get(fmt.Sprintf("http://www.airnowapi.org/aq/observation/zipCode/current/?format=application/json&API_KEY=%s&zipCode=%s&distance=%d", a, zipcode, radiusMiles))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var ro []rawObservation
	d := json.NewDecoder(r.Body)
	err = d.Decode(&ro)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character '<'") {
			// this happens if the user doesn't provide API key, but could be for other reasons such as rate limiting
			return nil, ErrBadServerResponse
		}
		return nil, err
	}

	obs := make([]*Observation, len(ro))
	for i, o := range ro {
		t, err := time.Parse("2006-01-02 15 MST", fmt.Sprintf("%s %d %s", strings.TrimSpace(o.DateObserved), o.HourObserved, o.LocalTimeZone))
		if err != nil {
			return nil, err
		}
		o.Category.Color = CategoryColor[o.Category.Num]
		obs[i] = &Observation{
			Time:     t,
			Area:     o.ReportingArea,
			State:    o.StateCode,
			LatLng:   Point{o.Latitude, o.Longitude},
			AQI:      o.AQI,
			Category: o.Category,
		}
	}

	return obs[0], nil
}
