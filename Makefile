VERSION := 0.0.1
BINARY := rav

# -s and -w strip debug headers
build:
	GCO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
	-a \
	-ldflags="-extldflags=-static -X main.version=${VERSION}" \
	-o ${BINARY}-linux-arm64 .

test:
	REACT_APP_FOO=bar go run . --dir ./test-data
	@echo
	@grep ^REACT_APP_FOO= ./test-data/changed.txt
	@echo "^ should be  =bar"

