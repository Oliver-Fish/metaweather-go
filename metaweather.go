package metaweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	endpointURL              = "https://www.metaweather.com"                   //Base API url
	locationQueryURL         = endpointURL + "/api/location/search/?query="    //Takes city name i.e London
	locationQueryLattLongURL = endpointURL + "/api/location/search/?lattlong=" //Takes Lat and Long
	weatherLocationURL       = endpointURL + "/api/location/"                  //Takes woeid that you get from Location queries
)

//LocationData stores the data returned from LocationQuery endpoint,
//this provides us with the Woeid that lets us pull weather data
type LocationData struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Woeid        int64  `json:"woeid"`
	LattLong     string `json:"latt_long"`
}

//WeatherData stores the data returned from the weatherLocationURL endpoint
//This Contains all the weather data as well as sources for that data
type WeatherData struct {
	ConsolidatedWeather []ConsolidatedWeather `json:"consolidated_weather"`
	Time                time.Time             `json:"time"`
	SunRise             time.Time             `json:"sun_rise"`
	SunSet              time.Time             `json:"sun_set"`
	TimezoneName        string                `json:"timezone_name"`
	Parent              Parent                `json:"parent"`
	Sources             []Source              `json:"sources"`
	Title               string                `json:"title"`
	LocationType        string                `json:"location_type"`
	Woeid               int64                 `json:"woeid"`
	LattLong            string                `json:"latt_long"`
	Timezone            string                `json:"timezone"`
}

//ConsolidatedWeather stores all of the weather data returned from WeatherLocationURL
//This struct is nested in Weather data but also directly used in GetWeatherDate
type ConsolidatedWeather struct {
	ID                   int64     `json:"id"`
	WeatherStateName     string    `json:"weather_state_name"`
	WeatherStateAbbr     string    `json:"weather_state_abbr"`
	WindDirectionCompass string    `json:"wind_direction_compass"`
	Created              time.Time `json:"created"`
	ApplicableDate       string    `json:"applicable_date"`
	MinTemp              float64   `json:"min_temp"`
	MaxTemp              float64   `json:"max_temp"`
	TheTemp              float64   `json:"the_temp"`
	WindSpeed            float64   `json:"wind_speed"`
	WindDirection        float64   `json:"wind_direction"`
	AirPressure          float64   `json:"air_pressure"`
	Humidity             int64     `json:"humidity"`
	Visibility           float64   `json:"visibility"`
	Predictability       int64     `json:"predictability"`
}

//Parent stores the data about our query to the WeatherData endpoint
//This struct is not directly used and nested directly into WeatherData
type Parent struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Woeid        int64  `json:"woeid"`
	LattLong     string `json:"latt_long"`
}

//Source contains the data of all weather sources used for our query to WeatherData endpoint
//This sturct is not directly used and nested directly into WeatherData
type Source struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	CrawlRate int64  `json:"crawl_rate"`
}

//GetLocation takes a location string either the name of a place or the long and Lat location of a place
//This can return multiple locations
func GetLocation(loc string) ([]LocationData, error) {
	var lDat []LocationData
	if strings.Contains(loc, ",") { //We Found a comma so this is a latlong
		err := getJSONData(locationQueryLattLongURL+loc, &lDat)
		if err != nil {
			return nil, err
		}
	} else { //No comma was found so this is a string of text instead (city,country etc)
		err := getJSONData(locationQueryURL+loc, &lDat)
		if err != nil {
			return nil, err
		}
	}
	return lDat, nil
}

//GetWeather takes a woeid string can get this from a LocationData type, this will return all weather data including the source of said data
func GetWeather(woeid string) (WeatherData, error) {
	var wDat WeatherData
	err := getJSONData(weatherLocationURL+woeid, &wDat)
	if err != nil {
		return nil, err
	}
	return wDat, nil
}

//GetWeatherDate takes a woeid string can get this from a LocationData type, and a time.Time for the date this will return only weather data and no source data
func GetWeatherDate(woeid string, date time.Time) ([]ConsolidatedWeather, error) { //Data string should be a go time/date object?
	var wDat []ConsolidatedWeather
	y := date.Year()
	m := int(date.Month())
	d := date.Day()
	dateString := strconv.Itoa(y) + "/" + strconv.Itoa(m) + "/" + strconv.Itoa(d)
	err := getJSONData(weatherLocationURL+woeid+"/"+dateString, &wDat)
	if err != nil {
		return nil, err
	}
	return wDat, nil
}

//getJSONData is a helper function to prevent me doing the below error checking each time I need to json unmarshal on API endpoints
func getJSONData(url string, out interface{}) error {
	fmt.Println(url)
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	dat, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &out)
	if err != nil {
		return err
	}
	return nil
}
