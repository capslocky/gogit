package main

import (
	"fmt"
	"log"
	"os"

	"github.com/qingyunha/gogit"
)

func main() {
	var sha1dir string
	sha1dir = os.Getenv(gogit.DbEnvironment)
	if sha1dir != "" {
		fi, err := os.Stat(sha1dir)
		if err == nil && fi.Mode().IsDir() {
			return
		}
		log.Printf("DB_ENVIRONMENT set to bad directory %s", sha1dir)
	}

	log.Printf("defaulting to private storage area")
	sha1dir = gogit.DefaultDbEnvironment
	if err := os.MkdirAll(sha1dir, 0700); err != nil {
		log.Fatalf("%s", err)
	}

	for i := 0; i < 256; i++ {
		if err := os.Mkdir(fmt.Sprintf("%s/%02x", sha1dir, i), 0700); err != nil {
			if !os.IsExist(err) {
				log.Fatalf("%s", err)
			}
		}
	}

}
