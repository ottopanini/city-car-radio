package main

import (
	"bufio"
	"fmt"
	"io/fs"
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
	stations := getSourceStations(openFile("redirectradio.ini"))

	replaceUrlsByRedirects(&stations)

	fileNameRadio := "radio.ini"
	content := generateRadioIniFileContent(openFile(fileNameRadio), stations)
	writeFile(fileNameRadio, content)
}

func replaceUrlsByRedirects(stations *map[string]*station) {
	// for each station
	// open url and parse for redirect url
	for station := range *stations {
		url, err := getRedirectUrl((*stations)[station].url)
		if err == nil {
			(*stations)[station].url = url
		}
	}
}

func openFile(fileName string) fs.File {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening src file: ", err)
		return nil
	}
	defer file.Close()

	return file
}

func writeFile(fileName string, content string) {
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing file: ", err)
		return
	}
}

func generateRadioIniFileContent(file fs.File, stations map[string]*station) string {
	scanner := bufio.NewScanner(file)
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

	return lines
}

func getSourceStations(file fs.File) map[string]*station {
	// read url from file in format: station="url|name"...
	//fileNameRadio := "radio.ini"
	stations := make(map[string]*station, 0)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for i := 0; scanner.Scan(); i++ {
		station, err := getStation(scanner.Text())
		if err != nil {
			fmt.Printf("Error parsing line(%d): %s\n", i, err)
			continue
		}
		stations[station.name] = station
	}

	return stations
}

func getStation(line string) (*station, error) {
	re := regexp.MustCompile(`station="([^|]*)\|([^"]*)"`)
	split := re.FindStringSubmatch(line)

	if len(split) != 3 || split[1] == "" || split[2] == "" {
		return nil, fmt.Errorf("invalid line: %s", line)
	} else {
		return &station{url: split[1], name: split[2]}, nil
	}
}

func getRedirectUrl(url string) (string, error) {
	get, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	if get.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", get.Status)
	}

	return get.Request.URL.String(), nil
}
