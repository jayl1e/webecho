run:
	go run -trimpath main.go
docker:
	docker buildx build --platform=linux/amd64,linux/arm64 -t jayl1e/webecho --push .
docker-slim:
	docker buildx build --platform=linux/amd64,linux/arm64 -t jayl1e/webecho:slim --push . -f slim.Dockerfile
build:
	go build -trimpath -o webecho .
build-dbg:
	go build -trimpath -gcflags="-N -l" -o webecho . 
build-slim:
	CGO_ENABLED=0 go build -trimpath -o webecho .