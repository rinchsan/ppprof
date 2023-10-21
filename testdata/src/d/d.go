package d

import (
	"fmt"
	_ "net/http/pprof"
	"runtime"
)

func main() { // want "should set up pprof at the beginning of main"
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	fmt.Println("ppprof")
	fmt.Println("ppprof")
}
