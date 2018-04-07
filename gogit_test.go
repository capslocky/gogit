package gogit

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestA(*testing.T) {
	os.RemoveAll(".dircache")

	log.Println("======= make")
	cmd := exec.Command("make")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)

	log.Println("======= init-db")
	cmd = exec.Command("./init-db")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)

	log.Println("======= update-cache")
	cmd = exec.Command("./update-cache", "gogit.go")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)

	log.Println("======= wrtie-tree")
	cmd = exec.Command("./write-tree")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)

	log.Println("======= commit-tree")
	tree := strings.TrimSpace(string(output))
	cmd = exec.Command("./commit-tree", tree)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "init commit\n")
	}()
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)

	log.Println("======= cat-file")
	commit := strings.TrimSpace(string(output))
	cmd = exec.Command("./cat-file", commit)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)
}
