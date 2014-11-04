
.PHONY: all
all:
	go build -tags debug

.PHONY: release
release:
	go build 

clean:
	go clean
