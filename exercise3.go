package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(comma("1234.5678"))
}
func comma(s string) string { //refactored
	var buf bytes.Buffer
	start := len(s)
	f := strings.Index(s, ".")
	if f > 0 {
		start = f
	}
	for i := start; i >= 0; i-- {
		if i%3 == 0 && i != 0 {
			buf.WriteString("," + s[start-i:start-i+1])
		} else if i >= 1 {
			buf.WriteString(s[start-i : start-i+1])
		}
	}
	if f > 0 {
		buf.WriteString(s[start:])
	}
	return buf.String()
}
