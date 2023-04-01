# City Car Radio
Replaces custom radio stations in City Car Drivings radio.ini which need redirects to get the actual stream url.

# Build

```bash
env GOOS=windows GOARCH=amd64 go build -o cityradio.exe
```

# Run

In the user folder with the radio.ini file must be another file named redirectradio.ini with urls to get the redirect 
urls. The name of the stations must match the name in the radio.ini file.

```bat
cityradio.exe
```
