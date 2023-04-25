// Static analyzer.
// Usage:
//
//	./cmd/staticlint/osexit -osexit ./cmd/shortener/main.go
package main

import (
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/pkg/osexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"honnef.co/go/tools/staticcheck"
)

func getSA() []*analysis.Analyzer {
	list := make([]*analysis.Analyzer, 0)
	for _, a := range staticcheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "SA") {
			list = append(list, a.Analyzer)
		}
	}
	return list
}

func main() {
	analyzers := []*analysis.Analyzer{
		printf.Analyzer,
		shift.Analyzer,
		shadow.Analyzer,
		bools.Analyzer,
		assign.Analyzer,
		httpresponse.Analyzer,
		osexit.Analyzer,
	}
	analyzers = append(analyzers, getSA()...)

	multichecker.Main(analyzers...)
}
