package shorter_test

import (
	"fmt"

	"github.com/Vasily-van-Zaam/ushortener/pkg/shorter"
)

func Example() {
	sh2, _ := shorter.New([]string{
		"0", "1",
	})
	out1 := sh2.Convert("123456789")
	fmt.Println(out1)

	out2 := sh2.UnConvert(out1)
	fmt.Println(out2)

	sh3, _ := shorter.New([]string{
		"A", "B", "C",
	})
	out3 := sh3.Convert("123456789")
	fmt.Println(out3)

	out4 := sh3.UnConvert(out3)
	fmt.Println(out4)

	// Output:
	// 111010110111100110100010101
	// 123456789
	// CCBCBACCACACBCCAA
	// 123456789
}
