.PHONY: build test clean helm-template helm-lint e2e rock-build

build:
	go build -o _output/bin/fancni ./cmd/fancni/

test:
	go test ./... -v -count=1

clean:
	rm -rf _output/

rock-build:
	rockcraft pack

helm-template:
	helm template fancni deploy/helm/fancni/

helm-lint:
	helm lint deploy/helm/fancni/

e2e:
	bash tests/e2e/test-e2e.sh
