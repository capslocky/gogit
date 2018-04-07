package main

import (
	"fmt"
	"os"

	"github.com/qingyunha/gogit"
)

func main() {
	hex := os.Args[1]
	typ, content, _ := gogit.Sha1ReadFile(hex)
	fmt.Printf("type: %s\n===\n%s", typ, content)
}
