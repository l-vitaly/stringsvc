String Service
==============

### Install Vendor

Install `https://github.com/Masterminds/glide`

``` bash
$ cd $GOPATH/src/github.com/l-vitaly/stringsvc && glide install
```

### Install:

``` bash
$ cd $GOPATH/src/github.com/l-vitaly/stringsvc && glide install
$ go install github.com/l-vitaly/stringsvc
$ $GOPATH/bin/stringsvc
```

### Build

``` bash
$ go build github.com/l-vitaly/stringsvc/cmd/stringsvc
$ ./stringsvc
```

### Arguments

| Argument      | Description                   | Default
|---------------|-------------------------------|-----------------------|
| debug.addr    | Debug address<sup>*</sup>     | 0.0.0.0:62101 
| zipkin.addr   | Zipkin address<sup>*</sup>    | 0.0.0.0:9411
| consul.addr   | Consul address<sup>*</sup>    | 0.0.0.0:8500
| addr          | Service address<sup>*</sup>   | 0.0.0.0:62001

**address this is `host:port` string*

### Test

``` bash
$ go test github.com/l-vitaly/stringsvc/stringsvctest 
```
