default: scooter-spotter
all: scooter-spotter

scooter-spotter: fetch-dependencies main.go $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -ldflags '-s -w' -o $@ .

clean:
	rm -f bin/upx
	rm -f scooter-spotter

.PHONY: default all clean

fetch-golangci-lint:
	# brew install golangci-lint
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.14.0

fetch-upx: bin/upx
	# brew install upx
	cd tmp && \
		wget https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz && \
		tar -xf upx-3.95-amd64_linux.tar.xz
	cp tmp/upx-3.95-amd64_linux/upx bin/
	chmod +x bin/upx

fetch-dependencies:
	@go mod vendor -v

build-compressed: scooter-spotter
	@bin/upx scooter-spotter

.PHONY: build-compressed fetch-dependencies

docker-build: Dockerfile
	docker build -t kakkoyun/scooter-spotter:latest .

docker-push: docker-build
	docker push kakkoyun/scooter-spotter:latest

.PHONY: docker-build docker-push

check:
	golangci-lint run -v --enable-all -D gochecknoglobals

.PHONY: check
