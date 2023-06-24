package osexit_test

import (
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/pkg/osexit"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexit.Analyzer, "./...")
}
