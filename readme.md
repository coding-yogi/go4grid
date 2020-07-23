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

* ### start
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

* ### terminate
    ```
    ./go4grid terminate --help
    cleans up selenium 4 grid hub and nodes

    Usage:
    go4grid terminate [flags]

    Flags:
    -h, --help               help for terminate
        --namespace string   kube namespace (default "default")
    ```

* ### status
    ```
    ./go4grid status --help
    gets the current state of selenium 4 grid

    Usage:
    go4grid status [flags]

    Flags:
    -h, --help               help for status
        --namespace string   kube namespace (default "default")
    ```

    Sample output:
    ```
    |--------------------------------|-----------|------|---------------------------|-----------------------------|
    |              NAME              | NAMESPACE | PODS |          CREATED          |            IMAGE            |
    |--------------------------------|-----------|------|---------------------------|-----------------------------|
    | go4grid-selenium4-hub          | default   | 1/1  | 2020-07-23T10:33:38+08:00 | selenium/hub:4.0.0          |
    | go4grid-selenium4-node-chrome  | default   | 1/1  | 2020-07-23T10:34:18+08:00 | selenium/node-chrome:4.0.0  |
    | go4grid-selenium4-node-firefox | default   | 1/1  | 2020-07-23T10:34:18+08:00 | selenium/node-firefox:4.0.0 |
    |--------------------------------|-----------|------|---------------------------|-----------------------------|
    ```

## Scaling grid
For scaling grid same `start` command can be used. Go4Grid will analyze the current state of Grid and will scale up or down the nodes as needed