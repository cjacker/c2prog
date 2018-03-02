all:
	go get
	go build -o c2prog *.go

install:
	go install

.PHONY: all
