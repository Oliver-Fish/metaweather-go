package metaweather

import (
	"reflect"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	meta := New()
	tmeta := reflect.TypeOf(meta).String()
	if tmeta != "*metaweather.Client" {
		t.Errorf("Expected type *metaweather.Client got %q", tmeta)
	}
}

func TestBaseURL(t *testing.T) {
	tString := "localhost:8080"
	meta := New(BaseURL(tString))
	if meta.baseURL != tString {
		t.Errorf("Expected value %q got %q", tString, meta.baseURL)
	}
}

func TestGetJSONData(t *testing.T) {
	tString := "localhost:9999"
	meta := New(BaseURL(tString))
	t.Run("API Down", func(t *testing.T) {
		err := meta.getJSONData("/notfound", nil)
		if err == nil {
			t.Error("Expected error got nil ")
		}
	})
}

func TestGetLocation(t *testing.T) {
	testQueries := []struct {
		testName string
		query    string
	}{
		{"Single City", "Amsterdam"},
		{"Multi-Return", "A"},
		{"Invalid", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
	}
	meta := New()
	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := meta.GetLocation(tt.query)
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
	tString := "localhost:9999"
	meta = New(BaseURL(tString))
	t.Run("API Down", func(t *testing.T) {
		_, err := meta.GetLocation("")
		if err == nil {
			t.Error("Expected error got nil ")
		}
	})
}

func TestGetLocationLattLong(t *testing.T) {
	testQueries := []struct {
		testName string
		latt     string
		long     string
	}{
		{"Bristol", "51.453732", "-2.591560"},
	}
	meta := New()
	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := meta.GetLocationLattLong(tt.latt, tt.long)
			if err != nil {
				t.Fatal(err)
			}
			if len(dat) == 0 {
				t.Errorf("No values returned for %q - %q,%q", tt.testName, tt.latt, tt.long)
			}
			for _, v := range dat {
				if v.Title == "" {
					t.Error("No value for field title, JSON didn't unmarshal correctly")
				}
			}

		})
	}
	tString := "localhost:9999"
	meta = New(BaseURL(tString))
	t.Run("API Down", func(t *testing.T) {
		_, err := meta.GetLocationLattLong("", "")
		if err == nil {
			t.Error("Expected error got nil ")
		}
	})
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

	meta := New()
	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := meta.GetWeather(tt.query)
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
	tString := "localhost:9999"
	meta = New(BaseURL(tString))
	t.Run("API Down", func(t *testing.T) {
		_, err := meta.GetWeather("")
		if err == nil {
			t.Error("Expected error got nil ")
		}
	})
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
	meta := New()
	for _, tt := range testQueries {
		t.Run(tt.testName, func(t *testing.T) {
			dat, err := meta.GetWeatherDate(tt.query, curTime)
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
	tString := "localhost:9999"
	meta = New(BaseURL(tString))
	t.Run("API Down", func(t *testing.T) {
		_, err := meta.GetWeatherDate("", time.Now())
		if err == nil {
			t.Error("Expected error got nil")
		}
	})
}
