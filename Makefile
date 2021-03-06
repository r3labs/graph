install:
	go install -v

build:
	go build -v ./...

lint:
	go vet ./...

test:
	go test -v ./... --cover

deps: dev-deps
	go get -u github.com/r3labs/diff

dev-deps:
	go get -u github.com/smartystreets/goconvey/convey

clean:
	go clean

dist-clean:
	rm -rf pkg src bin
