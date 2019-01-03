LOGFWD_PATH ?= src/github.com/v3io/logfwd
LOGFWD_TAG ?= latest
LOGFWD_BUILD_COMMAND ?= CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -o $(GOPATH)/bin/logfwd $(GOPATH)/$(LOGFWD_PATH)/cmd/logfwd/main.go

.PHONY: all
all: lint build
	@echo Done.

.PHONY: build
build:
	docker build --tag=logfwd:$(LOGFWD_TAG) .

.PHONY: ensure-gopath bin
bin:
	$(LOGFWD_BUILD_COMMAND)

.PHONY: lint
lint: ensure-gopath
	@echo Installing linters...
	go get -u gopkg.in/alecthomas/gometalinter.v2
	@$(GOPATH)/bin/gometalinter.v2 --install

	@echo Linting...
	@$(GOPATH)/bin/gometalinter.v2 \
		--deadline=300s \
		--disable-all \
		--enable-gc \
		--enable=deadcode \
		--enable=goconst \
		--enable=gofmt \
		--enable=golint \
		--enable=gosimple \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=misspell \
		--enable=staticcheck \
		--enable=unconvert \
		--enable=varcheck \
		--enable=vet \
		--enable=vetshadow \
		--enable=errcheck \
		--exclude="_test.go" \
		--exclude="comment on" \
		--exclude="error should be the last" \
		--exclude="should have comment" \
		./cmd/... ./pkg/...

	@echo Done.

.PHONY: vet
vet:
	go vet ./cmd/...
	go vet ./pkg/...

.PHONY: test
test:
	go test -v ./cmd/...
	go test -v ./pkg/...

.PHONY: ensure-gopath
ensure-gopath:
ifndef GOPATH
	$(error GOPATH must be set)
endif
