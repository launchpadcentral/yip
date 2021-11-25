# yip 

yip is a command line tool used for doing in place modifications to a yaml object.

Prerequisite: 
  * Install [go](https://golang.org/doc/install)

Installation: (for Mac/Darwin)
  1. `git clone git@github.com:launchpadcentral/yip.git`
  2. `cd yip`  
  3. `mv main.go yip.go`
  4. `export GOBIN=/usr/local/bin`
  5. `go install yip.go`
  
Verify via
  * `yip --version`
  
Repository Packages:
  * set proper GOOS, GOARCH and Release #
  * Linux: GOOS=linux GOARCH=amd64 go build -o /tmp/yip-1.0.0-linux-amd64
  * file saved to yip/tmp/yip-1.0.0-linux-amd64
