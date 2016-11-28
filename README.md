String Service
==============

### Sync vendor

``` bash
$ cd ./stringsvc
$ govendor sync
```

### Install:

``` bash
$ go install github.com/l-vitaly/stringsvc/cmd/stringsvc
$ stringsvc
```

### Build

``` bash
$ go build github.com/l-vitaly/stringsvc/cmd/stringsvc
$ ./stringsvc
```


### Arguments

| Argument      | Description                 | Default
|---------------|-----------------------------|-----------------------|
| configName    | Set config file name        | stringsvc-config
| host          | Listen host name            | 0.0.0.0
| port          | Listen port                 | 8082

### Test

``` bash
$ go test github.com/l-vitaly/stringsvc/stringsvc_test 
```
