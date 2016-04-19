package main

import (
	"flag"
	"fmt"
	"github.com/lafolle/flen"
	"os"
)

var (
	pkg                          string
	bucketSize                   int
	inclTests                    bool
	lenLowerLimit, lenUpperLimit int
)

// init sets clas and flag package.
func init() {
	flag.BoolVar(&inclTests, "t", false, "include tests files")
	flag.IntVar(&bucketSize, "bs", 5, "bucket size (natural number)")
	flag.IntVar(&lenLowerLimit, "l", 0, "min length (inclusive)")
	flag.IntVar(&lenUpperLimit, "u", flen.Sentinel, "max length (exclusive)")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: flen <pkg> [options]\n")
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	pkg = flag.Args()[0]

	if pkg == "" {
		flag.Usage()
		return
	}

	flenOptions := &flen.Options{
		IncludeTests: inclTests,
		BucketSize:   bucketSize,
	}
	flens, err := flen.GenerateFuncLens(pkg, flenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	if rangeAsked(lenLowerLimit, lenUpperLimit) {
		zeroLenFuncs := flens.GetZeroLenFuncs()
		if len(zeroLenFuncs) > 0 {
			fmt.Println("0 len funcs")
			zeroLenFuncs.Print()
		}

		extImplFuncs := flens.GetExternallyImplementedFuncs()
		if len(extImplFuncs) > 0 {
			fmt.Println("Externally implemented funcs")
			extImplFuncs.Print()
		}
		flens.DisplayHistogram()
	} else {
		rangeFlens := flens.Query(lenLowerLimit, lenUpperLimit)
		if len(rangeFlens) > 0 {
			fmt.Printf("Functions with length in range [%d, %d)\n", lenLowerLimit, lenUpperLimit)
			rangeFlens.Print()
		}
	}

	os.Exit(0)
}

func rangeAsked(ll, ul int) bool {
	return ll != 0 || ul != flen.Sentinel
}
