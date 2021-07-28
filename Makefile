VERSION := 0.0.1
BINARY := rav

# -s and -w strip debug headers
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build \
	-a -ldflags '-extldflags "-static"' \
	-o ${BINARY}-linux-amd64 .

test:
	REACT_APP_FOO=bar \
	REACT_APP_WAZ=baz \
	  go run . --dir ./test-data
	@echo
	@grep ^REACT_APP_FOO= ./test-data/changed.txt
	@echo "^ should be  =bar"
	@grep ^REACT_APP_WAZ= ./test-data/changed.txt
	@echo "^ should be  =baz"

