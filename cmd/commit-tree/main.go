package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/qingyunha/gogit"
)

func main() {
	var parents []string
	if len(os.Args) < 2 {
		fmt.Printf("usage:\n  commit-tree <sha1> [-p <sha1>]* < changelog\n")
		os.Exit(1)
	}
	tree := os.Args[1]
	for i := 2; i < len(os.Args); i += 2 {
		if os.Args[i] != "-p" {
			fmt.Printf("usage:\n  commit-tree <sha1> [-p <sha1>]* < changelog\n")
			os.Exit(1)
		}
		parents = append(parents, os.Args[i+1])
	}
	var buffer bytes.Buffer
	buffer.Write([]byte(fmt.Sprintf("tree %s\n", tree)))
	for _, p := range parents {
		buffer.Write([]byte(fmt.Sprintf("parent %s\n", p)))
	}
	realdate := time.Now().Format(time.UnixDate)
	u, _ := user.Current()
	realname := u.Username
	realemail := fmt.Sprintf("%s@%s", realname, "localhost")

	name := os.Getenv("COMMITTER_NAME")
	if name == "" {
		name = realname
	}
	email := os.Getenv("COMMITTER_EMAIL")
	if email == "" {
		email = realemail
	}
	date := os.Getenv("COMMITTER_DATE")
	if date == "" {
		date = realdate
	}
	buffer.Write([]byte(fmt.Sprintf("author %s <%s> %s\n", name, email, date)))
	buffer.Write([]byte(fmt.Sprintf("committer %s <%s> %s\n\n", realname, realemail, realdate)))
	reader := bufio.NewReader(os.Stdin)
	comment, _ := reader.ReadBytes('\n')
	buffer.Write(comment)

	meta := []byte(fmt.Sprintf("commit %d\000", buffer.Len()))
	sha1 := gogit.Sha1WriteFile(bytes.Join([][]byte{meta, buffer.Bytes()}, []byte{}))
	fmt.Printf("%0x\n", sha1)
}
