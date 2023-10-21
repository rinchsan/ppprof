package c

import (
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

	_ = "ppprof"
}
