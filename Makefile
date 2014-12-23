
.PHONY: all
all:
	go build -tags debug

.PHONY: release
release:
	go build  -ldflags '-s'

clean:
	go clean
