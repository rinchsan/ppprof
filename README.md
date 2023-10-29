# ppprof

Set u**p** **pprof** by **ppprof**

## Run

```shell
go run github.com/rinchsan/ppprof/cmd/ppprof@latest -fix main.go
go run github.com/rinchsan/gosimports/cmd/gosimports@latest -w main.go
```

## Example

```go
package main

import "fmt"

func main() {
    fmt.Println("hello")
}
```

↓

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof"
    "runtime"
)

func main() {
    runtime.SetBlockProfileRate(1)
    runtime.SetMutexProfileFraction(1)
    go func() {
        log.Fatal(http.ListenAndServe("localhost:6060", nil))
    }()

    fmt.Println("hello")
}
```
