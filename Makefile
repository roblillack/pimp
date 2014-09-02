BASE := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT := $(notdir $(BASE:/=))
GOPATH := $(BASE:/=)/ext:$(BASE:/=)
ENV := env GOPATH=${GOPATH} PATH=${BASE}ext/bin:$$PATH

build:
	@${ENV} go build

clean:
	rm -rf bin/ ext/bin/ ext/pkg/ ./${PROJECT}
