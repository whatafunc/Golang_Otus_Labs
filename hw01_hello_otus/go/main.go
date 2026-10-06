package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

var greetingMsg = "Hello, OTUS!"

func reverseGreeting(origMsg string) string {
	origMsg = reverse.String(origMsg) // reverse the order of symbols in a passed greeeting msg
	return origMsg
}

func main() {
	greetingMsg = reverseGreeting(greetingMsg)
	fmt.Println(greetingMsg) // output the result
}
