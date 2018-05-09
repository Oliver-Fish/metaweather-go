package metaweather

import (
	"testing"
	"time"
)

func TestGetLocation(t *testing.T) {
	testQueries := []struct {
		testName string
		query    string
	}{
		{"Single City", "Amsterdam"},
		{"Multi-Return", "A"},
		{"Invalid", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
		{"Lat Long Query", "36.96,-122.02"},
	}

	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := GetLocation(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if tt.testName != "Invalid" {
				if len(dat) == 0 {
					t.Errorf("No values returned for %q - %q", tt.testName, tt.query)
				}
				for _, v := range dat {
					if v.Title == "" {
						t.Error("No value for field title, JSON didn't unmarshal correctly")
					}
				}
			} else {
				if len(dat) != 0 {
					t.Errorf("values returned for %q - %q", tt.testName, tt.query) //This shouldn't have any values
				}
			}

		})
	}
}

func TestGetWeather(t *testing.T) {
	testQueries := []struct {
		testName string
		query    string
	}{
		{"London", "44418"},
		{"New York", "2459115"},
		{"Amsterdam", "727232"},
		{"Invalid", "00000"},
	}

	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := GetWeather(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if tt.testName != "Invalid" {
				if tt.testName != dat.Title {
					t.Errorf("invalid value expected %q got %q", tt.testName, dat.Title)
				}
			} else {
				if dat.Title != "" {
					t.Errorf("values returned for %q - %q", tt.testName, tt.query) //This shouldn't have any values
				}
			}
		})
	}
}

func TestGetWeatherDate(t *testing.T) {
	testQueries := []struct {
		testName string
		query    string
	}{
		{"London", "44418"},
		{"New York", "2459115"},
		{"Amsterdam", "727232"},
	}
	curTime := time.Now()
	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := GetWeatherDate(tt.query, curTime)
			if err != nil {
				t.Fatal(err)
			}
			if len(dat) == 0 {
				t.Error("No Values returned from api")
			}
			for _, v := range dat {
				if v.WeatherStateName == "" {
					t.Error("No value for field WeatherStateName, JSON didn't unmarshal correctly")
				}
			}
		})
	}
}
