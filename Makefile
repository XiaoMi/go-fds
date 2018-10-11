.PHONY: test
test:
	go test ./fds/ -v -cover -coverprofile cover.out

lint:
	gofmt -s -w ./fds/
	goimports -w ./fds/
	golint ./fds/
	go vet ./fds/
