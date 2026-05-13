.PHONY: build test clean docker-build docker-build-cni docker-build-init helm-template helm-lint

build:
	go build -o _output/bin/fancni ./cmd/fancni/

test:
	go test ./... -v -count=1

clean:
	rm -rf _output/

docker-build-cni:
	docker build -t fancni:latest -f deploy/docker/Dockerfile.cni .

docker-build-init:
	docker build -t fancni-init:latest -f deploy/docker/Dockerfile.init .

docker-build: docker-build-cni docker-build-init

helm-template:
	helm template fancni deploy/helm/fancni/

helm-lint:
	helm lint deploy/helm/fancni/
