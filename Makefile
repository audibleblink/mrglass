.PHONY: all clean

DIR = release
OUT = ${DIR}/{{.OS}}_{{.Arch}}
GOFLAGS = -trimpath -ldflags="-s -w -buildid="
OSARCH ?= "linux/amd64 linux/arm windows/amd64 darwin/amd64"

all:
	gox -osarch ${OSARCH} ${GOFLAGS} -output ${OUT}

clean: 
	rm -rf ${DIR}

