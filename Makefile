GO_BUILD=GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w"
RELEASE_VERSION = "0.1.0"
DOCKER_HUB_USER = "iguaziodocker"

all: lint bin image
	@echo Done.

bin: ensure-gopath
	$(GO_BUILD) -o logfwd cmd/logfwd/main.go

image:
	docker build --rm --tag $(DOCKER_HUB_USER)/logfwd:$(RELEASE_VERSION) .

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
check-gopath:
ifndef GOPATH
    $(error GOPATH must be set)
endif
