package e

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Println("ppprof")
}
