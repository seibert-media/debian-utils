install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/create_debian_package/create_debian_package.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/extract_zip/extract_zip.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/update_apt_source_list/update_apt_source_list.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/update_available_apt_source_list/update_available_apt_source_list.go
test:
	GO15VENDOREXPERIMENT=1 go test `glide novendor`
check:
	golint ./...
	errcheck -ignore '(Close|Write)' ./...
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor
