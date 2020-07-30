APP = mrglass
DIR = release

FLAGS = -trimpath -ldflags="-s -w -buildid="
PLATFORMS ?= linux windows darwin

os=$(word 1, $@)

all: ${PLATFORMS}

release: ${PLATFORMS}
	@tar caf release.tar.gz ${DIR}
	@rm -rf release

${PLATFORMS}:
	GOOS=${os} go build ${FLAGS} -o ${DIR}/${APP}-${os} 

clean: 
	rm -rf ${DIR}*

.PHONY: clean release
