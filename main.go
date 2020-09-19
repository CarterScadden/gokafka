package main

import (
	"fmt"
	"strings"
)

func main() {
	messages := [3]string{
		"carter",
		"tanner",
		"michael",
	}

	for _, s := range messages {
		scream(s)
	}
}

func scream(m string) {
	fmt.Printf("%v\n", strings.ToUpper(m))
}
