package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/qingyunha/gogit"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("read-tree <key>")
	}
	typ, content, err := gogit.Sha1ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if typ != "tree" {
		log.Fatalf("need a tree, got a %s", typ)
	}
	for len(content) != 0 {
		i := bytes.Index(content, []byte{0})
		if len(content) < i+20 {
			log.Fatalf("corrupt tree")
		}
		var mode os.FileMode
		var path string
		fmt.Sscanf(string(content[:i]), "%o %s", &mode, &path)
		sha1 := content[i+1 : i+21]
		content = content[i+21:]
		fmt.Printf("%o %s %x\n", mode, path, sha1)
	}
}
