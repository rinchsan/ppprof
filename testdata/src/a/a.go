package a // want "should import net/http/pprof"

import "fmt"

func main() { // want "should set up pprof at the beginning of main"
	fmt.Println("ppprof")
	fmt.Println("ppprof")
	fmt.Println("ppprof")
	fmt.Println("ppprof")
}
