ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

SAMPLE_DIR:=$(ROOT_DIR)/sample

ZWALLET=zwallet

.PHONY: $(ZWALLET)

zwallet-test:
	go test -v -tags bn256 ./...

$(ZWALLET):
	go build -v -tags bn256 -o $@
	cp $(ZWALLET) $(SAMPLE_DIR)/zwallet

gomod-download:
	go mod download -json

gomod-clean:
	go clean -i -r -x -modcache  ./...

install: zwallet | zwallet-test

clean: gomod-clean
	@rm -rf $(ROOT_DIR)/$(ZWALLET)

help:
	@echo "Environment: "
	@echo "\tGOPATH=$(GOPATH)"
	@echo "\tGOROOT=$(GOROOT)"
	@echo ""
	@echo "Supported commands:"
	@echo "\tmake help              - display environment and make targets"
	@echo ""
	@echo "Install"
	@echo "\tmake install           - build, test and install the wallet cli"
	@echo "\tmake zwallet           - installs the wallet cli"
	@echo "\tmake zwallet-test      - run zwallet test"
	@echo ""
	@echo "Clean:"
	@echo "\tmake clean             - deletes all build output files"
	@echo "\tmake gomod-download    - download the go modules"
	@echo "\tmake gomod-clean       - clean the go modules"




