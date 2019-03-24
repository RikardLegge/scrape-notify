# scrape-notify

A simple tool to send notifications when a website is update and does not support RSS feeds. 
The tool currently only has one integration, namely www.chalemersstudentbostader.se.

## Build environment
Tested using go 1.12.1 but should work with most releases since the tool ha no external dependencies.

### Run 
`go run main/main.go`

### Test
`go test ./...`

### Build docker container
`docker build . -t scrape-notify`
