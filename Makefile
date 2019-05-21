GOCMD = go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get
FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
FLAGS_LINUX = $(FLAGS_SHARED) GOOS=linux
FLAGS_MAC = $(FLAGS_SHARED) GOOS=darwin
FLAGS_WIN = $(FLAGS_SHARED) GOOS=windows

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

all: install

test: deps install
	$(FLAG_MODULE) go test ./...

test-build: deps build

deps:
	$(FLAG_MODULE) $(GOGET) github.com/spf13/cobra/cobra@7547e83; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/ginkgo/ginkgo@505cc35; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/gomega@ce690c5; \
	$(FLAG_MODULE) $(GOGET) github.com/golang/protobuf/protoc-gen-go@347cf4a; \
  $(FLAG_MODULE) $(GOGET) golang.org/x/tools/cmd/goimports

version:
	echo "package bhlclone" > version.go
	echo "" >> version.go
	echo "const Version = \"$(VERSION)"\" >> version.go
	echo "const Build = \"$(DATE)\"" >> version.go

build: version
	cd bhlclone; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(GOBUILD)

install: version
	cd bhlclone; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(GOINSTALL)

release: version
	cd bhlclone; \
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(GOBUILD); \
	tar zcf /tmp/bhlclone-$(VER)-linux.tar.gz bhlclone; \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(GOBUILD); \
	tar zcf /tmp/bhlclone-$(VER)-mac.tar.gz bhlclone; \
	$(GOCLEAN); \
	$(FLAGS_WIN) $(GOBUILD); \
	zip -9 /tmp/bhlclone-$(VER)-win-64.zip bhlclone.exe; \
	$(GOCLEAN);