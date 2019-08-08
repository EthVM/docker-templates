.SILENT :
.PHONY : deps install dist-clean docker

TAG:=0.1.0
LDFLAGS:=-X main.VERSION=$(TAG)

all: install

deps:
	echo "Ensuring dependencies..."
	dep ensure

install:
	echo "Building docker-templates"
	go install -ldflags "$(LDFLAGS)"

docker:
	echo "Building docker image..."
	docker build -t ethvm/docker-templates:"$(TAG)" .
	docker build -t ethvm/docker-templates:latest .

docker-push:
	echo "Pushing docker image to registry..."
	docker push ethvm/docker-templates:"$(TAG)"
	docker push ethvm/docker-templates:latest

dist-clean:
	rm -rf dist
	rm -f docker-templates-*.tar.gz

dist: dist-clean deps
	mkdir -p dist/alpine-linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/amd64/docker-templates
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux/amd64/docker-templates
	mkdir -p dist/linux/386 && GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/linux/386/docker-templates
	mkdir -p dist/linux/armel && GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o dist/linux/armel/docker-templates
	mkdir -p dist/linux/armhf && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -o dist/linux/armhf/docker-templates
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/docker-templates

release: dist
	tar -cvzf dockerize-alpine-linux-amd64-$(TAG).tar.gz -C dist/alpine-linux/amd64 docker-templates
	tar -cvzf dockerize-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 docker-templates
	tar -cvzf dockerize-linux-386-$(TAG).tar.gz -C dist/linux/386 docker-templates
	tar -cvzf dockerize-linux-armel-$(TAG).tar.gz -C dist/linux/armel docker-templates
	tar -cvzf dockerize-linux-armhf-$(TAG).tar.gz -C dist/linux/armhf docker-templates
	tar -cvzf dockerize-darwin-amd64-$(TAG).tar.gz -C dist/darwin/amd64 docker-templates
