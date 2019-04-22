export BUILD_PACKAGE = github.com/void616/ethlevel
export BUILD_TAGS = 

OUTPUT_DIR = ./build/bin

# Windows: use Docker to build with CGO
ifeq ($(OS),Windows_NT)
    CGO_DOCKER := 1
else
    CGO_DOCKER := 0
endif
ifndef GOPATH
	$(error GOPATH is undefined)
endif

.PHONY: build

all: build dockerize

build_clean:
	rm -rf ./build/bin/* | true
	mkdir -p ./build/bin | true

build: build_clean maingo

maingo:
	@{ \
	export APP=ethlevel ;\
	export BUILD_APP=main.go ;\
	export BUILD_OS=linux ;\
	export BUILD_ARCH=amd64 ;\
	export BUILD_CGO=1 ;\
	if [ "$$BUILD_CGO" != "" ]; then export BUILD_CGO=1; fi ;\
	if [ "$$BUILD_OS" == "windows" ]; then APPEXT=.exe; fi ;\
	export BUILD_OUTFILE=$${APP}_$${BUILD_ARCH}$${APPEXT} ;\
	export BUILD_OUTDIR=$(OUTPUT_DIR)/$${APP}_$${BUILD_OS} ;\
	\
	if [ "$$BUILD_CGO" != "" ] && [ "$(CGO_DOCKER)" == "1" ]; then \
		echo "Building $$BUILD_APP via Docker" ;\
		docker build -t gobuild_with_docker -f ./scripts/gobuild_with_docker . ;\
		./scripts/gobuild_with_docker.sh ;\
	else \
		echo "Building $$BUILD_APP" ;\
		./scripts/gobuild.sh ;\
	fi ;\
	echo "Done" ;\
	}

dockerize:
	docker build -t ethlevel -f ./build/linux_amd64.dockerfile ./build