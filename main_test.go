package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestGetRedirectUrlShouldReturnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	_, err := getRedirectUrl(server.URL)

	if err == nil || err.Error() != "404 Not Found" {
		t.Errorf("Expected \"test error\" error, got %v", err)
	}
}

func TestGetRedirectUrlShouldReturnRedirectUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://icecast.omroep.nl/radio1-bb-mp3", http.StatusFound)
	}))
	defer server.Close()

	redirectUrl, err := getRedirectUrl(server.URL)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if redirectUrl != "https://icecast.omroep.nl/radio1-bb-mp3" {
		t.Errorf("Expected \"https://icecast.omroep.nl/radio1-bb-mp3\", got %v", redirectUrl)
	}
}

// Table Test
func TestGetStation(t *testing.T) {
	tests := []struct {
		spec           string
		expectedResult *station
		expectedError  error
		given          string
	}{
		{
			spec:           "should split correctly",
			expectedResult: &station{name: "Test", url: "http://testurl.com/test?g=test"},
			expectedError:  nil,
			given:          "station=\"http://testurl.com/test?g=test|Test\"",
		},
		{
			spec:           "should fail on missing station parts in given string",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid line: http://testurl.com/test?g=test"),
			given:          "http://testurl.com/test?g=test",
		},
		{
			spec:           "should fail on missing delimiter in given station string",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid line: station=\"http://testurl.com/test?g=test test\""),
			given:          "station=\"http://testurl.com/test?g=test test\"",
		},
		{
			spec:           "should fail on missing quotation marks in given station string",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid line: station=http://testurl.com/test?g=test|test"),
			given:          "station=http://testurl.com/test?g=test|test",
		},
		{
			spec:           "should fail on missing url in given station string",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid line: station=\"|Test\""),
			given:          "station=\"|Test\"",
		},
		{
			spec:           "should fail on missing station name in given station string",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid line: station=\"http://testurl.com/test?g=test|\""),
			given:          "station=\"http://testurl.com/test?g=test|\"",
		},
	}

	for _, test := range tests {
		result, err := getStation(test.given)

		if !((test.expectedResult == nil || reflect.DeepEqual(*test.expectedResult, *result)) &&
			(test.expectedError == nil || test.expectedError.Error() == err.Error())) {
			t.Errorf("%s: %s", test.spec, err)
		}
		fmt.Printf("passed: %s\n", test.spec)
	}
}

func TestGetSourceStations(t *testing.T) {
	fsMap := fstest.MapFS{
		"stations.txt": &fstest.MapFile{
			Data: []byte(`station="http://fritz.de/livemp3_s|Fritz"
			station="http://www.rockantenne.de/livemp3_s|Rockantenne"
			station="RadioEins"`),
		},
	}

	open, _ := fsMap.Open("stations.txt")

	stations := getSourceStations(open)

	if len(stations) != 2 {
		t.Errorf("Expected 3 stations, got %v", len(stations))
	}

	if stations["Fritz"] == nil {
		t.Errorf("Expected \"Fritz\" station")
	}

	if stations["Fritz"].url != "http://fritz.de/livemp3_s" {
		t.Errorf("Expected \"http://fritz.de/livemp3_s\", got %v", stations["Fritz"].url)
	}

	if stations["Rockantenne"] == nil {
		t.Errorf("Expected \"Rockantenne\" sration")
	}

	if stations["Rockantenne"].url != "http://www.rockantenne.de/livemp3_s" {
		t.Errorf("Expected \"http://www.rockantenne.de/livemp3_s\", got %v", stations["Rockantenne"].url)
	}

	if stations["RadioEins"] != nil {
		t.Errorf("Didn't expect \"RadioEins\" station")
	}
}

func TestReplaceUrlsByRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "https://icecast.omroep.nl/radio1-bb-mp3", http.StatusFound)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	stations := map[string]*station{
		"Redirect": &station{name: "Redirect", url: server.URL + "/redirect"},
		"NotFound": &station{name: "NotFound", url: server.URL},
	}

	replaceUrlsByRedirects(&stations)

	if stations["Redirect"].url != "https://icecast.omroep.nl/radio1-bb-mp3" {
		t.Errorf("Expected \"https://icecast.omroep.nl/radio1-bb-mp3\", got %v", stations["Redirect"].url)
	}

	if stations["NotFound"].url != server.URL {
		t.Errorf("Expected \"%v\", got %v", server.URL, stations["NotFound"].url)
	}
}
