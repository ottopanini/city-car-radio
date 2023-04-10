package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
