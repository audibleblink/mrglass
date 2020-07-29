GOX := $(GOPATH)/bin/gox

DIR = release
OUT = ${DIR}/{{.OS}}_{{.Arch}}

GOFLAGS = -trimpath -ldflags="-s -w -buildid="
OSARCH ?= "linux/amd64 linux/arm windows/amd64 darwin/amd64"

all: $(GOX)
	gox -osarch ${OSARCH} ${GOFLAGS} -output ${OUT}

release: all
	@tar caf release.tar.gz ./release
	@rm -rf release

clean: 
	rm -rf ${DIR}

$(GOX):
	go get -u github.com/audibleblink/gox

.PHONY: all clean release
