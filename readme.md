# Go4Grid
Golang based commadline tool to spin up Selenium 4 Grid on Kubernetes within seconds

## Why Go4Grid?
* Setup, Teardown, Scale up/down grid with single command
* No need of handling deployment files
* Most of tasks are asynchronous and run concurrently

## Building from source
* Clone repo
* Run `go get -u ./... && go build -u ./...`

## Using go4grid
Currently available commands

### start
 ```
./go4grid start --help        
start up selenium 4 grid hub and nodes

Usage:
  go4grid start [flags]

Flags:
      --chrome int32       number of chrome nodes (default 1)
      --firefox int32      number of firefox nodes (default 1)
  -h, --help               help for start
      --namespace string   kube namespace (default "default")
 ```

### terminate
 ```
 ./go4grid terminate --help
cleans up selenium 4 grid hub and nodes

Usage:
  go4grid terminate [flags]

Flags:
  -h, --help               help for terminate
      --namespace string   kube namespace (default "default")
 ```