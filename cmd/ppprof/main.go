package main

import (
	"github.com/rinchsan/ppprof"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(ppprof.Analyzer) }
