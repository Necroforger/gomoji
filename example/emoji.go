package main

/*
Usage:
	go run emoji.go ":computer: :heart:"
*/

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Necroforger/gomoji"
)

func main() {
	flag.Parse()
	fmt.Println(gomoji.Format(strings.Join(flag.Args(), "")))
}
