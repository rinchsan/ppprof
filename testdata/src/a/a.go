package a

func main() { // want "should set up pprof at the beginning of main"
	_ = "ppprof"
}

type A struct{}

func (A) main() {
}
