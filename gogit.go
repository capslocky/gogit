package gogit

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	DbEnvironment        = "SHA1_FILE_DIRECTORY"
	DefaultDbEnvironment = ".dircache/objects"

	CacheSignature = 0x44495243
)

var (
	ActiveCache Cache
)

type CacheEntry struct {
	Mtime time.Time
	Mode  os.FileMode
	Size  int64
	Name  string
	Sha1  [20]byte
}

type Cache struct {
	Signature uint
	Version   int
	Sha1      [20]byte

	Entries map[string]CacheEntry
}

func NewCache(path string) *Cache {
	if path == "" {
		path = ".dircache/index"
	}
	cache := &Cache{}
	cache.Entries = make(map[string]CacheEntry)
	cache.readCache(path)
	if cache.Signature == 0 {
		cache.Signature = CacheSignature
		cache.Version = 1
	}
	return cache
}

func (c *Cache) WriteCache() {
	f, err := os.OpenFile(".dircache/index.lock", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		log.Fatalf("unable to crate new cachefile: %s", err)
	}
	defer f.Close()

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(c)
	if err != nil {
		log.Fatalf("encode error: %s", err)
	}

	_, err = f.Write(buffer.Bytes())
	if err != nil {
		log.Fatalf("write error: %s", err)
	}
	if err := os.Rename(".dircache/index.lock", ".dircache/index"); err != nil {
		log.Fatalf("rename error: %s", err)
	}
}

func (c *Cache) readCache(path string) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatalf("%s", err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	err = dec.Decode(c)
	if err != nil {
		log.Fatalf("decode error: %s", err)
	}
}

func (c *Cache) Add(path string) error {
	f, err := os.Open(path)
	if err != nil {
		delete(c.Entries, path)
		log.Print(err)
		return nil
	}
	defer f.Close()

	var e CacheEntry
	fi, err := f.Stat()
	if err != nil {
		log.Print(err)
		return err
	}

	var buf bytes.Buffer
	meta := []byte(fmt.Sprintf("blob %d\000", fi.Size()))
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	buf.Write(meta)
	buf.Write(b)
	sha1 := Sha1WriteFile(buf.Bytes())

	e.Name = path
	e.Size = fi.Size()
	e.Mtime = fi.ModTime()
	e.Mode = fi.Mode()
	e.Sha1 = sha1
	c.Entries[e.Name] = e
	return nil
}

func Sha1WriteFile(data []byte) [20]byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data)
	w.Close()
	sha1 := sha1.Sum(buf.Bytes())
	fname := fmt.Sprintf("%s/%02x/%02x", DefaultDbEnvironment, sha1[0], sha1[1:])
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalf("sha1WriteFile error: %s", err)
	}
	defer f.Close()
	if _, err := f.Write(buf.Bytes()); err != nil {
		log.Fatalf("sha1WriteFile error: %s", err)
	}
	return sha1
}

func Sha1ReadFile(hex string) (typ string, content []byte, err error) {
	fname := fmt.Sprintf("%s/%s/%s", DefaultDbEnvironment, hex[0:2], hex[2:])
	f, err := os.Open(fname)
	if err != nil {
		log.Print(err)
		return "", nil, err
	}
	var buf bytes.Buffer
	r, err := zlib.NewReader(f)
	if err != nil {
		log.Print(err)
		return "", nil, err
	}
	io.Copy(&buf, r)
	size := 0
	fmt.Fscanf(&buf, "%s %d", &typ, &size)
	content = buf.Bytes()
	i := bytes.Index(content, []byte{0})
	if size != len(content[i+1:]) {
		log.Fatalf("invalid data %d %d", size, len(content[i:]))
	}
	return typ, content[i+1:], err
}
