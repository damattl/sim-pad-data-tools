go build -o sim-pad-data-tools-m1 ./src
GOOS=windows GOARCH=amd64 go build -o sim-pad-data-tools.exe ./src
GOOS=darwin GOARCH=amd64 go build -o sim-pad-data-tools-darwin-amd64 ./src
