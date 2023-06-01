// Package for converting a number to a string from one bit depth to another.
package shorter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Main structure.
type shorter struct {
	symbols []string
	bit     int64
}

// New create universal.
func New(symbols []string) (*shorter, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols")
	}
	if len(symbols) < 2 {
		return nil, fmt.Errorf("at least two characters")
	}
	checkList := make([]string, 0, len(symbols))
	for _, s := range symbols {
		for _, c := range checkList {
			if s == c {
				return nil, fmt.Errorf("characters must be non-repetitive")
			}
		}
		checkList = append(checkList, s)
	}
	return &shorter{
		symbols: symbols,
		bit:     int64(len(symbols)),
	}, nil
}

// Create converter int to string59.
//
/* List symbols:
"$", "1", "2", "3", "4", "5", "6", "7", "8", "9",
"a", "A", "b", "B", "c", "C", "d", "D", "e", "E",
"f", "F", "g", "G", "h", "H", "i", "j", "J", "k",
"K", "L", "m", "M", "n", "N", "o", "p", "P", "q",
"Q", "r", "R", "s", "S", "t", "T", "u", "U", "v",
"V", "w", "W", "x", "X", "y", "Y", "z", "Z",.
*/
func NewShorter59() *shorter {
	var (
		symbols = []string{
			"$", "1", "2", "3", "4", "5", "6", "7", "8", "9",
			"a", "A", "b", "B", "c", "C", "d", "D", "e", "E",
			"f", "F", "g", "G", "h", "H", "i", "j", "J", "k",
			"K", "L", "m", "M", "n", "N", "o", "p", "P", "q",
			"Q", "r", "R", "s", "S", "t", "T", "u", "U", "v",
			"V", "w", "W", "x", "X", "y", "Y", "z", "Z",
		}
		bit = int64(len(symbols))
	)
	return &shorter{
		symbols: symbols,
		bit:     bit,
	}
}

// Function verssion 1.
// @Depricated.
func (s *shorter) ToStringV1(id int64) string {
	var list []int64
	result := []string{}
start:
	list = []int64{}
	if id < s.bit {
		return s.symbols[id]
	}
	for {
		n := id % s.bit
		list = append(list, n)
		if n == id {
			break
		}
		id -= s.bit
	}
	a := list[len(list)-1]
	result = append(result, s.symbols[a])

	id = int64(len(list[:len(list)-1]))
	if id >= s.bit {
		goto start
	}
	result = append(result, s.symbols[id])

	i := 0
	j := len(result) - 1
	for i < j {
		result[i], result[j] = result[j], result[i]
		i++
		j--
	}
	return strings.Join(result, "")
}

// New version.
func (s *shorter) Convert(str string) string {
	id, _ := strconv.ParseInt(str, 0, 64)
	result := []string{}
start:

	if id < s.bit {
		return s.symbols[id]
	}
	n := id % s.bit
	id /= s.bit
	result = append(result, s.symbols[n])
	if id >= s.bit {
		goto start
	}
	result = append(result, s.symbols[id])

	i := 0
	j := len(result) - 1
	for i < j {
		result[i], result[j] = result[j], result[i]
		i++
		j--
	}
	return strings.Join(result, "")
}

// Un Converter.
// sum m*b^n.
func (s *shorter) UnConvert(id string) string {
	if id == "" {
		return ""
	}
	list := strings.Split(id, "")
	listIndex := make([]int64, len(list))
	resultSum := make([]int64, 0)
	for i, v := range list {
		for n, z := range s.symbols {
			if z == v {
				listIndex[i] = int64(n)
				break
			}
		}
	}
	count := 0
start:
	l := listIndex[count+1:]
	a := listIndex[count]
	ex := len(l)
	resultSum = append(resultSum, a*int64(math.Pow(float64(s.bit), float64(ex))))
	if ex != 0 {
		count++
		goto start
	}
	var res int64
	for _, v := range resultSum {
		res += v
	}
	return fmt.Sprint(res)
}
