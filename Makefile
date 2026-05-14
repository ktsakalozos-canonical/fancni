.PHONY: build test clean docker-build docker-build-cni docker-build-init helm-template helm-lint e2e rock-build

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

rock-build:
	rockcraft pack

helm-template:
	helm template fancni deploy/helm/fancni/

helm-lint:
	helm lint deploy/helm/fancni/

e2e:
	bash tests/e2e/test-e2e.sh
