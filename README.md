# City Car Radio
Replaces custom radio stations in **City Car Driving** radio.ini which need redirects to get the actual stream url.

# Build

```bash
env GOOS=windows GOARCH=amd64 go build -o cityradio.exe
```

# Run

In the user folder with the radio.ini file must be another file named redirectradio.ini with urls to get the redirect 
urls. The name of the stations must match the name in the radio.ini file.  

Example (redirectradio.ini):  
```ini
station="http://fritz.de/livemp3_s|Fritz" 
station="http://www.rockantenne.de/livemp3_s|Rockantenne" # TLS
station="http://radiioeins.de/livemp3_s|RadioEins" # not reachable
```
Example (radio.ini):
```ini
station="something|Fritz"
station="http://www.rockantenne.de/livemp3_s|Rockantenne"
station="http://radiioeins.de/livemp3_s|RadioEins"
station="http://fritz.de/livemp3_s|Fritz2"
```
And then run the binary:
```bat
cityradio.exe
```
Example result (radio.ini):
```ini
station="http://d131.rndfnk.com/ard/rbb/fritz/live/mp3/48/stream.mp3?cid=01FC1WK1H1KDC5EBKW5Z191BW4&sid=2Nq7ZZ2uDctGS94qKfIEBRAhPCZ&token=bl5ZLHZIZ0jbwt4qIVJADvl5njKyn9UIwcFHhEGpTno&tvf=XFbh8S78URdkMTMxLnJuZGZuay5jb20|Fritz"
station="https://www.rockantenne.de/livemp3_s|Rockantenne"
station="http://radiioeins.de/livemp3_s|RadioEins"
station="http://fritz.de/livemp3_s|Fritz2"
```