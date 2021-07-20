# ISecL aas-manager

This library is used to populate all users and permissions in Authentication and Authorization Service.

## System Requirements
- RHEL 8.1
- Proxy settings if applicable

## Software requirements
- git
- `go` version `go1.14.4`

# Step By Step Build Instructions

### Install `go` version `go1.14.4`
The `aas-manager` was validated with the latest version `1.14.4` of `go`. It is recommended that you use 1.14.4 version of `go`. More recent versions may introduce compatibility issues. You can use the following to install `go`.

```shell
wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
tar -xzf go1.14.4.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

## Build aas-manager

- Git clone the aas-manager
- Run scripts to build the aas-manager

```shell
git clone https://github.com/intel-secl/intel-secl.git
cd tools/aas-manager
go build ./...
```

# Links
https://01.org/intel-secl/
