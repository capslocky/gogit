package main

import (
	"bytes"
	"fmt"

	"github.com/qingyunha/gogit"
)

func main() {
	cache := gogit.NewCache("")
	var buffer bytes.Buffer
	for name, entry := range cache.Entries {
		buffer.Write([]byte(fmt.Sprintf("%o %s\000", entry.Mode, name)))
		buffer.Write(entry.Sha1[:])
	}
	meta := []byte(fmt.Sprintf("tree %d\000", buffer.Len()))
	sha1 := gogit.Sha1WriteFile(bytes.Join([][]byte{meta, buffer.Bytes()}, []byte{}))
	fmt.Printf("%0x\n", sha1)
}
