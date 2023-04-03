package service

import "log"

var list []int64

func Convert62(id int64) []int64 {
	symbols := []string{
		"_", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "A", "b", "B", "c", "C", "d", "D", "e", "E",
		"f", "F", "g", "G", "h", "H", "i", "j", "J", "k",
		"K", "L", "m", "M", "n", "N", "o", "p", "P", "q",
		"Q", "r", "R", "s", "S", "t", "T", "u", "U", "v",
		"V", "w", "W", "x", "X", "y", "Y", "z", "Z",
	}
	l := len(symbols)
	al := id % int64(l)
	list = append(list, al)
	next := id - int64(l)
	log.Println(al, next, l)
	if next >= 0 {
		Convert62(next + 1)
	}

	return list
}
