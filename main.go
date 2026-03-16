package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/manuelarte/testcomments/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
