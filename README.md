# cfportmap-plugin
A CF cli plugin to map routes to non-standard HTTP ports

## Installation
1. git clone the repo to your desktop
1. In the repo, run `go build` to compile a binary
1. run `cf install-plugin <path-to-binary>`

## Usage

### Map a route to non-standard app port
```
cf map-route-port my-app example.com apphost 8888
```