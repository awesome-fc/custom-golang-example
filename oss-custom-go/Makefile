build-img:
	docker build -t fc-custom-go-runtime-bot -f build-image/Dockerfile build-image

build: build-img
	docker run --rm -it -v $$(pwd):/tmp -w /tmp/code fc-custom-go-runtime-bot bash -c "go mod vendor && GOARCH=amd64 go build -o /tmp/code/bootstrap /tmp/code"
	chmod +x ./code/bootstrap

deploy: build
	s deploy -y
