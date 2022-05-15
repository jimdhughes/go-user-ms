.PHONY: proto docker docker-image docker-package

docker-image:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker: docker-image docker-package
docker-package:
	docker build -t go-users-ms -f Dockerfile.scratch .
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto

