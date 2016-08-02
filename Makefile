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
	docker run

clean:
	@rm gadget

test:
	@echo "==> Testing gadget ..."
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs -n1 go get -d
	go test ./...

PHONY: all format test
