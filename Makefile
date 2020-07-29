GOX := $(GOPATH)/bin/gox

DIR = release
OUT = ${DIR}/{{.OS}}_{{.Arch}}

GOFLAGS = -trimpath -ldflags="-s -w -buildid="
OSARCH ?= "linux/amd64 linux/arm windows/amd64 darwin/amd64"

release: all
	@tar caf release.tar.gz ./release
	@rm -rf release

all: $(GOX)
	gox -osarch ${OSARCH} ${GOFLAGS} -output ${OUT}

clean: 
	rm -rf ${DIR} release.tar.gz

$(GOX):
	go get -u github.com/audibleblink/gox

.PHONY: all clean release
