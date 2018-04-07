PROG=update-cache show-diff init-db write-tree read-tree commit-tree cat-file
# DEP=$(shell git ls-files | grep go)
DEP=gogit.go

all: $(PROG)

install: $(PROG)
	install $(PROG) $(shell go env GOPATH)/bin/

init-db: $(DEP) cmd/init-db/main.go
	go build ./cmd/init-db

update-cache: $(DEP) cmd/update-cache/main.go
	go build ./cmd/update-cache

show-diff: $(DEP) cmd/show-diff/main.go
	go build ./cmd/show-diff

write-tree: $(DEP) cmd/write-tree/main.go
	go build ./cmd/write-tree

read-tree: $(DEP) cmd/read-tree/main.go
	go build ./cmd/read-tree

commit-tree: $(DEP) cmd/commit-tree/main.go
	go build ./cmd/commit-tree

cat-file: $(DEP) cmd/cat-file/main.go
	go build ./cmd/cat-file

clean:
	echo "clean..."
	rm $(PROG)
