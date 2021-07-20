# ISecL platform-info

This library is used to get the platform related information such as hardwareUUID, bios name, CPU flags etc.
cmd - contains the command line utility for the library.

## System Requirements
- RHEL 8.1
- Epel 8 Repo
- Proxy settings if applicable

## Software requirements
- git
- `go` version 1.14.4

# Step By Step Build Instructions

### Install `go` version 1.14.4
The `platform-info` requires Go version 1.14.4 that has support for `go modules`. The build was validated with the latest version 1.14.4 of `go`. It is recommended that you use 1.14.4 version of `go`. You can use the following to install `go`.
```shell
wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
tar -xzf go1.14.4.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

## Build platform-info

- Git clone the platform-info
- Run scripts to build the platform-info

```shell
git clone https://github.com/intel-secl/platform-info.git
cd platform-info
```
```shell
go build ./...
```

# Links
https://01.org/intel-secl/
