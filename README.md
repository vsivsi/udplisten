# udplisten

Simple Golang UDP Server which writes received packets to a file.

## Build
* `go build main.go -o udplisten`

## Run
* `--host`: default host is `127.0.0.1`
* `--port`: default port is 1234
* `--file`: If set, the received data will be appended to this file, otherwise it is written to stdout
* `--buffer`: default is 1500
