.PHONY: build
build:
	pkger
	go build .

.PHONY: install
install:
	pkger
	go install
