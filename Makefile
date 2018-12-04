docker-image:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker: docker-image docker-package
docker-package:
	docker build -t go-users-ms -f Dockerfile.scratch .

