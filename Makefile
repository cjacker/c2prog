all:
	go get
	go build -o c2prog *.go

arm:
	go get
	GOOS=linux GOARCH=arm go build -o c2prog *.go

install:
	go install

.PHONY: all arm
