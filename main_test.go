package main

import (
	"fmt"
	"net/http"
	netUrl "net/url"
	"testing"
)

type MockHhttpClient struct{}

func (c *MockHhttpClient) Get(url string) (*http.Response, error) {
	if url == "test_success" {
		return &http.Response{
			Request: &http.Request{
				URL: &netUrl.URL{
					Scheme: "https",
					Host:   "icecast.omroep.nl",
					Path:   "/radio1-bb-mp3",
				},
			},
		}, nil
	} else {
		return nil, fmt.Errorf("test error")
	}
}
func TestGetRedirectUrlShouldReturnError(t *testing.T) {
	Client = &MockHhttpClient{}
	_, err := getRedirectUrl("https://www.nporadio1.nl/live")

	if err == nil || err.Error() != "test error" {
		t.Errorf("Expected \"test error\" error, got %v", err)
	}
}

func TestGetRedirectUrlShouldReturnRedirectUrl(t *testing.T) {
	Client = &MockHhttpClient{}
	url, err := getRedirectUrl("test_success")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != "https://icecast.omroep.nl/radio1-bb-mp3" {
		t.Errorf("Expected \"https://icecast.omroep.nl/radio1-bb-mp3\", got %v", url)
	}
}
