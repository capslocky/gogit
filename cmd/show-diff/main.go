package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/qingyunha/gogit"
)

func main() {
	cache := gogit.NewCache("")
	for name, entry := range cache.Entries {
		fi, err := os.Stat(name)
		if err != nil {
			fmt.Printf("%s: %s\n", name, err)
			continue
		}
		if fi.ModTime().After(entry.Mtime) {
			fmt.Printf("%s: %x\n", name, entry.Sha1)
			showDiff(name, entry)
		} else {
			fmt.Printf("%s: ok\n", name)
		}
	}
	cache.WriteCache()
}

func showDiff(path string, e gogit.CacheEntry) {
	cmd := exec.Command("diff", "-u", "-", path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Print(err)
		return
	}
	go func() {
		defer stdin.Close()
		_, content, err := gogit.Sha1ReadFile(fmt.Sprintf("%0x", e.Sha1))
		if err != nil {
			log.Print(err)
		}
		stdin.Write(content)
	}()
	out, _ := cmd.CombinedOutput()
	fmt.Printf("%s\n", out)

}
