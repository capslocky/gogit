package main

import (
	"fmt"
	"os"

	"github.com/qingyunha/gogit"
)

func main() {
	cache := gogit.NewCache("")
	fmt.Printf("%+v\n", cache)
	for _, path := range os.Args[1:] {
		cache.Add(path)
	}
	fmt.Printf("%+v\n", cache)
	cache.WriteCache()
}
