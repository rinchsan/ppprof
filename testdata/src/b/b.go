package b // want "should import net/http/pprof"

import "fmt"

func main() int {
	fmt.Println("ppprof")
	return 1
}

type A struct{}

func (A) main() {
}
