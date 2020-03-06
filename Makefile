
GOLINT := $(GOPATH)/bin/golint
GOBINDATA := $(GOPATH)/bin/go-bindata

build: clean checks
	@echo "Compiling project"
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o output/deployer cmd/tenant/main.go
	@echo "Building docker image"
	docker build -t gcr.io/scp-engsrvperfdev-gcp/jaeger-deploy:experimental -f build/Dockerfile .

$(GOBINDATA):
	@echo "Installing go-bindata"
	go get -u github.com/go-bindata/go-bindata/v3/...

$(GOLINT):
	@echo "Installing golint"
	go get -u golang.org/x/lint/golint

clean:
	@echo "Cleaning the output directory"
	rm -rf output

checks: $(GOBINDATA) $(GOLINT)
	go-bindata -o utils/bindata.go -pkg utils -prefix "configs" configs
	golint ./...
	go vet ./...
	go fmt ./...

.PHONY: clean build
