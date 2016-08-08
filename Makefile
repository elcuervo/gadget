all:
	@mkdir -p bin/
	@echo "==> Installing dependencies"
	@go get -d -v ./...

format:
	@echo "==> Formating project ..."
	go fmt ./...

build:
	@echo "==> Building ..."
	@go build -o bin/gadget .

dist:
	@echo "==> Creating executables ..."
	@./dist.sh

container:
	docker build -f Dockerfile -t elcuervo/gadget .

create:
	#docker rmi -f gadget-builder
	docker build -t gadget-builder -f Dockerfile.build .
	docker run gadget-builder

clean:
	@rm gadget

test:
	@echo "==> Testing gadget ..."
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs -n1 go get -d
	go test ./...

PHONY: all format test
