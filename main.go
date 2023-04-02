package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type station struct {
	url  string
	name string
}

func (s station) toString() string {
	return `station="` + s.url + "|" + s.name + `"`
}

func main() {
	stations := getSourceStations()

	// for each station
	// open url and parse header for redirect url
	for station := range stations {
		url, err := getRedirectUrl(stations[station].url)
		if err == nil {
			stations[station].url = url
		}
	}

	writeStations(stations)
}

func writeStations(stations map[string]*station) {
	fileNameRadio := "radio.ini"

	open, err := os.Open(fileNameRadio)
	if err != nil {
		fmt.Println("Error opening src file: ", err)
		return
	}

	scanner := bufio.NewScanner(open)
	scanner.Split(bufio.ScanLines)
	var lines string
	for i := 0; scanner.Scan(); i++ {
		station, err := getStation(scanner.Text())
		curRedirectStation := stations[station.name]

		if i != 0 {
			lines = lines + "\n"
		}

		if curRedirectStation != nil && err == nil {
			lines = lines + curRedirectStation.toString()
		} else {
			lines = lines + scanner.Text()
		}
	}

	err = os.WriteFile(fileNameRadio, []byte(lines), 0644)
	if err != nil {
		fmt.Println("Error writing file: ", err)
		return
	}
}

func getSourceStations() map[string]*station {
	// read url from file in format: station="url|name"...
	//fileNameRadio := "radio.ini"
	fileNameRedirectRadio := "redirectradio.ini"
	stations := make(map[string]*station, 0)

	// read file
	open, err := os.Open(fileNameRedirectRadio)
	if err != nil {
		fmt.Println("Error opening src file: ", err)
		return stations
	}

	scanner := bufio.NewScanner(open)
	scanner.Split(bufio.ScanLines)
	for i := 0; scanner.Scan(); i++ {
		station, err := getStation(scanner.Text())
		if err != nil {
			fmt.Printf("Error parsing line(%d): %s\n", i, err)
			continue
		}
		stations[station.name] = &station
	}

	return stations
}

func getStation(line string) (station, error) {
	re := regexp.MustCompile(`station="([^|]*)\|([^"]*)"`)
	split := re.FindStringSubmatch(line)

	if len(split) != 3 {
		return station{}, fmt.Errorf("invalid line: %s", line)
	} else {
		return station{url: split[1], name: split[2]}, nil
	}
}

var (
	Client httpClient
)

func init() {
	Client = &http.Client{}
}

type httpClient interface {
	Get(url string) (*http.Response, error)
}

func getRedirectUrl(url string) (string, error) {
	get, err := Client.Get(url)
	if err != nil {
		return "", err
	}

	return get.Request.URL.String(), nil
}
