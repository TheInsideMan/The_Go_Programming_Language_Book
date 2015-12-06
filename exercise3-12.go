package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(annagram("something", "something"))
}

func annagram(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if strings.Index(b, a[i:i+1]) < 0 {
			return false
		}
	}
	return true
}
